package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/examples/gin/model"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"io/ioutil"
)

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
	config := *flag.String("config", "./config.toml", "Configuration file path")
	//eng := ant.New().AddRemoteProvider("etcd3", "http://127.0.0.1:2379", "/test.toml").Serve(app)
	eng := ant.New(config).Serve(app)
	result := model.Admin{}
	//ant.Db().Find(&result)
	ant.Db("mysql2").Table("s_admin").Find(&result)
	//
	alog.Debug("main", zap.Any("result", result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}

	eng.Close()
}
