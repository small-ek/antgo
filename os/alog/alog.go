package alog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Write Inherit zap log
var Write *zap.Logger

// Logs parameter structure
type Logs struct {
	Path        string //Save Path
	Level       string //Set log level,info debug warn
	MaxBackups  int    //Keep 30 backups, 300 by default
	MaxSize     int    //Each log file saves 10M, the default is 10M
	MaxAge      int    //7 days reserved, 180 days by default
	Compress    bool   //Whether to compress, no compression by default
	ServiceName string //Log service name, default antgo
	Format      string //Log format default console
	Console     bool   //Whether to output the console display
}

// New Default setting log
func New(path string) *Logs {
	return &Logs{
		Path:        path,
		MaxSize:     10,
		MaxBackups:  300,
		MaxAge:      180,
		Compress:    false,
		ServiceName: "antgo",
	}
}

// Register Set log
func (logs *Logs) Register() *zap.Logger {
	// Log split
	hook := lumberjack.Logger{
		Filename:   logs.Path,
		MaxSize:    logs.MaxSize,
		MaxBackups: logs.MaxBackups,
		MaxAge:     logs.MaxAge,
		Compress:   logs.Compress,
	}

	// Set log level
	// debug->info->warn->error
	//日志输出等级
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

	EncodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	LevelEncoder := zapcore.LowercaseLevelEncoder
	EncodeCaller := zapcore.FullCallerEncoder
	if logs.Format == "console" {
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

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   LevelEncoder, // Lowercase encoder
		EncodeTime:    EncodeTime,   // ISO8601 UTC 时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
		EncodeCaller: EncodeCaller, // Full path encoder
		EncodeName:   zapcore.FullNameEncoder,
	}

	var format zapcore.Encoder

	//json格式
	if logs.Format == "json" {
		format = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		format = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var console zapcore.WriteSyncer
	//输出控制台
	if logs.Console {
		console = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		console = zapcore.AddSync(&hook)
	}

	core := zapcore.NewCore(
		format,
		console,
		level,
	)
	// Open development mode, stack trace
	caller := zap.AddCaller()
	// Open document and line number
	development := zap.Development()
	// Set the initialization field, such as: add a server name
	filed := zap.Fields(zap.String("service_name", logs.ServiceName))
	// Construction log
	Write = zap.New(core, caller, development, filed)
	defer Write.Sync()
	return Write
}

// SetServiceName setting log<打印日志服务名称>
func (logs *Logs) SetServiceName(ServiceName string) *Logs {
	logs.ServiceName = ServiceName
	return logs
}

// SetConsole Whether to output the console display<是否控制台打印>
func (logs *Logs) SetConsole(console bool) *Logs {
	logs.Console = console
	return logs
}

// SetMaxAge setting log
func (logs *Logs) SetMaxAge(MaxAge int) *Logs {
	logs.MaxAge = MaxAge
	return logs
}

// SetMaxBackups setting log
func (logs *Logs) SetMaxBackups(MaxBackups int) *Logs {
	logs.MaxBackups = MaxBackups
	return logs
}

// SetMaxSize Set log maximum
func (logs *Logs) SetMaxSize(MaxSize int) *Logs {
	logs.MaxSize = MaxSize
	return logs
}

// SetLevel setting log level
func (logs *Logs) SetLevel(level string) *Logs {
	logs.Level = level
	return logs
}

// SetPath setting log path
func (logs *Logs) SetPath(path string) *Logs {
	logs.Path = path
	return logs
}

// SetFormat Log format default console
func (logs *Logs) SetFormat(format string) *Logs {
	logs.Format = format
	return logs
}

// SetCompress Do you need compression
func (logs *Logs) SetCompress(compress bool) *Logs {
	logs.Compress = compress
	return logs
}

func Debug(msg string, fields ...zap.Field) {
	Write.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Write.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Write.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Write.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Write.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Write.Fatal(msg, fields...)
}

// Sync 确保日志缓冲区内容写入输出
func Sync() error {
	return Write.Sync()
}
