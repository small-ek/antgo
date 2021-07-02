package aredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/small-ek/antgo/os/logs"
	"time"
)

//New parameter structure
type Client struct {
	Options       redis.Options
	Clients       *redis.Client
	Ctx           context.Context
	ClusterClient *redis.ClusterClient
	Mode          bool
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
		Mode:    true,
		Options: options,
		Clients: client,
		Ctx:     ctx,
	}
}

//NewClusterClient <Redis集群>
func NewClusterClient(Addrs []string, Password string) *Client {
	var ctx = context.Background()
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    Addrs,
		Password: Password,
		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		//RouteRandomly: true,
	})
	err := client.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		logs.Error(err.Error())
	}

	return &Client{
		Mode:          false,
		ClusterClient: client,
		Ctx:           ctx,
	}
}

//NewFailoverClient <Redis哨兵>
func NewFailoverClient(SentinelAddrs []string, MasterName, Password string, Db int) *Client {
	var ctx = context.Background()
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    MasterName,
		SentinelAddrs: SentinelAddrs,
		Password:      Password,
		DB:            Db,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		logs.Error(err.Error())
	}

	return &Client{
		Mode:    false,
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
	if c.Mode {
		defer c.Clients.Close()
	} else {
		defer c.ClusterClient.Close()
	}
}

//Ping <心跳>
func (c *Client) Ping() string {
	var pong string
	var err error

	if c.Mode {
		pong, err = c.Clients.Ping(c.Ctx).Result()
	} else {
		pong, err = c.Clients.Ping(c.Ctx).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return pong
}

//TTL<获取过期时间>
func (c *Client) TTL(key string) time.Duration {
	var result time.Duration
	var err error

	if c.Mode {
		result, err = c.Clients.TTL(c.Ctx, key).Result()
	} else {
		result, err = c.ClusterClient.TTL(c.Ctx, key).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return result
}
