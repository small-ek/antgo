package acron

import "go.uber.org/zap"

// zapCronLogger 适配器将 zap.Logger 转为 cron.Logger 接口
// zapCronLogger adapts zap.Logger to cron.Logger interface
type zapCronLogger struct {
	logger *zap.Logger
}

// Info 实现 cron.Logger 接口的 Info 方法
// Info logs informational messages from cron
func (l *zapCronLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Infow(msg, keysAndValues...)
}

// Error 实现 cron.Logger 接口的 Error 方法，记录错误及上下文
// Error logs errors recovered or occurred within cron
func (l *zapCronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := []interface{}{"error", err}
	fields = append(fields, keysAndValues...)
	l.logger.Sugar().Errorw(msg, fields...)
}
