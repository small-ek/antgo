package agin

import (
	"bytes"
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

// 使用更快的 jsoniter 配置以提升性能（生产可选）
// 如果有兼容性顾虑，可改回 ConfigCompatibleWithStandardLibrary
var json = jsoniter.ConfigFastest

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

// buffer 池（保留），但回收时会判断 cap，避免被超大 buffer 长期占用
var apiBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 8192)) // 8KB 初始
	},
}

// 控制 buffer 返回池的阈值（超过则丢弃）
const bufferRetainCapThreshold = 64 * 1024 // 64KB

func putBackBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	// 如果底层容量过大，丢弃以避免池被撑爆
	if buf.Cap() > bufferRetainCapThreshold {
		return
	}
	buf.Reset()
	apiBufferPool.Put(buf)
}

// 异步日志写入相关（非必须，但高并发时能显著降低阻塞）
type logEntry struct {
	level  int // 0=info,1=warn,2=error
	msg    string
	fields []zap.Field
}

var (
	logChan       chan logEntry
	logWorkerOnce sync.Once
)

// 启动异步日志工作协程，缓冲可调
func startLogWorker() {
	logWorkerOnce.Do(func() {
		// 缓冲大小可按机器/负载调整
		logChan = make(chan logEntry, 16384)
		go func() {
			for e := range logChan {
				switch e.level {
				case 2:
					alog.Write.Error(e.msg, e.fields...)
				case 1:
					alog.Write.Warn(e.msg, e.fields...)
				default:
					alog.Write.Info(e.msg, e.fields...)
				}
			}
		}()
	})
}

// 尝试异步入队，入队失败（channel 满）则回退到同步写入（保证日志不会全部丢失）
func enqueueLog(level int, msg string, fields []zap.Field) {
	// ensure worker started
	startLogWorker()
	le := logEntry{level: level, msg: msg, fields: fields}
	select {
	case logChan <- le:
		// queued
	default:
		// fallback to sync write to avoid silent丢失
		switch level {
		case 2:
			alog.Write.Error(msg, fields...)
		case 1:
			alog.Write.Warn(msg, fields...)
		default:
			alog.Write.Info(msg, fields...)
		}
	}
}

// ----------------- Logger 中间件（保留原结构，尽量少改动） -----------------
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

		// 跳过指定路径
		if exactSkipPaths[currentPath] {
			c.Next()
			return
		}
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
		// 注意：这里不再 defer 直接放回池（避免中途被提前 Put）
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

		// 预分配字段切片，减少扩容
		logFields := make([]zap.Field, 0, 16)
		prepared := prepareLogFieldsWithSlice(
			c,
			statusCode,
			path,
			startTime,
			endTime,
			requestBody,
			headerWhitelist,
			enableRequestBody,
		)
		logFields = append(logFields, prepared...)

		// 记录响应体（限制大小） / Record response body (with size limit)
		responseBody := w.body.Bytes()

		if enableResponseBody {
			var parsedBody interface{}
			// 尝试解析 JSON
			if err := json.Unmarshal(responseBody, &parsedBody); err != nil {
				parsedBody = string(responseBody)
			}
			logFields = append(logFields, zap.Any("response_body", parsedBody))
		}

		// 将 buffer 放回池（受阈值控制）
		putBackBuffer(buffer)

		// 异步写日志（channel 满时会回退为同步写）
		switch {
		case statusCode >= 500:
			enqueueLog(2, "HTTP Server Error", logFields)
		case statusCode >= 400:
			enqueueLog(1, "HTTP Client Error", logFields)
		default:
			enqueueLog(0, "HTTP Access Log", logFields)
		}
	}
}

// prepareLogFieldsWithSlice: 返回 zap.Field 切片（便于预分配）
// 把原 prepareLogFields 拆分成返回 slice 的版本，减少中间分配
func prepareLogFieldsWithSlice(
	c *gin.Context,
	status int,
	path string,
	startTime, endTime time.Time,
	requestBody []byte,
	headerWhitelist []string,
	enableRequestBody bool,
) []zap.Field {
	logFields := make([]zap.Field, 0, 16)

	// 基础字段 / Basic fields
	logFields = append(logFields,
		zap.Int("status", status),
		zap.String("path", path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
		zap.String("latency", endTime.Sub(startTime).String()),
	)

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
		parsedBody, err := parseRequestLogBody(c, requestBody, c.ContentType())
		if err != nil {
			// 解析错误写到 error logger（避免影响主日志链路）
			alog.Write.Error("parseLogBody failed", zap.Error(err))
		}
		logFields = append(logFields, zap.Any("request_body", parsedBody))
	}

	return logFields
}

