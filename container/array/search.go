package array

//SearchStrings 搜索切片
func SearchStrings(list []string, str string) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchInts 搜索切片
func SearchInts(list []int, str int) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchFloat64s 搜索切片
func SearchFloat64s(list []float64, str float64) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchFloat32s 搜索切片
func SearchFloat32s(list []float32, str float32) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}

//SearchInterfaces 搜索切片
func SearchInterfaces(list []interface{}, str interface{}) int {
	for i := 0; i < len(list); i++ {
		var value = list[i]
		if value == str {
			return i
		}
	}
	return -1
}
