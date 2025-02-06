package serve

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/gin_middleware"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/os/config"
	"io/ioutil"
)

// LoadSrv Load Api service<加载API服务>
func LoadSrv() {
	gin.ForceConsoleColor()

	configPath := flag.String("config", "./config/config.toml", "Configuration file path")

	flag.Parse()

	eng := ant.New(*configPath).Serve(load())

	defer eng.Close()
}

func load() *gin.Engine {
	var app = gin.New()
	//开发者模式
	if config.GetBool("system.debug") == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}
	app.Use(gin_middleware.Recovery()).Use(gin_middleware.Logger())

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	app.POST("/error", func(c *gin.Context) {
		var test = []int{1, 2, 3}
		fmt.Println(test[4])
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	return app

}
