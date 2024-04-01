package aredis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

// SetExpiration 设置过期时间
// SetExpiration<毫秒>
func (c *ClientRedis) SetExpiration(key string, expiration ...int64) error {
	var ex time.Duration
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0]) * time.Millisecond
	} else {
		ex = time.Duration(0)
	}

	if c.Mode {
		return c.Clients.Expire(c.Ctx, key, ex).Err()
	}
	return c.ClusterClient.Expire(c.Ctx, key, ex).Err()
}

// ZSet 有序集合add
func (c *ClientRedis) ZSet(key string, members ...redis.Z) error {
	if c.Mode {
		return c.Clients.ZAdd(c.Ctx, key, members...).Err()
	}
	return c.ClusterClient.ZAdd(c.Ctx, key, members...).Err()
}

// ZIncrBy 增加元素的分数
func (c *ClientRedis) ZIncrBy(key string, increment float64, member string) (float64, error) {
	if c.Mode {
		return c.Clients.ZIncrBy(c.Ctx, key, increment, member).Result()
	}
	return c.ClusterClient.ZIncrBy(c.Ctx, key, increment, member).Result()
}

// ZIncrBy 获取有序集合中指定排名范围内的成员列表及其分数
func (c *ClientRedis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	if c.Mode {
		return c.Clients.ZRevRangeWithScores(c.Ctx, key, start, stop).Result()
	}
	return c.ClusterClient.ZRevRangeWithScores(c.Ctx, key, start, stop).Result()
}

// ZIncrBy 照分数范围删除有序集合中的成员
func (c *ClientRedis) ZRemRangeByScore(key, min, max string) (int64, error) {
	if c.Mode {
		return c.Clients.ZRemRangeByScore(c.Ctx, key, min, max).Result()
	}
	return c.ClusterClient.ZRemRangeByScore(c.Ctx, key, min, max).Result()
}

// ZRemRangeByRank 按照排名范围删除有序集合中的成员
func (c *ClientRedis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	if c.Mode {
		return c.Clients.ZRemRangeByRank(c.Ctx, key, start, stop).Result()
	}
	return c.ClusterClient.ZRemRangeByRank(c.Ctx, key, start, stop).Result()
}

// ZRange 返回集合中某个索引范围的元素，根据分数从小到大排序
func (c *ClientRedis) ZRange(key string, start, stop int64) ([]string, error) {
	if c.Mode {
		return c.Clients.ZRange(c.Ctx, key, start, stop).Result()
	}
	return c.ClusterClient.ZRange(c.Ctx, key, start, stop).Result()
}

// ZRevRange 返回集合中某个索引范围的元素，根据分数从小到大排序
func (c *ClientRedis) ZRevRange(key string, start, stop int64) ([]string, error) {
	if c.Mode {
		return c.Clients.ZRevRange(c.Ctx, key, start, stop).Result()
	}
	return c.ClusterClient.ZRevRange(c.Ctx, key, start, stop).Result()
}

// ZRevRank 获取有序集合中指定成员的排名（从低到高）
func (c *ClientRedis) ZRank(key string, member string) (int64, error) {
	if c.Mode {
		return c.Clients.ZRank(c.Ctx, key, member).Result()
	}
	return c.ClusterClient.ZRank(c.Ctx, key, member).Result()
}

// ZRevRank 获取有序集合中指定成员的排名（从高到低）
func (c *ClientRedis) ZRevRank(key string, member string) (int64, error) {
	if c.Mode {
		return c.Clients.ZRevRank(c.Ctx, key, member).Result()
	}
	return c.ClusterClient.ZRevRank(c.Ctx, key, member).Result()
}

// ZCard 获取有序集合的元素数量
func (c *ClientRedis) ZCard(key string) (int64, error) {
	if c.Mode {
		return c.Clients.ZCard(c.Ctx, key).Result()
	}
	return c.ClusterClient.ZCard(c.Ctx, key).Result()
}

// ZCard 获取有序集合指定成员的分数
func (c *ClientRedis) ZScore(key, member string) (float64, error) {
	if c.Mode {
		return c.Clients.ZScore(c.Ctx, key, member).Result()
	}
	return c.ClusterClient.ZScore(c.Ctx, key, member).Result()
}

// ZCount 获取有序集合指定分数范围内的成员数量
func (c *ClientRedis) ZCount(key, min, max string) (int64, error) {
	if c.Mode {
		return c.Clients.ZCount(c.Ctx, key, min, max).Result()
	}
	return c.ClusterClient.ZCount(c.Ctx, key, min, max).Result()
}

// ZRem 删除有序集合中的指定成
func (c *ClientRedis) ZRem(key string, members ...interface{}) (int64, error) {
	if c.Mode {
		return c.Clients.ZRem(c.Ctx, key, members).Result()
	}
	return c.ClusterClient.ZRem(c.Ctx, key, members).Result()
}

// ZRangeWithScores 获取有序集合指定排名范围内的成员及其分数
func (c *ClientRedis) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	if c.Mode {
		return c.Clients.ZRangeWithScores(c.Ctx, key, start, stop).Result()
	}
	return c.ClusterClient.ZRangeWithScores(c.Ctx, key, start, stop).Result()
}

// ZInterStore 计算多个有序集合的交集，并存储到新的有序集合中
func (c *ClientRedis) ZInterStore(destination string, store *redis.ZStore) (int64, error) {
	if c.Mode {
		return c.Clients.ZInterStore(c.Ctx, destination, store).Result()
	}
	return c.ClusterClient.ZInterStore(c.Ctx, destination, store).Result()
}

// ZUnionStore 计算多个有序集合的并集，并存储到新的有序集合中
func (c *ClientRedis) ZUnionStore(destination string, store *redis.ZStore) (int64, error) {
	if c.Mode {
		return c.Clients.ZUnionStore(c.Ctx, destination, store).Result()
	}
	return c.ClusterClient.ZUnionStore(c.Ctx, destination, store).Result()
}

// ZRangeByScoreWithScores 根据opt条件升序查询
func (c *ClientRedis) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	if c.Mode {
		return c.Clients.ZRangeByScoreWithScores(c.Ctx, key, opt).Result()
	}
	return c.ClusterClient.ZRangeByScoreWithScores(c.Ctx, key, opt).Result()
}

// ZRevRangeByScoreWithScores 根据opt条件降序查询
func (c *ClientRedis) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	if c.Mode {
		return c.Clients.ZRevRangeByScoreWithScores(c.Ctx, key, opt).Result()
	}
	return c.ClusterClient.ZRevRangeByScoreWithScores(c.Ctx, key, opt).Result()
}
