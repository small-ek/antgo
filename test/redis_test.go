package test

import (
	"context"
	"github.com/small-ek/antgo/aredis"
	"log"
	"testing"
)

var ctx = context.Background()

func TestRedis(t *testing.T) {
	client := aredis.New("127.0.0.1:6379", "", 0)
	client.Set("key", "value")
	log.Println(client.Get("key"))
	log.Println(client.TTL("key"))
	client.Push("list_test", "message1")
	client.Push("list_test", "message2")
	log.Println(client.GetList("list_test"))
}
