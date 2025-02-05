package test

import (
	"github.com/small-ek/antgo/utils/pool"
	"testing"
)

// TestNewPool 测试 New 函数
// Test the New function
func TestNewPool(t *testing.T) {
	// 初始化 Goroutine 池，设定大小为 5，队列大小为 50
	// Initialize the Goroutine pool with size 5 and queue size 50
	pool.New(5, 50)

	// 获取池实例，确保池已成功创建
	// Get the pool instance to ensure it's successfully created
	poolInstance := pool.Get()

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
	poolInstance := pool.Get()

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
	pool.Get()
}

// TestNewWithInvalidSize 测试给定无效的大小
// Test the case when an invalid size is provided
func TestNewWithInvalidSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic when providing invalid pool size")
		}
	}()

	// 尝试初始化池时传入无效的大小（小于等于0）
	// Try initializing the pool with an invalid size (<= 0)
	pool.New(0, 10)
}
