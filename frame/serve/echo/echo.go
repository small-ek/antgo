package gin

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/serve"
)

// Gin structure value is a Gin GoAdmin adapter.
type Echo struct {
	serve.BaseAdapter
	ctx echo.Context
	app *echo.Echo
}

func init() {
	ant.Register(&Echo{})
}

// Name implements the method Adapter.Name.
func (e *Echo) Name() string {
	return "echo"
}

// SetApp implements the method Adapter.Use.
func (e *Echo) SetApp(app interface{}) error {
	var (
		eng *echo.Echo
		ok  bool
	)
	if eng, ok = app.(*echo.Echo); !ok {
		return errors.New("gin adapter SetApp: wrong parameter")
	}
	e.app = eng
	return nil
}
