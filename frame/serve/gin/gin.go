package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/frame/serve"
	"github.com/small-ek/antgo/os/alog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Gin structure value is a Gin GoAdmin adapter.
type Gin struct {
	serve.BaseAdapter
	ctx *gin.Context
	app *gin.Engine
	Srv *http.Server
}

func init() {
	ant.Register(new(Gin))
}

// Name implements the method Adapter.Name.
func (gins *Gin) Name() string {
	return "gin"
}

// SetApp implements the method Adapter.Use.
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

// Run http service<不加载配置服务>
func (eng *Gin) Run(addr string) {
	eng.Srv = &http.Server{
		Addr:    ":" + addr,
		Handler: eng.app,
	}
	fmt.Printf("  PID: %d \n", os.Getpid())
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://127.0.0.1" + eng.Srv.Addr)

	go func() {
		// 服务连接
		if err := eng.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return
}

// Close http service<关闭当前一些服务>
func (eng *Gin) Close() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	alog.Warn("Exit service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := eng.Srv.Shutdown(ctx); err != nil {
		alog.Error("Server Shutdown:" + err.Error())
	}

	return
}
