package agin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
)

const (
	maxBodySize = 4 << 20 // 4MB
	bodyKey     = "recovery_body"
)

// Recovery 捕获 panic 并记录结构化日志
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheRequestBody(c)
		defer handlePanic(c)
		c.Next()
	}
}

// cacheRequestBody 预读取并缓存请求体
func cacheRequestBody(c *gin.Context) {
	if c.Request == nil || c.Request.Body == nil {
		return
	}

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodySize))
	if err != nil {
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	c.Set(bodyKey, body)
}

// handlePanic 统一处理 panic
func handlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		stack := debug.Stack()
		fields := buildLogFields(c, err, stack)
		alog.Write.Error("Recovery from panic", fields...)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// buildLogFields 构建日志字段
func buildLogFields(c *gin.Context, err interface{}, stack []byte) []zap.Field {
	return []zap.Field{
		zap.String("client_ip", c.ClientIP()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Any("headers", c.Request.Header),
		zap.Any("query", c.Request.URL.Query()),
		zap.Any("request_body", parseRequestBody(c)),
		zap.Any("panic", err),
		zap.String("panic_at", extractPanicLocation(stack)), // 新增:精准定位
		zap.Strings("stack", SplitStack(stack)),             // 改进:数组形式
		zap.String("request_id", getRequestID(c)),
	}
}

// extractPanicLocation 提取业务代码中的 panic 位置
func extractPanicLocation(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))

	for i := 0; i < len(lines)-1; i++ {
		line := strings.TrimSpace(string(lines[i]))
		nextLine := strings.TrimSpace(string(lines[i+1]))

		// 跳过空行
		if line == "" || nextLine == "" {
			continue
		}

		// 检查下一行(文件路径)是否为框架代码
		if isFrameworkCode(nextLine) {
			continue
		}

		// 检查当前行是否为函数调用(通常以包名开头或包含括号)
		if strings.Contains(line, "(") || strings.Contains(nextLine, ".go:") {
			return fmt.Sprintf("%s at %s", line, nextLine)
		}
	}

	return "unknown"
}

// isFrameworkCode 判断是否为框架/标准库代码
func isFrameworkCode(line string) bool {
	return strings.Contains(line, "runtime/") || // 标准库
		strings.Contains(line, "vendor/") || // vendor 模式
		strings.Contains(line, "/pkg/mod/") // modules 模式
}

// SplitStack 分割堆栈为数组(仅保留关键帧)
func SplitStack(stack []byte) []string {
	lines := bytes.Split(stack, []byte("\n"))
	var result []string
	for i := 0; i < len(lines) && i < 20; i++ { // 限制 10 帧
		if line := strings.TrimSpace(string(lines[i])); line != "" {
			result = append(result, line)
		}
	}
	return result
}

// parseRequestBody 解析请求体
func parseRequestBody(c *gin.Context) interface{} {
	// 优先使用缓存的请求体
	if v, ok := c.Get(bodyKey); ok {
		if body, ok := v.([]byte); ok {
			return parseBodyByContentType(c, body)
		}
	}

	// 降级方案:实时读取
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodySize))
	if err != nil {
		return fmt.Sprintf("read error: %v", err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return parseBodyByContentType(c, body)
}

// parseBodyByContentType 根据 Content-Type 解析请求体
func parseBodyByContentType(c *gin.Context, body []byte) interface{} {
	contentType := c.ContentType()

	switch {
	case strings.Contains(contentType, "application/json"):
		result, err := parseJSONBody(body)
		if err != nil {
			return fmt.Sprintf("json error: %v", err)
		}
		return result

	case strings.Contains(contentType, "form-data"),
		strings.Contains(contentType, "x-www-form-urlencoded"):
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return parseFormData(c)

	default:
		return string(body)
	}
}

// parseFormData 解析表单数据
func parseFormData(c *gin.Context) interface{} {
	if err := c.Request.ParseForm(); err != nil {
		return fmt.Sprintf("form error: %v", err)
	}
	return c.Request.PostForm
}

// getRequestID 获取请求 ID
func getRequestID(c *gin.Context) string {
	return c.GetHeader("X-Request-Id")
}
