package acron

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"              // 任务调度库 cron scheduler library
	"github.com/small-ek/antgo/crypto/auuid" // 生成唯一请求 ID UUID generator for request IDs
	"github.com/small-ek/antgo/net/httpx/middleware/agin"
	"go.uber.org/zap" // 结构化日志库 structured logging
)

// requestIDKey 定义上下文中的键，用于请求跟踪
// requestIDKey is the key for storing request ID in context
const requestIDKey string = "request_id"

// ContextJob 支持上下文的定时任务接口，Run 方法接收 context 并返回 error
// ContextJob is a cron job interface accepting a context and returning an error
type ContextJob interface {
	Run(ctx context.Context) error
}

// JobMeta 存储任务的元数据信息
// JobMeta stores metadata information for each job
type JobMeta struct {
	EntryID cron.EntryID  // Cron 内部任务ID Internal cron entry ID
	Spec    string        // Cron 时间表达式 Cron schedule expression
	Timeout time.Duration // 任务超时时间 Job timeout duration
	Type    string        // 任务类型 (func/job) Job type (func/job)
}

// Crontab 定时任务管理器，支持全局上下文、超时、日志及并发计数
// Crontab manages cron jobs with root context, per-job timeout, logging, and concurrency tracking
type Crontab struct {
	ctx     context.Context    // 根上下文，用于派生子任务 root context for child jobs
	cron    *cron.Cron         // robfig/cron 实例 cron scheduler instance
	ids     map[string]JobMeta // 任务ID到元数据的映射 Map of task IDs to job metadata
	mu      sync.RWMutex       // 保护ids和cron实例的互斥锁 Mutex for protecting ids and cron instance
	logger  *zap.Logger        // Zap 日志实例 zap logger for structured logging
	timeout time.Duration      // 默认任务超时时间 default timeout for each job
	running int32              // 并发执行任务计数 counter for concurrently running jobs
}

// New 创建并返回 Crontab 实例
// New creates a Crontab with given root context, zap logger, and default timeout
func New(ctx context.Context, logger *zap.Logger, defaultTimeout time.Duration) *Crontab {

	// 支持秒级的 cron 解析器 parser supporting seconds
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom |
			cron.Month | cron.DowOptional | cron.Descriptor,
	)

	cronLogger := &zapCronLogger{logger: logger}
	c := cron.New(
		cron.WithParser(parser),       // 使用自定义解析器 use custom parser
		cron.WithLocation(time.Local), // 本地时区 local timezone
		cron.WithChain(
			cron.SkipIfStillRunning(cronLogger), // 如果上次任务未完成，则跳过 skip if still running
			cron.Recover(cronLogger),            // 任务 panic 时自动恢复 recover from panic
		),
	)

	return &Crontab{
		ctx:     ctx,
		cron:    c,
		ids:     make(map[string]JobMeta),
		logger:  logger,
		timeout: defaultTimeout,
	}
}

// Start 启动调度器
// Start starts the cron scheduler
func (c *Crontab) Start() {
	c.cron.Start()
	c.logger.Info("cron scheduler started")
}

// Stop 停止调度并返回停止后的上下文，等待所有正在运行任务完成
// Stop stops scheduler and returns a context that ends when all jobs complete
func (c *Crontab) Stop() context.Context {
	c.logger.Info("stopping cron scheduler gracefully")
	return c.cron.Stop()
}

// StopWithTimeout 限时停止，超时后强制结束
// StopWithTimeout stops scheduler and forcibly ends after duration d
func (c *Crontab) StopWithTimeout(d time.Duration) {
	c.logger.Info("stopping cron with timeout", zap.Duration("timeout", d))
	ctx := c.Stop()
	select {
	case <-ctx.Done():
		c.logger.Info("cron stopped gracefully")
	case <-time.After(d):
		c.logger.Warn("cron forced stop after timeout")
	}
}

// RunningTasks 返回当前正在执行的任务数量
// RunningTasks returns number of jobs currently executing
func (c *Crontab) RunningTasks() int {
	return int(atomic.LoadInt32(&c.running))
}

// AddFunc 添加无上下文的简单任务
// AddFunc registers a simple function without context
func (c *Crontab) AddFunc(id, spec string, f func()) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查任务ID是否已存在
	// Check if task ID already exists
	if _, exists := c.ids[id]; exists {
		return errors.New("task id exists")
	}

	eid, err := c.cron.AddFunc(spec, f)
	if err != nil {
		c.logger.Error("AddFunc failed", zap.String("id", id), zap.Error(err))
		return err
	}

	// 存储任务元数据
	// Store job metadata
	c.ids[id] = JobMeta{
		EntryID: eid,
		Spec:    spec,
		Type:    "func",
	}

	c.logger.Info("task registered",
		zap.String("id", id),
		zap.String("spec", spec),
		zap.String("type", "func"))
	return nil
}

