package gomap

// StrArrayToMap convert a slice to a map
func StrArrayToMap[T any](array []T, mapFn func(T) string) map[string]T {
	m := make(map[string]T)
	if len(array) == 0 {
		return m
	}
	for _, item := range array {
		key := mapFn(item)
		if len(key) > 0 {
			m[key] = item
		}
	}
	return m
}

// StrArrayToMap convert a slice to a Multimap
func StrArrayToMultiMap[T any](array []T, mapFn func(T) string) *ArrayMultiMap[string, T] {
	m := NewArrayMultiMap[string, T]()
	if len(array) == 0 {
		return m
	}
	for _, item := range array {
		key := mapFn(item)
		if len(key) > 0 {
			m.Put(key, item)
		}
	}
	return m
}
