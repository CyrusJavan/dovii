package keyvaluestore_test

import (
	"testing"

	"github.com/CyrusJavan/dovii/src/keyvaluestore"
	"github.com/stretchr/testify/assert"
)

func TestBasicMemoryGet(t *testing.T) {
	var db keyvaluestore.KeyValueStore = make(keyvaluestore.BasicMemory)

	key := "test"
	value := "value"

	err := db.Set(key, value)
	assert.NoError(t, err)

	got, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestBasicMemoryGetKeyNotFound(t *testing.T) {
	var db keyvaluestore.KeyValueStore = make(keyvaluestore.BasicMemory)

	key := "test"
	value := "value"

	err := db.Set(key, value)
	assert.NoError(t, err)

	_, err = db.Get(key + " ")
	assert.Error(t, err)
}

func TestBasicMemorySet(t *testing.T) {
	var db keyvaluestore.KeyValueStore = make(keyvaluestore.BasicMemory)

	key := "test"
	value := "value"

	err := db.Set(key, value)
	assert.NoError(t, err)
}

func TestBasicFileGet(t *testing.T) {
	var db keyvaluestore.KeyValueStore
	var err error
	db, err = keyvaluestore.NewBasicFile()
	assert.NoError(t, err)

	key := "test"
	value := "value"

	err = db.Set(key, value)
	assert.NoError(t, err)

	got, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestBasicFileGetKeyNotFound(t *testing.T) {
	var db keyvaluestore.KeyValueStore
	var err error
	db, err = keyvaluestore.NewBasicFile()
	assert.NoError(t, err)

	key := "test"
	value := "value"

	err = db.Set(key, value)
	assert.NoError(t, err)

	_, err = db.Get(key + " ")
	assert.Error(t, err)
}

func TestBasicFileSet(t *testing.T) {
	var db keyvaluestore.KeyValueStore
	var err error
	db, err = keyvaluestore.NewBasicFile()
	assert.NoError(t, err)

	key := "test"
	value := "value"

	err = db.Set(key, value)
	assert.NoError(t, err)
}

var result string

func BenchmarkBasicMemoryGetSet(b *testing.B) {
	var db keyvaluestore.KeyValueStore = make(keyvaluestore.BasicMemory)
	key := "key"
	value := "value"
	var r string
	for n := 0; n < b.N; n++ {
		db.Set(key+string(n), value)
		r, _ = db.Get(key + string(n))
	}
	result = r
}

func BenchmarkBasicFileGetSet(b *testing.B) {
	var db keyvaluestore.KeyValueStore
	db, _ = keyvaluestore.NewBasicFile()
	key := "key"
	value := "value"
	var r string
	for n := 0; n < b.N; n++ {
		db.Set(key+string(n), value)
		r, _ = db.Get(key + string(n))
	}
	result = r
}
