package serve

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	g errgroup.Group
)

var Engine *gin.Engine

type Option struct {
	Server *http.Server
}

// 运行服务(Run the service)
func (this *Option) Run() *Option {
	gin.ForceConsoleColor()

	go func() {
		// 服务连接
		if err := this.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("服务连接失败,可能端口冲突" + err.Error())
		}
	}()

	if err := g.Wait(); err != nil {
		log.Println("启动失败,可能端口冲突请修改配置端口" + err.Error())
	}

	fmt.Println("  App running at:")
	fmt.Println("  -Local: http://" + this.Server.Addr)
	return this
}

// 服务等待,多服务情况在最后等待(Service waiting)
func (this *Option) Wait() *Option {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := this.Server.Shutdown(ctx); err != nil {
		log.Println("服务器关闭:" + err.Error())
	}
	log.Println("服务器退出")
	return this
}
