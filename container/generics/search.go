package generics

func Search[T comparable](list []T, key T) int {
	for i, v := range list {
		if v == key {
			return i
		}
	}
	return -1
}
