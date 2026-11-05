// Package gofig provides a thread-safe, modular configuration container
// for managing hierarchical key-value pairs. It enforces boundary safety
// by deep-copying inbound and outbound data maps, isolating internal state.
package gofig

import (
	"errors"
	"sync"
)

// ErrNilProvider is a sentinel error returned by NewWithProvider when
// a nil Provider source is supplied. It signals that configuration
// data cannot be loaded from a non-existent provider.
var ErrNilProvider = errors.New("gofig: nil provider")

// Gofig is the central configuration manager that stores and manages
// hierarchical key-value pairs. It provides thread-safe access guarded
// by an internal RWMutex and enforces boundary safety across all operations.
type Gofig struct {
	rwMutex sync.RWMutex
	kvMap   map[string]any
}

// New returns a new empty Gofig.
func New() *Gofig {
	return &Gofig{kvMap: make(map[string]any)}
}

// NewWithMap returns a new Gofig loaded with a deep copy of the given
// key-value pairs.
func NewWithMap(keyValues map[string]any) (*Gofig, error) {
	var gfg *Gofig = New()
	var mapCopy map[string]any = DeepCopyMap(keyValues)
	mergeMaps(gfg.kvMap, mapCopy)
	return gfg, nil
}

// NewWithProvider returns a new Gofig loaded with a deep copy of the
// key-value pairs read from the given Provider.
func NewWithProvider(provider Provider) (*Gofig, error) {
	if provider == nil {
		return nil, ErrNilProvider
	}

	var keyValues map[string]any
	var err error
	keyValues, err = provider.Read()
	if err != nil {
		return nil, err
	}

	return NewWithMap(keyValues)
}

// NewWithGofig returns a new Gofig loaded with a deep copy of the given
// Gofig's configuration.
func NewWithGofig(other *Gofig) (*Gofig, error) {
	if other == nil {
		return New(), nil
	}

	return NewWithMap(other.Raw())
}

// Get returns the value at the given path parts, reporting whether the
// key was found. It performs a read-locked traversal over the nested
// configuration maps.
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

// Set stores a value at the given path parts, creating missing intermediate
// maps as needed. It enforces boundary safety by deep-copying map values
// and returns the receiver to enable method chaining.
func (gfg *Gofig) Set(parts []string, value any) *Gofig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.ensureMapInitialized()

	if len(parts) == 0 {
		if newConf, ok := value.(map[string]any); ok {
			gfg.kvMap = DeepCopyMap(newConf)
		}
		return gfg
	}

	setPath(gfg.kvMap, parts, value)

	return gfg
}

// Delete removes the configuration value at the given path parts.
// When called with no path parts, it clears the entire configuration map
// and returns the receiver for fluent chaining.
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

// Merge deep-merges the configuration of another Gofig into this one.
// It snapshots the source container under a read lock and merges nested
// maps recursively, returning the receiver for method chaining.
func (gfg *Gofig) Merge(other *Gofig) *Gofig {
	if other == nil {
		return gfg
	}

	var otherMap map[string]any
	other.rwMutex.RLock()
	otherMap = DeepCopyMap(other.kvMap)
	other.rwMutex.RUnlock()

	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()

	gfg.ensureMapInitialized()
	mergeMaps(gfg.kvMap, otherMap)

	return gfg
}

// MergeMap deep-merges the given key-value pairs into the configuration.
// It creates a deep copy of the input map to preserve boundary safety
// and returns the receiver for method chaining.
func (gfg *Gofig) MergeMap(keyValues map[string]any) *Gofig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.ensureMapInitialized()
	mergeMaps(gfg.kvMap, DeepCopyMap(keyValues))

	return gfg
}

// Raw returns a complete deep copy of the whole configuration map.
// It guarantees outbound boundary safety so that caller mutations do
// not alter internal Gofig state.
func (gfg *Gofig) Raw() map[string]any {
	gfg.rwMutex.RLock()
	defer gfg.rwMutex.RUnlock()

	var raw map[string]any = DeepCopyMap(gfg.kvMap)
	if raw == nil {
		raw = make(map[string]any)
	}

	return raw
}

// ensureMapInitialized allocates the internal configuration map if it is
// uninitialized and returns the receiver.
func (gfg *Gofig) ensureMapInitialized() *Gofig {
	if gfg.kvMap == nil {
		gfg.kvMap = make(map[string]any)
	}
	return gfg
}

// setPath walks path segments and stores value at the leaf, creating missing
// intermediate maps. Map values are deep-copied to prevent aliasing.
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
		kvMap[last] = DeepCopyMap(valueMap)
		return
	}
	kvMap[last] = value
}

// deletePath walks path segments and deletes the leaf key. Missing paths
// are ignored without error.
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

// getPath walks path segments through nested maps and returns the leaf
// value along with a boolean indicating whether the key was found.
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

// mergeMaps deep-merges src into dst: nested maps are merged recursively,
// while every other value overwrites dst.
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
