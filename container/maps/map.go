package maps

import (
	"sync"
)

//Map parameter structure
type Map struct {
	Map  map[string]interface{} //
	lock *sync.RWMutex          // 加锁
}

//New ...
func New() *Map {
	return &Map{Map: make(map[string]interface{})}
}

//Set ...
func (get *Map) Set(key string, value interface{}) {
	get.lock.Lock()
	defer get.lock.Unlock()
	get.Map[key] = value
}

//Get ...
func (get *Map) Get(key string) interface{} {
	get.lock.Lock()
	defer get.lock.Unlock()
	_, err := get.Map[key]
	if err {
		return nil
	}
	return get.Map[key]
}

//GetOrSet ...
func (get *Map) GetOrSet(key string, value interface{}) interface{} {
	get.lock.Lock()
	defer get.lock.Unlock()
	_, err := get.Map[key]
	if err {
		get.Map[key] = value
		return value
	}
	return get.Map[key]
}

//Count ...
func (get *Map) Count() int {
	get.lock.Lock()
	defer get.lock.Unlock()
	return len(get.Map)
}

//Delete ...
func (get *Map) Delete(key string) bool {
	get.lock.Lock()
	defer get.lock.Unlock()
	delete(get.Map, key)
	_, err := get.Map[key]
	if err {
		return false
	}
	return true
}

//LockFunc locks writing by callback function <f>
func (get *Map) LockFunc(f func(Map map[string]interface{})) *Map {
	get.lock.Lock()
	defer get.lock.Unlock()

	f(get.Map)
	return get
}

//ReadLockFunc locks writing by callback function <f>
func (get *Map) ReadLockFunc(f func(Map map[string]interface{})) *Map {
	get.lock.RLock()
	defer get.lock.RUnlock()

	f(get.Map)
	return get
}
