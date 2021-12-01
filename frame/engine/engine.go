package engine

import (
	"bytes"
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/os/logger"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	template2 "html/template"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

)

// Engine is the core component of antgo.
type Engine struct {
	announceLock sync.Once
}

// Default return the default engine instance.
func Default() *Engine {
	engine = &Engine{

	}
	return engine
}

// Use enable the adapter.
func (eng *Engine) Use(router interface{}) error {
	if eng.Adapter == nil {
		emptyAdapterPanic()
	}

	eng.Services.Add(auth.InitCSRFTokenSrv(eng.DefaultConnection()))
	eng.initSiteSetting()
	eng.initJumpNavButtons()
	eng.initPlugins()

	printInitMsg(language.Get("initialize success"))

	return eng.Adapter.Use(router, eng.PluginList)
}

// Register set default adapter of engine.
func Register(ada adapter.WebFrameWork) {
	if ada == nil {
		emptyAdapterPanic()
	}
	defaultAdapter = ada
}
