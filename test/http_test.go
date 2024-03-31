package test

import (
	"flag"
	"github.com/small-ek/antgo/net/ahttp"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/ants"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	var (
		count int
		//mutex sync.RWMutex
	)
	//var wg sync.WaitGroup
	RegisterLog()

	//numWorkers := 1000

	ants.InitPool(10)
	defer ants.NewPool.Release()
	var result, err = ahttp.Client().Debug(true).SetBody(map[string]interface{}{"test": "test"}).SetHeader(map[string]string{"Content-Type": "text/html; charset=utf-8", "test2": "test2"}).SetCookie(map[string]string{"test": "test"}).Get("http://127.0.0.1:3000/user")

	var _, _ = ahttp.Client().Debug(true).SetBody(map[string]interface{}{"test": "test"}).SetHeader(map[string]string{"Content-Type": "text/html; charset=utf-8", "test2": "test2"}).SetCookie(map[string]string{"test": "test"}).Post("http://127.0.0.1:3000/user")
	alog.Write.Info("123", zap.String("result", conv.String(result)), zap.Error(err), zap.Int("num", count))

	//for i := 0; i < numWorkers; i++ {
	//	wg.Add(1)
	//	//time.Sleep(100 * time.Nanosecond)
	//	_ = ants.NewPool.Submit(func() {
	//		defer wg.Done()
	//		mutex.Lock()
	//		count++
	//		mutex.Unlock()
	//
	//		//alog.Write.Info("123", zap.Int("num", count))
	//
	//		var result, err = ahttp.Client().Debug(true).SetBody(map[string]interface{}{"test": "test"}).SetCookie(map[string]string{"test": "test"}).SetHeader(map[string]string{"test": "test", "Content-Type": "application/json"}).Get("http://127.0.0.1:3000/user")
	//		alog.Write.Info("123", zap.String("result", conv.String(result)), zap.Error(err), zap.Int("num", count))
	//		time.Sleep(1000)
	//	})
	//}
	//wg.Wait()
	//
	//fmt.Println("run go num: ", ants.NewPool.Running())
}

func RegisterLog() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	alog.New("./log/ek2.log").SetServiceName("api").Register()
}
