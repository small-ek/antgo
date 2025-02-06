package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Gin 结构体是Gin框架的适配器，用于集成Gin到自定义框架中
// Gin struct is an adapter for Gin framework, integrating Gin into the custom framework.
type Gin struct {
	serve.BaseAdapter
	ctx *gin.Context
	app *gin.Engine
	Srv *http.Server
}

func init() {
	ant.Register(new(Gin))
}

// Name 返回当前适配器名称
// Name returns the name of the current adapter.
func (g *Gin) Name() string {
	return "gin"
}

// SetApp 设置并验证Gin引擎实例
// SetApp sets and validates the Gin engine instance.
func (g *Gin) SetApp(app interface{}) error {
	engine, ok := app.(*gin.Engine)
	if !ok {
		return errors.New("gin adapter SetApp: invalid parameter type")
	}

	// 设置生产模式提升性能
	// Set production mode to enhance performance
	gin.SetMode(gin.ReleaseMode)
	g.app = engine
	return nil
}

// Run 启动HTTP服务（不加载配置服务）
// Run starts the HTTP server (without loading configuration service)
func (g *Gin) Run(addr string) {
	// 初始化HTTP服务器配置
	// Initialize HTTP server configuration
	g.Srv = &http.Server{
		Addr:    ":" + addr,
		Handler: g.app,
	}

	// 输出服务启动信息
	// Print service startup information
	alog.Write.Info("Service started",
		zap.Int("pid", os.Getpid()),
		zap.String("address", "http://127.0.0.1"+g.Srv.Addr),
	)

	// 启动异步HTTP服务
	// Start asynchronous HTTP service
	go func() {
		if err := g.Srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			alog.Write.Fatal("Server startup failed", zap.Error(err))
		}
	}()
}

// Close 实现优雅的HTTP服务关闭
// Close implements graceful shutdown of the HTTP service
func (g *Gin) Close() {
	// 创建带缓冲的信号通道
	// Create buffered signal channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// 阻塞等待关闭信号
	// Block waiting for shutdown signal
	<-quit

	// 创建带超时的上下文
	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行优雅关闭
	// Perform graceful shutdown
	if err := g.Srv.Shutdown(ctx); err != nil {
		alog.Write.Error("Server shutdown error", zap.Error(err))
	} else {
		alog.Write.Info("Service shutdown", zap.Int("pid", os.Getpid()))
	}
}
