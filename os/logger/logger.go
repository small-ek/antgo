package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/small-ek/antgo/conv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"runtime"
	"time"
)

//Write Inherit zap log
var Write *zap.Logger

//Create 默认调用
var Create Log

//New Log parameter structure
type Log struct {
	Path        string //Save Path
	Level       string //Set log level,info debug warn
	MaxBackups  int    //Keep 30 backups, 300 by default
	MaxSize     int    //Each log file saves 10M, the default is 10M
	MaxAge      int    //7 days reserved, 30 days by default
	Compress    bool   //Whether to compress, no compression by default
	ServiceName string //Log service name, default Ginp
	Write       *zap.Logger
}

//Default setting log
func Default(path string) *Log {
	return &Log{
		Path:        path,
		MaxSize:     3,
		MaxBackups:  500,
		MaxAge:      30,
		Compress:    true,
		ServiceName: "antgo",
	}
}

//SetServiceName setting log
func (log *Log) SetServiceName(ServiceName string) *Log {
	log.ServiceName = ServiceName
	return log
}

//SetMaxAge setting log
func (log *Log) SetMaxAge(MaxAge int) *Log {
	log.MaxAge = MaxAge
	return log
}

//SetMaxBackups setting log
func (log *Log) SetMaxBackups(MaxBackups int) *Log {
	log.MaxBackups = MaxBackups
	return log
}

//SetMaxSize Set log maximum
func (log *Log) SetMaxSize(MaxSize int) *Log {
	log.MaxSize = MaxSize
	return log
}

//SetPath setting log path
func (log *Log) SetPath(path string) *Log {
	log.Path = path
	return log
}

//Register Set log
func (log *Log) Register() *zap.Logger {
	// Log split
	hook := lumberjack.Logger{
		Filename:   log.Path,
		MaxSize:    log.MaxSize,
		MaxBackups: log.MaxBackups,
		MaxAge:     log.MaxAge,
		Compress:   log.Compress,
	}
	WriteSyncer := zapcore.AddSync(&hook)
	// Set log level
	// debug->info->warn->error
	var level zapcore.Level
	switch log.Level {
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
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // Lowercase encoder
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, // ISO8601 UTC 时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
		EncodeCaller: zapcore.FullCallerEncoder, // Full path encoder
		EncodeName:   zapcore.FullNameEncoder,
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
	filed := zap.Fields(zap.String("serviceName", log.ServiceName))
	// Construction log
	log.Write = zap.New(core, caller, development, filed)
	defer log.Write.Sync()
	return log.Write
}

//Record 记录
func (log *Log) Record(msg string, data []byte) {
	log.Write.Info(msg, zap.ByteString("record", data))
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

//Errors Error printing<错误打印>
func Errors(msg string, args ...interface{}) {
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

//Error Error print log<捕获异常>
func Error(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		log.Println(file + ":" + conv.String(line) + "<" + err.Error() + ">")
		FormateLog([]interface{}{file, line, f.Name()}).Sugar().Error(err.Error())
	}
}
