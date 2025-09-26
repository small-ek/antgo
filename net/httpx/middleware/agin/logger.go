package agin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/small-ek/antgo/net/httpx"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// json is the shared instance for JSON serialization/deserialization.
// Using this shared instance prevents redundant initialization.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	skipMethods := config.GetStringSlice("log.skip_methods")  // 跳过日志记录的方法 / methods to skip logging
	skipPaths := config.GetStringSlice("log.skip_paths")      // 跳过日志的路由路径
	enableRequestBody := config.GetBool("log.request_body")   // 是否启用请求体Body
	enableResponseBody := config.GetBool("log.response_body") // 是否启用Debug日志 / enable debug logs

	// 转换跳过方法为map提高查询效率 / Convert skip methods to map for faster lookup
	skipMethodsMap := make(map[string]bool, len(skipMethods))
	for _, m := range skipMethods {
		skipMethodsMap[strings.ToUpper(m)] = true
	}

	// 预处理跳过路径：分离精确匹配和前缀匹配
	exactSkipPaths := make(map[string]bool)
	prefixSkipPaths := []string{}
	for _, path := range skipPaths {
		if path == "" {
			continue
		}
		if strings.HasSuffix(path, "/*") {
			// 前缀匹配：移除末尾的"/*"
			prefix := strings.TrimSuffix(path, "/*")
			if prefix != "" {
				prefixSkipPaths = append(prefixSkipPaths, prefix)
			}
		} else {
			// 精确匹配
			exactSkipPaths[path] = true
		}
	}

	return func(c *gin.Context) {
		startTime := time.Now()
		currentPath := c.Request.URL.Path

		// 跳过指定HTTP方法 / Skip specified HTTP methods
		if skipMethodsMap[c.Request.Method] {
			c.Next()
			return
		}

		// 2. 跳过指定路径（新增逻辑）
		// 精确匹配检查
		if exactSkipPaths[currentPath] {
			c.Next()
			return
		}
		// 前缀匹配检查
		for _, prefix := range prefixSkipPaths {
			if strings.HasPrefix(currentPath, prefix) {
				c.Next()
				return
			}
		}

		// 读取请求体（限制大小） / Read request body (with size limit)
		maxSize := httpx.CalculateMaxSize(c.Request.ContentLength)
		requestBody, newRC, err := httpx.ReadBody(c.Request.Body, maxSize)
		if err != nil {
			alog.Write.Error("Read request body failed", zap.Error(err))
		}
		// 重新构造 c.Request.Body 以便后续的中间件或处理函数使用
		c.Request.Body = newRC
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
			enableRequestBody,
		)

		// 记录响应体（限制大小） / Record response body (with size limit)
		responseBody := w.body.Bytes()

		if enableResponseBody {
			var parsedBody interface{}

			// 尝试解析 JSON
			if err := json.Unmarshal(responseBody, &parsedBody); err != nil {
				// 不是 JSON，直接作为字符串保存
				parsedBody = string(responseBody)
			}
			logFields = append(logFields, zap.Any("response_body", parsedBody))
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
	enableRequestBody bool,
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
		logFields = append(logFields, zap.String("request_id", values))
	}

	// 请求体处理 / Process request body
	if enableRequestBody {
		parsedBody, err := parseRequestLogBody(requestBody, c.ContentType())
		if err != nil {
			alog.Write.Error("parseLogBody failed", zap.Error(err))
		}
		logFields = append(logFields, zap.Any("request_body", parsedBody))
	}

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

// parseRequestLogBody 解析请求体或者参数，支持多种 Content-Type 格式
func parseRequestLogBody(body []byte, contentType string) (interface{}, error) {
	var parsedBody interface{}

	switch {
	case strings.Contains(contentType, "application/json"):
		if err := json.Unmarshal(body, &parsedBody); err != nil {
			return string(body), err
		}
	case strings.Contains(contentType, "application/xml"), strings.Contains(contentType, "text/xml"):
		var xmlData map[string]interface{}
		if err := xml.Unmarshal(body, &xmlData); err != nil {
			return string(body), err
		}
		parsedBody = xmlData
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		formData, err := url.ParseQuery(string(body))
		if err != nil {
			return string(body), err
		}
		parsedBody = formData
	case strings.Contains(contentType, "multipart/form-data"):
		parsedBody = "multipart/form-data: not parsed (binary/form boundary)"
	case strings.Contains(contentType, "application/octet-stream"):
		parsedBody = "binary data (not parsed)"
	case strings.HasPrefix(contentType, "image/"):
		parsedBody = fmt.Sprintf("binary image data (%s)", contentType)
	case strings.HasPrefix(contentType, "video/"):
		parsedBody = fmt.Sprintf("binary video data (%s)", contentType)
	case strings.HasPrefix(contentType, "audio/"):
		parsedBody = fmt.Sprintf("binary audio data (%s)", contentType)
	case strings.Contains(contentType, "text/plain"):
		parsedBody = string(body)
	default:
		parsedBody = string(body) // 原始输出
	}

	return parsedBody, nil
}
