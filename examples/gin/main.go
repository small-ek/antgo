package main

import (
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = ioutil.Discard

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	eng := ant.Default().SetConfig("config.toml").Use(app).Run(&http.Server{
		Addr:              ":8080",
		Handler:           app,
		ReadTimeout:       240 * time.Second, //设置秒的读超时
		WriteTimeout:      240 * time.Second, //设置秒的写超时
		ReadHeaderTimeout: 60 * time.Second,  //读取头超时
		IdleTimeout:       120 * time.Second, //空闲超时
		MaxHeaderBytes:    2097152,           //请求头最大字节
	})
	eng.Close()
}
