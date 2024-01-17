package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/fiber"
	"log"
)

func main() {
	app := fiber.New()

	app.Get("/api/*", func(c fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	}).Name("api")
	log.Fatal(app.Listen(":3000"))
	//config := *flag.String("config", "./config.toml", "Configuration file path")
	eng := ant.New().Etcd([]string{"127.0.0.1:2379"}, "/test.toml", "", "").Serve(app)

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

	defer eng.Close()
}
