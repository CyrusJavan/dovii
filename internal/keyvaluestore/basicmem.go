package keyvaluestore

import "fmt"

// BasicMemory is the simplest possible in memory database utilizing just a
// regular map to store keys and values.
type BasicMemory map[string]string

// Get the value of the key
func (s BasicMemory) Get(key string) (string, error) {
	if val, ok := s[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Key not found")
}

// Set the key to the value
func (s BasicMemory) Set(key, value string) error {
	s[key] = value
	return nil
}
