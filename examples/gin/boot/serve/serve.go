package serve

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/examples/gin/boot/lang"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/gin_middleware"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/i18n"
	"github.com/small-ek/antgo/os/config"
	"io/ioutil"
	"os"
	"time"
)

// LoadSrv Load Api service<加载API服务>
func LoadSrv() {
	gin.ForceConsoleColor()

	configPath := flag.String("config", "./config/config.toml", "Configuration file path")

	flag.Parse()

	eng := ant.New(*configPath).AddFunc(func() {
		lang.Register()
	}).Serve(load())

	defer eng.Close()
}

func load() *gin.Engine {
	var app = gin.New()
	//开发者模式
	if config.GetBool("system.debug") == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}
	app.Use(gin_middleware.Recovery()).Use(gin_middleware.Logger()).Use(i18n.Middleware())

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	app.GET("/query", func(c *gin.Context) {
		var list = make([]map[string]interface{}, 0)
		ant.Db().Table("sys_menu").Find(&list)
		c.JSON(200, gin.H{
			"list": list,
		})
	})

	app.POST("/error", func(c *gin.Context) {
		var test = []int{1, 2, 3}
		fmt.Println(test[4])
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	app.GET("/lang", func(c *gin.Context) {
		testDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
		c.JSON(200, gin.H{
			"正常翻译":  i18n.T(c, "common.hello", "World"),
			"正常翻译2": i18n.T(c, "common.welcome"),
			"正常翻译3": i18n.T(c, "nested.key"),
			"复数1":   i18n.TPlural(c, 1, "plural.apple"),
			"复数2":   i18n.TPlural(c, 3, "plural.apple"),
			"日期":    i18n.TDate(c, testDate),
		})
	})

	app.GET("/pid", func(c *gin.Context) {
		c.String(200, fmt.Sprintf("%d", os.Getpid()))
	})
	return app

}
