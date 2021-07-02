package aredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/small-ek/antgo/os/logs"
	"time"
)

//New parameter structure
type Client struct {
	Options redis.Options
	Clients *redis.Client
	Ctx     context.Context
}

//New setting redis
func New(Addr, Password string, DB int) *Client {
	var ctx = context.Background()
	options := redis.Options{
		Addr:     Addr,     //Address
		Password: Password, // no password set
		DB:       DB,       // use default DB
	}
	client := redis.NewClient(&options)
	_, err := client.Ping(ctx).Result()

	if err != nil {
		logs.Error(err.Error())
	}
	return &Client{
		Options: options,
		Clients: client,
		Ctx:     ctx,
	}
}

//SetOptions <修改配置>
func (c *Client) SetOptions(Options *redis.Options) *Client {
	c.Options = *Options
	return c
}

//Close <关闭>
func (c *Client) Close() {
	defer c.Clients.Close()
}

//Ping <心跳>
func (c *Client) Ping() string {
	pong, err := c.Clients.Ping(c.Ctx).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return pong
}

//TTL<获取过期时间>
func (c *Client) TTL(key string) time.Duration {
	result, err := c.Clients.TTL(c.Ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return result
}
