package aredis

// AddSet value<修改集合>
func (c *ClientRedis) AddSet(key string, members ...interface{}) int64 {
	var count int64

	if c.Mode {
		count = c.Clients.SAdd(c.Ctx, key, members).Val()
	} else {
		count = c.ClusterClient.SAdd(c.Ctx, key, members).Val()
	}

	return count
}

// GetSetLength value<获取集合长度>
func (c *ClientRedis) GetSetLength(key string) int64 {
	var count int64

	if c.Mode {
		count = c.Clients.SCard(c.Ctx, key).Val()
	} else {
		count = c.ClusterClient.SCard(c.Ctx, key).Val()
	}

	return count
}

// GetSet value<获取集合>
func (c *ClientRedis) GetSet(key string) []string {
	var result []string
	if c.Mode {
		result = c.Clients.SMembers(c.Ctx, key).Val()
	} else {
		result = c.ClusterClient.SMembers(c.Ctx, key).Val()
	}

	return result
}

// GetSet value<获取集合>
func (c *ClientRedis) RemoveSet(key string, members ...interface{}) int64 {
	var result int64

	if c.Mode {
		result = c.Clients.SRem(c.Ctx, key, members).Val()
	} else {
		result = c.ClusterClient.SRem(c.Ctx, key, members).Val()
	}

	return result
}
