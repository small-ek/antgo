package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/os/alog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io"
	"net/url"
	"sync"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
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

		// 获取请求数据
		requestBody := buffer.Bytes()
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		// 获取 response 内容
		w := &responseBodyWriter{body: buffer, ResponseWriter: c.Writer}
		c.Writer = w
		// 设置开始时间
		startTime := time.Now()
		c.Next()

		// 开始记录日志的逻辑
		endTime := time.Now()
		responStatus := c.Writer.Status()
		// 请求URL
		path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
		logFields := []zap.Field{
			zap.Int("status", responStatus),
			zap.String("path", path),
			zap.String("method", c.Request.Method),
			zap.Any("header", c.Request.Header),
			zap.String("ip", c.ClientIP()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", endTime.Sub(startTime).String()),
		}
		isNotOptions := c.Request.Method != "OPTIONS"
		if isNotOptions {
			logFields = append(logFields, zap.Any("request-body", requestBody))
			logFields = append(logFields, zap.Any("response-body", w.body.String()))
		}

		logger := alog.Write
		if responStatus > 400 && responStatus <= 499 {
			// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意
			logger.Warn("HTTP Warning "+cast.ToString(responStatus), logFields...)
		} else {
			logger.Debug("HTTP Access Log", logFields...)
		}
		api.Put(buffer)
		c.Request.Body.Close()
	}
}
