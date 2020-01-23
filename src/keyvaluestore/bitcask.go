package keyvaluestore

import (
	"math/rand"
	"time"

	"github.com/prologic/bitcask"
)

// Bitcask is
type Bitcask struct {
	bc *bitcask.Bitcask
}

// Get does
func (b Bitcask) Get(key string) (string, error) {
	v, err := b.bc.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

// Set does
func (b Bitcask) Set(key, value string) error {
	return b.bc.Put([]byte(key), []byte(value))
}

// NewBitcask creates
func NewBitcask(uniqueDB bool) (Bitcask, error) {
	fn := "/tmp/dovii_bitcask"
	if uniqueDB {
		fn += randString(8)
	}
	db, _ := bitcask.Open(fn)
	return Bitcask{db}, nil
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randString(length int) string {
	return stringWithCharset(length, charset)
}
