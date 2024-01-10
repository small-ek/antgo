package ant

import (
	"context"
	"flag"
	"fmt"
	"github.com/small-ek/antgo/db/adb"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/alog"
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
	Config       config.ConfigStr
	announceLock sync.Once
}

// defaultAdapter is the default adapter.
var defaultAdapter serve.WebFrameWork

// New return the default engine instance.
func New(configPath ...string) *Engine {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()

	if len(configPath) > 0 {
		err := config.New(configPath...).Regiter()
		if err != nil {
			panic(err)
		}
		loadApp()
	}

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
	fmt.Printf("  PID: %d \n", os.Getpid())
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://127.0.0.1" + eng.Srv.Addr)
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return eng
}

// defaultServer
func defaultServer(app http.Handler) *http.Server {
	addr := config.GetString("system.address")
	if addr == "" {
		addr = ":8080"
	}

	return &http.Server{
		Addr:    addr,
		Handler: app,
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

	fmt.Printf("  PID: %d \n", os.Getpid())
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://127.0.0.1" + eng.Srv.Addr)
	return eng
}

// Close signal<关闭服务操作>
func (eng *Engine) Close(f ...func()) *Engine {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	Log().Warn("Exit service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connections := config.GetMaps("connections")

	if len(connections) > 0 {
		defer adb.Close()
	}

	if err := eng.Srv.Shutdown(ctx); err != nil {
		Log().Error("Server Shutdown:" + err.Error())
	}

	if len(f) > 0 {
		f[0]()
	}
	return eng
}

// GetServer Get http service<获取http服务>
func (eng *Engine) GetServer() *http.Server {
	return eng.Srv
}

// SetConfig Modify the configuration path<修改配置路径>
func (eng *Engine) SetConfig(filePath ...string) *Engine {
	err := config.New(filePath...).Regiter()
	if err != nil {
		panic(err)
	}

	loadApp()
	return eng
}

// AddSecureRemoteProvider.<添加远程连接>
func (eng *Engine) AddRemoteProvider(provider, endpoint, path string) *Engine {
	err := config.New().AddRemoteProvider(provider, endpoint, path)
	if err != nil {
		panic(err)
	}
	loadApp()
	return eng
}

// SetLog Modify log path.<修改日志路径>
func (eng *Engine) SetLog(filePath string) *Engine {
	alog.New(filePath).Register()
	return eng
}

// loadApp.<加载应用>
func loadApp() {
	if config.Config != nil {
		//加载默认配置
		initLog()
		adb.InitDb()
	}
}
