package keyvalue

import "github.com/slava-y2k/gofig/maps"

// KeyValue is a simple Provider that supply values from provided map
// deepCopy - makes a deep copy on every read.
type KeyValue struct {
	keyValues map[string]any
	deepCopy  bool
}

// New returns a KeyValue serving the given key-value pairs.
// A nil map will force to create an empty map
func New(keyValues map[string]any) *KeyValue {
	if keyValues == nil {
		keyValues = make(map[string]any)
	}

	return &KeyValue{keyValues: keyValues}
}

// WithDeepCopy enables or disabled deep copy.
func (kv *KeyValue) WithDeepCopy(isEnabled bool) *KeyValue {
	kv.deepCopy = isEnabled
	return kv
}

// Read returns the underlying map. Refer to deepCopy for the details.
func (kv *KeyValue) Read() (map[string]any, error) {
	if !kv.deepCopy {
		return kv.keyValues, nil
	}
	return maps.DeepCopyMap(kv.keyValues), nil
}
