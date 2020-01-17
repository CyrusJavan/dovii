package keyvaluestore

// KeyValueStore interface
type KeyValueStore interface {
	Get(string) (string, error)
	Set(string, string) error
}
