// Package keyvalue provides an in-memory configuration provider.
package keyvalue

import "github.com/slava-y2k/gofig"

// KeyValue is an in-memory Provider over a caller-supplied key-value map;
// copying on Read is opt-in through WithDeepCopy.
type KeyValue struct {
	keyValues map[string]any
	deepCopy  bool
}

func New(keyValues map[string]any) *KeyValue {
	if keyValues == nil {
		keyValues = make(map[string]any)
	}

	return &KeyValue{keyValues: keyValues}
}

// WithDeepCopy switches Read between the served map itself and a fresh copy.
// Copying is off by default: no allocation, but the served map is exposed.
func (kv *KeyValue) WithDeepCopy(isEnabled bool) *KeyValue {
	kv.deepCopy = isEnabled
	return kv
}

func (kv *KeyValue) Read() (map[string]any, error) {
	if !kv.deepCopy {
		return kv.keyValues, nil
	}
	return gofig.DeepCopyMap(kv.keyValues), nil
}
