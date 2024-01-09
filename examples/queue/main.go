package main

import (
	"github.com/small-ek/antgo/container/queue"
	"log"
	"time"
)

func main() {
	go func() {
		var dm = queue.New()
		//添加任务
		dm.AddTask(time.Now().Add(time.Second*10), "testJob", func(args ...interface{}) {
			log.Println(args)
			dm.Stop()
		}, []interface{}{1, 2, 3})
		//什么时候结束
		time.AfterFunc(time.Second*30, func() {
			dm.Stop()
		})
		dm.Start()
	}()
}
