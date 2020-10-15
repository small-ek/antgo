package container

import (
	"sync"
)

type Map struct {
	Map  map[string]interface{}
	lock sync.RWMutex // 加锁
}

// New ...
func NewMap() *Map {
	return &Map{Map: make(map[string]interface{})}
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
	_, err := this.Map[key]
	if err {
		return nil
	} else {
		return this.Map[key]
	}
}

// GetOrSet ...
func (this *Map) GetOrSet(key string, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, err := this.Map[key]
	if err {
		this.Map[key] = value
		return value
	} else {
		return this.Map[key]
	}
}

// Count ...
func (this *Map) Count() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return len(this.Map)
}

// Delete ...
func (this *Map) Delete(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.Map, key)
	_, err := this.Map[key]
	if err {
		return false
	} else {
		return true
	}
}
