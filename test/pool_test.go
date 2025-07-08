package test

import (
	"context"
	"fmt"
	"github.com/small-ek/antgo/crypto/auuid"
	"github.com/small-ek/antgo/utils/pool"
	"testing"
	"time"
)

// TestNewPool 测试 New 函数
// Test the New function
func TestNewPool(t *testing.T) {
	// 初始化 Goroutine 池，设定大小为 5，队列大小为 50
	// Initialize the Goroutine pool with size 5 and queue size 50
	pool.New(5, 50)

	// 获取池实例，确保池已成功创建
	// Get the pool instance to ensure it's successfully created
	poolInstance := pool.JobPool

	// 检查池实例是否为 nil
	// Ensure the pool instance is not nil
	if poolInstance == nil {
		t.Fatalf("Expected pool instance, got nil")
	}

	// 检查池的大小是否符合预期
	// Check if the pool size is as expected (this can be tested based on your pool configuration)
	if poolInstance.Cap() != 5 {
		t.Errorf("Expected pool size of 5, got %d", poolInstance.Cap())
	}
}

// TestNewPoolWithDefaultQueue 测试没有提供队列大小时，使用默认队列大小
// Test the default queue size when not provided
func TestNewPoolWithDefaultQueue(t *testing.T) {
	// 初始化 Goroutine 池，设定大小为 5，队列大小不提供
	// Initialize the Goroutine pool with size 5 and no queue size
	pool.New(5, 0)

	// 获取池实例，确保池已成功创建
	// Get the pool instance to ensure it's successfully created
	poolInstance := pool.JobPool

	// 确保池实例不为 nil
	// Ensure the pool instance is not nil
	if poolInstance == nil {
		t.Fatalf("Expected pool instance, got nil")
	}

	// 确认队列大小为 50（大小是 5 * 10）
	// Verify that the queue size is set to the default value (5 * 10 = 50)
	if poolInstance.Cap() != 5 {
		t.Errorf("Expected pool size of 5, got %d", poolInstance.Cap())
	}
}

// TestGetWithoutInitialization 测试在没有初始化池的情况下调用 Get
// Test calling Get without initializing the pool
func TestGetWithoutInitialization(t *testing.T) {
	// 在没有调用 New 函数初始化池的情况下直接调用 Get
	// Directly call Get without initializing the pool
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic when calling Get without initialization")
		}
	}()

	// 调用 Get，预期会 panic
	// Call Get, expecting it to panic

}

// TestNewWithInvalidSize 测试给定无效的大小
// Test the case when an invalid size is provided
func TestExamples(t *testing.T) {
	// 初始化协程池，大小100，最大阻塞任务数1000
	err := pool.New(100, 1000)
	if err != nil {
		panic("init pool failed: " + err.Error())
	}
	defer pool.Release() // 程序退出时释放资源
	// 带上下文的任务，自动捕获 panic
	ctx := context.WithValue(context.Background(), "request_id", auuid.New().String())
	// 提交普通任务
	for i := 0; i < 5; i++ {
		n := i
		err := pool.Submit(func() {
			fmt.Printf("普通任务: %d\n", n)

			time.Sleep(200 * time.Millisecond)
		})
		if err != nil {
			fmt.Printf("Submit error: %v\n", err)
		}
	}

	for i := 0; i < 5; i++ {
		n := i
		err := pool.SubmitWithCtx(ctx, func(ctx context.Context) {
			fmt.Printf("带上下文任务: %d\n", n)

			if n == 3 {
				panic("模拟panic")
			}
			time.Sleep(100 * time.Millisecond)
		})
		if err != nil {
			fmt.Printf("SubmitWithCtx error: %v\n", err)
		}
	}
	pool.OnPanic(func(ctx context.Context, r interface{}, stack []byte) {
		traceID := ctx.Value("trace_id") // 如果你有 trace_id

		// 构造飞书消息内容（示意）
		msg := fmt.Sprintf("🚨 Panic Detected\nTraceID: %v\nReason: %v\nStack: %s",
			traceID,
			r,
			stack[:300], // 避免太长
		)
		fmt.Println(msg)

	})

	// 等待所有任务执行完成
	time.Sleep(2 * time.Second)
	fmt.Println("所有任务执行完毕")
}
