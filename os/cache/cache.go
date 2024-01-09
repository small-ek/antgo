package cache

import (
	"github.com/coocood/freecache"
	"github.com/small-ek/antgo/crypto/hash"
	"github.com/small-ek/antgo/utils/conv"
)

const (
	//Sets the start memory size
	cacheSize = 1024 * 1024
	//the entry will not be written to the cache. expireSeconds <= 0 means no expire,
	cacheExpire = 0
)

var cache = freecache.NewCache(cacheSize)

// Get cached data
func Get(key string) []byte {
	//判断是否有缓存
	h := hash.Sha1(key)
	getData, _ := cache.Get([]byte(h))

	return getData
}

// Set the cache data
func Set(key string, value interface{}, expire ...int) {
	//判断是否有缓存
	h := hash.Sha1(conv.String(key))

	if len(expire) > 0 {
		_ = cache.Set([]byte(h), conv.Bytes(value), expire[0])
	}
	_ = cache.Set([]byte(h), conv.Bytes(value), cacheExpire)
}

// Sets cache data based on value
func Sets(value interface{}, expire ...int) {
	h := hash.Sha1(conv.String(value))

	if len(expire) > 0 {
		_ = cache.Set([]byte(h), conv.Bytes(value), expire[0])
	}
	_ = cache.Set([]byte(h), conv.Bytes(value), cacheExpire)
}

// GetOrSet returns existing value or if record doesn't exist
func GetOrSet(key string, value interface{}, expire ...int) []byte {
	h := hash.Sha1(key)
	if len(expire) > 0 {
		var result, _ = cache.GetOrSet(conv.Bytes(h), conv.Bytes(value), expire[0])
		return result
	}
	var result, _ = cache.GetOrSet(conv.Bytes(h), conv.Bytes(value), cacheExpire)
	return result
}

// Remove Delete the cache
func Remove(key string) bool {
	h := hash.Sha1(key)
	result := cache.Del([]byte(h))
	return result
}

// Clear the cache
func Clear() bool {
	cache.Clear()
	return true
}
