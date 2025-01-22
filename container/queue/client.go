package queue

import (
	"github.com/hibiken/asynq"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"sync"
	"time"
)

// ClientConfig Redis client configuration parameters
// Redis客户端配置参数
type ClientConfig struct {
	Addr         string        // Redis服务器地址 Host:Port
	Password     string        // Redis认证密码
	DB           int           // Redis数据库编号 0-15
	PoolSize     int           // 连接池大小（建议为最大并发数2倍）
	DialTimeout  time.Duration // 连接建立超时时间
	ReadTimeout  time.Duration // 读取操作超时时间
	WriteTimeout time.Duration // 写入操作超时时间
}

// TaskOption Functional options pattern for task configuration
// 任务配置的功能选项模式
type TaskOption func(*taskOptions)

type taskOptions struct {
	delay     time.Duration // 延迟执行时间
	maxRetry  int           // 最大重试次数
	queue     string        // 队列名称
	timeout   time.Duration // 任务处理超时
	deadline  time.Time     // 任务截止时间
	uniqueTTL time.Duration // 唯一任务锁定时长
}

// AsyncClient Asynchronous task queue client
// 异步任务队列客户端
type AsyncClient struct {
	client   *asynq.Client        // Asynq核心客户端
	redisOpt asynq.RedisClientOpt // Redis连接配置
	config   ClientConfig         // 客户端配置
	logger   *zap.Logger          // 日志记录器
	mu       sync.RWMutex         // 读写锁（保证线程安全）
}

var (
	instance *AsyncClient // 单例实例
	once     sync.Once    // 单例控制
)

// NewClient Create a singleton client instance (Thread-safe)
// 创建单例客户端实例（线程安全）
func NewClient(cfg ClientConfig, opts ...ClientOption) *AsyncClient {
	once.Do(func() {
		// 设置默认选项
		options := clientOptions{
			logger: zap.NewNop(), // 默认无日志
		}
		for _, opt := range opts {
			opt(&options)
		}

		// 初始化Redis配置
		redisOpt := asynq.RedisClientOpt{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}

		// 创建客户端实例
		instance = &AsyncClient{
			client:   asynq.NewClient(redisOpt),
			redisOpt: redisOpt,
			config:   cfg,
			logger:   options.logger,
		}

		// 连接健康检查
		if err := instance.Ping(); err != nil {
			instance.logger.Fatal("Redis connection failed",
				zap.Error(err),
				zap.Any("config", cfg))
		}
	})
	return instance
}

// ClientOption Client configuration options
// 客户端配置选项
type ClientOption func(*clientOptions)

type clientOptions struct {
	logger *zap.Logger
}

// WithLogger Custom logger configuration
// 自定义日志配置
func WithLogger(logger *zap.Logger) ClientOption {
	return func(o *clientOptions) {
		o.logger = logger
	}
}

// 任务选项配置函数
// WithDelay Set task delay execution time
// 设置任务延迟执行时间
func WithDelay(delay time.Duration) TaskOption {
	return func(opts *taskOptions) {
		opts.delay = delay
	}
}

// WithMaxRetry Set maximum retry attempts
// 设置最大重试次数
func WithMaxRetry(maxRetry int) TaskOption {
	return func(opts *taskOptions) {
		opts.maxRetry = maxRetry
	}
}

// WithQueue Specify queue name
// 指定队列名称
func WithQueue(queue string) TaskOption {
	return func(opts *taskOptions) {
		opts.queue = queue
	}
}

// WithTimeout Set task processing timeout
// 设置任务处理超时
func WithTimeout(timeout time.Duration) TaskOption {
	return func(opts *taskOptions) {
		opts.timeout = timeout
	}
}

// WithDeadline Set task deadline
// 设置任务截止时间
func WithDeadline(deadline time.Time) TaskOption {
	return func(opts *taskOptions) {
		opts.deadline = deadline
	}
}

