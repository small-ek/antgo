package generics

import (
	"fmt"
	"sync"
)

type Array[T comparable] struct {
	Slice []T
	lock  sync.RWMutex // 加锁
}

func New[T comparable]() *Array[T] {
	return &Array[T]{Slice: make([]T, 0)}
}

func (a *Array[T]) Append(value T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.Slice = append([]T(a.Slice), value)
}

func (a *Array[T]) Len() int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return len([]T(a.Slice))
}

func (a *Array[T]) List() []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.Slice
}

//Insert Array
func (a *Array[T]) Insert(index int, value T) (err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if index < 0 || index >= len([]T(a.Slice)) {
		err = fmt.Errorf("invalid Delete index %d", index)
		return
	}

	var reset = make([]T, 0)
	prefix := append(reset, a.Slice[index:]...)
	a.Slice = append([]T(a.Slice[0:index]), value)
	a.Slice = append([]T(a.Slice), prefix...)

	return
}

func (a *Array[T]) Delete(index int) (t T, err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if index < 0 || index >= len([]T(a.Slice)) {
		err = fmt.Errorf("invalid Delete index %d", index)
		return
	}
	t = a.Slice[index]
	a.Slice = append([]T(a.Slice[:index]), a.Slice[index+1:]...)
	return
}

func (a *Array[T]) Set(index int, value T) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	if index < 0 || index >= len([]T(a.Slice)) {
		return false
	}
	a.Slice[index] = value
	return true
}

func (a *Array[T]) Get(index int) (t T, err error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if index < 0 || index >= len([]T(a.Slice)) {
		err = fmt.Errorf("invalid Get index %d", index)
		return
	}
	return a.Slice[index], nil
}

func (a *Array[T]) Search(value T) int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for i, v := range a.Slice {
		if v == value {
			return i
		}
	}
	return -1
}

func (a *Array[T]) Clear() {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.Slice = make([]T, 0)
}

func (a *Array[T]) LockFunc(f func(array []T)) *Array[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	f(a.Slice)
	return a
}

//extension
func (a *Array[T]) Map(f func(T) any) *Array[any] {
	a.lock.Lock()
	defer a.lock.Unlock()

	arrayU := make([]any, len([]T(a.Slice)))
	if len([]T(a.Slice)) > 0 {
		for _, v := range a.Slice {
			arrayU = append(arrayU, f(v))
		}
	}

	return &Array[any]{
		arrayU,
		sync.RWMutex{},
	}
}
