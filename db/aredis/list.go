package aredis

// GetList value<获取列表长度>
func (c *ClientRedis) GetListLength(key string) (lens int64) {
	if c.Mode {
		lens = c.Clients.LLen(c.Ctx, key).Val()
	} else {
		lens = c.ClusterClient.LLen(c.Ctx, key).Val()
	}

	return
}

// GetList value<获取列表>
func (c *ClientRedis) GetList(key string) []string {
	lens := c.GetListLength(key)
	var list []string
	if c.Mode {
		list = c.Clients.LRange(c.Ctx, key, 0, lens-1).Val()
	} else {
		list = c.ClusterClient.LRange(c.Ctx, key, 0, lens-1).Val()
	}

	return list
}

// GetListIndex value<返回名称为key的list中index位置的元素>
func (c *ClientRedis) GetListIndex(key string, index int64) (list string) {
	if c.Mode {
		list = c.Clients.LIndex(c.Ctx, key, index).Val()
	} else {
		list = c.Clients.LIndex(c.Ctx, key, index).Val()
	}
	return
}

// SetList value<修改列表>
func (c *ClientRedis) SetList(key string, index int64, value interface{}) error {
	if c.Mode {
		return c.Clients.LSet(c.Ctx, key, index, value).Err()
	}
	return c.ClusterClient.LSet(c.Ctx, key, index, value).Err()
}

// RemoveList value<删除列表>
// count 参数表示删除多少个key中的list
func (c *ClientRedis) RemoveList(key string, value interface{}, count ...int64) error {
	var counts int64 = 0
	if len(count) > 0 {
		counts = count[0]
	}

	if c.Mode {
		return c.Clients.LRem(c.Ctx, key, counts, value).Err()
	}
	return c.ClusterClient.LRem(c.Ctx, key, counts, value).Err()
}

// RemoveListLeft value<返回并删除名称为key的list中的首元素>
func (c *ClientRedis) RemoveListLeft(key string) error {
	if c.Mode {
		return c.Clients.LPop(c.Ctx, key).Err()
	}
	return c.ClusterClient.LPop(c.Ctx, key).Err()
}

// RemoveListRight value<返回并删除名称为key的list中的尾元素>
func (c *ClientRedis) RemoveListRight(key string) error {
	if c.Mode {
		return c.Clients.LPop(c.Ctx, key).Err()
	}
	return c.ClusterClient.LPop(c.Ctx, key).Err()
}

// PushList value<添加>
func (c *ClientRedis) PushList(key string, value interface{}) error {
	if c.Mode {
		return c.Clients.RPush(c.Ctx, key, value).Err()
	}
	return c.ClusterClient.RPush(c.Ctx, key, value).Err()

}
