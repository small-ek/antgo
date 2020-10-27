package serve

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

//Group ...
var Group errgroup.Group

//Engine ...
var Engine *gin.Engine

//New ...
type New struct {
	Server *http.Server
}

//Default service parameters
func Default(router *gin.Engine, address string) *New {
	return &New{
		&http.Server{
			Addr:              address,
			ReadTimeout:       240 * time.Second, //设置秒的读超时
			WriteTimeout:      240 * time.Second, //设置秒的写超时
			ReadHeaderTimeout: 60 * time.Second,  //读取头超时
			IdleTimeout:       120 * time.Second, //空闲超时
			MaxHeaderBytes:    2097152,           //请求头最大字节
			Handler:           router,
		},
	}
}

//Run the service
func (get *New) Run() *New {
	gin.ForceConsoleColor()
	Group.Go(func() error {
		return get.Server.ListenAndServe()
	})
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://" + get.Server.Addr)
	return get
}

//Wait Service waiting, multi-service situation waiting at the end
func Wait() {
	if err := Group.Wait(); err != nil {
		log.Fatal(err)
	}
}
