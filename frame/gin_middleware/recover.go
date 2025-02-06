package gin_middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

const (
	maxBodySize2      = 1024 * 4096 // 4MB 最大请求体限制
	stackTraceSkipNum = 4           // 跳过栈帧数量以精简日志
)

// Recovery 异常恢复中间件，捕获 panic 并记录结构化日志
// Recovery middleware to capture panic and log structured context
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handlePanic(c)
		c.Next()
	}
}

// handlePanic 统一处理 panic 的逻辑
// Unified panic handling logic
func handlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		// 获取精简的调用栈
		// Get simplified stack trace
		stack := debug.Stack()
		//shortStack := shortenStackTrace(stack, stackTraceSkipNum)

		// 构建请求上下文日志字段
		// Build request context log fields
		logFields := buildLogFields(c, err, stack)

		// 记录结构化日志
		// Log structured context
		alog.Write.Error("Recovery from panic", logFields...)

		// 终止请求并返回 500 状态码
		// Abort request and return 500 status
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// buildLogFields 构建日志字段
// Build log fields with request context
func buildLogFields(c *gin.Context, err interface{}, stack []byte) []zap.Field {
	request := c.Request

	// 获取安全的请求体解析
	// Safely parse request body
	bodyParams := parseRequestBody(c)

	// 处理请求路径
	// Process request path
	path, _ := url.QueryUnescape(request.URL.RequestURI())

	return []zap.Field{
		zap.String("client_ip", c.ClientIP()),
		zap.String("method", request.Method),
		zap.String("path", path),
		zap.Any("query_params", c.Request.URL.Query()),
		zap.Reflect("body_params", bodyParams),
		zap.String("panic", fmt.Sprintf("%v", err)),
		zap.Any("stack", stack),
		zap.String("x-request-id", getRequestID(c)),
	}
}

// parseRequestBody 安全解析请求体
// Safely parse request body with protection
func parseRequestBody(c *gin.Context) interface{} {
	// 根据 Content-Type 决定解析策略
	// Parse strategy based on Content-Type
	contentType := c.ContentType()

	// 限制读取的请求体大小
	// Limit request body reading size
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodySize2))
	if err != nil {
		return fmt.Sprintf("read body error: %v", err)
	}

	// 恢复请求体供后续处理
	// Restore body for subsequent processing
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	switch {
	case strings.Contains(contentType, "application/json"):
		return parseJSONBody(body)
	case strings.Contains(contentType, "form-data"),
		strings.Contains(contentType, "x-www-form-urlencoded"):
		return parseFormData(c)
	default:
		return string(body)
	}
}

// parseJSONBody 解析 JSON 格式请求体
// Parse JSON request body
func parseJSONBody(body []byte) interface{} {
	var params map[string]interface{}
	if err := json.Unmarshal(body, &params); err != nil {
		return string(body) // 返回原始内容如果解析失败
	}
	return params
}

// parseFormData 解析表单数据
// Parse form data
func parseFormData(c *gin.Context) interface{} {
	if err := c.Request.ParseForm(); err != nil {
		return "parse form error"
	}
	return c.Request.PostForm
}

// getRequestID 获取请求唯一标识
// Get request unique identifier
func getRequestID(c *gin.Context) string {
	if id := c.GetHeader("X-Request-Id"); id != "" {
		return id
	}
	return ""
}
