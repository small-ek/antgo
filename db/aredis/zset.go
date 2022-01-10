package aredis

import (
	"github.com/go-redis/redis/v8"
)

//AddSet value<修改集合>
func (c *Client) AddZset(key string, value []*redis.Z) int64 {
	var err error
	var count int64
	if c.Mode {
		count, err = c.Clients.ZAddNX(c.Ctx, key, value...).Result()
	} else {
		count, err = c.ClusterClient.ZAddNX(c.Ctx, key, value...).Result()
	}

	if err != nil {
		panic(err)
	}
	return count
}

//GetSetLength value<获取集合长度>
func (c *Client) GetZsetLength(key string) int64 {
	var err error
	var count int64

	if c.Mode {
		count, err = c.Clients.SCard(c.Ctx, key).Result()
	} else {
		count, err = c.ClusterClient.SCard(c.Ctx, key).Result()
	}

	if err != nil {
		panic(err)
	}
	return count
}

//GetZsetScore value<获取集合>
func (c *Client) GetZsetMember(key string, score string) float64 {
	var result float64
	var err error
	if c.Mode {
		result, err = c.Clients.ZScore(c.Ctx, key, score).Result()
	} else {
		result, err = c.ClusterClient.ZScore(c.Ctx, key, score).Result()
	}
	if err != nil {
		panic(err)
	}
	return result
}

//GetZsetScore value<获取集合>
func (c *Client) GetZsetScore(key string, member string) int64 {
	var result int64
	var err error
	if c.Mode {
		result, err = c.Clients.ZRank(c.Ctx, key, member).Result()
	} else {
		result, err = c.ClusterClient.ZRank(c.Ctx, key, member).Result()
	}

	if err != nil {
		panic(err)
	}
	return result
}

//GetZsetRange value<获取有序集合>
func (c *Client) GetZsetRange(key string, start, stop int64) []string {
	var result []string
	var err error
	if c.Mode {
		result, err = c.Clients.ZRange(c.Ctx, key, start, stop).Result()
	} else {
		result, err = c.ClusterClient.ZRange(c.Ctx, key, start, stop).Result()
	}

	if err != nil {
		panic(err)
	}
	return result
}

//GetZsetRange value<返回有序集合指定区间内的成员分数从高到低>
func (c *Client) GetZsetRevRange(key string, start, stop int64) []string {
	var result []string
	var err error
	if c.Mode {
		result, err = c.Clients.ZRevRange(c.Ctx, key, start, stop).Result()
	} else {
		result, err = c.ClusterClient.ZRevRange(c.Ctx, key, start, stop).Result()
	}

	if err != nil {
		panic(err)
	}
	return result
}

//GetZsetRange value<返回有序集合指定区间内的成员分数从高到低>
func (c *Client) GetZsetRangeByScore(key string, opt *redis.ZRangeBy) []string {
	var result []string
	var err error
	if c.Mode {
		result, err = c.Clients.ZRangeByScore(c.Ctx, key, opt).Result()
	} else {
		result, err = c.ClusterClient.ZRangeByScore(c.Ctx, key, opt).Result()
	}

	if err != nil {
		panic(err)
	}
	return result
}

//RemoveZset value<获取集合>
func (c *Client) RemoveZset(key string, members ...interface{}) int64 {
	var result int64
	var err error
	if c.Mode {
		result, err = c.Clients.SRem(c.Ctx, key, members).Result()
	} else {
		result, err = c.Clients.SRem(c.Ctx, key, members).Result()
	}

	if err != nil {
		panic(err)
	}
	return result
}
