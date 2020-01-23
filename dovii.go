package dovii

import (
	"fmt"

	"github.com/CyrusJavan/dovii/src/keyvaluestore"
)

// KVStore is how consumers of dovii will access the underlying data store
type KVStore struct {
	engine keyvaluestore.KeyValueStore
}

// NewKVStore takes in functional options and returns the KVStore
func NewKVStore(options ...func(*KVStore)) (*KVStore, error) {
	var store KVStore

	for _, option := range options {
		option(&store)
	}

	return &store, nil
}

// BasicMemEngine is the simplest in memory engine
func BasicMemEngine(store *KVStore) {
	store.engine = make(keyvaluestore.BasicMemory)
}

// BasicFileEngine is the simplest (slowest) on disk engine
func BasicFileEngine(store *KVStore) {
	bf, err := keyvaluestore.NewBasicFile()
	if err != nil {
		panic(err)
	}
	store.engine = bf
}

// BitcaskEngine uses a go implementation of bitcask
func BitcaskEngine(store *KVStore) {
	bc, err := keyvaluestore.NewBitcask(true)
	if err != nil {
		panic(err)
	}
	store.engine = bc
}

// Get accesses the underlying db engine to return a value
func (store *KVStore) Get(key string) (string, error) {
	if store.engine == nil {
		return "", fmt.Errorf("no dovii engine has been intitialized")
	}
	return store.engine.Get(key)
}

// Set accesses the underlying db engine to set a value
func (store *KVStore) Set(key, value string) error {
	if store.engine == nil {
		return fmt.Errorf("no dovii engine has been intitialized")
	}
	return store.engine.Set(key, value)
}
