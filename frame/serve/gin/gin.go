package engine

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/engine"
	"github.com/small-ek/antgo/frame/serve"
)

// Gin structure value is a Gin GoAdmin adapter.
type Gin struct {
	serve.BaseAdapter
	ctx *gin.Context
	app *gin.Engine
}

func init() {
	engine.Register(new(Gin))
}

// Use implements the method Adapter.Use.
func (gins *Gin) SetApp(app interface{}) error {
	var (
		eng *gin.Engine
		ok  bool
	)
	if eng, ok = app.(*gin.Engine); !ok {
		return errors.New("gin adapter SetApp: wrong parameter")
	}
	gins.app = eng
	return nil
}
