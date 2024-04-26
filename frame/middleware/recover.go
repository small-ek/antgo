package middleware

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
	"time"
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
						alog.Write.Error("ioutil.ReadAll error",
							zap.Any("err", err),
							zap.Any("stack", debug.Stack()),
						)
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}
					// 把刚刚读出来的再写进去其他地方使用没有
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
					//解析Body
					if request.Header["Content-Type"] != nil && request.Header["Content-Type"][0] == "application/json" {
						json.Unmarshal(body, &requestBody)
					} else if len(body) > 0 {
						bodyList := strings.Split(string(body), "&")
						for i := 0; i < len(bodyList); i++ {
							value := strings.Split(bodyList[i], "=")
							if len(value) >= 2 {
								requestBody[value[0]] = value[1]
							}

						}
					}
				}
				// 请求URL
				path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
				// 请求类型
				method := request.Method
				// 请求IP
				ip := c.ClientIP()

				alog.Write.Error("Recovery from panic",
					zap.Any("ip", ip),
					zap.Time("time", time.Now()),
					zap.Any("path", path),
					zap.Any("request", requestBody),
					zap.Any("method", method),
					zap.Any("header", request.Header),
					zap.Any("err", err),
					zap.Any("stack", debug.Stack()),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		// 继续往下处理
		c.Next()
	}
}
