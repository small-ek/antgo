package pool

import (
	"github.com/panjf2000/ants/v2"
	"sync"
	"sync/atomic"
)

var (
	pool atomic.Value // 使用 atomic.Value 来存储 Goroutine 池实例，简化代码
	once sync.Once    // 确保只初始化一次
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
	if queue <= 0 {
		queue = size * 10 // 默认队列大小是 Goroutine 池大小的 10 倍
	}

	// 使用 sync.Once 确保初始化池只进行一次
	// Use sync.Once to ensure the pool is initialized only once
	once.Do(func() {
		// 创建 Goroutine 池实例
		// Create a new Goroutine pool instance
		p, err := ants.NewPool(size, ants.WithPreAlloc(true), ants.WithMaxBlockingTasks(queue))
		if err != nil {
			panic("failed to create pool: " + err.Error()) // 如果池创建失败，抛出异常
		}
		pool.Store(p) // 使用 atomic.Value.Store 原子存储池实例
	})
}

// Get 获取 Goroutine 池实例
// Get retrieves the Goroutine pool instance
func Get() *ants.Pool {
	// 原子加载池实例
	// Atomically load the pool instance
	p := pool.Load()
	if p != nil {
		return p.(*ants.Pool) // 返回池实例，进行类型断言
	}

	// 如果池未初始化，抛出异常
	// Panic if the pool is not initialized
	panic("pool uninitialized: call New() first")
}
