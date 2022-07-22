package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"io/ioutil"
)

type Test struct {
	Name string
}

func main() {
	app := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	gin.DefaultWriter = ioutil.Discard

	app.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	app.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	config := *flag.String("config", "config.toml", "Configuration file path")
	eng := ant.Default().SetConfig(config).Serve(app)

	//result := model.Admin{}
	//ant.Db().Find(&result)
	////ant.Db().Table("s_admin").Find(&result)
	//
	//ant.Log().Info(conv.String(result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}

	eng.Close()
}
