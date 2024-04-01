package aredis

// AddSet value<添加元素到集合>
func (c *ClientRedis) AddSet(key string, members ...interface{}) error {
	if c.Mode {
		return c.Clients.SAdd(c.Ctx, key, members).Err()
	}
	return c.ClusterClient.SAdd(c.Ctx, key, members).Err()
}

// GetSetLength value<获取集合长度>
func (c *ClientRedis) GetSetLength(key string) (int64, error) {
	if c.Mode {
		return c.Clients.SCard(c.Ctx, key).Result()
	}
	return c.ClusterClient.SCard(c.Ctx, key).Result()
}

// SIsMember value<检查元素是否存在>
func (c *ClientRedis) SIsMember(key string, member interface{}) (bool, error) {
	if c.Mode {
		return c.Clients.SIsMember(c.Ctx, key, member).Result()
	}
	return c.ClusterClient.SIsMember(c.Ctx, key, member).Result()
}

// GetSet value<获取集合的所有成员>
func (c *ClientRedis) GetSet(key string) ([]string, error) {
	if c.Mode {
		return c.Clients.SMembers(c.Ctx, key).Result()
	}
	return c.ClusterClient.SMembers(c.Ctx, key).Result()
}

// SRandMember value<随机获取集合的一个成员>
func (c *ClientRedis) SRandMember(key string) (string, error) {
	if c.Mode {
		return c.Clients.SRandMember(c.Ctx, key).Result()
	}
	return c.ClusterClient.SRandMember(c.Ctx, key).Result()
}

// SPop value<移除并返回集合中的一个随机元素>
func (c *ClientRedis) SPop(key string) (string, error) {
	if c.Mode {
		return c.Clients.SPop(c.Ctx, key).Result()
	}
	return c.ClusterClient.SPop(c.Ctx, key).Result()
}

// SInter value<计算多个集合的交集>
func (c *ClientRedis) SInter(keys ...string) ([]string, error) {
	if c.Mode {
		return c.Clients.SInter(c.Ctx, keys...).Result()
	}
	return c.ClusterClient.SInter(c.Ctx, keys...).Result()
}

// SUnion value<计算多个集合的并集>
func (c *ClientRedis) SUnion(keys ...string) ([]string, error) {
	if c.Mode {
		return c.Clients.SUnion(c.Ctx, keys...).Result()
	}
	return c.ClusterClient.SUnion(c.Ctx, keys...).Result()
}

// SDiff value<计算多个集合的差集>
func (c *ClientRedis) SDiff(keys ...string) ([]string, error) {
	if c.Mode {
		return c.Clients.SDiff(c.Ctx, keys...).Result()
	}
	return c.ClusterClient.SDiff(c.Ctx, keys...).Result()
}

// SRem value<移除元素>
func (c *ClientRedis) SRem(key string, members ...interface{}) (int64, error) {
	if c.Mode {
		return c.Clients.SRem(c.Ctx, key, members).Result()
	}
	return c.ClusterClient.SRem(c.Ctx, key, members).Result()
}
