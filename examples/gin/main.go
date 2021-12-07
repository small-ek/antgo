package main

import (
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/engine"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	app := gin.New()
	//app.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	eng := engine.Default()
	if err := eng.Use(app); err != nil {
		panic(err)
	}

	go func() {
		_ = app.Run(":9033")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	app.Run(":9000") // 监听并在 0.0.0.0:8080 上启动服务
}
