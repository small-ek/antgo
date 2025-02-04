package acron

import (
	"errors"
	"github.com/robfig/cron/v3"
	"sync"
)

// Crontab crontab 管理器 / Crontab manager
type Crontab struct {
	inner   *cron.Cron              // Cron 实例 / Cron instance
	ids     map[string]cron.EntryID // 存储任务ID与Cron Entry ID的映射 / Store task IDs and Cron Entry IDs
	mu      sync.RWMutex            // 读写锁 / Read-write mutex
	running bool                    // 记录是否正在运行 / Flag to track whether cron engine is running
}

// New 创建新的 crontab / Create a new crontab
func New() *Crontab {
	// 设定支持秒级的 Cron 表达式解析器 / Set up a parser for second-level cron expressions
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

	return &Crontab{
		inner: cron.New(cron.WithParser(secondParser), cron.WithChain()), // 初始化 Cron 实例 / Initialize Cron instance
		ids:   make(map[string]cron.EntryID),                             // 初始化任务ID映射表 / Initialize task ID map
	}
}

// IDs 返回所有有效的 cron 任务ID / Return all valid cron task IDs
func (c *Crontab) IDs() []string {
	c.mu.RLock()         // 只读锁 / Acquire read lock
	defer c.mu.RUnlock() // 释放读锁 / Release read lock

	validIDs := make([]string, 0, len(c.ids)) // 存储有效的ID / Store valid IDs
	invalidIDs := make([]string, 0)           // 存储无效的ID / Store invalid IDs

	// 遍历所有任务ID，检查其有效性 / Check the validity of each task ID
	for sid, eid := range c.ids {
		if e := c.inner.Entry(eid); e.ID != eid { // 如果任务无效 / If the task is invalid
			invalidIDs = append(invalidIDs, sid) // 添加到无效ID列表 / Add to invalid IDs
			continue
		}
		validIDs = append(validIDs, sid) // 添加到有效ID列表 / Add to valid IDs
	}

	// 清理无效ID / Clean up invalid IDs
	if len(invalidIDs) > 0 {
		c.mu.Lock()         // 升级为写锁 / Upgrade to write lock
		defer c.mu.Unlock() // 释放写锁 / Release write lock
		for _, id := range invalidIDs {
			delete(c.ids, id) // 删除无效ID / Delete invalid IDs
		}
	}

	return validIDs // 返回所有有效的ID / Return all valid IDs
}

// Start 启动 cron 引擎 / Start the cron engine
func (c *Crontab) Start() {
	c.inner.Start()  // 启动 Cron / Start Cron
	c.running = true // 标记为运行中 / Set the running flag to true
}

// Stop 停止 cron 引擎 / Stop the cron engine
func (c *Crontab) Stop() {
	c.inner.Stop()    // 停止 Cron / Stop Cron
	c.running = false // 标记为停止 / Set the running flag to false
}

// IsRunning 检查 Cron 引擎是否正在运行 / Check if the cron engine is running
func (c *Crontab) IsRunning() bool {
	return c.running // 返回运行状态 / Return the running status
}

// DelByID 根据ID删除一个 cron 任务 / Remove a cron task by its ID
func (c *Crontab) DelByID(id string) {
	c.mu.Lock()         // 获取写锁 / Acquire write lock
	defer c.mu.Unlock() // 释放写锁 / Release write lock

	eid, ok := c.ids[id] // 获取任务ID对应的Cron Entry ID / Get the Cron Entry ID for the task ID
	if !ok {
		return // 如果ID不存在，直接返回 / Return if ID does not exist
	}
	c.inner.Remove(eid) // 移除任务 / Remove the task
	delete(c.ids, id)   // 删除任务ID映射 / Delete the task ID mapping
}

// AddByID 根据ID添加一个 cron 任务 / Add a cron task by its ID
// id 为唯一标识符 / id is unique
// spec 为 cron 表达式 / spec is the cron expression
func (c *Crontab) AddByID(id string, spec string, cmd cron.Job) error {
	c.mu.Lock()         // 获取写锁 / Acquire write lock
	defer c.mu.Unlock() // 释放写锁 / Release write lock

	if _, ok := c.ids[id]; ok { // 检查任务ID是否已存在 / Check if task ID already exists
		return errors.New("crontab id exists") // 如果已存在，返回错误 / Return error if ID already exists
	}
	eid, err := c.inner.AddJob(spec, cmd) // 添加Cron任务 / Add the cron task
	if err != nil {
		return err // 返回错误 / Return error
	}
	c.ids[id] = eid // 保存任务ID和Cron Entry ID映射 / Save the task ID and Cron Entry ID mapping
	return nil      // 返回成功 / Return success
}

// AddByFunc 根据ID和函数添加一个 cron 任务 / Add a cron task using a function
func (c *Crontab) AddByFunc(id string, spec string, f func()) error {
	c.mu.Lock()         // 获取写锁 / Acquire write lock
	defer c.mu.Unlock() // 释放写锁 / Release write lock

	if _, ok := c.ids[id]; ok { // 检查任务ID是否已存在 / Check if task ID already exists
		return errors.New("crontab id exists") // 如果已存在，返回错误 / Return error if ID already exists
	}
	eid, err := c.inner.AddFunc(spec, f) // 添加Cron任务 / Add the cron task
	if err != nil {
		return err // 返回错误 / Return error
	}
	c.ids[id] = eid // 保存任务ID和Cron Entry ID映射 / Save the task ID and Cron Entry ID mapping
	return nil      // 返回成功 / Return success
}

// IsExists 检查 cron 任务是否存在 / Check if a cron task exists by job ID
func (c *Crontab) IsExists(jid string) bool {
	c.mu.RLock()           // 获取读锁 / Acquire read lock
	defer c.mu.RUnlock()   // 释放读锁 / Release read lock
	_, exist := c.ids[jid] // 检查任务ID是否存在 / Check if the task ID exists
	return exist           // 返回结果 / Return the result
}
