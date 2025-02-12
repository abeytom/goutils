package gomap

type ArrayMultiMap[K comparable, V any] struct {
	data map[K][]V
}

// NewArrayMultiMap creates a new ArrayMultiMap instance
func NewArrayMultiMap[K comparable, V any]() *ArrayMultiMap[K, V] {
	return &ArrayMultiMap[K, V]{data: make(map[K][]V)}
}

// Put adds a value to the multimap under the specified key
func (m *ArrayMultiMap[K, V]) Put(key K, value V) {
	m.data[key] = append(m.data[key], value)
}

// Get retrieves the values associated with a key
func (m *ArrayMultiMap[K, V]) Get(key K) []V {
	return m.data[key]
}

// Remove deletes a specific value from a key
func (m *ArrayMultiMap[K, V]) Remove(key K, eq func(v V) bool) {
	values := m.data[key]
	for i, v := range values {
		if eq(v) {
			m.data[key] = append(values[:i], values[i+1:]...)
			return
		}
	}
}

// RemoveAll removes all values associated with a key
func (m *ArrayMultiMap[K, V]) RemoveAll(key K) {
	delete(m.data, key)
}

// Keys returns all keys in the multimap
func (m *ArrayMultiMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}
