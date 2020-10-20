package array

import "sync"

//Array parameter structure
type Array struct {
	Slice []interface{}
	lock  *sync.RWMutex // 加锁
}

//New Array
func New() *Array {
	return &Array{Slice: make([]interface{}, 0)}
}

//Append Set Array
func (get *Array) Append(value interface{}) {
	get.lock.Lock()
	defer get.lock.Unlock()
	get.Slice = append(get.Slice, value)
}

//Len Count Array
func (get *Array) Len() int {
	get.lock.Lock()
	defer get.lock.Unlock()
	return len(get.Slice)
}

//List Array
func (get *Array) List() []interface{} {
	get.lock.Lock()
	defer get.lock.Unlock()
	return get.Slice
}

//InsertAfter Array
func (get *Array) InsertAfter(index int, value interface{}) []interface{} {
	get.lock.Lock()
	defer get.lock.Unlock()

	var reset = make([]interface{}, 0)
	prefix := append(reset, get.Slice[index:]...)
	get.Slice = append(get.Slice[0:index], value)
	get.Slice = append(get.Slice, prefix...)
	return get.Slice
}

//Delete Array
func (get *Array) Delete(index int) []interface{} {
	get.lock.Lock()
	defer get.lock.Unlock()

	get.Slice = append(get.Slice[:index], get.Slice[index+1:]...)
	return get.Slice
}

//Set Array
func (get *Array) Set(index int, value interface{}) {
	get.lock.Lock()
	defer get.lock.Unlock()

	get.Slice[index] = value
}

//Get Array
func (get *Array) Get(index int) interface{} {
	get.lock.Lock()
	defer get.lock.Unlock()
	return get.Slice[index]
}

//Search Array
func (get *Array) Search(value interface{}) int {
	get.lock.Lock()
	defer get.lock.Unlock()
	for i := 0; i < len(get.Slice); i++ {
		if get.Slice[i] == value {
			return i
		}
	}
	return 0
}

//Clear Array
func (get *Array) Clear() {
	get.lock.Lock()
	defer get.lock.Unlock()
	get.Slice = make([]interface{}, 0)
}

//LockFunc locks writing by callback function <f>
func (get *Array) LockFunc(f func(array []interface{})) *Array {
	get.lock.Lock()
	defer get.lock.Unlock()

	f(get.Slice)
	return get
}

//ReadLockFunc locks writing by callback function <f>
func (get *Array) ReadLockFunc(f func(array []interface{})) *Array {
	get.lock.RLock()
	defer get.lock.RUnlock()

	f(get.Slice)
	return get
}
