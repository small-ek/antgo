package adb

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func New(zapLogger *zap.Logger) Logger {
	return Logger{
		ZapLogger:                 zapLogger,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             200 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

func (l Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}

	alog.Write.Sugar().Infof(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	alog.Write.Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	alog.Write.Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	sql = formatSQLForLog(sql)
	logFields := []zap.Field{
		zap.String("line", utils.FileWithLineNum()),
		zap.String("sql", sql),
		zap.Float64("duration", float64(elapsed.Nanoseconds())/1e6),
		zap.Int64("rows", rows),
	}

	requestID := ctx.Value("request_id")
	if requestID != nil {
		logFields = append(logFields, zap.String("request_id", requestID.(string)))
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.ZapLogger.Error("sql_error", logFields...)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		l.ZapLogger.Warn("sql_warn", logFields...)
	case l.LogLevel >= gormlogger.Info:
		l.ZapLogger.Info("sql_info", logFields...)
	}
}

func formatSQLForLog(sql string) string {
	var builder strings.Builder
	builder.Grow(len(sql))

	inSingleQuote := false
	inDoubleQuote := false
	lastSpace := false

	writeSpace := func() {
		if builder.Len() > 0 && !lastSpace {
			builder.WriteByte(' ')
			lastSpace = true
		}
	}

	for i := 0; i < len(sql); i++ {
		ch := sql[i]

		if inSingleQuote {
			switch ch {
			case '\'':
				builder.WriteByte(ch)
				lastSpace = false
				if i+1 < len(sql) && sql[i+1] == '\'' {
					i++
					builder.WriteByte(sql[i])
					continue
				}
				inSingleQuote = false
			case '\r', '\n', '\t':
				builder.WriteByte(' ')
				lastSpace = false
			default:
				builder.WriteByte(ch)
				lastSpace = false
			}
			continue
		}

		switch ch {
		case '\'':
			inSingleQuote = true
			builder.WriteByte(ch)
			lastSpace = false
		case '"':
			if inDoubleQuote && i+1 < len(sql) && sql[i+1] == '"' {
				i++
				continue
			}
			inDoubleQuote = !inDoubleQuote
		case '\r', '\n', '\t', ' ':
			writeSpace()
		default:
			builder.WriteByte(ch)
			lastSpace = false
		}
	}

	return strings.TrimSpace(builder.String())
}
