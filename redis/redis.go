package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ctx = context.Background()

type Option struct {
	Addr     string //地址
	Password string //no password set
	DB       int    //use default DB
	Clients  *redis.Client
}

func Default(Addr, Password string, DB int) *Option {
	return &Option{
		Addr:     Addr,
		Password: Password, // no password set
		DB:       DB,       // use default DB
	}
}
func (this *Option) Client() Option {
	client := redis.NewClient(&redis.Options{
		Addr:     this.Addr,
		Password: this.Password, // no password set
		DB:       this.DB,       // use default DB
	})
	_, err := client.Ping(ctx).Result()

	if err != nil {
		log.Println("Redis加载失败:" + err.Error())
	}
	return Option{
		Clients: client,
	}
}

//获取值
func Get(key string) string {
	var Option = new(Option).Client().Clients
	var result, err = Option.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

//修改Redis
func Set(key string, value interface{}, expiration ...int64) {
	var ex time.Duration
	var Option = new(Option).Client().Clients
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0])
	} else {
		ex = time.Duration(0)
	}

	err := Option.Set(ctx, key, value, ex).Err()
	if err != nil {
		log.Println("redis修改值错误或者配置不正确:" + err.Error())
	}
}
