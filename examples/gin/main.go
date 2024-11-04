package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/gin_middleware"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/os/config"
	"io/ioutil"
)

func main() {
	app := gin.New()
	if config.GetBool("system.debug") == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}
	app.Use(requestid.New()).Use(gin_middleware.Recovery()).Use(gin_middleware.Logger())

	if config.GetBool("system.cors") == true {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"*"}
		corsConfig.AllowHeaders = []string{"*"}
		app.Use(cors.New(corsConfig))
	}

	app.GET("/", func(c *gin.Context) {
		c.String(200, config.GetString("casbin.path"))
	})

	app.GET("/test2", func(c *gin.Context) {
		var test = []string{"11", "22"}
		fmt.Println(test[3])
		c.String(200, config.GetString("casbin.path"))
	})
	app.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	configPath := flag.String("config", "./config.toml", "Configuration file path")
	eng := ant.New(*configPath).Serve(app)

	defer eng.Close()
}
