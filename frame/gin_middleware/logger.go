package gin_middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// responseBodyWriter 用于捕获响应体 / responseBodyWriter for capturing response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应数据并捕获 / Write response data and capture
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

var apiBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096)) // 4KB初始缓冲区 / 4KB initial buffer
	},
}

// Logger 记录HTTP请求日志中间件
// Logger middleware for recording HTTP request logs
func Logger() gin.HandlerFunc {
	// 初始化配置 / Initialize configuration
	headerWhitelist := config.GetStringSlice("log.header_whitelist")
	maxBodySize := config.GetInt("log.max_body_size")         // 最大请求体大小 / max request body size
	skipMethods := config.GetStringSlice("log.skip_methods")  // 跳过日志记录的方法 / methods to skip logging
	enableResponseBody := config.GetBool("log.response_body") // 是否启用Debug日志 / enable debug logs

	// 转换跳过方法为map提高查询效率 / Convert skip methods to map for faster lookup
	skipMethodsMap := make(map[string]bool, len(skipMethods))
	for _, m := range skipMethods {
		skipMethodsMap[strings.ToUpper(m)] = true
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		// 跳过指定HTTP方法 / Skip specified HTTP methods
		if skipMethodsMap[c.Request.Method] {
			c.Next()
			return
		}

		// 读取请求体（限制大小） / Read request body (with size limit)
		var requestBody []byte
		if c.Request.Body != nil {
			if maxBodySize > 0 {
				requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, int64(maxBodySize)))
			} else {
				requestBody, _ = io.ReadAll(c.Request.Body)
			}
			// 重新构造 c.Request.Body 以便后续的中间件或处理函数使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		// 获取响应体缓冲区 / Get response body buffer
		buffer := apiBufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		defer apiBufferPool.Put(buffer)

		// 包装响应写入器 / Wrap response writer
		w := &responseBodyWriter{
			body:           buffer,
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		// 处理请求 / Process request
		c.Next()
		endTime := time.Now()

		// 准备日志字段 / Prepare log fields
		statusCode := c.Writer.Status()
		path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
		logFields := prepareLogFields(
			c,
			statusCode,
			path,
			startTime,
			endTime,
			requestBody,
			headerWhitelist,
			maxBodySize,
		)

		// 记录响应体（限制大小） / Record response body (with size limit)
		responseBody := w.body.Bytes()
		if maxBodySize > 0 && len(responseBody) > maxBodySize {
			responseBody = responseBody[:maxBodySize]
		}
		if enableResponseBody {
			logFields = append(logFields, zap.ByteString("response_body", responseBody))
		}

		// 根据状态码记录不同级别日志 / Log different levels based on status code
		switch {
		case statusCode >= 500:
			alog.Write.Error("HTTP Server Error", logFields...)
		case statusCode >= 400:
			alog.Write.Warn("HTTP Client Error", logFields...)
		default:
			alog.Write.Info("HTTP Access Log", logFields...)
		}
	}
}

// prepareLogFields 准备日志字段
// prepareLogFields prepares log fields
func prepareLogFields(
	c *gin.Context,
	status int,
	path string,
	startTime, endTime time.Time,
	requestBody []byte,
	headerWhitelist []string,
	maxBodySize int,
) []zap.Field {
	// 基础字段 / Basic fields
	logFields := []zap.Field{
		zap.Int("status", status),
		zap.String("path", path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
		zap.String("latency", endTime.Sub(startTime).String()),
	}

	// 错误信息 / Error messages
	if len(c.Errors) > 0 {
		logFields = append(logFields, zap.Strings("errors", c.Errors.Errors()))
	}

	// 请求头处理 / Process headers
	if len(headerWhitelist) > 0 && headerWhitelist[0] == "*" {
		logFields = append(logFields, zap.Any("headers", c.Request.Header))
	} else {
		logFields = append(logFields, zap.Any("headers", filterHeaders(c.Request.Header, headerWhitelist)))
	}

	// 单独记录X-Request-Id / Record X-Request-Id separately
	if values := c.GetHeader("X-Request-Id"); values != "" {
		logFields = append(logFields, zap.String("x_request_id", values))
	}

	// 请求体处理 / Process request body
	if len(requestBody) > maxBodySize {
		requestBody = requestBody[:maxBodySize]
	}
	logFields = append(logFields, zap.ByteString("request_body", requestBody))

	return logFields
}

// filterHeaders 基于白名单过滤请求头
// filterHeaders filters headers based on whitelist
func filterHeaders(headers http.Header, whitelist []string) map[string][]string {
	filtered := make(map[string][]string)
	for _, key := range whitelist {
		// 跳过已单独处理的字段 / Skip separately processed fields
		if strings.EqualFold(key, "X-Request-Id") {
			continue
		}

		if values := headers.Values(key); len(values) > 0 {
			filtered[key] = values
		}
	}
	return filtered
}
