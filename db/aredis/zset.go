package aredis

import (
	"github.com/redis/go-redis/v9"
)

// AddSet value<修改集合>
func (c *ClientRedis) AddZset(key string, value ...redis.Z) int64 {

	var count int64
	if c.Mode {
		count = c.Clients.ZAddNX(c.Ctx, key, value...).Val()
	} else {
		count = c.ClusterClient.ZAddNX(c.Ctx, key, value...).Val()
	}

	return count
}

// GetSetLength value<获取集合长度>
func (c *ClientRedis) GetZsetLength(key string) int64 {
	var count int64

	if c.Mode {
		count = c.Clients.SCard(c.Ctx, key).Val()
	} else {
		count = c.ClusterClient.SCard(c.Ctx, key).Val()
	}

	return count
}

// GetZsetScore value<获取集合>
func (c *ClientRedis) GetZsetMember(key string, score string) float64 {
	var result float64

	if c.Mode {
		result = c.Clients.ZScore(c.Ctx, key, score).Val()
	} else {
		result = c.ClusterClient.ZScore(c.Ctx, key, score).Val()
	}

	return result
}

// GetZsetScore value<获取集合>
func (c *ClientRedis) GetZsetScore(key string, member string) int64 {
	var result int64

	if c.Mode {
		result = c.Clients.ZRank(c.Ctx, key, member).Val()
	} else {
		result = c.ClusterClient.ZRank(c.Ctx, key, member).Val()
	}

	return result
}

// GetZsetRange value<获取有序集合>
func (c *ClientRedis) GetZsetRange(key string, start, stop int64) []string {
	var result []string

	if c.Mode {
		result = c.Clients.ZRange(c.Ctx, key, start, stop).Val()
	} else {
		result = c.ClusterClient.ZRange(c.Ctx, key, start, stop).Val()
	}

	return result
}

// GetZsetRange value<返回有序集合指定区间内的成员分数从高到低>
func (c *ClientRedis) GetZsetRevRange(key string, start, stop int64) []string {
	var result []string

	if c.Mode {
		result = c.Clients.ZRevRange(c.Ctx, key, start, stop).Val()
	} else {
		result = c.ClusterClient.ZRevRange(c.Ctx, key, start, stop).Val()
	}

	return result
}

// GetZsetRange value<返回有序集合指定区间内的成员分数从高到低>
func (c *ClientRedis) GetZsetRangeByScore(key string, opt *redis.ZRangeBy) []string {
	var result []string

	if c.Mode {
		result = c.Clients.ZRangeByScore(c.Ctx, key, opt).Val()
	} else {
		result = c.ClusterClient.ZRangeByScore(c.Ctx, key, opt).Val()
	}

	return result
}

// RemoveZset value<获取集合>
func (c *ClientRedis) RemoveZset(key string, members ...interface{}) int64 {
	var result int64

	if c.Mode {
		result = c.Clients.SRem(c.Ctx, key, members).Val()
	} else {
		result = c.Clients.SRem(c.Ctx, key, members).Val()
	}

	return result
}
