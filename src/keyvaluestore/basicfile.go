package keyvaluestore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// BasicFile is a naive implementation of file based storage
type BasicFile struct {
	r *bufio.Reader
	w *bufio.Writer
}

// Get does
func (bf BasicFile) Get(key string) (string, error) {
	contents, err := ioutil.ReadFile("/tmp/doviibf")
	if err != nil {
		return "", err
	}
	var m map[string]string
	err = json.Unmarshal(contents, &m)
	if err != nil {
		return "", err
	}
	if val, ok := m[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Key not found")
}

// Set does
func (bf BasicFile) Set(key, value string) error {
	contents, err := ioutil.ReadFile("/tmp/doviibf")
	if err != nil {
		return err
	}
	var m map[string]string
	err = json.Unmarshal(contents, &m)
	if err != nil {
		return err
	}
	m[key] = value
	contents, err = json.Marshal(m)
	if err != nil {
		return err
	}
	ioutil.WriteFile("/tmp/doviibf", contents, 0666)
	return nil
}

// NewBasicFile creates the db file if it doesn't exist or opens the existing
// db file
func NewBasicFile() (BasicFile, error) {
	return BasicFile{}, nil
}
