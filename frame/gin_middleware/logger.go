package gin_middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

var apiBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

var headerCache = sync.Map{}

// Logger records request logs
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := apiBufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		defer apiBufferPool.Put(buffer)

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{body: buffer, ResponseWriter: c.Writer}
		c.Writer = w

		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		responseStatus := c.Writer.Status()
		path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
		logFields := prepareLogFields(c, responseStatus, path, startTime, endTime, requestBody)

		if c.Request.Method != "OPTIONS" {
			logFields = append(logFields, zap.String("response_body", w.body.String()))

			if responseStatus > 400 && responseStatus <= 499 {
				alog.Write.Warn("HTTP Warning "+cast.ToString(responseStatus), logFields...)
			} else {
				alog.Write.Debug("HTTP Access log", logFields...)
			}
		}
	}
}

// prepareLogFields prepares the log fields
func prepareLogFields(c *gin.Context, status int, path string, startTime, endTime time.Time, requestBody []byte) []zap.Field {
	logFields := []zap.Field{
		zap.Int("status", status),
		zap.String("path", path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.String("duration_ms", endTime.Sub(startTime).String()),
	}

	// Always log X-Request-Id separately to ensure it's unique each time
	if values, ok := c.Request.Header["X-Request-Id"]; ok {
		logFields = append(logFields, zap.Strings("x_request_id", values))
	}

	// 请求头白名单
	headerWhitelist := config.GetStringSlice("log.header_whitelist")
	if len(headerWhitelist) > 0 && headerWhitelist[0] == "*" {
		logFields = append(logFields, zap.Any("header", c.Request.Header))
	} else {
		logFields = append(logFields, zap.Any("header", filterHeaders(c.Request.Header, headerWhitelist)))
	}

	logFields = append(logFields, zap.ByteString("request_body", requestBody))
	return logFields
}

// filterHeaders filters the headers based on the whitelist
func filterHeaders(headers http.Header, whitelist []string) map[string][]string {
	cacheKey := generateCacheKey(headers, whitelist)

	if cachedHeaders, ok := headerCache.Load(cacheKey); ok {
		return cachedHeaders.(map[string][]string)
	}

	filteredHeaders := make(map[string][]string, len(whitelist))
	for _, key := range whitelist {
		if key == "X-Request-Id" {
			continue
		}
		if values, ok := headers[key]; ok {
			filteredHeaders[key] = values
		}
	}

	headerCache.Store(cacheKey, filteredHeaders)
	return filteredHeaders
}

// generateCacheKey generates a unique key for the cache based on headers and whitelist
func generateCacheKey(headers http.Header, whitelist []string) string {
	var builder strings.Builder
	for _, key := range whitelist {
		if key == "X-Request-Id" {
			continue
		}
		builder.WriteString(key)
		if values, ok := headers[key]; ok {
			for _, value := range values {
				builder.WriteString(value)
			}
		}
	}
	return builder.String()
}
