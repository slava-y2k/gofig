package keyvalue

import "github.com/slava-y2k/gofig"

// KeyValue is a simple Provider that supplies values from the provided map.
// When deepCopy is enabled, Read makes a deep copy on every read.
type KeyValue struct {
	keyValues map[string]any
	deepCopy  bool
}

// New returns a KeyValue serving the given key-value pairs.
func New(keyValues map[string]any) *KeyValue {
	if keyValues == nil {
		keyValues = make(map[string]any)
	}

	return &KeyValue{keyValues: keyValues}
}

// WithDeepCopy enables or disables deep copy.
func (kv *KeyValue) WithDeepCopy(isEnabled bool) *KeyValue {
	kv.deepCopy = isEnabled
	return kv
}

// Read returns the underlying map, deep-copied when configured.
func (kv *KeyValue) Read() (map[string]any, error) {
	if !kv.deepCopy {
		return kv.keyValues, nil
	}
	return gofig.DeepCopyMap(kv.keyValues), nil
}
