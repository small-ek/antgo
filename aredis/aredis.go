package aredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/small-ek/antgo/os/logs"
	"time"
)

var ctx = context.Background()

//New parameter structure
type Client struct {
	Options *redis.Options
	Clients *redis.Client
}

//New setting redis
func New(Addr, Password string, DB int) *Client {
	options := &redis.Options{
		Addr:     Addr,     //Address
		Password: Password, // no password set
		DB:       DB,       // use default DB
	}
	client := redis.NewClient(options)
	_, err := client.Ping(ctx).Result()

	if err != nil {
		logs.Error(err.Error())
	}
	return &Client{
		Options: options,
		Clients: client,
	}
}

//SetOptions <修改配置>
func (c *Client) SetOptions(Options *redis.Options) *Client {
	c.Options = Options
	return c
}

//Close <关闭>
func (c *Client) Close() {
	defer c.Clients.Close()
}

//Get value
func (c *Client) Get(key string) string {
	var result, err = c.Clients.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

//Set value<设置读写>
//expiration<毫秒>
func (c *Client) Set(key string, value interface{}, expiration ...int64) {
	var ex time.Duration
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0])
	} else {
		ex = time.Duration(0)
	}
	err := c.Clients.Set(ctx, key, value, ex).Err()
	if err != nil {
		logs.Error(err.Error())
	}
}

//SetNX value<不存在才设置>
//expiration<毫秒>
func (c *Client) SetNX(key string, value interface{}, expiration ...int64) bool {
	var ex time.Duration
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0])
	} else {
		ex = time.Duration(0)
	}
	result, err := c.Clients.SetNX(ctx, key, value, ex).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return result
}

//GetList value<获取列表>
func (c *Client) GetListLength(key string) int64 {
	lens, err := c.Clients.LLen(ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return lens
}

//GetList value<获取列表>
func (c *Client) GetList(key string) []string {
	lens := c.GetListLength(key)
	list, err := c.Clients.LRange(ctx, key, 0, lens-1).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return list
}

//Push value<添加>
func (c *Client) Push(key string, value interface{}) {
	err := c.Clients.RPush(ctx, key, value).Err()
	if err != nil {
		logs.Error(err.Error())
	}
}

//Ping <心跳>
func (c *Client) Ping() string {
	pong, err := c.Clients.Ping(ctx).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return pong
}

//TTL<获取过期时间>
func (c *Client) TTL(key string) time.Duration {
	result, err := c.Clients.TTL(ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return result
}
