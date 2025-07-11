package agin

import (
	"context"
	"github.com/gin-gonic/gin"
)

// WithContextRequestID 设置请求 ID上下文
// WithContextRequestID sets the request ID context
func WithContextRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("x-request-id")
		c.Set("request_id", requestID)

		// 将参数存入标准库的 context.Context
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx) // 更新请求的上下文
		c.Next()
	}
}
