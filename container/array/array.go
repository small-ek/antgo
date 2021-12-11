package array

import "sync"

//Array parameter structure
type Array struct {
	Slice []interface{}
	lock  sync.RWMutex // 加锁
}

//New Array
func New() *Array {
	return &Array{Slice: make([]interface{}, 0)}
}

//Append Set Array
func (a *Array) Append(value interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Slice = append(a.Slice, value)
}

//Len Count Array
func (a *Array) Len() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return len(a.Slice)
}

//List Array
func (a *Array) List() []interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.Slice
}

//InsertAfter Array
func (a *Array) InsertAfter(index int, value interface{}) []interface{} {
	a.lock.Lock()
	defer a.lock.Unlock()

	var reset = make([]interface{}, 0)
	prefix := append(reset, a.Slice[index:]...)
	a.Slice = append(a.Slice[0:index], value)
	a.Slice = append(a.Slice, prefix...)
	return a.Slice
}

//Delete Array
func (a *Array) Delete(index int) []interface{} {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.Slice = append(a.Slice[:index], a.Slice[index+1:]...)
	return a.Slice
}

//Set Array
func (a *Array) Set(index int, value interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.Slice[index] = value
}

//Get Array
func (a *Array) Get(index int) interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.Slice[index]
}

//Search Array
func (a *Array) Search(value interface{}) int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for i := 0; i < len(a.Slice); i++ {
		if a.Slice[i] == value {
			return i
		}
	}
	return 0
}

//Clear Array
func (a *Array) Clear() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Slice = make([]interface{}, 0)
}

//LockFunc locks writing by callback function <f>
func (a *Array) LockFunc(f func(array []interface{})) *Array {
	a.lock.Lock()
	defer a.lock.Unlock()

	f(a.Slice)
	return a
}

//ReadLockFunc locks writing by callback function <f>
func (a *Array) ReadLockFunc(f func(array []interface{})) *Array {
	a.lock.RLock()
	defer a.lock.RUnlock()

	f(a.Slice)
	return a
}
