package array

import "sync"

//Array parameter structure
type Array struct {
	Slice []interface{}
	lock  *sync.RWMutex // 加锁
}

// New Array
func New() *Array {
	return &Array{Slice: make([]interface{}, 0)}
}

// Append Set Array
func (this *Array) Append(value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Slice = append(this.Slice, value)
}

// Len Count Array
func (this *Array) Len() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return len(this.Slice)
}

// List Array
func (this *Array) List() []interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Slice
}

// InsertAfter Array
func (this *Array) InsertAfter(index int, value interface{}) []interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	var reset = make([]interface{}, 0)
	prefix := append(reset, this.Slice[index:]...)
	this.Slice = append(this.Slice[0:index], value)
	this.Slice = append(this.Slice, prefix...)
	return this.Slice
}

// Delete Array
func (this *Array) Delete(index int) []interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.Slice = append(this.Slice[:index], this.Slice[index+1:]...)
	return this.Slice
}

// Set Array
func (this *Array) Set(index int, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.Slice[index] = value
}

// Get Array
func (this *Array) Get(index int) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Slice[index]
}

// Search Array
func (this *Array) Search(value interface{}) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	for i := 0; i < len(this.Slice); i++ {
		if this.Slice[i] == value {
			return i
		}
	}
	return 0
}

// Clear Array
func (this *Array) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Slice = make([]interface{}, 0)
}

// LockFunc locks writing by callback function <f>
func (this *Array) LockFunc(f func(array []interface{})) *Array {
	this.lock.Lock()
	defer this.lock.Unlock()

	f(this.Slice)
	return this
}

// ReadLockFunc locks writing by callback function <f>
func (this *Array) ReadLockFunc(f func(array []interface{})) *Array {
	this.lock.RLock()
	defer this.lock.RUnlock()

	f(this.Slice)
	return this
}
