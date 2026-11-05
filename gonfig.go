// Package gonfig provides hierarchical key-value configuration.
package gonfig

import (
	"errors"
	"sync"
)

// ErrNilProvider separates a wiring mistake from a source failure.
var ErrNilProvider = errors.New("gonfig: nil provider")

// Gonfig is a concurrency-safe hierarchical key-value container. Its zero value is ready to use.
type Gonfig struct {
	rwMutex sync.RWMutex
	kvMap   map[string]any
}

// New returns an empty Gonfig for symmetry with the NewWithXXX family.
// The zero value is already usable, so it is optional.
func New() *Gonfig {
	return &Gonfig{kvMap: make(map[string]any)}
}

// NewWithMap never fails; the error result only keeps signatures uniform.
func NewWithMap(keyValues map[string]any) (*Gonfig, error) {
	var gfg *Gonfig = New()
	var mapCopy map[string]any = DeepCopyMap(keyValues)
	mergeMaps(gfg.kvMap, mapCopy)
	return gfg, nil
}

// NewWithProvider copies the result no matter the provider's own policy,
// so a provider's aliasing choice cannot reach the new container.
func NewWithProvider(provider Provider) (*Gonfig, error) {
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

// NewWithGonfig snapshots another Gonfig; nil yields an empty container,
// so optional parents need no call-site guard.
func NewWithGonfig(other *Gonfig) (*Gonfig, error) {
	if other == nil {
		return New(), nil
	}

	return NewWithMap(other.Raw())
}

// Get returns live configuration values: no parts give the whole map,
// and a map-valued key comes back by reference, so do not mutate results.
func (gfg *Gonfig) Get(parts ...string) (any, bool) {
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

// Set writes along the key path parts. With no parts it replaces the root;
// only a map is accepted, because the root must stay a map.
func (gfg *Gonfig) Set(parts []string, value any) *Gonfig {
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

// Delete removes a key by path parts;
// missing paths and non-map segments are ignored, so callers need no existence check.
func (gfg *Gonfig) Delete(parts ...string) *Gonfig {
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

// Merge copies the source before locking the receiver; lock order therefore cannot deadlock.
func (gfg *Gonfig) Merge(other *Gonfig) *Gonfig {
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

// MergeMap lets later sources win on conflicts and recursively merges nested maps;
// the argument is copied, so the caller keeps it.
func (gfg *Gonfig) MergeMap(keyValues map[string]any) *Gonfig {
	gfg.rwMutex.Lock()
	defer gfg.rwMutex.Unlock()
	gfg.ensureMapInitialized()
	mergeMaps(gfg.kvMap, DeepCopyMap(keyValues))

	return gfg
}

// Raw returns a deep copy, never nil; a snapshot is safe to mutate. Empty state needs no nil check.
func (gfg *Gonfig) Raw() map[string]any {
	gfg.rwMutex.RLock()
	defer gfg.rwMutex.RUnlock()

	var raw map[string]any = DeepCopyMap(gfg.kvMap)
	if raw == nil {
		raw = make(map[string]any)
	}

	return raw
}

// ensureMapInitialized allocates the root map on demand; callers hold the write lock.
func (gfg *Gonfig) ensureMapInitialized() *Gonfig {
	if gfg.kvMap == nil {
		gfg.kvMap = make(map[string]any)
	}
	return gfg
}

// setPath replaces a non-map value that blocks the path;
// a nested write therefore discards whatever that position held before.
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

// deletePath stops at a non-map segment; setPath would replace it instead.
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

// getPath stops at a non-map segment and reports the key as not found.
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

// mergeMaps does not copy src; callers pass a collection they already own.
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