// filterHeaders 基于白名单过滤请求头（保持原来行为）
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

// parseRequestLogBody 是统一入口，根据 Content-Type 分发给不同的解析器
func parseRequestLogBody(c *gin.Context, body []byte, contentType string) (interface{}, error) {
	// 轻量 Content-Type 处理：常见类型使用更快判断
	ct := strings.ToLower(strings.TrimSpace(contentType))

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return parseJSONBody(body)
	case strings.HasPrefix(ct, "application/xml"), strings.HasPrefix(ct, "text/xml"):
		return parseXMLBody(body)
	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		return parseFormURLEncodedBody(body)
	case strings.HasPrefix(ct, "multipart/form-data"):
		return parseMultipartFormDataFromContext(c)
	case strings.HasPrefix(ct, "application/octet-stream"):
		return "binary data (not parsed)", nil
	case strings.HasPrefix(ct, "image/"):
		return fmt.Sprintf("binary image data (%s)", contentType), nil
	case strings.HasPrefix(ct, "video/"):
		return fmt.Sprintf("binary video data (%s)", contentType), nil
	case strings.HasPrefix(ct, "audio/"):
		return fmt.Sprintf("binary audio data (%s)", contentType), nil
	case strings.HasPrefix(ct, "text/plain"):
		return string(body), nil
	default:
		// 未知类型，直接输出原始字符串
		return string(body), nil
	}
}

// JSON 解析
func parseJSONBody(body []byte) (interface{}, error) {
	if len(body) == 0 {
		return nil, nil
	}
	var parsed interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return string(body), err
	}
	return parsed, nil
}

// XML 解析
func parseXMLBody(body []byte) (interface{}, error) {
	if len(body) == 0 {
		return nil, nil
	}
	// 为保持稳定性，默认返回原始字符串（可以再加长度限制或尝试解析小 payload）
	return string(body), nil
}

// x-www-form-urlencoded 解析
func parseFormURLEncodedBody(body []byte) (interface{}, error) {
	if len(body) == 0 {
		return nil, nil
	}
	formData, err := url.ParseQuery(string(body))
	if err != nil {
		return string(body), err
	}
	return formData, nil
}

// multipart/form-data 解析（生产版）
// 仅输出普通字段与文件名；对大体量请求会选择省略解析以保护内存
func parseMultipartFormDataFromContext(c *gin.Context) (interface{}, error) {
	formFields := make(map[string]interface{})

	// 若 ContentLength 明确且非常大，直接省略解析（避免占用内存）
	const multipartOmitThreshold = 10 << 20 // 10MB
	if c.Request.ContentLength > multipartOmitThreshold {
		return "multipart/form-data: omitted (too large)", nil
	}

	// 解析表单，限制内存临时解析大小（用于 small-form 提取文本字段）
	// 这里使用 2MB 内存阈值用于解析表单字段和文件 metadata
	const multipartParseMem = 2 << 20 // 2MB
	if err := c.Request.ParseMultipartForm(multipartParseMem); err != nil {
		// 解析失败时不阻塞主逻辑，返回简短提示
		return "multipart/form-data: parse failed", err
	}

	form := c.Request.MultipartForm
	if form == nil {
		return "multipart/form-data: empty form", nil
	}

	// 普通字段
	for key, vals := range form.Value {
		if len(vals) > 0 {
			formFields[key] = vals[0]
		}
	}

	// 文件字段：只打印文件名，避免内存占用
	for key, files := range form.File {
		// 构造文件名切片（避免 fmt.Sprintf 带来的额外分配）
		names := make([]string, 0, len(files))
		for _, f := range files {
			names = append(names, f.Filename)
		}
		formFields[key] = map[string]interface{}{
			"files": names,
			"count": len(names),
		}
	}

	return formFields, nil
}