// WithUnique Set unique task lock TTL (Prevent duplicate tasks)
// 设置唯一任务锁定时间（防止重复任务）
func WithUnique(ttl time.Duration) TaskOption {
	return func(opts *taskOptions) {
		opts.uniqueTTL = ttl
	}
}

// Enqueue Add task to queue (Thread-safe)
// 添加任务到队列（线程安全）
func (c *AsyncClient) Enqueue(taskType string, payload interface{}, opts ...TaskOption) (*asynq.TaskInfo, error) {
	// 序列化任务数据
	data, err := conv.ToJSON(payload)
	if err != nil {
		c.logger.Error("Payload serialization failed",
			zap.String("task_type", taskType),
			zap.Error(err))
		return nil, err
	}

	// 创建基础任务
	task := asynq.NewTask(taskType, data)

	// 处理任务选项
	options := processTaskOptions(opts)

	// 加锁保证线程安全
	c.mu.Lock()
	defer c.mu.Unlock()

	// 提交任务到队列
	info, err := c.client.Enqueue(task, options...)
	if err != nil {
		c.logger.Error("Task enqueue failed",
			zap.String("task_type", taskType),
			zap.Error(err))
		return nil, err
	}

	c.logger.Info("Task enqueued successfully",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue))
	return info, nil
}

// processTaskOptions Convert custom options to Asynq options
// 将自定义选项转换为Asynq原生选项
func processTaskOptions(opts []TaskOption) []asynq.Option {
	var options taskOptions
	for _, opt := range opts {
		opt(&options)
	}

	var asynqOpts []asynq.Option
	// 延迟执行
	if options.delay > 0 {
		asynqOpts = append(asynqOpts, asynq.ProcessIn(options.delay))
	}
	// 最大重试
	if options.maxRetry > 0 {
		asynqOpts = append(asynqOpts, asynq.MaxRetry(options.maxRetry))
	}
	// 指定队列
	if options.queue != "" {
		asynqOpts = append(asynqOpts, asynq.Queue(options.queue))
	}
	// 处理超时
	if options.timeout > 0 {
		asynqOpts = append(asynqOpts, asynq.Timeout(options.timeout))
	}
	// 截止时间
	if !options.deadline.IsZero() {
		asynqOpts = append(asynqOpts, asynq.Deadline(options.deadline))
	}
	// 唯一任务锁
	if options.uniqueTTL > 0 {
		asynqOpts = append(asynqOpts, asynq.Unique(options.uniqueTTL))
	}

	return asynqOpts
}

// Ping Check Redis connection status
// 检查Redis连接状态
func (c *AsyncClient) Ping() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client.Ping()
}

// Close Gracefully close connection
// 优雅关闭连接
func (c *AsyncClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.client.Close(); err != nil {
		c.logger.Error("Connection close failed",
			zap.Error(err),
			zap.Any("config", c.config))
		return err
	}
	c.logger.Info("Connection closed cleanly")
	return nil
}

// HealthCheck Connection health monitor
// 连接健康监控
func (c *AsyncClient) HealthCheck() error {
	if err := c.Ping(); err != nil {
		c.logger.Warn("Redis connection unhealthy",
			zap.Error(err),
			zap.String("addr", c.config.Addr))
		return err
	}
	c.logger.Debug("Redis connection healthy")
	return nil
}

// GetQueueInfo Get queue information (Support dynamic inspection)
// 获取队列信息（支持动态检查）
func (c *AsyncClient) GetQueueInfo(queue string) (*asynq.QueueInfo, error) {
	inspector := asynq.NewInspector(c.redisOpt)
	defer inspector.Close()

	info, err := inspector.GetQueueInfo(queue)
	if err != nil {
		c.logger.Error("Get queue info failed",
			zap.String("queue", queue),
			zap.Error(err))
		return nil, err
	}

	c.logger.Debug("Queue status",
		zap.String("queue", queue),
		zap.Int("size", info.Size),
		zap.Int("active_workers", info.Active))
	return info, nil
}
