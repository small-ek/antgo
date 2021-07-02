package aredis

import (
	"github.com/small-ek/antgo/os/logs"
	"time"
)

//Get value
func (c *Client) Get(key string) string {
	var result string
	var err error
	if c.Mode {
		result, err = c.Clients.Get(c.Ctx, key).Result()
	} else {
		result, err = c.ClusterClient.Get(c.Ctx, key).Result()
	}

	if err != nil {
		return ""
	}
	return result
}

//Set value<设置读写>
//expiration<毫秒>
func (c *Client) Set(key string, value interface{}, expiration ...int64) error {
	var ex time.Duration
	if len(expiration) > 0 {
		ex = time.Duration(expiration[0])
	} else {
		ex = time.Duration(0)
	}

	if c.Mode {
		return c.Clients.Set(c.Ctx, key, value, ex).Err()
	} else {
		return c.ClusterClient.Set(c.Ctx, key, value, ex).Err()
	}

}

//SetNX value<不存在才设置>
//expiration<毫秒>
func (c *Client) SetNX(key string, value interface{}, expiration ...int64) bool {
	var ex time.Duration
	var result bool
	var err error

	if len(expiration) > 0 {
		ex = time.Duration(expiration[0])
	} else {
		ex = time.Duration(0)
	}

	if c.Mode {
		result, err = c.Clients.SetNX(c.Ctx, key, value, ex).Result()
	} else {
		result, err = c.ClusterClient.SetNX(c.Ctx, key, value, ex).Result()
	}

	if err != nil {
		logs.Error(err.Error())
	}
	return result
}
