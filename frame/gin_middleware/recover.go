package gin_middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

// Recovery Catch the exception and write it to the log(捕获异常并且写入日志)
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				request := c.Request
				var body []byte
				requestBody := make(map[string]interface{})

				if request.Body != nil {
					body, err = ioutil.ReadAll(request.Body)
					if err != nil {
						alog.Write.Error("ReadAll error",
							zap.Error(err.(error)), // 类型断言
							zap.Any("stack", debug.Stack()),
						)
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}
					// 把刚刚读出来的再写进去其他地方使用没有
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
					//解析Body
					if request.Header["Content-Type"] != nil && request.Header["Content-Type"][0] == "application/json" {
						if err2 := json.Unmarshal(body, &requestBody); err != nil {
							alog.Write.Error("json.Unmarshal", zap.Error(err2))
						}
					} else if len(body) > 0 {
						for _, pair := range strings.Split(string(body), "&") {
							if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
								requestBody[kv[0]] = kv[1]
							}
						}
					}
				}
				// 请求URL
				path, _ := url.QueryUnescape(c.Request.URL.RequestURI())

				logFields := []zap.Field{
					zap.Any("ip", c.ClientIP()),
					zap.Any("path", path),
					zap.Any("request", requestBody),
					zap.Any("method", request.Method),
					zap.Any("header", request.Header),
					zap.Any("err", err),
					zap.Any("stack", debug.Stack()),
				}

				// Always log X-Request-Id separately to ensure it's unique each time
				if values, ok := c.Request.Header["X-Request-Id"]; ok {
					logFields = append(logFields, zap.Strings("X-Request-Id", values))
				}
				alog.Write.Error("Recovery from panic",
					logFields...,
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		// 继续往下处理
		c.Next()
	}
}
