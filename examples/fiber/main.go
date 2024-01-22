package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/small-ek/antgo/examples/gin/model"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/fiber"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
)

func main() {
	app := fiber.New()

	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		alog.Info("1212", zap.String("121223", "343434"))
		return c.SendString(msg) // => ✋ register
	}).Name("api")
	config := *flag.String("config", "./config.toml", "Configuration file path")
	eng := ant.New(config).AddFunc(func() {
		result := model.Admin{}
		ant.Db("mysql2").Table("admin").Find(&result)
		alog.Info("result", zap.String("12", conv.String(result)))
	}).Serve(app)

	defer eng.Close(func() {
		alog.Info("close", zap.String("121223", "343434"))
	})
	//eng := ant.New(config).Serve(app)

	//alog.Info("main", zap.Any("result", result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}
	//defer eng.Close()
}
