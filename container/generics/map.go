package generics


import (
	"sync"
)

type MapT[K comparable, V any] struct {
	Map  map[K]V      //
	lock sync.RWMutex // 加锁
}

func NewT[K comparable, V any]() *MapT[K, V] {
	return &MapT[K, V]{Map: make(map[K]V)}
}

func (m *MapT[K, V]) Set(key K, value V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.Map[key] = value
}

func (m *MapT[K, V]) Get(key K) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v, ok := m.Map[key]
	return v, ok
}

func (m *MapT[K, V]) GetOrSet(key K, value V) V {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, ok := m.Map[key]
	if !ok {
		m.Map[key] = value
		return value
	}
	return v
}

func (m *MapT[K, V]) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(map[K]V(m.Map))
}

func (m *MapT[K, V]) Delete(key K) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.Map, key)
	_, ok := m.Map[key]
	return !ok
}

func (m *MapT[K, V]) LockFunc(f func(Map map[K]V)) *MapT[K, V] {
	m.lock.Lock()
	defer m.lock.Unlock()

	f(m.Map)
	return m
}

