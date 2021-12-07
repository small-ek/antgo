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
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = ioutil.Discard

	app := gin.Default()
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	eng := engine.Default()
	if err := eng.Use(app); err != nil {
		panic(err)
	}
	
	app.Run(":9033")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
}
