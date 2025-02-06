package ant

import (
	"flag"
	"github.com/small-ek/antgo/db/adb"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
	"log"
	"net/http"
)

// Engine is the core component of antgo.
// Engine 是 antgo 的核心组件.
type Engine struct {
	Adapter serve.WebFrameWork // Web framework adapter for handling HTTP requests.
	// 用于处理 HTTP 请求的 Web 框架适配器.
	Srv *http.Server // HTTP server instance.
	// HTTP 服务实例.
	Config config.ConfigStr // Configuration structure.
	// 配置结构体.
	port string // Custom port (if provided).
	// 自定义端口（如果提供）.
}

// defaultAdapter holds the registered web framework adapter.
// defaultAdapter 保存注册的 Web 框架适配器.
var defaultAdapter serve.WebFrameWork

// New creates and returns a new Engine instance.
// Optionally accepts one or more configuration file paths.
// New 创建并返回一个新的 Engine 实例。
// 可选传入一个或多个配置文件路径.
func New(configPath ...string) *Engine {
	// Set detailed log flags.
	// 设置详细的日志标记.
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()

	if len(configPath) > 0 {
		if err := config.New(configPath...).Register(); err != nil {
			panic(err)
		}
		loadApp()
	}

	return &Engine{
		Adapter: defaultAdapter,
	}
}

// AddConfig adds an additional configuration file.
// AddConfig 添加额外的配置文件.
func (eng *Engine) AddConfig(configPath string) *Engine {
	if configPath != "" {
		if err := config.AddConfigFile(configPath); err != nil {
			panic(err)
		}
	}
	return eng
}

// AddFunc executes one or more initialization functions during Engine setup.
// AddFunc 在 Engine 初始化期间执行一个或多个函数.
func (eng *Engine) AddFunc(funcs ...func()) *Engine {
	for _, f := range funcs {
		f()
	}
	return eng
}

// SetPort sets the port on which the Engine will listen.
// SetPort 设置 Engine 监听的端口.
func (eng *Engine) SetPort(port string) *Engine {
	eng.port = port
	return eng
}

// Register registers the default web framework adapter.
// Register 注册默认的 Web 框架适配器.
func Register(ada serve.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}

// Serve starts the HTTP server with the provided application.
// Serve 使用指定的应用程序启动 HTTP 服务.
func (eng *Engine) Serve(app interface{}) *Engine {
	if eng.Adapter == nil {
		panic("adapter is nil")
	}

	// Retrieve the address from configuration; default to "8888" if not set.
	// 从配置中获取地址，如果未设置则默认为 "8888".
	addr := config.GetString("system.address")
	if addr == "" {
		addr = "8888"
	}
	// Override with custom port if provided.
	// 如果提供了自定义端口，则覆盖配置中的端口.
	if eng.port != "" {
		addr = eng.port
	}

	// Set the application for the adapter.
	// 为适配器设置应用程序.
	if err := eng.Adapter.SetApp(app); err != nil {
		panic(err)
	}

	// Run the adapter on the specified address.
	// 在指定地址上启动适配器.
	eng.Adapter.Run(addr)
	return eng
}

// Close gracefully shuts down the Engine and its associated resources.
// Optionally executes a callback function after closing resources.
// Close 优雅地关闭 Engine 及其相关资源。
// 如果提供了回调函数，则在关闭资源后执行该函数.
func (eng *Engine) Close(callbacks ...func()) *Engine {
	eng.Adapter.Close()

	// Close database connections if configured.
	// 如果配置了数据库连接，则关闭它们.
	if connections := config.GetMaps("connections"); len(connections) > 0 {
		adb.Close()
	}

	// Execute the first callback function if provided.
	// 如果提供了回调函数，则执行第一个.
	if len(callbacks) > 0 {
		callbacks[0]()
	}
	return eng
}

// SetConfig resets the Engine's configuration using new configuration files.
// SetConfig 使用新的配置文件重置 Engine 的配置.
func (eng *Engine) SetConfig(filePath ...string) *Engine {
	if err := config.New(filePath...).Register(); err != nil {
		panic(err)
	}
	loadApp()
	return eng
}

// AddRemoteProvider adds a remote configuration provider to the Engine.
// AddRemoteProvider 添加远程配置提供者.
func (eng *Engine) AddRemoteProvider(provider, endpoint, path string) *Engine {
	if err := config.New().AddRemoteProvider(provider, endpoint, path); err != nil {
		panic(err)
	}
	loadApp()
	return eng
}

// Etcd configures etcd as the configuration backend.
// Etcd 使用 etcd 作为配置后端.
func (eng *Engine) Etcd(hosts, paths []string, username, pwd string) *Engine {
	if len(hosts) > 0 && len(paths) > 0 {
		if err := config.New().Etcd3(hosts, paths, username, pwd); err != nil {
			panic(err)
		}
		loadApp()
	}
	return eng
}

// SetLog sets the log file path and registers the logging system.
// SetLog 设置日志文件路径并注册日志系统.
func (eng *Engine) SetLog(filePath string) *Engine {
	alog.New(filePath).Register()
	return eng
}

// loadApp initializes application components such as logging, database connections, and Redis.
// loadApp 初始化应用组件，例如日志、数据库连接和 Redis.
func loadApp() {
	if config.Config != nil {
		initLog() // Initialize logging system.
		// 初始化日志系统.
		adb.InitDb(config.GetMaps("connections")) // Initialize database connections.
		// 初始化数据库连接.
		initRedis() // Initialize Redis if configured.
		// 如果配置了 Redis，则初始化 Redis.
	}
}
