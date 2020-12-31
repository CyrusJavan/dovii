package keyvaluestore_test

import (
	"testing"

	"github.com/CyrusJavan/dovii/internal/keyvaluestore"
	"github.com/stretchr/testify/assert"
)

func TestBitcaskGet(t *testing.T) {
	var db keyvaluestore.KeyValueStore
	var err error
	db, err = keyvaluestore.NewBitcask(true)
	assert.NoError(t, err)

	key := "test"
	value := "value"

	err = db.Set(key, value)
	assert.NoError(t, err)

	got, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestBitcaskGetKeyNotFound(t *testing.T) {
	var db keyvaluestore.KeyValueStore
	var err error
	db, err = keyvaluestore.NewBitcask(true)
	assert.NoError(t, err)

	key := "test"
	value := "value"

	err = db.Set(key, value)
	assert.NoError(t, err)

	_, err = db.Get(key + " ")
	assert.Error(t, err)
}

func TestBitcaskSet(t *testing.T) {
	var db keyvaluestore.KeyValueStore
	var err error
	db, err = keyvaluestore.NewBitcask(true)
	assert.NoError(t, err)

	key := "test"
	value := "value"

	err = db.Set(key, value)
	assert.NoError(t, err)
}
