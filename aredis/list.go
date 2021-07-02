package aredis

import (
	"github.com/small-ek/antgo/os/logs"
)

//GetList value<获取列表长度>
func (c *Client) GetListLength(key string) int64 {
	lens, err := c.Clients.LLen(c.Ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return lens
}

//GetList value<获取列表>
func (c *Client) GetList(key string) []string {
	lens := c.GetListLength(key)
	list, err := c.Clients.LRange(c.Ctx, key, 0, lens-1).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return list
}

//GetListIndex value<返回名称为key的list中index位置的元素>
func (c *Client) GetListIndex(key string, index int64) string {
	list, err := c.Clients.LIndex(c.Ctx, key, index).Result()
	if err != nil {
		logs.Error(err.Error())
	}
	return list
}

//SetList value<修改列表>
func (c *Client) SetList(key string, index int64, value interface{}) error {
	err := c.Clients.LSet(c.Ctx, key, index, value).Err()
	return err
}

//RemoveList value<删除列表>
//count 参数表示删除多少个key中的list
func (c *Client) RemoveList(key string, value interface{}, count ...int64) error {
	var counts int64 = 0
	if len(count) > 0 {
		counts = count[0]
	}
	err := c.Clients.LRem(c.Ctx, key, counts, value).Err()
	return err
}

//RemoveListLeft value<返回并删除名称为key的list中的首元素>
func (c *Client) RemoveListLeft(key string) error {
	err := c.Clients.LPop(c.Ctx, key).Err()
	return err
}

//RemoveListRight value<返回并删除名称为key的list中的尾元素>
func (c *Client) RemoveListRight(key string) error {
	err := c.Clients.LPop(c.Ctx, key).Err()
	return err
}

//PushList value<添加>
func (c *Client) PushList(key string, value interface{}) error {
	err := c.Clients.RPush(c.Ctx, key, value).Err()
	return err
}
