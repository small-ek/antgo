package gin_cors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors Cross-domain request
func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")                                      //允许访问所有域
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE,OPTIONS") //允许请求类型
	c.Header("Access-Control-Allow-Credentials", "true")                              //服务器是否接受浏览器发送的Cookie
	c.Header("Connection", "keep-alive")                                              //可以使一次TCP连接为同意用户的多次请求服务,提高了响应速度。
	c.Header("Access-Control-Max-Age", "3600")                                        //多少秒以后再次OPTIONS.默认60分钟
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Strict-Transport-Security", " max-age=63072000; includeSubdomains; preload")
	//放行所有OPTIONS方法
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	//处理请求
	c.Next()
}
