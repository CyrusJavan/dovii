package keyvaluestore_test

import (
	"testing"

	"github.com/CyrusJavan/dovii/src/keyvaluestore"
	"github.com/stretchr/testify/assert"
)

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
