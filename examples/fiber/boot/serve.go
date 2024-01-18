package boot

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/small-ek/antgo/frame/ant"
)

func Serve() {
	app := fiber.New()

	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	}).Name("api")
	config := *flag.String("config", "./config.toml", "Configuration file path")
	eng := ant.New(config).Serve(app)
	defer eng.Close()
}
