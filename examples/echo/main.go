package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/echo"
	"net/http"
)

func main() {

	app := echo.New()
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	app.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	config := *flag.String("config", "config.toml", "Configuration file path")
	eng := ant.Default().SetConfig("./", config).Serve(app)

	//result := model.Admin{}
	//ant.Db().Find(&result)
	//ant.Db().Table("s_admin").Find(&result)

	//ant.Log().Info(conv.String(result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}

	eng.Close()
}
