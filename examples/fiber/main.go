package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/fiber"
)

func main() {
	app := fiber.New()

	app.Get("/api/*", func(c fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	}).Name("api")
	config := *flag.String("config", "./config.toml", "Configuration file path")
	ant.New(config).Serve(app)

	//eng := ant.New(config).Serve(app)
	//result := model.Admin{}
	//ant.Db().Find(&result)
	//ant.Db("mysql2").Table("s_admin").Find(&result)
	//
	//alog.Info("main", zap.Any("result", result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}
	//defer eng.Close()
}
