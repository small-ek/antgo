package acron

import (
	"context"
	"go.uber.org/zap"
	"sync/atomic"
	"testing"
	"time"
)

const (
	// 调度表达式与超时时间常量
	simplePeriod   = "@every 200ms"
	rescheduleOld  = "@every 500ms"
	rescheduleNew  = "@every 100ms"
	healthInterval = 10 * time.Millisecond

	// 最大等待时长，用于 waitFor 辅助函数
	maxWait = 1 * time.Second
)

// waitFor 在 timeout 期限内轮询 count 指针是否达到期望值 n
func waitFor(count *int32, n int32, timeout time.Duration) bool {
	deadline := time.After(timeout)
	for {
		select {
		case <-deadline:
			return false
		default:
			if atomic.LoadInt32(count) >= n {
				return true
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// mockJob 用于测试 ContextJob：执行后睡眠 execDelay 并返回 returnErr
type mockJob struct {
	execDelay time.Duration
	returnErr error
	called    *int32
}

func (m *mockJob) Run(ctx context.Context) error {
	atomic.AddInt32(m.called, 1)
	select {
	case <-time.After(m.execDelay):
		return m.returnErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

// panicJob 用于测试 Panic 恢复，直接 panic
type panicJob struct {
	called *int32
}

func (p *panicJob) Run(ctx context.Context) error {
	atomic.AddInt32(p.called, 1)
	panic("test panic")
}

func TestAddFuncAndRemove(t *testing.T) {
	logger := zap.NewNop()
	ct := New(context.Background(), logger, time.Second)
	// 确保测试结束时优雅停止并等待完成
	defer func() {
		done := ct.Stop()
		<-done.Done()
	}()

	var called int32
	if err := ct.AddFunc("simple", simplePeriod, func() {
		atomic.AddInt32(&called, 1)
	}); err != nil {
		t.Fatalf("AddFunc 失败: %v", err)
	}

	ct.Start()

	// 等待至少执行一次
	if !waitFor(&called, 1, maxWait) {
		t.Errorf("期望 simple 调用至少 1 次, 实际: %d", called)
	}

	// 移除任务后，计数重置并验证不再执行
	ct.Remove("simple")
	atomic.StoreInt32(&called, 0)
	time.Sleep(250 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Errorf("Remove 后不应再执行, 但执行了 %d 次", called)
	}
}

func TestAddJobWithTimeout_Success(t *testing.T) {
	logger := zap.NewNop()
	ct := New(context.Background(), logger, 500*time.Millisecond)
	defer func() {
		done := ct.Stop()
		<-done.Done()
	}()

	var callCount int32
	job := &mockJob{execDelay: 100 * time.Millisecond, returnErr: nil, called: &callCount}

	if err := ct.AddJobWithTimeout("jobSuccess", simplePeriod, job, 300*time.Millisecond); err != nil {
		t.Fatalf("AddJobWithTimeout 失败: %v", err)
	}

	ct.Start()

	if !waitFor(&callCount, 1, maxWait) {
		t.Errorf("预期 jobSuccess 至少执行一次, 实际: %d", callCount)
	}
}

func TestAddJobWithTimeout_Timeout(t *testing.T) {
	logger := zap.NewNop()
	ct := New(context.Background(), logger, 100*time.Millisecond)
	defer func() {
		done := ct.Stop()
		<-done.Done()
	}()

	var callCount int32
	// execDelay 超过 timeout，确保超时分支
	job := &mockJob{execDelay: 500 * time.Millisecond, returnErr: nil, called: &callCount}

	if err := ct.AddJobWithTimeout("jobTimeout", simplePeriod, job, 100*time.Millisecond); err != nil {
		t.Fatalf("AddJobWithTimeout 设置超时失败: %v", err)
	}

	ct.Start()

	if !waitFor(&callCount, 1, maxWait) {
		t.Errorf("预期 jobTimeout 执行一次 (超时后返回), 实际: %d", callCount)
	}
}

func TestReschedule(t *testing.T) {
	logger := zap.NewNop()
	ct := New(context.Background(), logger, time.Second)
	defer func() {
		done := ct.Stop()
		<-done.Done()
	}()

	var called int32
	if err := ct.AddFunc("res", rescheduleOld, func() {
		atomic.AddInt32(&called, 1)
	}); err != nil {
		t.Fatalf("AddFunc 失败: %v", err)
	}

	ct.Start() // 先启动调度器，确保任务开始计时

	// 等待旧规则至少执行一次（500ms周期）
	if !waitFor(&called, 1, maxWait) {
		t.Fatalf("等待初始执行失败")
	}

	// 重新调度到更短间隔
	if err := ct.Reschedule("res", rescheduleNew); err != nil {
		t.Fatalf("Reschedule 失败: %v", err)
	}

	// 重置计数器，专注测量重新调度后的执行
	atomic.StoreInt32(&called, 0)

	// 等待新规则执行（100ms周期，应快速触发）
	if !waitFor(&called, 2, maxWait) {
		t.Errorf("预期 reschedule 后至少执行 2 次, 实际: %d", called)
	}
}

func TestPanicRecovery(t *testing.T) {
	logger := zap.NewNop()
	ct := New(context.Background(), logger, time.Second)
	defer func() {
		done := ct.Stop()
		<-done.Done()
	}()

	var called int32
	job := &panicJob{called: &called}

	if err := ct.AddJobWithTimeout("panicJob", simplePeriod, job, time.Second); err != nil {
		t.Fatalf("AddJobWithTimeout 失败: %v", err)
	}

	ct.Start()

	// 等待至少执行一次并 panic 恢复
	if !waitFor(&called, 1, maxWait) {
		t.Errorf("预期 panicJob 至少执行一次, 实际: %d", called)
	}
}

func TestContextCancelStopsJob(t *testing.T) {
	logger := zap.NewNop()
	// 创建一个可取消的根上下文
	rootCtx, cancel := context.WithCancel(context.Background())
	ct := New(rootCtx, logger, time.Second)

	var called int32
	job := &mockJob{execDelay: 500 * time.Millisecond, returnErr: nil, called: &called}
	if err := ct.AddJobWithTimeout("cancelJob", simplePeriod, job, time.Second); err != nil {
		t.Fatalf("AddJobWithTimeout 失败: %v", err)
	}

	ct.Start()

	// 等待第一次调用
	if !waitFor(&called, 1, maxWait) {
		t.Fatalf("预期 cancelJob 执行一次, 实际: %d", called)
	}

	// 取消根上下文并等待调度器停止
	cancel()
	done := ct.Stop()
	select {
	case <-done.Done():
		// 确保没有 panic
	case <-time.After(500 * time.Millisecond):
		t.Fatal("调度器在根上下文取消后未能及时停止")
	}
}
