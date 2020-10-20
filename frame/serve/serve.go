package serve

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

var (
	Group errgroup.Group //Group ...
)

var Engine *gin.Engine //Engine ...

//New ...
type New struct {
	Server *http.Server
}

//Default service parameters
func Default(router *gin.Engine, port string) *New {
	return &New{
		&http.Server{
			Addr:              "127.0.0.1:" + port,
			ReadTimeout:       120 * time.Second, //设置秒的读超时
			WriteTimeout:      120 * time.Second, //设置秒的写超时
			ReadHeaderTimeout: 60 * time.Second,  //读取头超时
			IdleTimeout:       120 * time.Second, //空闲超时
			MaxHeaderBytes:    2097152,
			Handler:           router,
		},
	}
}

//Run the service
func (this *New) Run() *New {
	gin.ForceConsoleColor()
	Group.Go(func() error {
		return this.Server.ListenAndServe()
	})
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://" + this.Server.Addr)
	return this
}

//Wait Service waiting, multi-service situation waiting at the end
func Wait() {
	if err := Group.Wait(); err != nil {
		log.Fatal(err)
	}
}
