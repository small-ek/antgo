package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/examples/gin/model"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"io/ioutil"
)

func main() {
	app := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	gin.DefaultWriter = ioutil.Discard

	app.GET("/", func(c *gin.Context) {

		c.String(200, config.GetString("casbin.path"))
	})
	app.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	//config := *flag.String("config", "./config.toml", "Configuration file path")
	//eng := ant.New().Etcd([]string{"127.0.0.1:2379"}, "/test.toml", "", "").Serve(app)
	config := *flag.String("config", "./config.toml", "Configuration file path")
	eng := ant.New(config).Serve(app)
	//eng := ant.New(config).Serve(app)
	result := model.Admin{}
	ant.Db("mysql2").Table("s_admin").Find(&result)
	ant.Log().Info("result", zap.String("12", conv.String(result)))
	//alog.Info("main", zap.Any("result", result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}

	defer eng.Close()
}
