package pool

import (
	"context"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
)

const (
	defaultQueueMultiplier = 10
)

var (
	jobPool          *ants.Pool
	jobPoolLock      sync.RWMutex
	userPanicHandler func(ctx context.Context, r interface{}, stack []byte)
)

// JobPool 暴露全局 Goroutine 池实例（只读）
var JobPool *ants.Pool

// New 初始化 Goroutine 池（线程安全）
// size: Goroutine 数量, queueSize: 最大阻塞任务数
// 重复初始化返回错误
func New(size, queueSize int) error {
	if size <= 0 {
		return errors.New("pool size must be positive")
	}

	jobPoolLock.Lock()
	defer jobPoolLock.Unlock()

	if jobPool != nil {
		return errors.New("pool already initialized")
	}

	if queueSize < size {
		queueSize = size * defaultQueueMultiplier
	}

	pool, err := ants.NewPool(
		size,
		ants.WithPreAlloc(true),
		ants.WithNonblocking(false),
		ants.WithMaxBlockingTasks(queueSize),
	)
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}

	jobPool = pool
	JobPool = pool
	return nil
}

// Submit 提交任务到 Goroutine 池（线程安全）
func Submit(task func()) error {
	jobPoolLock.RLock()
	defer jobPoolLock.RUnlock()

	if jobPool == nil {
		return errors.New("pool not initialized")
	}
	return jobPool.Submit(task)
}

// SubmitWithCtx 提交带上下文的任务，自动捕获 panic
func SubmitWithCtx(ctx context.Context, task func(ctx context.Context)) error {
	wrappedTask := func() {
		defer handlePanic(ctx)
		task(ctx)
	}

	jobPoolLock.RLock()
	defer jobPoolLock.RUnlock()

	if jobPool == nil {
		return errors.New("pool not initialized")
	}
	return jobPool.Submit(wrappedTask)
}

// Release 释放资源（线程安全）
func Release() {
	jobPoolLock.Lock()
	defer jobPoolLock.Unlock()

	if jobPool != nil {
		jobPool.Release()
		jobPool = nil
		JobPool = nil
	}
}

// OnPanic 设置用户自定义 panic 处理函数（可用于上报监控、报警等）
func OnPanic(handler func(ctx context.Context, r interface{}, stack []byte)) {
	userPanicHandler = handler
}

func handlePanic(ctx context.Context) {
	if r := recover(); r != nil {
		stack := debug.Stack()
		if alog.Write != nil {
			// 默认日志记录
			alog.Error(ctx, "goroutine panic recovered",
				zap.Any("recover", r),
				zap.Any("stack", debug.Stack()),
			)
		}

		// 用户自定义处理器
		if userPanicHandler != nil {
			userPanicHandler(ctx, r, stack)
		}
	}
}
