package ant

import (
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/os/logger"
)

// GetConfig Get configuration content<获取配置>
func GetConfig(name string) *config.Config {
	return config.Decode().Get(name)
}

//initConfigLog Set log according to the configuration file<初始化配置日志>
func initConfigLog() {
	cfg := config.Decode()
	logPath := cfg.Get("log.path").String()

	if cfg.Get("log.switch").Bool() == true && logPath != "" {
		logger.Default(logPath).
			SetLevel(cfg.Get("log.level").String()).
			SetServiceName(cfg.Get("log.service_name").String()).
			SetMaxSize(cfg.Get("log.max_size").Int()).
			SetMaxAge(cfg.Get("log.max_age").Int()).
			SetMaxBackups(cfg.Get("log.max_backups").Int()).
			SetFormat(cfg.Get("log.format").String()).
			SetConsole(cfg.Get("log.console").Bool()).
			SetCompress(cfg.Get("log.compress").Bool()).
			Register()
	}
}
