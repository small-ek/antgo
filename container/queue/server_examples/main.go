package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/small-ek/antgo/container/queue"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/gin_middleware"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/os/config"
	"go.uber.org/zap"
	"io/ioutil"
)

func main() {
	// 初始化zap日志 (需要在应用层初始化)
	// Initialize zap logger (should be initialized at application level)
	gin.ForceConsoleColor()
	app := gin.New()
	if config.GetBool("system.debug") == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}
	app.Use(gin_middleware.Recovery()).Use(gin_middleware.Logger())
	configPath := flag.String("config", "./config.toml", "Configuration file path")

	app.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	eng := ant.New(*configPath).Serve(app)

	defer eng.Close()
	// 服务配置 Service configuration
	cfg := &queue.ServiceConfig{
		RedisAddress:  "127.0.0.1:6379",
		RedisPassword: "",
		RedisDB:       1,
		Concurrency:   20,
		Logger:        ant.Log(),
		Queues: map[string]int{
			"critical": 6,
			"default":  5,
		},
	}

	// 创建服务实例 Create service instance
	svc := queue.NewService(cfg)

	// 注册任务处理器 Register task handlers
	svc.RegisterHandler("order:process", handleOrderTask)

	// 启动服务 Start service
	if err := svc.Start(); err != nil {
		zap.L().Fatal("Failed to start service", zap.Error(err))
	}

}

func handleOrderTask(ctx context.Context, task *asynq.Task) error {

	zap.L().Info("Processing email task", zap.ByteString("payload", task.Payload()))
	// 实现业务逻辑...
	return nil
}
