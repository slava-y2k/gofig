package gofig

import (
	"errors"
	"sync"
)

// ErrNilProvider is returned by NewWithProvider when the given provider
// is nil.
var ErrNilProvider = errors.New("gofig: nil provider")

// Gofig is a key-value pair manager which provides
type Gofig struct {
	rwMutex sync.RWMutex
	kvMap   map[string]any
}

func New() *Gofig {
	return &Gofig{kvMap: make(map[string]any)}
}

// NewWithMap creates a new Gofig loaded with a deep copy of the given
// key-value pairs.
func NewWithMap(keyValues map[string]any) (*Gofig, error) {
	var gfg *Gofig = New()
	mergeMaps(gfg.kvMap, DeepCopyMap(keyValues))

	return gfg, nil
}

// NewWithProvider creates a new Gofig populated from the supplied
// Provider.
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

// NewWithGofig creates a new Gofig loaded with a deep copy of the other
// Gofig's configuration.
func NewWithGofig(other *Gofig) (*Gofig, error) {
	if other == nil {
		return New(), nil
	}

	return NewWithMap(other.Raw())
}

// Get returns the value at the given path parts, reporting whether the
// key was found.
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
// intermediate maps.
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

// Delete removes the value at the given path parts. Missing paths are
// ignored.
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

// MergeMap deep-merges the given key-value pairs into the
// configuration.
func (gfg *Gofig) MergeMap(keyValues map[string]any) *Gofig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.ensureMapInitialized()
	mergeMaps(gfg.kvMap, DeepCopyMap(keyValues))

	return gfg
}

// Raw returns a deep copy of the whole configuration map.
func (gfg *Gofig) Raw() map[string]any {
	gfg.rwMutex.RLock()
	defer gfg.rwMutex.RUnlock()

	var raw map[string]any = DeepCopyMap(gfg.kvMap)
	if raw == nil {
		raw = make(map[string]any)
	}

	return raw
}

// ensureMapInitialized makes sure kvMap is non-nil.
func (gfg *Gofig) ensureMapInitialized() *Gofig {
	if gfg.kvMap == nil {
		gfg.kvMap = make(map[string]any)
	}
	return gfg
}

// setPath walks parts and stores value at the leaf, creating missing
// intermediate maps.
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
