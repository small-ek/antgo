package ant

import (
	"flag"
	"github.com/small-ek/antgo/db/adb"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/utils/ants"
	"log"
	"net/http"
)

// Engine is the core component of antgo.
type Engine struct {
	Adapter serve.WebFrameWork
	Srv     *http.Server
	Config  config.ConfigStr
	port    string
}

// defaultAdapter is the default adapter.
var defaultAdapter serve.WebFrameWork

// New return the default engine instance.
func New(configPath ...string) *Engine {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()

	if len(configPath) > 0 {
		err := config.New(configPath...).Register()
		if err != nil {
			panic(err)
		}
		loadApp()
	}

	return &Engine{
		Adapter: defaultAdapter,
	}
}

func (eng *Engine) AddConfig(configPath string) *Engine {
	if len(configPath) > 0 {
		err := config.AddConfigFile(configPath)
		if err != nil {
			panic(err)
		}
	}
	return eng
}

// AddFunc Add function execution
func (eng *Engine) AddFunc(f ...func()) *Engine {
	if len(f) > 0 {
		for i := 0; i < len(f); i++ {
			f[i]()
		}
	}
	return eng
}

// SetPort Set Port
func (eng *Engine) SetPort(port string) *Engine {
	eng.port = port
	return eng
}

// Register the default adapter.<服务注册>
func Register(ada serve.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}

// Serve http service<默认服务加载>
func (eng *Engine) Serve(app interface{}) *Engine {
	if eng.Adapter == nil {
		panic("adapter is nil")
	}

	addr := config.GetString("system.address")
	if addr == "" {
		addr = "8081"
	}
	if eng.port != "" {
		addr = eng.port
	}
	if err := eng.Adapter.SetApp(app); err != nil {
		panic(err)
	}
	eng.Adapter.Run(addr)
	return eng
}

// Close signal<关闭服务操作>
func (eng *Engine) Close(f ...func()) *Engine {
	eng.Adapter.Close()
	connections := config.GetMaps("connections")

	if len(connections) > 0 {
		defer adb.Close()
	}
	if len(f) > 0 {
		f[0]()
	}
	return eng
}

// SetConfig Modify the configuration path<修改配置路径>
func (eng *Engine) SetConfig(filePath ...string) *Engine {
	err := config.New(filePath...).Register()
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
func (eng *Engine) Etcd(host, path []string, username, pwd string) *Engine {
	if len(host) > 0 && len(host) > 0 {
		err := config.New().Etcd3(host, path, username, pwd)
		if err != nil {
			panic(err)
		}
		loadApp()
	}
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
		adb.InitDb(config.GetMaps("connections"))
		initRedis()

		ants.InitPool(config.GetInt("system.max_pool_count"))
	}
}
