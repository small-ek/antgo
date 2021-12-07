package engine

import (
	"github.com/small-ek/antgo/frame/serve"
	"log"
	"sync"
)

// Engine is the core component of antgo.
type Engine struct {
	Adapter      serve.WebFrameWork
	announceLock sync.Once
}

var engine *Engine
var defaultAdapter serve.WebFrameWork

// Default return the default engine instance.
func Default() *Engine {
	log.Println(defaultAdapter)
	engine = &Engine{
		Adapter: defaultAdapter,
	}
	return engine
}

// Use enable the adapter.
func (eng *Engine) Use(router interface{}) error {
	log.Println(eng)
	log.Println(eng.Adapter == nil)
	if eng.Adapter == nil {
		panic("adapter is nil")
	}

	return eng.Adapter.SetApp(router)
}

// Register set default adapter of engine.
func Register(ada serve.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}
