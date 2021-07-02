package aredis

import "github.com/small-ek/antgo/os/logs"

//AddSet value<修改集合>
func (c *Client) AddSet(key string, members ...interface{}) int64 {
	count, err := c.Clients.SAdd(c.Ctx, key, members).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return count
}

//GetSetLength value<获取集合长度>
func (c *Client) GetSetLength(key string) int64 {
	count, err := c.Clients.SCard(c.Ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return count
}

//GetSet value<获取集合>
func (c *Client) GetSet(key string) []string {
	result, err := c.Clients.SMembers(c.Ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return result
}

//GetSet value<获取集合>
func (c *Client) RemoveSet(key string, members ...interface{}) int64 {
	result, err := c.Clients.SRem(c.Ctx, key, members).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return result
}

//SetDiff value<差集>
func (c *Client) SetDiff(keys ...string) []string {
	result, err := c.Clients.SDiff(c.Ctx, keys[0]).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return result
}
