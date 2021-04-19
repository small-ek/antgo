package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/small-ek/antgo/os/logs"
	"time"
)

var ctx = context.Background()

//New parameter structure
type New struct {
	Addr     string //Address
	Password string //no password set
	DB       int    //use default DB
	Clients  *redis.Client
}

//Default setting redis
func Default(Addr, Password string, DB int) *New {
	return &New{
		Addr:     Addr,     //Address
		Password: Password, // no password set
		DB:       DB,       // use default DB
	}
}

//Client ...
func (get *New) Client() New {
	client := redis.NewClient(&redis.Options{
		Addr:     get.Addr,     //Address
		Password: get.Password, // no password set
		DB:       get.DB,       // use default DB
	})
	_, err := client.Ping(ctx).Result()

	if err != nil {
		logs.Error(err.Error())
	}
	return New{
		Clients: client,
	}
}

//Get value
func Get(key string) string {
	var Option = new(New).Client().Clients
	var result, err = Option.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

//Set value
func Set(key string, value interface{}, expiration ...int64) {
	var ex time.Duration
	var Option = new(New).Client().Clients
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0])
	} else {
		ex = time.Duration(0)
	}

	err := Option.Set(ctx, key, value, ex).Err()
	if err != nil {
		logs.Error(err.Error())
	}
}
