package aredis

import (
	"time"
)

// Get value
func (c *ClientRedis) Get(key string) string {
	var result string

	if c.Mode {
		result = c.Clients.Get(c.Ctx, key).Val()
	} else {
		result = c.ClusterClient.Get(c.Ctx, key).Val()
	}
	return result
}

// Set value<设置读写>
// expiration<毫秒>
func (c *ClientRedis) Set(key string, value interface{}, expiration ...int64) error {
	var ex time.Duration
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0]) * time.Millisecond
	} else {
		ex = time.Duration(0)
	}

	if c.Mode {
		return c.Clients.Set(c.Ctx, key, value, ex).Err()
	} else {
		return c.ClusterClient.Set(c.Ctx, key, value, ex).Err()
	}

}

// Remove value<删除数据>
func (c *ClientRedis) Remove(key string) (int64, error) {
	return c.Clients.Del(c.Ctx, key).Result()
}

// SetNX value<不存在才设置>
// expiration<毫秒>
func (c *ClientRedis) SetNX(key string, value interface{}, expiration ...int64) bool {
	var ex time.Duration
	var result bool

	if len(expiration) > 0 {
		ex = time.Duration(expiration[0]) * time.Millisecond
	} else {
		ex = time.Duration(0)
	}

	if c.Mode {
		result = c.Clients.SetNX(c.Ctx, key, value, ex).Val()
	} else {
		result = c.ClusterClient.SetNX(c.Ctx, key, value, ex).Val()
	}

	return result
}
