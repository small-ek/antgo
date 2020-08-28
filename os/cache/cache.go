package cache

import (
	"github.com/coocood/freecache"
	"github.com/small-ek/ginp/conv"
	"github.com/small-ek/ginp/crypto/sha256"
)

const (
	//Sets the start memory size
	cacheSize   = 1024 * 1024
	cacheExpire = 1000
)

var cache = freecache.NewCache(cacheSize)

//Get cached data
func Get(key string) []byte {
	//判断是否有缓存
	var hash = sha256.Create(key)
	getData, _ := cache.Get([]byte(hash))

	return getData
}

//Set the cache data
func Set(key string, value interface{}, expire ...int) {
	//判断是否有缓存
	var hash = sha256.Create(conv.String(key))

	if len(expire) > 0 {
		cache.Set([]byte(hash), conv.StructToBytes(value), expire[0])
	}
	cache.Set([]byte(hash), conv.StructToBytes(value), cacheExpire)
}

func GetOrSet(key string, value interface{}, expire ...int) []byte {
	var hash = sha256.Create(key)
	if len(expire) > 0 {
		var result, _ = cache.GetOrSet(conv.Bytes(hash), conv.Bytes(value), expire[0])
		return result
	}
	var result, _ = cache.GetOrSet(conv.Bytes(hash), conv.Bytes(value), cacheExpire)
	return result
}

//Delete the cache
func Remove(key string) bool {
	var hash = sha256.Create(key)
	result := cache.Del([]byte(hash))
	return result
}

//Clear the cache
func Clear() bool {
	cache.Clear()
	return true
}
