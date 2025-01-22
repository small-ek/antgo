package array

import (
	"errors"
	"sync"
)

// ConcurrentArray is a generic thread-safe array.
// ConcurrentArray 是一个泛型线程安全数组。
type ConcurrentArray[T comparable] struct {
	data []T          // Underlying typed data storage. 底层类型化数据存储。
	lock sync.RWMutex // Read-write lock for concurrency control. 读写锁用于并发控制。
}

// New creates a new ConcurrentArray with initial capacity.
// New 创建一个具有初始容量的 ConcurrentArray。
// Parameters:
//
//	capacity - Initial capacity to reduce reallocations. 初始容量以减少内存重分配。
func New[T comparable](capacity int) *ConcurrentArray[T] {
	return &ConcurrentArray[T]{
		data: make([]T, 0, capacity),
	}
}

// Append adds an element to the end of the array.
// Append 向数组末尾添加元素。
func (ca *ConcurrentArray[T]) Append(element T) {
	ca.lock.Lock()
	defer ca.lock.Unlock()
	ca.data = append(ca.data, element)
}

// Len returns the number of elements in the array.
// Len 返回数组中元素的数量。
func (ca *ConcurrentArray[T]) Len() int {
	ca.lock.RLock()
	defer ca.lock.RUnlock()
	return len(ca.data)
}

// List returns a copy of the underlying data to prevent external modification.
// List 返回底层数据的副本以防止外部修改。
func (ca *ConcurrentArray[T]) List() []T {
	ca.lock.RLock()
	defer ca.lock.RUnlock()
	copied := make([]T, len(ca.data))
	copy(copied, ca.data)
	return copied
}

// Insert inserts an element at the specified index.
// Insert 在指定索引处插入元素。
// Returns error if index is out of bounds.
// 如果索引越界则返回错误。
func (ca *ConcurrentArray[T]) Insert(index int, element T) error {
	ca.lock.Lock()
	defer ca.lock.Unlock()

	if index < 0 || index > len(ca.data) { // Allow inserting at len(data) (append)
		return errors.New("index out of bounds")
	}

	// Optimized insertion with single allocation
	newData := make([]T, 0, len(ca.data)+1)
	newData = append(newData, ca.data[:index]...)
	newData = append(newData, element)
	newData = append(newData, ca.data[index:]...)
	ca.data = newData

	return nil
}

// Delete removes the element at the specified index.
// Delete 删除指定索引处的元素。
// Returns error if index is invalid.
// 如果索引无效则返回错误。
func (ca *ConcurrentArray[T]) Delete(index int) error {
	ca.lock.Lock()
	defer ca.lock.Unlock()

	if index < 0 || index >= len(ca.data) {
		return errors.New("index out of bounds")
	}

	ca.data = append(ca.data[:index], ca.data[index+1:]...)
	return nil
}

// Set updates the element at the specified index.
// Set 更新指定索引处的元素。
// Returns error if index is invalid.
// 如果索引无效则返回错误。
func (ca *ConcurrentArray[T]) Set(index int, element T) error {
	ca.lock.Lock()
	defer ca.lock.Unlock()

	if index < 0 || index >= len(ca.data) {
		return errors.New("index out of bounds")
	}

	ca.data[index] = element
	return nil
}

// Get returns the element at the specified index.
// Get 返回指定索引处的元素。
// Returns zero value and error if index is invalid.
// 如果索引无效则返回零值和错误。
func (ca *ConcurrentArray[T]) Get(index int) (T, error) {
	ca.lock.RLock()
	defer ca.lock.RUnlock()

	var zero T
	if index < 0 || index >= len(ca.data) {
		return zero, errors.New("index out of bounds")
	}
	return ca.data[index], nil
}

// Search returns the first index of the target element.
// Search 返回目标元素的第一个索引。
// Returns -1 if not found.
// 如果未找到则返回 -1。
func (ca *ConcurrentArray[T]) Search(target T) int {
	ca.lock.RLock()
	defer ca.lock.RUnlock()

	for i, v := range ca.data {
		if v == target {
			return i
		}
	}
	return -1
}

// Clear removes all elements from the array.
// Clear 清空数组中的所有元素。
func (ca *ConcurrentArray[T]) Clear() {
	ca.lock.Lock()
	defer ca.lock.Unlock()
	ca.data = nil
}

// WithWriteLock executes a function with write lock held.
// WithWriteLock 在持有写锁的情况下执行函数。
func (ca *ConcurrentArray[T]) WithWriteLock(fn func([]T)) {
	ca.lock.Lock()
	defer ca.lock.Unlock()
	fn(ca.data)
}

// WithReadLock executes a function with read lock held.
// WithReadLock 在持有读锁的情况下执行函数。
func (ca *ConcurrentArray[T]) WithReadLock(fn func([]T)) {
	ca.lock.RLock()
	defer ca.lock.RUnlock()
	fn(ca.data)
}
