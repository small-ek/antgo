package search

// SearchString 搜索切片
func SearchString(list []string, key string) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}

// SearchInt 搜索切片
func SearchInt(list []int, key int) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}

// SearchInt32 搜索切片
func SearchInt32(list []int32, key int32) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}

// SearchInt64 搜索切片
func SearchInt64(list []int64, key int64) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}

// SearchFloat32 搜索切片
func SearchFloat32(list []float32, key float32) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}

// SearchFloat64 搜索切片
func SearchFloat64(list []float64, key float64) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}

// SearchInterface 搜索切片
func SearchInterface(list []interface{}, key interface{}) int {
	for i, _ := range list {
		if list[i] == key {
			return i
		}
	}
	return -1
}
