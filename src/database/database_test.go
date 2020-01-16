package database_test

import (
	"testing"

	"github.com/CyrusJavan/dovii/src/database"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	db := make(database.MapDatabase)

	key := "test"
	value := "value"

	err := db.Set(key, value)
	assert.NoError(t, err)

	got, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestGetKeyNotFound(t *testing.T) {
	db := make(database.MapDatabase)

	key := "test"
	value := "value"

	err := db.Set(key, value)
	assert.NoError(t, err)

	_, err = db.Get(key + " ")
	assert.Error(t, err)
}

func TestSet(t *testing.T) {
	db := make(database.MapDatabase)

	key := "test"
	value := "value"

	err := db.Set(key, value)
	assert.NoError(t, err)
}
