package gofig

import (
	"errors"
	"sync"

	"github.com/slava-y2k/gofig/maps"
)

// ErrNilProvider is returned by Load when the given provider is nil.
var ErrNilProvider = errors.New("gofig: nil provider")

// Gofig is a key-value pair manager. It loads key-value pairs from
// Providers, merges them into a single map, and provides methods to
// read, add, and remove keys.
//
// Get, Set, and Delete address nested maps through path parts: Get
// and Delete take the parts as a variadic list, Set as a slice.
//
// Methods that cannot fail (Merge, MergeMap, Set, Delete) return the
// Gofig itself, enabling a fluent interface.
type Gofig struct {
	rwMutex sync.RWMutex
	kvMap   map[string]any
}

// New returns an empty Gofig using "." as the key path delimiter.
func New() *Gofig {
	return &Gofig{kvMap: make(map[string]any)}
}

// Load reads the provider and merges its output into the
// configuration.
func (gfg *Gofig) Load(provider Provider) error {
	if provider == nil {
		return ErrNilProvider
	}

	var keyValues map[string]any
	var err error
	keyValues, err = provider.Read()
	if err != nil {
		return err
	}

	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.init()
	mergeMaps(gfg.kvMap, maps.DeepCopyMap(keyValues))

	return nil
}

// Get returns the value at the given path parts. The second result
// reports whether the key was found. Without parts, Get returns the
// whole configuration.
func (gfg *Gofig) Get(parts ...string) (any, bool) {
	gfg.rwMutex.RLock()
	defer gfg.rwMutex.RUnlock()

	if gfg.kvMap == nil {
		return nil, false
	}
	if len(parts) == 0 {
		return gfg.kvMap, true
	}

	return getPath(gfg.kvMap, parts)
}

// Set stores value at the given path parts, creating missing
// intermediate maps. Without parts, Set replaces the whole
// configuration when value is a map[string]any.
func (gfg *Gofig) Set(parts []string, value any) *Gofig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.init()

	if len(parts) == 0 {
		if newConf, ok := value.(map[string]any); ok {
			gfg.kvMap = maps.DeepCopyMap(newConf)
		}
		return gfg
	}

	setPath(gfg.kvMap, parts, value)

	return gfg
}

// Merge deep-merges the configuration of another Gofig into this one:
// nested maps are merged recursively, every other value overwrites the
// existing one. A nil Gofig is ignored.
func (gfg *Gofig) Merge(other *Gofig) *Gofig {
	if other == nil {
		return gfg
	}

	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.init()
	mergeMaps(gfg.kvMap, other.snapshot())

	return gfg
}

// MergeMap deep-merges the given key-value pairs into the
// configuration: nested maps are merged recursively, every other value
// overwrites the existing one.
func (gfg *Gofig) MergeMap(keyValues map[string]any) *Gofig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.init()
	mergeMaps(gfg.kvMap, maps.DeepCopyMap(keyValues))

	return gfg
}

// Delete removes the value at the given path parts. Without parts,
// Delete clears the whole configuration. Missing paths are ignored.
func (gfg *Gofig) Delete(parts ...string) *Gofig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()

	if len(parts) == 0 {
		gfg.kvMap = make(map[string]any)
		return gfg
	}

	if gfg.kvMap != nil {
		deletePath(gfg.kvMap, parts)
	}

	return gfg
}

// Raw returns a deep copy of the whole configuration map, so callers
// cannot mutate the internal state. Nested maps are copied; all other
// values are shared.
func (gfg *Gofig) Raw() map[string]any {
	var raw map[string]any = gfg.snapshot()
	if raw == nil {
		raw = make(map[string]any)
	}

	return raw
}

// snapshot returns a deep copy of the configuration map, reading it
// under the read lock.
func (gfg *Gofig) snapshot() map[string]any {
	gfg.rwMutex.RLock()
	defer gfg.rwMutex.RUnlock()

	return maps.DeepCopyMap(gfg.kvMap)
}

// init lazily initializes the configuration map so that the zero
// value of Gofig is usable.
func (gfg *Gofig) init() {
	if gfg.kvMap == nil {
		gfg.kvMap = make(map[string]any)
	}
}

// setPath walks parts and stores value at the leaf, creating missing
// intermediate maps. A non-map value on the path is overwritten with
// a new map.
func setPath(kvMap map[string]any, parts []string, value any) {
	for _, part := range parts[:len(parts)-1] {
		next, ok := kvMap[part].(map[string]any)
		if !ok {
			next = make(map[string]any)
			kvMap[part] = next
		}
		kvMap = next
	}

	last := parts[len(parts)-1]
	if valueMap, ok := value.(map[string]any); ok {
		kvMap[last] = maps.DeepCopyMap(valueMap)
		return
	}
	kvMap[last] = value
}

// deletePath walks parts and deletes the leaf key. Missing paths are
// ignored.
func deletePath(kvMap map[string]any, parts []string) {
	for _, part := range parts[:len(parts)-1] {
		next, ok := kvMap[part].(map[string]any)
		if !ok {
			return
		}
		kvMap = next
	}
	delete(kvMap, parts[len(parts)-1])
}

// getPath walks parts and returns the leaf value and whether it was
// found.
func getPath(kvMap map[string]any, parts []string) (any, bool) {
	for _, part := range parts[:len(parts)-1] {
		next, ok := kvMap[part].(map[string]any)
		if !ok {
			return nil, false
		}
		kvMap = next
	}
	value, ok := kvMap[parts[len(parts)-1]]

	return value, ok
}

// mergeMaps deep-merges src into dst: nested maps are merged
// recursively, every other value overwrites dst.
func mergeMaps(dst, src map[string]any) {
	for key, value := range src {
		if srcMap, ok := value.(map[string]any); ok {
			if dstMap, ok := dst[key].(map[string]any); ok {
				mergeMaps(dstMap, srcMap)
				continue
			}
		}
		dst[key] = value
	}
}
