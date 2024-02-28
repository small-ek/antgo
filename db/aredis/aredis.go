package aredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Client parameter structure
type ClientRedis struct {
	Options       redis.Options
	Clients       *redis.Client
	Ctx           context.Context
	ClusterClient *redis.ClusterClient
	Mode          bool
}

var once sync.Once
var Client map[string]*ClientRedis

// New setting redis
func New(list []map[string]any) map[string]*ClientRedis {
	once.Do(func() {
		if Client == nil {
			Client = make(map[string]*ClientRedis)
		}
		var ctx = context.Background()
		if len(list) > 0 {
			for i := 0; i < len(list); i++ {
				row := list[i]
				Addr, Password, DB, Name := row["address"].(string), row["password"].(string), conv.Int(row["db"]), row["name"].(string)
				if Addr != "" && DB >= 0 && Name != "" {
					options := redis.Options{
						Addr:     Addr,     //Address
						Password: Password, // no password set
						DB:       DB,       // use default DB
					}
					client := redis.NewClient(&options)
					_, err := client.Ping(ctx).Result()
					if err != nil {
						alog.Panic("redis error:", zap.Error(err))
					}

					Client[Name] = &ClientRedis{
						Mode:    true,
						Options: options,
						Clients: client,
						Ctx:     ctx,
					}
				}
			}
		}

	})
	return Client
}

// NewClusterClient <Redis集群>
func NewClusterClient(Addrs []string, Password string) *ClientRedis {
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
		alog.Panic("NewFailoverClient", zap.Error(err))
	}

	return &ClientRedis{
		Mode:          false,
		ClusterClient: client,
		Ctx:           ctx,
	}
}

// NewFailoverClient <Redis哨兵>
func NewFailoverClient(SentinelAddrs []string, MasterName, Password string, Db int) *ClientRedis {
	var ctx = context.Background()
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    MasterName,
		SentinelAddrs: SentinelAddrs,
		Password:      Password,
		DB:            Db,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		alog.Panic("NewFailoverClient", zap.Error(err))
	}

	return &ClientRedis{
		Mode:    false,
		Clients: client,
		Ctx:     ctx,
	}
}

// SetOptions <修改配置>
func (c *ClientRedis) SetOptions(Options *redis.Options) *ClientRedis {
	c.Options = *Options
	return c
}

// Close <关闭>
func (c *ClientRedis) Close() {
	if c.Mode {
		defer c.Clients.Close()
	} else {
		defer c.ClusterClient.Close()
	}
}

// Ping <心跳>
func (c *ClientRedis) Ping() string {
	var pong string

	if c.Mode {
		pong = c.Clients.Ping(c.Ctx).Val()
	} else {
		pong = c.ClusterClient.Ping(c.Ctx).Val()
	}

	return pong
}

// TTL<获取过期时间>
func (c *ClientRedis) TTL(key string) time.Duration {
	var result time.Duration

	if c.Mode {
		result = c.Clients.TTL(c.Ctx, key).Val()
	} else {
		result = c.ClusterClient.TTL(c.Ctx, key).Val()
	}

	return result
}
