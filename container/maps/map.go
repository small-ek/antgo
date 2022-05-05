package maps

import (
	"sync"
)

//Map parameter structure
type Map struct {
	Map  map[interface{}]interface{} //
	lock sync.RWMutex                // 加锁
}

//NewMap
func NewMap() *Map {
	return &Map{Map: make(map[interface{}]interface{})}
}

//Set ...
func (m *Map) Set(key interface{}, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Map[key] = value
}

//Get ...
func (m *Map) Get(key interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	v, ok := m.Map[key]
	if !ok {
		return nil
	}
	return v
}

//GetOrSet ...
func (m *Map) GetOrSet(key interface{}, value interface{}) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.Map[key]
	if !ok {
		m.Map[key] = value
		return value
	}
	return v
}

//Count ...
func (m *Map) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.Map)
}

//Delete ...
func (m *Map) Delete(key interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.Map, key)
}

//LockFunc locks writing by callback function <f>
func (m *Map) LockFunc(f func(Map map[interface{}]interface{})) *Map {
	m.lock.Lock()
	defer m.lock.Unlock()

	f(m.Map)
	return m
}

//ReadLockFunc locks writing by callback function <f>
func (m *Map) ReadLockFunc(f func(Map map[interface{}]interface{})) *Map {
	m.lock.RLock()
	defer m.lock.RUnlock()

	f(m.Map)
	return m
}
