package engine

import (
	"github.com/small-ek/antgo/frame"
	"sync"
)

// Engine is the core component of antgo.
type Engine struct {
	Adapter      frame.WebFrameWork
	announceLock sync.Once
}
var engine *Engine
var defaultAdapter frame.WebFrameWork

// Default return the default engine instance.
func Default() *Engine {
	engine = &Engine{
		Adapter:    defaultAdapter,
	}
	return engine
}

// Use enable the adapter.
func (eng *Engine) Use(router interface{}) error {
	if eng.Adapter == nil {
		panic("adapter is nil")
	}

	return eng.Adapter.Use(router)
}

// Register set default adapter of engine.
func Register(ada frame.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}
