package array

//SearchString 搜索切片
func SearchString(list []string, str string) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchInt 搜索切片
func SearchInt(list []int, str int) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchInt32 搜索切片
func SearchInt32(list []int32, str int32) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchInt64 搜索切片
func SearchInt64(list []int64, str int64) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchFloat32 搜索切片
func SearchFloat32(list []float32, str float32) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchFloat64 搜索切片
func SearchFloat64(list []float64, str float64) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchInterface 搜索切片
func SearchInterface(list []interface{}, str interface{}) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}
