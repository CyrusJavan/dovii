package keyvaluestore

// KeyValueStore defines what our service will provide. Right now it is the
// bare minimum of a key-value store.
type KeyValueStore interface {
	Get(string) (string, error)
	Set(string, string) error
}
