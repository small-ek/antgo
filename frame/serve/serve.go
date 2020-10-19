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
	Group errgroup.Group
)

var Engine *gin.Engine

type Option struct {
	Server *http.Server
}

func Default(router *gin.Engine, port string) *Option {
	return &Option{
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

// 运行服务(Run the service)
func (this *Option) Run() *Option {
	gin.ForceConsoleColor()
	Group.Go(func() error {
		return this.Server.ListenAndServe()
	})
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://" + this.Server.Addr)
	return this
}

// 服务等待,多服务情况在最后等待(Service waiting)
func Wait() {
	if err := Group.Wait(); err != nil {
		log.Fatal(err)
	}
}
