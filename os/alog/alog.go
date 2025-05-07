package alog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

// Write is the global logger instance that will be used to log messages
var Write *zap.Logger

// 预创建带Skip的专用Logger（单例）
var wrappedLogger *zap.Logger

// Logs represents the configuration options for logging
type Logs struct {
	Path        string // Log file save path
	Level       string // Log level (info, debug, warn)
	MaxBackups  int    // Number of backups to keep (default 300)
	MaxSize     int    // Maximum size of each log file (default 10MB)
	MaxAge      int    // Maximum age of log files (default 180 days)
	Compress    bool   // Whether to compress log files (default false)
	ServiceName string // Service name for logging (default "antgo")
	Format      string // Log format ("console" or "json")
	Console     bool   // Whether to output logs to the console
}

// New creates a new Logs instance with default settings
// 创建一个带有默认设置的Logs实例
func New(path string) *Logs {
	return &Logs{
		Path:        path,
		MaxSize:     10,      // Default max size 10MB
		MaxBackups:  300,     // Default max backups 300
		MaxAge:      180,     // Default max age 180 days
		Compress:    false,   // Default no compression
		ServiceName: "antgo", // Default service name
	}
}

// Register configures and initializes the logger based on the Logs settings
// 注册并根据Logs设置配置和初始化日志器
func (logs *Logs) Register() *zap.Logger {
	// Log file rotation
	hook := lumberjack.Logger{
		Filename:   logs.Path,
		MaxSize:    logs.MaxSize,
		MaxBackups: logs.MaxBackups,
		MaxAge:     logs.MaxAge,
		Compress:   logs.Compress,
	}

	// Set log level based on user input
	var level zapcore.Level
	switch logs.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}

	// Customize time format
	EncodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	// Customize level encoding
	LevelEncoder := zapcore.LowercaseLevelEncoder
	EncodeCaller := zapcore.FullCallerEncoder
	if logs.Format == "console" {
		// Console log format customization
		EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
		}
		LevelEncoder = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + level.CapitalString() + "]")
		}
		EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(caller.String())
		}
	}

	// Encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   LevelEncoder,
		EncodeTime:    EncodeTime,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000) // Milliseconds precision
		},
		EncodeCaller: EncodeCaller,
		EncodeName:   zapcore.FullNameEncoder,
	}

	// Choose encoder based on format type (JSON or Console)
	var format zapcore.Encoder
	if logs.Format == "json" {
		format = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		format = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Multi-write syncer for console and file output
	var console zapcore.WriteSyncer
	if logs.Console {
		console = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		console = zapcore.AddSync(&hook)
	}

	// Create core logger with encoder, output writer, and log level
	core := zapcore.NewCore(format, console, level)

	// Add caller information and stack traces for development
	caller := zap.AddCaller()

	development := zap.Development()

	// Include custom service name in logs
	filed := zap.Fields(zap.String("service_name", logs.ServiceName))

	// Construct the logger
	Write = zap.New(core, caller, development, filed)
	defer Write.Sync() // Ensure logs are flushed
	wrappedLogger = Write.WithOptions(zap.AddCallerSkip(1))
	return Write
}

// SetServiceName sets the service name for logs
// 设置日志的服务名称
func (logs *Logs) SetServiceName(ServiceName string) *Logs {
	logs.ServiceName = ServiceName
	return logs
}

// SetConsole sets whether logs should be output to the console
// 设置日志是否输出到控制台
func (logs *Logs) SetConsole(console bool) *Logs {
	logs.Console = console
	return logs
}

// SetMaxAge sets the maximum age of log files
// 设置日志文件的最大存储天数
func (logs *Logs) SetMaxAge(MaxAge int) *Logs {
	logs.MaxAge = MaxAge
	return logs
}

// SetMaxBackups sets the number of backups to keep
// 设置保存的备份数量
func (logs *Logs) SetMaxBackups(MaxBackups int) *Logs {
	logs.MaxBackups = MaxBackups
	return logs
}

// SetMaxSize sets the maximum size of each log file
// 设置每个日志文件的最大大小
func (logs *Logs) SetMaxSize(MaxSize int) *Logs {
	logs.MaxSize = MaxSize
	return logs
}

// SetLevel sets the log level (debug, info, warn, error, etc.)
// 设置日志级别（debug, info, warn, error等）
func (logs *Logs) SetLevel(level string) *Logs {
	logs.Level = level
	return logs
}

// SetPath sets the path where log files should be stored
// 设置日志文件存储路径
func (logs *Logs) SetPath(path string) *Logs {
	logs.Path = path
	return logs
}

// SetFormat sets the log format (console or json)
// 设置日志格式（console或json）
func (logs *Logs) SetFormat(format string) *Logs {
	logs.Format = format
	return logs
}

// SetCompress sets whether log files should be compressed
// 设置日志文件是否需要压缩
func (logs *Logs) SetCompress(compress bool) *Logs {
	logs.Compress = compress
	return logs
}

// Debug logs a message at the Debug level
// Debug级别日志记录
func Debug(msg string, fields ...zap.Field) {
	wrappedLogger.Debug(msg, fields...)
}

// Info logs a message at the Info level
// Info级别日志记录
func Info(msg string, fields ...zap.Field) {
	wrappedLogger.Info(msg, fields...)
}

// Warn logs a message at the Warn level
// Warn级别日志记录
func Warn(msg string, fields ...zap.Field) {
	wrappedLogger.Warn(msg, fields...)
}

// Error logs a message at the Error level
// Error级别日志记录
func Error(msg string, fields ...zap.Field) {
	wrappedLogger.Error(msg, fields...)
}

// Panic logs a message at the Panic level
// Panic级别日志记录
func Panic(msg string, fields ...zap.Field) {
	wrappedLogger.Panic(msg, fields...)
}

// Fatal logs a message at the Fatal level
// Fatal级别日志记录
func Fatal(msg string, fields ...zap.Field) {
	wrappedLogger.Fatal(msg, fields...)
}

// Sync ensures that all buffered log entries are written
// 确保所有缓存的日志条目被写入
func Sync() error {
	return Write.Sync()
}

// WithRequestID adds the request ID to the logger
func WithRequestID(requestID string) *zap.Logger {
	if requestID != "" {
		return Write.With(zap.String("request_id", requestID))
	}
	return Write
}
