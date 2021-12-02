package main

import (
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/engine"
)

func main() {
	app := gin.New()
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	engine.Default().Use(app)
	app.Run(":9000") // 监听并在 0.0.0.0:8080 上启动服务
}
