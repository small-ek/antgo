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

//Server ...
type Serve struct {
	Server *http.Server
}

//New service parameters
func New(router *gin.Engine, address ...string) *Serve {
	addr := ":8080"
	if len(address) > 0 {
		addr = address[0]
	}
	return &Serve{
		&http.Server{
			Addr:              addr,
			ReadTimeout:       240 * time.Second, //设置秒的读超时
			WriteTimeout:      240 * time.Second, //设置秒的写超时
			ReadHeaderTimeout: 60 * time.Second,  //读取头超时
			IdleTimeout:       120 * time.Second, //空闲超时
			MaxHeaderBytes:    2097152,           //请求头最大字节
			Handler:           router,
		},
	}
}

//SetMaxHeaderBytes
func (s *Serve) SetMaxHeaderBytes(MaxHeaderBytes int) *Serve {
	s.Server.MaxHeaderBytes = MaxHeaderBytes
	return s
}

//Run the service
func (s *Serve) Run() *Serve {
	gin.ForceConsoleColor()
	Group.Go(func() error {
		return s.Server.ListenAndServe()
	})
	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://" + s.Server.Addr)
	return s
}

//Wait Service waiting, multi-service situation waiting at the end
func Wait() {
	if err := Group.Wait(); err != nil {
		log.Fatal(err)
	}
}
