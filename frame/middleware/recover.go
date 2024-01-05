package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
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
					body, _ = ioutil.ReadAll(request.Body)
					// 把刚刚读出来的再写进去其他地方使用没有
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
					//解析Body
					if request.Header["Content-Type"] != nil && request.Header["Content-Type"][0] == "application/json" {
						json.Unmarshal(body, &requestBody)
					} else if len(body) > 0 {
						bodyList := strings.Split(string(body), "&")
						for i := 0; i < len(bodyList); i++ {
							value := strings.Split(bodyList[i], "=")
							requestBody[value[0]] = value[1]
						}
					}
				}
				// 请求URL
				path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
				// 请求类型
				method := request.Method
				// 请求IP
				ip := c.ClientIP()

				alog.Write.Error("错误报错",
					zap.Any("ip", ip),
					zap.Any("path", path),
					zap.Any("request", requestBody),
					zap.Any("method", method),
					zap.Any("header", request.Header),
					zap.Any("err", err),
					zap.Any("stack", debug.Stack()),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Println(err)
				debug.PrintStack()
			}
		}()

		// 继续往下处理
		c.Next()
	}
}
