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
	"sync"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

var api = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 4096))
	},
}

// Logger 记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := api.Get().(*bytes.Buffer)
		buffer.Reset()
		defer func() {
			if buffer != nil {
				api.Put(buffer)
				buffer = nil
			}
			c.Request.Body.Close()
		}()

		requestBody := buffer.Bytes()
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
		// 请求URL
		path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
		logFields := []zap.Field{
			zap.Int("status", responseStatus),
			zap.String("path", path),
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("duration_ms", endTime.Sub(startTime).String()),
		}

		if c.Request.Method != "OPTIONS" {
			logFields = append(logFields, zap.Any("request_body", requestBody))
			logFields = append(logFields, zap.Any("response_body", w.body.String()))
			logFields = append(logFields, zap.Any("header", filterHeaders(c.Request.Header, config.GetStringSlice("log.header_whitelist"))))

			if responseStatus > 400 && responseStatus <= 499 {
				alog.Write.Warn("HTTP Warning "+cast.ToString(responseStatus), logFields...)
			} else {
				alog.Write.Debug("HTTP Access log", logFields...)
			}
		}

		api.Put(buffer)
		buffer = nil
	}
}

// filterHeaders filters the headers based on the whitelist 过滤请求头白名单
func filterHeaders(headers http.Header, whitelist []string) map[string][]string {
	filteredHeaders := make(map[string][]string)
	for _, key := range whitelist {
		if values, ok := headers[key]; ok {
			filteredHeaders[key] = values
		}
	}
	return filteredHeaders
}
