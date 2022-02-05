package ant

import (
	"context"
	"flag"
	"fmt"
	"github.com/small-ek/antgo/db/adb"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/os/logger"
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
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	return &Engine{
		Adapter: defaultAdapter,
	}
}

// Register the default adapter.<服务注册>
func Register(ada serve.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}

// Use enable the adapter.<引用组件>
func (eng *Engine) Use(router interface{}) *Engine {
	if eng.Adapter == nil {
		panic("adapter is nil")
	}
	if err := eng.Adapter.SetApp(router); err != nil {
		panic("gin adapter SetApp: wrong parameter")
	}
	return eng
}

// Run http service<不加载配置服务>
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
	addr := GetConfig("system.address").String()
	if addr == "" {
		addr = ":8080"
	}

	return &http.Server{
		Addr:              addr,
		Handler:           app,
		ReadTimeout:       240 * time.Second, //设置秒的读超时
		WriteTimeout:      240 * time.Second, //设置秒的写超时
		ReadHeaderTimeout: 60 * time.Second,  //读取头超时
		IdleTimeout:       120 * time.Second, //空闲超时
		MaxHeaderBytes:    2097152,           //请求头最大字节
	}
}

// Serve http service<默认服务加载>
func (eng *Engine) Serve(app http.Handler) *Engine {
	eng.Use(app)
	eng.Srv = defaultServer(app)

	go func() {
		if err := eng.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://127.0.0.1" + eng.Srv.Addr)
	return eng
}

// Close signal<关闭服务操作>
func (eng *Engine) Close() *Engine {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.Decode()
	connections := cfg.Get("connections").Maps()
	if len(connections) > 0 {
		defer adb.Close()
	}

	if err := eng.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return eng
}

// GetServer Get http service<获取http服务>
func (eng *Engine) GetServer() *http.Server {
	return eng.Srv
}

// SetConfig Modify the configuration path<修改配置路径>
func (eng *Engine) SetConfig(filePath string) *Engine {
	config.SetPath(filePath)
	//加载默认配置
	initConfigLog()
	adb.InitDb()
	return eng
}

// SetLog Modify log path.<修改日志路径>
func (eng *Engine) SetLog(filePath string) *Engine {
	logger.Default(filePath).Register()
	return eng
}
