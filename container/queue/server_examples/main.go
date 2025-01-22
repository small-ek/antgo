package main

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/small-ek/antgo/container/queue"
	"go.uber.org/zap"
)

func main() {
	// 初始化zap日志 (需要在应用层初始化)
	// Initialize zap logger (should be initialized at application level)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// 服务配置 Service configuration
	cfg := &queue.ServiceConfig{
		RedisAddress:  "127.0.0.1:6379",
		RedisPassword: "",
		RedisDB:       1,
		Concurrency:   20,
	}

	// 创建服务实例 Create service instance
	svc := queue.NewService(cfg)

	// 注册任务处理器 Register task handlers
	svc.RegisterHandler("order:process", handleOrderTask)

	// 启动服务 Start service
	if err := svc.Start(); err != nil {
		zap.L().Fatal("Failed to start service", zap.Error(err))
	}
	defer svc.Shutdown()
}

func handleOrderTask(ctx context.Context, task *asynq.Task) error {
	fmt.Println("------------------")
	fmt.Println(task.Payload())
	zap.L().Info("Processing email task", zap.ByteString("payload", task.Payload()))
	// 实现业务逻辑...
	return nil
}
