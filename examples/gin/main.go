package main

import (
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"io/ioutil"
)

func main() {
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = ioutil.Discard

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	eng := ant.Default().SetConfig("config.toml").SetLog("./log/ek.log").Serve(app)
	ant.Log().Info("1212")

	eng.Close()
}
