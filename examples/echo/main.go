package main

import (
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

	eng := ant.Default().SetConfig("config.toml").Serve(app)
	eng.Close()
}
