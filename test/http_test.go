package test

import (
	"flag"
	"github.com/small-ek/antgo/net/ahttp"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"log"
	"sync"
	"testing"
)

func TestHttp(t *testing.T) {
	var http = ahttp.Client()
	//var result, err = http.Debug().SetFile("test.jpg", "file").SetBody(map[string]interface{}{"name": "123.jpg"}).PostForm("http://127.0.0.1:102/upload_file")
	//log.Println(http)
	//
	//log.Println(string(result))
	//log.Println(err)
	//var result, err = http.Debug().SetHeader(map[string]string{"Content-Type": "text/html; charset=utf-8"}).Get("https://www.baidu.com/")
	//fmt.Println(result)
	//fmt.Println(err)
	RegisterLog()
	var wg sync.WaitGroup
	numWorkers := 1000
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			var result, err = http.Debug().SetHeader(map[string]string{"Content-Type": "text/html; charset=utf-8"}).Get("https://www.baidu.com/")

			var result2, err2 = http.Debug().Get("https://www.baidu.com/")

			alog.Write.Info("123", zap.Any("result", result), zap.Any("result2", result2),
				zap.Error(err), zap.Error(err2))
			// 在这里执行对 JwtManager 实例的操作
			// 例如，调用 Encrypt 或 Decode 方法
		}()
	}

	wg.Wait()
}

func RegisterLog() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	alog.New("./log/ek2.log").SetServiceName("api").Register()
}
