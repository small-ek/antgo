package ant

import (
	"context"
	"fmt"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Engine is the core component of antgo.
type Engine struct {
	Adapter      serve.WebFrameWork
	Srv          *http.Server
	Config       config.Config
	announceLock sync.Once
}

//defaultAdapter is the default adapter.
var defaultAdapter serve.WebFrameWork

// Default return the default engine instance.
func Default() *Engine {
	return &Engine{
		Adapter: defaultAdapter,
	}
}

// Register the default adapter.
func Register(ada serve.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}

// Use enable the adapter.
func (eng *Engine) Use(router interface{}) *Engine {
	if eng.Adapter == nil {
		panic("adapter is nil")
	}
	if err := eng.Adapter.SetApp(router); err != nil {
		panic("gin adapter SetApp: wrong parameter")
	}
	return eng
}

// Run http service
func (eng *Engine) Run(srv *http.Server) *Engine {
	eng.Srv = srv
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://127.0.0.1" + eng.Srv.Addr)
	return eng
}

//defaultServer
func defaultServer(app http.Handler) *http.Server {
	return &http.Server{
		Addr:              GetConfig("system.address").String(),
		Handler:           app,
		ReadTimeout:       240 * time.Second, //设置秒的读超时
		WriteTimeout:      240 * time.Second, //设置秒的写超时
		ReadHeaderTimeout: 60 * time.Second,  //读取头超时
		IdleTimeout:       120 * time.Second, //空闲超时
		MaxHeaderBytes:    2097152,           //请求头最大字节
	}
}

// Serve http service
func (eng *Engine) Serve(app http.Handler) *Engine {
	eng.Use(app)
	eng.Srv = defaultServer(app)
	go func() {
		// 服务连接
		if err := eng.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://127.0.0.1" + eng.Srv.Addr)
	return eng
}

// Close signal
func (eng *Engine) Close() *Engine {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := eng.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return eng
}

// GetServer Get http service
func (eng *Engine) GetServer() *http.Server {
	return eng.Srv
}

// SetConfig Modify the configuration path
func (eng *Engine) SetConfig(filePath string) *Engine {
	config.SetPath(filePath)
	return eng
}

// GetConfig Get configuration content
func GetConfig(name string) *config.Config {
	return config.Decode().Get(name)
}
