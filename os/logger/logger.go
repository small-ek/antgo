package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Write Inherit zap log
var Write *zap.Logger

//Log 默认调用
var Create New

//New Log parameter structure
type New struct {
	Path        string //Save Path
	Level       string //Set log level,info debug warn
	MaxBackups  int    //Keep 30 backups, 300 by default
	MaxSize     int    //Each log file saves 10M, the default is 10M
	MaxAge      int    //7 days reserved, 30 days by default
	Compress    bool   //Whether to compress, no compression by default
	ServiceName string //Log service name, default Ginp
}

//Default setting log
func Default(path string) *New {
	return &New{
		Path:        path,
		MaxSize:     10,
		MaxBackups:  300,
		MaxAge:      30,
		Compress:    false,
		ServiceName: "Ginp",
	}
}

//SetServiceName setting log
func (get *New) SetServiceName(ServiceName string) *New {
	get.ServiceName = ServiceName
	return get
}

//SetMaxAge setting log
func (get *New) SetMaxAge(MaxAge int) *New {
	get.MaxAge = MaxAge
	return get
}

//SetMaxBackups setting log
func (get *New) SetMaxBackups(MaxBackups int) *New {
	get.MaxBackups = MaxBackups
	return get
}

//SetMaxSize setting log
func (get *New) SetMaxSize(MaxSize int) *New {
	get.MaxSize = MaxSize
	return get
}

//SetPath setting log
func (get *New) SetPath(path string) *New {
	get.Path = path
	return get
}

//SetPath Set log path
func (get *New) Register() *zap.Logger {
	// Log split
	hook := lumberjack.Logger{
		Filename:   get.Path,
		MaxSize:    get.MaxSize,
		MaxBackups: get.MaxBackups,
		MaxAge:     get.MaxAge,
		Compress:   get.Compress,
	}
	WriteSyncer := zapcore.AddSync(&hook)
	// Set log level
	// debug->info->warn->error
	var level zapcore.Level
	switch get.Level {
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
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // Lowercase encoder
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // Full path encoder
		EncodeName:     zapcore.FullNameEncoder,
	}
	// Set log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		WriteSyncer,
		level,
	)
	// Open development mode, stack trace
	caller := zap.AddCaller()
	// Open document and line number
	development := zap.Development()
	// Set the initialization field, such as: add a server name
	filed := zap.Fields(zap.String("serviceName", get.ServiceName))
	// Construction log
	Write = zap.New(core, caller, development, filed)
	defer Write.Sync()
	return Write
}

//ToJsonData ...
func ToJsonData(args []interface{}) zap.Field {
	det := make([]string, 0)
	if len(args) > 0 {
		for _, v := range args {
			det = append(det, fmt.Sprintf("%+v", v))
		}
	}
	result := zap.Any("detail", det)
	return result
}

//FormateLog ...
func FormateLog(args []interface{}) *zap.Logger {
	log := Write.With(ToJsonData(args))
	return log
}

//Debug Debug printing<调试打印>
func Debug(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Debugf(msg)
}

//Error Error printing<错误打印>
func Error(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Errorf(msg)
}

//Warn Warning print<警告打印>
func Warn(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Warn(msg)
}

//Panic Abnormal printing<异常打印>
func Panic(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Panic(msg)
}

//Info Print log by default<默认情况下打印日志>
func Info(msg string, args ...interface{}) {
	FormateLog(args).Sugar().Infof(msg)
}

//AsyncInfo Asynchronous print log<异步打印日志>
func AsyncInfo(msg string, args ...interface{}) {
	go func() {
		FormateLog(args).Sugar().Infof(msg)
	}()
}
