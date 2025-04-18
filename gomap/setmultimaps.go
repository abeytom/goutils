package gomap

// SetMultimap holds unique values for each key
type SetMultimap struct {
	data map[string]map[string]struct{}
}

// NewSetMultimap creates and returns a new SetMultimap
func NewSetMultimap() *SetMultimap {
	return &SetMultimap{
		data: make(map[string]map[string]struct{}),
	}
}

// Add inserts a value into the set multimap under a given key
func (s *SetMultimap) Add(key, value string) {
	if _, exists := s.data[key]; !exists {
		s.data[key] = make(map[string]struct{})
	}
	s.data[key][value] = struct{}{}
}

// Remove deletes a value from the set multimap for a given key
func (s *SetMultimap) Remove(key, value string) {
	if _, exists := s.data[key]; exists {
		delete(s.data[key], value)
		if len(s.data[key]) == 0 {
			delete(s.data, key) // Remove key if no values left
		}
	}
}

// GetValues retrieves all values for a given key
func (s *SetMultimap) GetValues(key string) []string {
	values := []string{}
	if vals, exists := s.data[key]; exists {
		for val := range vals {
			values = append(values, val)
		}
	}
	return values
}

// Has checks if a value exists under a key
func (s *SetMultimap) Has(key, value string) bool {
	_, exists := s.data[key][value]
	return exists
}

// All returns the entire multimap
func (s *SetMultimap) All() map[string][]string {
	result := make(map[string][]string)
	for key, vals := range s.data {
		for val := range vals {
			result[key] = append(result[key], val)
		}
	}
	return result
}
