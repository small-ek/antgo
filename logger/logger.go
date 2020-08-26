package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Write *zap.Logger

type New struct {
	Path        string //保存路径
	Level       string //设置日志级别,info debug warn
	MaxBackups  int    //保留30个备份，默认300个
	MaxSize     int    //每个日志文件保存10M，默认 10M
	MaxAge      int    //保留7天，默认30天
	Compress    bool   //是否压缩,默认不压缩
	ServiceName string //日志服务名称,默认serviceName1
}

//默认设置日志
func Default(path string) *New {
	return &New{
		Path:        path,
		MaxSize:     10,
		MaxBackups:  300,
		MaxAge:      30,
		Compress:    false,
		ServiceName: "serviceName1",
	}
}

// logpath 日志文件路径
// loglevel 日志级别
func (this *New) Load() *zap.Logger {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   this.Path,       // 日志文件路径，默认 os.TempDir()
		MaxSize:    this.MaxSize,    // 每个日志文件保存10M，默认 100M
		MaxBackups: this.MaxBackups, // 保留30个备份，默认不限
		MaxAge:     this.MaxAge,     // 保留7天，默认不限
		Compress:   this.Compress,   // 是否压缩，默认不压缩
	}
	WriteSyncer := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	switch this.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		WriteSyncer,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("serviceName", this.ServiceName))
	// 构造日志
	Write = zap.New(core, caller, development, filed)
	defer Write.Sync()
	return Write
}

func ToJsonData(args []interface{}) zap.Field {
	det := make([]string, 0)
	if len(args) > 0 {
		for _, v := range args {
			det = append(det, fmt.Sprintf("%+v", v))
		}
	}
	zap := zap.Any("detail", det)
	return zap
}

func FormateLog(args []interface{}) *zap.Logger {
	log := Write.With(ToJsonData(args))
	return log
}

//调试打印 debug
func Debug(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Debugf(msg)
}

//打印 error
func Error(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Errorf(msg)
}

//打印 warn
func Warn(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Warn(msg)
}

//打印 Panic报错并记录
func Panic(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Panic(msg)
}

//默认打印 info
func Info(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Infof(msg)
}

//异步打印 info
func AsyncInfo(msg string, args ...interface{}) {
	go func() {
		FormateLog(args).Sugar().Infof(msg)
	}()
}
