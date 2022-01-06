package aredis

import "github.com/small-ek/antgo/os/logs"

//AddSet value<修改集合>
func (c *Client) AddSet(key string, members ...interface{}) int64 {
	var count int64
	var err error
	if c.Mode {
		count, err = c.Clients.SAdd(c.Ctx, key, members).Result()
	} else {
		count, err = c.ClusterClient.SAdd(c.Ctx, key, members).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return count
}

//GetSetLength value<获取集合长度>
func (c *Client) GetSetLength(key string) int64 {
	var count int64
	var err error
	if c.Mode {
		count, err = c.Clients.SCard(c.Ctx, key).Result()
	} else {
		count, err = c.ClusterClient.SCard(c.Ctx, key).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return count
}

//GetSet value<获取集合>
func (c *Client) GetSet(key string) []string {
	var result []string
	var err error
	if c.Mode {
		result, err = c.Clients.SMembers(c.Ctx, key).Result()
	} else {
		result, err = c.ClusterClient.SMembers(c.Ctx, key).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return result
}

//GetSet value<获取集合>
func (c *Client) RemoveSet(key string, members ...interface{}) int64 {
	var result int64
	var err error
	if c.Mode {
		result, err = c.Clients.SRem(c.Ctx, key, members).Result()
	} else {
		result, err = c.ClusterClient.SRem(c.Ctx, key, members).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return result
}
