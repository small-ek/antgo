package engine

import (
	"github.com/small-ek/antgo/frame/serve"
	"sync"
)

// Engine is the core component of antgo.
type Engine struct {
	Adapter      serve.WebFrameWork
	announceLock sync.Once
}

//defaultAdapter is the default adapter of engine.
var defaultAdapter serve.WebFrameWork

// Default return the default engine instance.
func Default() *Engine {
	return &Engine{
		Adapter: defaultAdapter,
	}
}

// Use enable the adapter.
func (eng *Engine) Use(router interface{}) error {
	if eng.Adapter == nil {
		panic("adapter is nil")
	}

	return eng.Adapter.Use(router)
}

// Register set default adapter of engine.
func Register(ada serve.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}
