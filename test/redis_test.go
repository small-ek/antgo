package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/small-ek/antgo/db/aredis"
	"log"
	"testing"
)

var ctx = context.Background()

func TestRedis(t *testing.T) {
	conn := aredis.New("127.0.0.1:6379", "", 0)
	//conn := aredis.NewClusterClient([]string{"127.0.0.1:6379", "127.0.0.1:6379"}, "")
	conn.Set("key", "value")
	log.Println(conn.Get("key"))
	log.Println(conn.TTL("key"))
	//client.PushList("list_test", "message1")
	conn.PushList("list_test", "message2")
	conn.SetList("list_test", 2, "message3")
	//client.RemoveList("list_test", "message1", 1)
	//log.Println(conn.GetListIndex("list_test", 1))
	log.Println(conn.GetList("list_test"))

	conn.AddSet("set_test", "111", "222", "77")
	log.Println(conn.GetSet("set_test"))
	conn.AddSet("set_test2", "111", "222", "3333", "444")
	//交集
	log.Println("交集")
	log.Println(conn.Clients.SInter(conn.Ctx, "set_test", "set_test2").Result())
	//并集
	log.Println("并集")
	log.Println(conn.Clients.SUnion(conn.Ctx, "set_test", "set_test2").Result())
	//差集
	log.Println("差集")
	log.Println(conn.Clients.SDiff(conn.Ctx, "set_test", "set_test2").Result())

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
	IncrByXX := redis.NewScript(`
        if redis.call("GET", KEYS[1]) ~= false then
            return redis.call("INCRBY", KEYS[1], ARGV[1])
        end
        return false
    `)

	n, err := IncrByXX.Run(conn.Ctx, conn.Clients, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)

	err = conn.Clients.Set(conn.Ctx, "xx_counter", "40", 0).Err()
	if err != nil {
		panic(err)
	}

	n, err = IncrByXX.Run(conn.Ctx, conn.Clients, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)
}

func TestRedis2(t *testing.T) {
	conn := aredis.New("127.0.0.1:6379", "", 0)
	log.Println(conn.Remove("123"))
}
