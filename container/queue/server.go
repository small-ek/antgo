package queue

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"time"
)

// ServiceConfig contains configuration for Asynq service
// ServiceConfig 包含Asynq服务配置
type ServiceConfig struct {
	RedisAddress  string        // Redis server address Redis服务器地址
	RedisPassword string        // Redis auth password Redis认证密码
	RedisDB       int           // Redis database number Redis数据库编号
	Concurrency   int           // Number of concurrent workers 并发工作协程数
	RetryStrategy RetryStrategy // Custom retry strategy 自定义重试策略
}

// RetryStrategy defines interface for custom retry strategies
// RetryStrategy 定义自定义重试策略接口
type RetryStrategy interface {
	GetDelay(retryCount int, err error, task *asynq.Task) time.Duration
}

// TaskHandler defines task handler function signature
// TaskHandler 定义任务处理函数签名
type TaskHandler func(ctx context.Context, task *asynq.Task) error

// Service represents the Asynq service instance
// Service 表示Asynq服务实例
type Service struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	config *ServiceConfig
}

// DefaultRetryStrategy provides production-ready retry logic
// DefaultRetryStrategy 提供生产环境就绪的重试策略
type DefaultRetryStrategy struct{}

func (d *DefaultRetryStrategy) GetDelay(retryCount int, err error, task *asynq.Task) time.Duration {
	switch {
	case retryCount <= 3:
		return time.Duration(retryCount) * 20 * time.Second
	case retryCount == 4:
		return 1 * time.Hour
	case retryCount == 5:
		return 5 * time.Hour
	default:
		return time.Duration((retryCount-4)*4+5) * time.Hour
	}
}

// NewService creates a new Asynq service instance
// NewService 创建新的Asynq服务实例
func NewService(cfg *ServiceConfig) *Service {
	// Set default values 设置默认值
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 10
	}
	if cfg.RedisDB == 0 {
		cfg.RedisDB = 1
	}
	if cfg.RetryStrategy == nil {
		cfg.RetryStrategy = &DefaultRetryStrategy{}
	}

	return &Service{
		config: cfg,
		mux:    asynq.NewServeMux(),
	}
}

// Start begins processing tasks
// Start 开始处理任务
func (s *Service) Start() error {
	// Validate configuration 配置校验
	if s.config.RedisAddress == "" {
		return fmt.Errorf("redis address is required")
	}

	// Create Redis connection options 创建Redis连接配置
	redisOpt := asynq.RedisClientOpt{
		Addr:     s.config.RedisAddress,
		Password: s.config.RedisPassword,
		DB:       s.config.RedisDB,
	}

	// Configure server 配置服务端
	s.server = asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: s.config.Concurrency,
			RetryDelayFunc: func(retryCount int, err error, task *asynq.Task) time.Duration {
				return s.config.RetryStrategy.GetDelay(retryCount, err, task)
			},
			Logger: s, // Implement asynq.Logger interface
		},
	)

	// Start processing 开始处理任务
	if err := s.server.Run(s.mux); err != nil {
		zap.L().Error("Failed to start Asynq server", zap.Error(err))
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// RegisterHandler registers task handler for given type
// RegisterHandler 注册任务处理器
func (s *Service) RegisterHandler(taskType string, handler TaskHandler) {
	s.mux.HandleFunc(taskType, handler)
}

// Shutdown gracefully stops the service
// Shutdown 优雅停止服务
func (s *Service) Shutdown() {
	if s.server != nil {
		s.server.Stop()
		s.server.Shutdown()
	}
}

// Implement asynq.Logger interface using zap
// 使用zap实现asynq.Logger接口

func (s *Service) Debug(args ...interface{}) {
	zap.L().Debug(fmt.Sprint(args...))
}

func (s *Service) Info(args ...interface{}) {
	zap.L().Info(fmt.Sprint(args...))
}

func (s *Service) Warn(args ...interface{}) {
	zap.L().Warn(fmt.Sprint(args...))
}

func (s *Service) Error(args ...interface{}) {
	zap.L().Error(fmt.Sprint(args...))
}

func (s *Service) Fatal(args ...interface{}) {
	zap.L().Fatal(fmt.Sprint(args...))
}
