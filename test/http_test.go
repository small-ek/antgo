package test

import (
	"flag"
	"fmt"
	"github.com/small-ek/antgo/net/ahttp"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/ants"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"log"
	"sync"
	"testing"
)

func TestHttp(t *testing.T) {
	RegisterLog()

	var (
		count int
		mutex sync.Mutex
	)
	var wg sync.WaitGroup

	numWorkers := 50000

	ants.InitPool(50000)
	defer ants.NewPool.Release()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		_ = ants.NewPool.Submit(func() {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			count++

			result, err := ahttp.Client().SetBody(map[string]interface{}{"test": "test"}).SetHeader(map[string]string{"test": "test", "Content-Type": "application/json"}).SetCookie(map[string]string{"test2": "test2"}).Debug(true).Post("http://127.0.0.1:3000/user")
			alog.Write.Info("123", zap.String("result", conv.String(result)), zap.Error(err), zap.Int("num", count))

		})

	}
	wg.Wait()

	fmt.Println("run go num: ", ants.NewPool.Running())
}

func RegisterLog() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	alog.New("./log/ek2.log").SetServiceName("api").Register()
}
