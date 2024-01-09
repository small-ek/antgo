package ant

import (
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
)

// GetConfig Get configuration content<获取配置>
func GetConfig(name string) any {
	return config.Get(name)
}

// initLog Set log according to the configuration file<初始化配置日志>
func initLog() {
	logPath := config.GetString("log.path")

	if config.GetBool("log.switch") == true && logPath != "" {
		alog.Default(logPath).
			SetLevel(config.GetString("log.level")).
			SetServiceName(config.GetString("log.service_name")).
			SetMaxSize(config.GetInt("log.max_size")).
			SetMaxAge(config.GetInt("log.max_age")).
			SetMaxBackups(config.GetInt("log.max_backups")).
			SetFormat(config.GetString("log.format")).
			SetConsole(config.GetBool("log.console")).
			SetCompress(config.GetBool("log.compress")).
			Register()
	}
}
