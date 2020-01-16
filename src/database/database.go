package database

import "fmt"

// Database interface
type Database interface {
	Get()
	Set()
}

// MapDatabase is the simplest possible in memory database
// utilizing just a regular map to store keys and values.
type MapDatabase map[string]string

// Get the value of the key
func (s MapDatabase) Get(key string) (string, error) {
	if val, ok := s[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Key not found")
}

// Set the key to the value
func (s MapDatabase) Set(key, value string) error {
	s[key] = value
	return nil
}
