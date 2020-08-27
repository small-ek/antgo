package cache

import (
	"github.com/small-ek/ginp/conv"
	"github.com/coocood/freecache"
	"github.com/small-ek/ginp/crypto/hashs"
)

type New struct {
	Key    string
	Group  string
	Expire int
}

const (
	//Sets the start memory size
	cacheSize = 1024 * 1024
)

var cache = freecache.NewCache(cacheSize)

//Get cached data
func (this *New) Get() []byte {
	//判断是否有缓存
	var hash = hashs.Sha256(this.Group + this.Key)
	getData, _ := cache.Get([]byte(hash))
	return getData
}

//Set the cache data
func (this *New) Set(data interface{}) {
	//判断是否有缓存
	if this.Expire == 0 {
		this.Expire = 10000
	}
	var hash = hashs.Sha256(this.Group + this.Key)
	go cache.Set([]byte(hash), conv.StructToBytes(data), this.Expire)
}

//Delete the cache
func Del(key []byte) bool {
	result := cache.Del(key)
	return result
}

//Clear the cache
func Clear() bool {
	cache.Clear()
	return true
}
