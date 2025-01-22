package queue

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"time"
)

// ServiceConfig contains configuration for Asynq service
type ServiceConfig struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	Concurrency   int
	RetryStrategy RetryStrategy
	Queues        map[string]int
	Logger        *zap.Logger // 新增Logger字段
}

type RetryStrategy interface {
	GetDelay(retryCount int, err error, task *asynq.Task) time.Duration
}

type TaskHandler func(ctx context.Context, task *asynq.Task) error

type Service struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	config *ServiceConfig
	logger *zap.Logger // 新增logger字段
}

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

func NewService(cfg *ServiceConfig) *Service {
	// 设置默认值
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 10
	}
	if cfg.RedisDB == 0 {
		cfg.RedisDB = 1
	}
	if cfg.RetryStrategy == nil {
		cfg.RetryStrategy = &DefaultRetryStrategy{}
	}
	if cfg.Queues == nil {
		cfg.Queues = map[string]int{"default": 1}
	}
	// 设置默认Logger（如果未提供）
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop()
	}

	return &Service{
		config: cfg,
		mux:    asynq.NewServeMux(),
		logger: cfg.Logger, // 初始化logger字段
	}
}

func (s *Service) Start() error {
	if s.config.RedisAddress == "" {
		return fmt.Errorf("redis address is required")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr:     s.config.RedisAddress,
		Password: s.config.RedisPassword,
		DB:       s.config.RedisDB,
	}

	s.server = asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: s.config.Concurrency,
			Queues:      s.config.Queues,
			RetryDelayFunc: func(retryCount int, err error, task *asynq.Task) time.Duration {
				return s.config.RetryStrategy.GetDelay(retryCount, err, task)
			},
			Logger: s, // 保持Logger接口实现
		},
	)

	if err := s.server.Run(s.mux); err != nil {
		s.logger.Error("Failed to start Asynq server", zap.Error(err)) // 使用传入的logger
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// 日志方法现在使用传入的logger
func (s *Service) Debug(args ...interface{}) {
	s.logger.Debug(fmt.Sprint(args...))
}

func (s *Service) Info(args ...interface{}) {
	s.logger.Info(fmt.Sprint(args...))
}

func (s *Service) Warn(args ...interface{}) {
	s.logger.Warn(fmt.Sprint(args...))
}

func (s *Service) Error(args ...interface{}) {
	s.logger.Error(fmt.Sprint(args...))
}

func (s *Service) Fatal(args ...interface{}) {
	s.logger.Fatal(fmt.Sprint(args...))
}

// 以下方法保持不变
func (s *Service) RegisterHandler(taskType string, handler TaskHandler) {
	s.mux.HandleFunc(taskType, handler)
}

func (s *Service) Shutdown() {
	if s.server != nil {
		s.server.Stop()
		s.server.Shutdown()
	}
}
