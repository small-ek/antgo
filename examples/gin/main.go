package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/db/adb/sql"
	models "github.com/small-ek/antgo/examples/gin/model"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/gin_middleware"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/utils/page"
	"io/ioutil"
)

func main() {
	app := gin.New()
	if config.GetBool("system.debug") == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}
	app.Use(gin_middleware.Recovery()).Use(gin_middleware.Logger())

	app.GET("/", func(c *gin.Context) {
		var list = []models.SysAdminUsers{}
		filters := []page.Filter{
			{Field: "username", Operator: "=", Value: "admin"},
			{Field: "username", Operator: "=", Value: ""},
		}
		err := ant.Db().Model(&models.SysAdminUsers{}).Scopes(
			sql.Filters(filters),
		).Find(&list).Error

		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, list)
	})

	app.GET("/test2", func(c *gin.Context) {
		var test = []string{"11", "22"}
		fmt.Println(test[3])
		c.String(200, config.GetString("casbin.path"))
	})
	app.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	configPath := flag.String("config", "./examples/gin/config.toml", "Configuration file path")
	eng := ant.New(*configPath).Serve(app)

	defer eng.Close()
}