// AddJobWithTimeout 添加带上下文支持和自定义超时的任务
// AddJobWithTimeout registers a ContextJob with timeout d
func (c *Crontab) AddJobWithTimeout(id, spec string, job ContextJob, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查任务ID是否已存在
	// Check if task ID already exists
	if _, exists := c.ids[id]; exists {
		return errors.New("task id exists")
	}

	// 包装实际运行逻辑，包含超时和恢复机制
	// Wrap execution logic with timeout and recovery
	run := func() {
		start := time.Now()
		atomic.AddInt32(&c.running, 1)        // 增加并发计数 increment counter
		defer atomic.AddInt32(&c.running, -1) // 完成后减少计数 decrement counter

		// 创建带超时的上下文
		// Create timeout context
		ctx, cancel := context.WithTimeout(c.ctx, d)
		defer cancel()

		// 为每次任务生成唯一请求 ID
		// Generate unique request ID for each execution
		reqID := auuid.New().String()
		ctx = context.WithValue(ctx, requestIDKey, reqID)

		// 准备日志字段
		// Prepare logging fields
		fields := []zap.Field{
			zap.String("task_id", id),
			zap.String("spec", spec),
			zap.String("request_id", reqID),
			zap.Duration("timeout", d),
		}
		c.logger.Debug("job start", fields...)

		// 错误处理通道
		// Error handling channel
		errChan := make(chan error, 1)

		// 并发执行任务，捕获panic
		// Execute job concurrently and capture panics
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// 捕获并解析堆栈，提取首个非 runtime/robfig 的用户代码帧
					stack := debug.Stack()

					c.logger.Error("panic recovered",
						append(fields,
							zap.Any("error", r),
							zap.Strings("stack", agin.SplitStack(stack)), // 改进:数组形式
						)...)
					errChan <- fmt.Errorf("panic: %v", r)
				}
			}()
			errChan <- job.Run(ctx) // 执行用户任务 execute user job
		}()

		// 监控超时和任务完成
		// Monitor timeout and job completion
		select {
		case err := <-errChan:
			dur := time.Since(start)
			fields = append(fields, zap.Duration("duration", dur))

			// 根据执行结果记录不同级别的日志
			// Log at different levels based on execution result
			switch {
			case err != nil:
				c.logger.Error("job failed", append(fields, zap.Error(err))...)
			case dur > d/2:
				c.logger.Warn("job completed near timeout", fields...)
			default:
				c.logger.Debug("job completed", fields...)
			}
		case <-ctx.Done():
			// 任务超时处理
			// Handle job timeout
			elapsed := time.Since(start)
			c.logger.Warn("job timeout", append(fields, zap.Duration("elapsed", elapsed), zap.Error(ctx.Err()))...)
		}
	}

	// 将包装后的函数添加到调度器
	// Add wrapped function to scheduler
	eid, err := c.cron.AddFunc(spec, run)
	if err != nil {
		c.logger.Error("AddJob failed", zap.String("id", id), zap.Error(err))
		return err
	}

	// 存储任务元数据
	// Store job metadata
	c.ids[id] = JobMeta{
		EntryID: eid,
		Spec:    spec,
		Timeout: d,
		Type:    "job",
	}

	c.logger.Info("context job registered",
		zap.String("id", id),
		zap.String("spec", spec),
		zap.Duration("timeout", d))
	return nil
}

// AddJob 使用默认超时添加上下文任务
// AddJob registers a ContextJob with default timeout
func (c *Crontab) AddJob(id, spec string, job ContextJob) error {
	return c.AddJobWithTimeout(id, spec, job, c.timeout)
}

// Remove 删除已注册的任务
// Remove cancels and removes a job by its id
func (c *Crontab) Remove(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if meta, exists := c.ids[id]; exists {
		c.cron.Remove(meta.EntryID)
		delete(c.ids, id)
		c.logger.Info("task removed",
			zap.String("id", id),
			zap.String("type", meta.Type))
	}
}

// IDs 返回所有注册任务 ID 列表
// IDs returns all registered task IDs
func (c *Crontab) IDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ids := make([]string, 0, len(c.ids))
	for id := range c.ids {
		ids = append(ids, id)
	}
	return ids
}

// JobStatus 获取任务状态信息
// JobStatus retrieves detailed status of a job
func (c *Crontab) JobStatus(id string) (cron.Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if meta, exists := c.ids[id]; exists {
		return c.cron.Entry(meta.EntryID), true
	}
	return cron.Entry{}, false
}

// Reschedule 重新调度现有任务
// Reschedule updates the schedule of an existing job
func (c *Crontab) Reschedule(id, newSpec string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	meta, exists := c.ids[id]
	if !exists {
		return errors.New("task not found")
	}

	// 获取现有任务函数
	// Get existing job function
	entry := c.cron.Entry(meta.EntryID)
	if entry.Job == nil {
		return errors.New("job function not found")
	}

	// 添加新调度的任务
	// Add job with new schedule
	newID, err := c.cron.AddFunc(newSpec, entry.Job.Run)
	if err != nil {
		return err
	}

	// 移除旧任务
	// Remove old job
	c.cron.Remove(meta.EntryID)

	// 先保存旧表达式，再更新
	oldSpec := meta.Spec
	meta.EntryID = newID
	meta.Spec = newSpec
	c.ids[id] = meta

	c.logger.Info("job rescheduled",
		zap.String("id", id),
		zap.String("old_spec", oldSpec),
		zap.String("new_spec", newSpec))
	return nil
}

// 健康检查协程 (可选)
// Health check goroutine (optional)
func (c *Crontab) StartHealthCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			running := c.RunningTasks()
			c.logger.Debug("system health check",
				zap.Int("running_tasks", running),
				zap.Int("total_jobs", len(c.IDs())))

			// 高负载预警
			// High load warning
			if running > 50 {
				c.logger.Warn("high system load",
					zap.Int("running_tasks", running))
			}
		}
	}()
}
