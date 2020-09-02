package container

import "sync"

type Map struct {
	Map  map[string]interface{}
	lock sync.RWMutex // 加锁
}

// Set ...
func (this *Map) Set(key string, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Map[key] = value
}

// Get ...
func (this *Map) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Map[key]
}
