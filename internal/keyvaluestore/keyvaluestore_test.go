package keyvaluestore_test

import (
	"testing"

	"github.com/CyrusJavan/dovii/internal/keyvaluestore"
)

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

func BenchmarkBitcaskGetSet(b *testing.B) {
	var db keyvaluestore.KeyValueStore
	db, _ = keyvaluestore.NewBitcask(true)
	key := "key"
	value := "value"
	var r string
	for n := 0; n < b.N; n++ {
		db.Set(key+string(n), value)
		r, _ = db.Get(key + string(n))
	}
	result = r
}
