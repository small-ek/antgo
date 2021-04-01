package test

import (
	"github.com/small-ek/antgo/os/cron"
	"log"
	"os"
	"testing"
)

type testTask struct {
}

func (t *testTask) Run() {
	log.Println("hello world2")

}
func TestCron(t *testing.T) {
	crontab := cron.Default()
	// 实现接口的方式添加定时任务
	task := &testTask{}
	log.Println(111)
	if err := crontab.AddByID("1", "*/5 * * * * ?", task); err != nil {
		log.Printf("error to add crontab task:%s", err)
		os.Exit(-1)
	}
	// 添加函数作为定时任务
	taskFunc := func() {
		log.Println("hello world")
	}
	if err := crontab.AddByFunc("2", "0 */1 * * * ?", taskFunc); err != nil {
		log.Printf("error to add crontab task:%s", err)
		os.Exit(-1)
	}
	crontab.Start()
	defer crontab.Stop()
	select {}
}
