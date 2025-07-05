package acron

import (
	"context"
)

// GetRequestID 从上下文中提取请求 ID
// GetRequestID retrieves request ID stored in context
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
