package test

import (
	"context"
	"github.com/small-ek/antgo/db/aredis"
	"github.com/small-ek/antgo/frame/ant"
	"log"
	"testing"
)

var ctx = context.Background()

func TestRedis(t *testing.T) {
	conn := aredis.New([]map[string]interface{}{
		{"name": "redis", "address": "localhost:6379", "db": 0},
	})

	//conn := aredis.NewClusterClient([]string{"127.0.0.1:6379", "127.0.0.1:6379"}, "")
	ant.Redis().Set("key:name:aaa", "value22")
	log.Println(ant.Redis().Get("key:name"))
	log.Println(ant.Redis().TTL("key"))
	//client.PushList("list_test", "message1")
	ant.Redis().PushList("list_test", "message2")
	ant.Redis().SetList("list_test", 2, "message3")
	//client.RemoveList("list_test", "message1", 1)
	//log.Println(conn.GetListIndex("list_test", 1))
	log.Println(ant.Redis().GetList("list_test"))

	ant.Redis().AddSet("set_test", "111", "222", "77")
	log.Println(ant.Redis().GetSet("set_test"))
	ant.Redis().AddSet("set_test2", "111", "222", "3333", "444")
	//交集
	log.Println("交集")
	log.Println(ant.Redis().Clients.SInter(conn["redis"].Ctx, "set_test", "set_test2").Result())
	//并集
	log.Println("并集")
	log.Println(ant.Redis().Clients.SUnion(conn["redis"].Ctx, "set_test", "set_test2").Result())
	//差集
	log.Println("差集")
	log.Println(ant.Redis().Clients.SDiff(conn["redis"].Ctx, "set_test", "set_test2").Result())

	//订阅
	//pubsub := conn.Clients.Subscribe(conn.Ctx, "subkey")
	//_, err := pubsub.Receive(conn.Ctx)
	//if err != nil {
	//	log.Fatal("pubsub.Receive")
	//}
	//log.Println("111")
	//ch := pubsub.Channel()
	//go time.AfterFunc(10*time.Second, func() {
	//	log.Println("Publish")
	//
	//	err = conn.Clients.Publish(conn.Ctx, "subkey", "test publish 1").Err()
	//	if err != nil {
	//		log.Fatal("redisdb.Publish", err)
	//	}
	//
	//	conn.Clients.Publish(conn.Ctx, "subkey", "test publish 2")
	//})
	//for msg := range ch {
	//	log.Println("recv channel:", msg.Channel, msg.Pattern, msg.Payload)
	//}
	//log.Println("222")

	//队列
	//IncrByXX := redis.NewScript(`
	//   if redis.call("GET", KEYS[1]) ~= false then
	//       return redis.call("INCRBY", KEYS[1], ARGV[1])
	//   end
	//   return false
	//`)
	//
	//n, err := IncrByXX.Run(conn.Ctx, conn.Clients, []string{"xx_counter"}, 2).Result()
	//fmt.Println(n, err)
	//
	//err = conn.Clients.Set(conn.Ctx, "xx_counter", "40", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//n, err = IncrByXX.Run(conn.Ctx, conn.Clients, []string{"xx_counter"}, 2).Result()
	//fmt.Println(n, err)
}
