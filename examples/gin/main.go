package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	app.Run(":9900") // 监听并在 0.0.0.0:8080 上启动服务
}
