package ant

import (
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
)

// GetConfig 获取配置内容 / Get configuration content
func GetConfig(name string) any {
	return config.Get(name)
}

// initLog 根据配置文件初始化日志 / initLog initializes logging according to the configuration file
func initLog() {
	// Check if logging is enabled / 检查是否启用了日志功能
	if !config.GetBool("log.switch") {
		return
	}

	// Retrieve the log path from configuration / 从配置中获取日志路径
	logPath := config.GetString("log.path")
	if logPath == "" {
		return // If log path is empty, skip log initialization / 如果日志路径为空则跳过日志初始化
	}

	// Cache all log related configuration values for better performance and readability
	// 缓存所有日志相关的配置项，提高性能和代码可读性
	level := config.GetString("log.level")
	serviceName := config.GetString("log.service_name")
	maxSize := config.GetInt("log.max_size")
	maxAge := config.GetInt("log.max_age")
	maxBackups := config.GetInt("log.max_backups")
	format := config.GetString("log.format")
	console := config.GetBool("log.console")
	compress := config.GetBool("log.compress")

	// Initialize the logger with the cached configuration values
	// 使用缓存的配置值初始化日志记录器
	logger := alog.New(logPath).
		SetLevel(level).
		SetServiceName(serviceName).
		SetMaxSize(maxSize).
		SetMaxAge(maxAge).
		SetMaxBackups(maxBackups).
		SetFormat(format).
		SetConsole(console).
		SetCompress(compress)

	// Register the logger / 注册日志记录器
	logger.Register()
}
