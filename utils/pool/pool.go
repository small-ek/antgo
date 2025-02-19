package pool

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

var (
	once    sync.Once // 确保只初始化一次
	JobPool *ants.Pool
)

// New 初始化 Goroutine 池
// New initializes the Goroutine pool
func New(size, queue int) {
	// 检查 Goroutine 池的大小是否大于0
	// Ensure the size is positive
	if size <= 0 {
		panic("size must be positive") // 如果大小无效，则抛出异常
	}

	// 如果队列大小无效，则使用默认值
	// Set a default queue size if the given queue size is invalid
	if queue <= size {
		queue = size * 10 // 默认队列大小是 Goroutine 池大小的 10 倍
	}

	// 使用 sync.Once 确保初始化池只进行一次
	// Use sync.Once to ensure the pool is initialized only once
	once.Do(func() {
		// 创建 Goroutine 池实例
		// Create a new Goroutine pool instance
		var err error
		JobPool, err = ants.NewPool(size, ants.WithPreAlloc(true), ants.WithMaxBlockingTasks(queue))
		if err != nil {
			panic("failed to create pool: " + err.Error()) // 如果池创建失败，抛出异常
		}
	})
}
