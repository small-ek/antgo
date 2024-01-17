package gin

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/serve"
	"log"
)

// Fiber structure value is a Gin GoAdmin adapter.
type Fiber struct {
	serve.BaseAdapter
	ctx fiber.Ctx
	app *fiber.App
}

func init() {

	ant.Register(&Fiber{})
}

// Name implements the method Adapter.Name.
func (e *Fiber) Name() string {
	return "fiber"
}

// SetApp implements the method Adapter.Use.
func (e *Fiber) SetApp(app interface{}) error {
	var (
		eng *fiber.App
		ok  bool
	)
	if eng, ok = app.(*fiber.App); !ok {
		return errors.New("gin adapter SetApp: wrong parameter")
	}
	e.app = eng
	return nil
}

// Run http service<加载服务>
func (f *Fiber) Run(Addr string) {
	log.Fatal(f.app.Listen(Addr))
	return
}

// Close http service<关闭当前服务>
func (eng *Fiber) Close() {
	return
}
