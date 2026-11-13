package gonfig

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

// Constructors.

// TestNewWithMapDoesNotAliasSource asserts the deep-copy invariant, which also
// proves the source contents are stored and served.
func TestNewWithMapDoesNotAliasSource(t *testing.T) {
	var src map[string]any = map[string]any{"db": map[string]any{"host": "localhost"}}

	var gfg *Gonfig
	var err error
	gfg, err = NewWithMap(src)
	if err != nil {
		t.Fatalf("NewWithMap() error = %v, want nil", err)
	}

	src["db"].(map[string]any)["host"] = "changed"

	assertValue(t, gfg, "localhost", true, "db", "host")
}

// TestNewWithMapNilIsEmpty covers the nil-input branch: no error, usable empty root.
func TestNewWithMapNilIsEmpty(t *testing.T) {
	var gfg *Gonfig
	var err error
	gfg, err = NewWithMap(nil)

	assertContainerContract(t, gfg, err, map[string]any{})
}

// stubProvider is a Provider returning a fixed result.
type stubProvider struct {
	keyValues map[string]any
}

func (sp stubProvider) Read() map[string]any { return sp.keyValues }

// TestNewWithProviderCopiesProviderMap asserts that a provider's aliasing
// choice cannot reach the new container; the read also proves the happy path.
func TestNewWithProviderCopiesProviderMap(t *testing.T) {
	var providerValues map[string]any = map[string]any{"db": map[string]any{"host": "localhost"}}
	var provider stubProvider = stubProvider{keyValues: providerValues}

	var gfg *Gonfig
	var err error
	gfg, err = NewWithProvider(provider)
	if err != nil {
		t.Fatalf("NewWithProvider() error = %v, want nil", err)
	}

	providerValues["db"].(map[string]any)["host"] = "changed"

	assertValue(t, gfg, "localhost", true, "db", "host")
}

// TestNewWithProviderNilMapIsEmpty covers the nil-output branch.
func TestNewWithProviderNilMapIsEmpty(t *testing.T) {
	var provider stubProvider = stubProvider{keyValues: nil}

	var gfg *Gonfig
	var err error
	gfg, err = NewWithProvider(provider)

	assertContainerContract(t, gfg, err, map[string]any{})
}

// TestNewWithProviderNilProvider covers the non-nil provider precondition.
func TestNewWithProviderNilProvider(t *testing.T) {
	var gfg *Gonfig
	var err error
	gfg, err = NewWithProvider(nil)

	assertConstructorError(t, gfg, err, ErrNilProvider)
}

// TestNewWithGonfigDoesNotAliasSource asserts the snapshot invariant: later
// source writes never reach the copy.
func TestNewWithGonfigDoesNotAliasSource(t *testing.T) {
	var src *Gonfig = New().Set([]string{"a"}, map[string]any{"b": 1})

	var gfg *Gonfig
	var err error
	gfg, err = NewWithGonfig(src)
	if err != nil {
		t.Fatalf("NewWithGonfig() error = %v, want nil", err)
	}

	src.Set([]string{"a", "b"}, 99)
	src.Set([]string{"new"}, true)

	assertValue(t, gfg, 1, true, "a", "b")
	assertValue(t, gfg, nil, false, "new")
}

// TestNewWithGonfigNilIsEmpty covers the nil-source branch, so optional
// parents need no call-site guard.
func TestNewWithGonfigNilIsEmpty(t *testing.T) {
	var gfg *Gonfig
	var err error
	gfg, err = NewWithGonfig(nil)

	assertContainerContract(t, gfg, err, map[string]any{})
}

// Get — RLock → [whole map: DeepCopyMap | getPath → deep copy]; the traversal
// branches live in the getPath pack.

// TestGetZeroValueReturnsNotFound covers Get's own nil-root branch.
func TestGetZeroValueReturnsNotFound(t *testing.T) {
	var gfg Gonfig

	assertValue(t, &gfg, nil, false)
	assertValue(t, &gfg, nil, false, "a")
}

// TestGetWholeConfigReturnsCopy covers the no-parts branch and the copy
// invariant: mutating the returned map never reaches the container.
func TestGetWholeConfigReturnsCopy(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a", "b"}, 1)

	var got any
	var ok bool
	got, ok = gfg.Get()

	var copied map[string]any
	var isMap bool
	copied, isMap = got.(map[string]any)
	if !ok || !isMap {
		t.Fatalf("Get() = %v, %v; want whole map, true", got, ok)
	}

	copied["a"].(map[string]any)["b"] = 99
	copied["new"] = true

	assertValue(t, gfg, 1, true, "a", "b")
	assertValue(t, gfg, nil, false, "new")
}

// TestGetMapValueReturnsCopy covers the path branch: a map-valued key comes
// back as a fresh copy.
func TestGetMapValueReturnsCopy(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, map[string]any{"b": 1})

	var got any
	var ok bool
	got, ok = gfg.Get("a")

	var copied map[string]any
	var isMap bool
	copied, isMap = got.(map[string]any)
	if !ok || !isMap {
		t.Fatalf("Get(a) = %v, %v; want map, true", got, ok)
	}

	copied["b"] = 99

	assertValue(t, gfg, 1, true, "a", "b")
}

// TestGetSliceReturnsCopy covers a slice-valued key: the returned slice keeps
// no backing array in common with the container.
func TestGetSliceReturnsCopy(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"list"}, []string{"a"})

	var got any
	var ok bool
	got, ok = gfg.Get("list")

	var copied []string
	var isSlice bool
	copied, isSlice = got.([]string)
	if !ok || !isSlice {
		t.Fatalf("Get(list) = %v, %v; want []string, true", got, ok)
	}

	copied[0] = "b"

	assertValue(t, gfg, []string{"a"}, true, "list")
}

// TestGetTypedMapReturnsCopy covers a typed map value: the returned map keeps
// no entries in common with the container.
func TestGetTypedMapReturnsCopy(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"cfg"}, map[string]string{"a": "1"})

	var got any
	var ok bool
	got, ok = gfg.Get("cfg")

	var copied map[string]string
	var isMap bool
	copied, isMap = got.(map[string]string)
	if !ok || !isMap {
		t.Fatalf("Get(cfg) = %v, %v; want map[string]string, true", got, ok)
	}

	copied["a"] = "2"

	assertValue(t, gfg, map[string]string{"a": "1"}, true, "cfg")
}

// TestGetKeysAreLiteral pins that parts are literal keys, not dot-split paths:
// a key containing a dot and the empty string are ordinary keys.
func TestGetKeysAreLiteral(t *testing.T) {
	var dotted *Gonfig = New().Set([]string{"a.b"}, 1)

	assertValue(t, dotted, 1, true, "a.b")
	assertValue(t, dotted, nil, false, "a", "b")

	var empty *Gonfig = New().Set([]string{""}, 1)

	assertValue(t, empty, 1, true, "")
}

// TestGetCreatesNoKeys pins the read-only invariant: failed lookups leave the
// container unchanged, so Get never creates intermediate maps.
func TestGetCreatesNoKeys(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a", "b"}, 1)
	var want map[string]any = map[string]any{"a": map[string]any{"b": 1}}

	assertValue(t, gfg, nil, false, "a", "missing")     // missing leaf
	assertValue(t, gfg, nil, false, "x", "y")           // missing branch
	assertValue(t, gfg, nil, false, "a", "b", "deeper") // blocked by a scalar

	assertRaw(t, gfg, want)
}

// getPath — literal traversal, stops at a non-map segment.

// TestGetPathReturnsDeepValue covers the deep-path hit.
func TestGetPathReturnsDeepValue(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": map[string]any{"b": map[string]any{"c": "deep"}}}

	var got any
	var ok bool
	got, ok = getPath(kvMap, []string{"a", "b", "c"})

	if !ok || got != "deep" {
		t.Errorf("getPath(a.b.c) = %v, %v; want deep, true", got, ok)
	}
}

// TestGetPathMissingReportsNotFound covers the miss branches: missing leaf,
// missing intermediate and a scalar blocking the path all report not found.
func TestGetPathMissingReportsNotFound(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": map[string]any{"b": 1}, "scalar": 2}

	var got any
	var ok bool
	got, ok = getPath(kvMap, []string{"a", "missing"})
	if ok {
		t.Errorf("getPath(a.missing) = %v, %v; want not found", got, ok)
	}

	got, ok = getPath(kvMap, []string{"x", "y"})
	if ok {
		t.Errorf("getPath(x.y) = %v, %v; want not found", got, ok)
	}

	got, ok = getPath(kvMap, []string{"scalar", "b"})
	if ok {
		t.Errorf("getPath(scalar.b) = %v, %v; want not found (non-map segment stops the walk)", got, ok)
	}
}

// TestGetPathStoredNilIsFound separates a present nil value from a missing key.
func TestGetPathStoredNilIsFound(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": nil}

	var got any
	var ok bool
	got, ok = getPath(kvMap, []string{"a"})

	if !ok || got != nil {
		t.Errorf("getPath(a) = %v, %v; want nil, true (a stored nil is present)", got, ok)
	}
}

// Set — Lock → ensureMapInitialized → [root: DeepCopyMap | setPath]; the path
// branches live in the setPath pack.

// TestZeroValueIsUsable covers the zero-value contract: the zero value needs
// no constructor and Set initializes storage.
func TestZeroValueIsUsable(t *testing.T) {
	var gfg Gonfig
	gfg.Set([]string{"a"}, 1) // must not panic on nil map

	assertValue(t, &gfg, 1, true, "a")
}

// TestSetEmptyKeyReplacesConfig covers the root-replacement branch: with no
// parts and a map argument, the root is replaced by a deep copy.
func TestSetEmptyKeyReplacesConfig(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"old"}, 1)
	var replacement map[string]any = map[string]any{"new": map[string]any{"leaf": 2}}
	gfg.Set(nil, replacement)

	replacement["new"].(map[string]any)["leaf"] = 99

	assertRaw(t, gfg, map[string]any{"new": map[string]any{"leaf": 2}})
}

// TestSetEmptyKeyWithNonMapIsNoop covers the root guard: the root must stay a
// map, so a non-map root value is ignored.
func TestSetEmptyKeyWithNonMapIsNoop(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"old"}, 1)
	gfg.Set(nil, 42)

	assertRaw(t, gfg, map[string]any{"old": 1})
}

// TestSetRootTypedNilMapLeavesEmptyConfig covers the root normalization: a
// typed nil root replacement leaves a usable empty root, so Get() and Raw()
// agree instead of reporting "not found" versus "empty".
func TestSetRootTypedNilMapLeavesEmptyConfig(t *testing.T) {
	var nilMap map[string]any
	var gfg *Gonfig = New().Set([]string{"a"}, 1)
	gfg.Set(nil, nilMap)

	assertRaw(t, gfg, map[string]any{})

	var got any
	var ok bool
	got, ok = gfg.Get()
	if !ok {
		t.Fatal("Get() not found, want an empty root map")
	}

	var root map[string]any
	var isMap bool
	root, isMap = got.(map[string]any)
	if !isMap || root == nil {
		t.Fatalf("Get() = %v, want a non-nil root map", got)
	}
}

// setPath — creates intermediates and writes the leaf.

// TestSetPathCreatesIntermediateMaps covers the creation branch and the
// pass-through of a scalar leaf.
func TestSetPathCreatesIntermediateMaps(t *testing.T) {
	var kvMap map[string]any = map[string]any{}
	setPath(kvMap, []string{"a", "b", "c"}, 1)

	var want map[string]any = map[string]any{"a": map[string]any{"b": map[string]any{"c": 1}}}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v", kvMap, want)
	}
}

// TestSetPathReplacesBlockingScalar covers the temporary invariant break where
// a scalar blocks the path: it is replaced and the old value is lost.
func TestSetPathReplacesBlockingScalar(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": "scalar"}
	setPath(kvMap, []string{"a", "b"}, 2)

	var want map[string]any = map[string]any{"a": map[string]any{"b": 2}}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v", kvMap, want)
	}
}

// TestSetPathOpensTypedNilMapOnPath covers the temporary invariant break where
// a typed nil map sits on the path: it counts as absent and is replaced.
func TestSetPathOpensTypedNilMapOnPath(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": map[string]any(nil)}
	setPath(kvMap, []string{"a", "b"}, 2) // must not panic

	var want map[string]any = map[string]any{"a": map[string]any{"b": 2}}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v", kvMap, want)
	}
}

// TestSetPathCopiesMapLeaf covers the leaf deep copy: the caller's map is not
// aliased at any map level.
func TestSetPathCopiesMapLeaf(t *testing.T) {
	var kvMap map[string]any = map[string]any{}
	var value map[string]any = map[string]any{"b": map[string]any{"c": 1}}
	setPath(kvMap, []string{"a"}, value)

	value["b"].(map[string]any)["c"] = 99

	if got := kvMap["a"].(map[string]any)["b"].(map[string]any)["c"]; got != 1 {
		t.Errorf("stored leaf = %v, want 1 (setPath must copy map values)", got)
	}
}

// TestSetPathTypedNilLeafBecomesEmptyMap covers the leaf normalization: a map
// argument that is a typed nil is stored as an empty map, so later writes work.
func TestSetPathTypedNilLeafBecomesEmptyMap(t *testing.T) {
	var kvMap map[string]any = map[string]any{}
	var nilMap map[string]any
	setPath(kvMap, []string{"a"}, nilMap)

	var stored map[string]any
	var isMap bool
	stored, isMap = kvMap["a"].(map[string]any)
	if !isMap {
		t.Fatalf("stored leaf = %T, want map[string]any", kvMap["a"])
	}
	if stored == nil {
		t.Error("stored leaf = nil map, want the typed nil normalized to an empty map")
	}
}

// TestSetPathCopiesNonMapValues covers the copy depth: a non-map leaf such as a
// slice is copied, so the caller keeps its own backing array.
func TestSetPathCopiesNonMapValues(t *testing.T) {
	var kvMap map[string]any = map[string]any{}
	var value []string = []string{"a"}
	setPath(kvMap, []string{"list"}, value)

	value[0] = "b"

	if got := kvMap["list"].([]string)[0]; got != "a" {
		t.Errorf("stored slice = %q, want an independent copy %q", got, "a")
	}
}

// ensureMapInitialized — callers hold the write lock.

// TestEnsureMapInitializedAllocatesNilRoot covers the nil-root branch: the
// zero value gets usable storage and the receiver comes back for chaining.
func TestEnsureMapInitializedAllocatesNilRoot(t *testing.T) {
	var gfg Gonfig

	var got *Gonfig = gfg.ensureMapInitialized()

	if got != &gfg {
		t.Errorf("ensureMapInitialized() = %p, want receiver %p", got, &gfg)
	}
	if gfg.kvMap == nil {
		t.Fatal("kvMap = nil, want an allocated root")
	}
}

// TestEnsureMapInitializedKeepsExistingRoot covers the no-op branch: a
// populated root is never reset.
func TestEnsureMapInitializedKeepsExistingRoot(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, 1)

	gfg.ensureMapInitialized()

	assertRaw(t, gfg, map[string]any{"a": 1})
}

// Delete — Lock → [empty parts: make | deletePath]; the traversal branches
// live in the deletePath pack.

// TestDeleteOnZeroValueIsNoop covers Delete's nil-root guard.
func TestDeleteOnZeroValueIsNoop(t *testing.T) {
	var gfg Gonfig

	gfg.Delete("a")

	assertRaw(t, &gfg, map[string]any{})
}

// TestDeleteEmptyKeyClearsConfig covers the root-clear branch.
func TestDeleteEmptyKeyClearsConfig(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, 1)

	gfg.Delete()

	assertRaw(t, gfg, map[string]any{})
}

// deletePath — stops at a non-map segment and removes the leaf only.

// TestDeletePathRemovesLeaf covers the leaf removal; the sibling survives.
func TestDeletePathRemovesLeaf(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": map[string]any{"b": 1, "c": 2}}
	deletePath(kvMap, []string{"a", "b"})

	var want map[string]any = map[string]any{"a": map[string]any{"c": 2}}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v", kvMap, want)
	}
}

// TestDeletePathKeepsEmptyParent pins the no-pruning invariant: removing the
// only child leaves the empty ancestors in place.
func TestDeletePathKeepsEmptyParent(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": map[string]any{"b": 1}}
	deletePath(kvMap, []string{"a", "b"})

	var want map[string]any = map[string]any{"a": map[string]any{}}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v (parents must not be pruned)", kvMap, want)
	}
}

// TestDeletePathMissingIsNoop covers the missing-path branch and the
// non-map stop: no panic, no replacement, no change.
func TestDeletePathMissingIsNoop(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": 1, "nested": map[string]any{"leaf": 2}}
	deletePath(kvMap, []string{"x", "y", "z"}) // missing path, must not panic
	deletePath(kvMap, []string{"a", "b"})      // blocked by a scalar

	var want map[string]any = map[string]any{"a": 1, "nested": map[string]any{"leaf": 2}}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v", kvMap, want)
	}
}

// TestDeletePathTypedNilMapIsNoop pins that a typed nil map on the path is
// safe: the walk descends into it and the delete is a no-op.
func TestDeletePathTypedNilMapIsNoop(t *testing.T) {
	var kvMap map[string]any = map[string]any{"a": map[string]any(nil)}
	deletePath(kvMap, []string{"a", "b"}) // must not panic

	var want map[string]any = map[string]any{"a": map[string]any(nil)}

	if !reflect.DeepEqual(kvMap, want) {
		t.Errorf("kvMap = %v, want %v", kvMap, want)
	}
}

// Merge — other.RLock → DeepCopyMap → other.RUnlock → Lock →
// ensureMapInitialized → mergeMaps; the union branches live in the mergeMaps pack.

// TestMergeMergesGonfigConfig pins the end-to-end wiring: snapshot, then union.
func TestMergeMergesGonfigConfig(t *testing.T) {
	var other *Gonfig = New().Set([]string{"db", "host"}, "remote").Set([]string{"db", "port"}, 3306)
	var gfg *Gonfig = New().Set([]string{"db"}, map[string]any{"host": "localhost", "user": "admin"})
	gfg.Merge(other)

	assertRaw(t, gfg, map[string]any{"db": map[string]any{"host": "remote", "port": 3306, "user": "admin"}})
}

// TestMergeDoesNotAliasSourceGonfig asserts the snapshot invariant: the source
// is copied before it is merged, so later source writes stay invisible.
func TestMergeDoesNotAliasSourceGonfig(t *testing.T) {
	var other *Gonfig = New().Set([]string{"a"}, map[string]any{"b": 1})
	var gfg *Gonfig = New()
	gfg.Merge(other)

	other.Set([]string{"a", "b"}, 99)
	other.Set([]string{"new"}, true)

	assertValue(t, gfg, 1, true, "a", "b")
	assertValue(t, gfg, nil, false, "new")
}

// TestMergeNilGonfigIsNoop covers the nil-source precondition.
func TestMergeNilGonfigIsNoop(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, 1)
	gfg.Merge(nil)

	assertRaw(t, gfg, map[string]any{"a": 1})
}

// TestMergeIntoZeroValue covers the zero-value receiver: storage is
// initialized and the source is applied.
func TestMergeIntoZeroValue(t *testing.T) {
	var gfg Gonfig
	var other *Gonfig = New().Set([]string{"db", "host"}, "remote")

	gfg.Merge(other)

	assertRaw(t, &gfg, map[string]any{"db": map[string]any{"host": "remote"}})
}

// TestMergeZeroValueSourceIsNoop covers the nil-source-root branch: copying a
// zero-value source yields nil, and merging it changes nothing.
func TestMergeZeroValueSourceIsNoop(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, 1)
	var other Gonfig

	gfg.Merge(&other)

	assertRaw(t, gfg, map[string]any{"a": 1})
}

// MergeMap — Lock → ensureMapInitialized → DeepCopyMap → mergeMaps.

// TestMergeDeepMergesNestedMaps pins the end-to-end wiring: the argument is
// copied and deeply united, the later source wins.
func TestMergeDeepMergesNestedMaps(t *testing.T) {
	var gfg *Gonfig = New()
	gfg.MergeMap(map[string]any{
		"app":  "gonfig",
		"db":   map[string]any{"host": "localhost", "port": 5432},
		"tags": []string{"a", "b"},
	})
	gfg.MergeMap(map[string]any{
		"app": "gonfig2",
		"db":  map[string]any{"port": 5433, "user": "admin"},
	})

	assertRaw(t, gfg, map[string]any{
		"app":  "gonfig2",
		"db":   map[string]any{"host": "localhost", "port": 5433, "user": "admin"},
		"tags": []string{"a", "b"},
	})
}

// TestMergeDoesNotAliasSourceMap asserts the argument is copied, so the caller
// keeps it.
func TestMergeDoesNotAliasSourceMap(t *testing.T) {
	var gfg *Gonfig = New()
	src := map[string]any{"db": map[string]any{"host": "localhost"}}
	gfg.MergeMap(src)

	src["db"].(map[string]any)["host"] = "changed"

	assertValue(t, gfg, "localhost", true, "db", "host")
}

// TestMergeMapNilSourceIsNoop covers the nil-argument precondition.
func TestMergeMapNilSourceIsNoop(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, 1)

	gfg.MergeMap(nil)

	assertRaw(t, gfg, map[string]any{"a": 1})
}

// TestMergeMapOnZeroValue covers the zero-value receiver: storage is
// initialized and the argument is applied.
func TestMergeMapOnZeroValue(t *testing.T) {
	var gfg Gonfig

	gfg.MergeMap(map[string]any{"a": 1})

	assertRaw(t, &gfg, map[string]any{"a": 1})
}

// mergeMaps — deep union; the caller passes a collection it already owns.

// TestMergeMapsDeepUnions covers the recursion: nested maps unite, the source
// wins on conflict, unrelated keys survive.
func TestMergeMapsDeepUnions(t *testing.T) {
	var dst map[string]any = map[string]any{"a": map[string]any{"x": 1, "y": 2}, "keep": true}
	mergeMaps(dst, map[string]any{"a": map[string]any{"y": 3, "z": 4}, "new": 5})

	var want map[string]any = map[string]any{
		"a":    map[string]any{"x": 1, "y": 3, "z": 4},
		"keep": true,
		"new":  5,
	}

	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst = %v, want %v", dst, want)
	}
}

// TestMergeMapsReplacesNonMapWithSourceMap covers the type-change branch:
// when the source value is a map and the destination value is not, the source wins.
func TestMergeMapsReplacesNonMapWithSourceMap(t *testing.T) {
	var dst map[string]any = map[string]any{"a": 1}
	mergeMaps(dst, map[string]any{"a": map[string]any{"b": 2}})

	var want map[string]any = map[string]any{"a": map[string]any{"b": 2}}

	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst = %v, want %v", dst, want)
	}
}

// TestMergeMapsReplacesTypedNilDstMap covers the merge-over-nil branch: a
// typed nil destination map is replaced before the recursion, so it cannot panic.
func TestMergeMapsReplacesTypedNilDstMap(t *testing.T) {
	var dst map[string]any = map[string]any{"a": map[string]any(nil)}
	mergeMaps(dst, map[string]any{"a": map[string]any{"b": 1}}) // must not panic

	var want map[string]any = map[string]any{"a": map[string]any{"b": 1}}

	if !reflect.DeepEqual(dst, want) {
		t.Errorf("dst = %v, want %v", dst, want)
	}
}

// TestMergeMapsDoesNotCopySource pins the documented no-copy contract: values
// merged in from the source are the caller's own collection.
func TestMergeMapsDoesNotCopySource(t *testing.T) {
	var src map[string]any = map[string]any{"a": map[string]any{"b": 1}}
	var dst map[string]any = map[string]any{}
	mergeMaps(dst, src)

	src["a"].(map[string]any)["b"] = 99

	if got := dst["a"].(map[string]any)["b"]; got != 99 {
		t.Errorf("merged value = %v, want the source's own map (mergeMaps must not copy src)", got)
	}
}

// Raw — RLock → DeepCopyMap → nil normalization; the copy depth lives in the
// DeepCopyMap pack.

// TestRawReturnsCopy asserts the deep-copy boundary: mutating the snapshot at
// top or nested level never reaches the container.
func TestRawReturnsCopy(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, map[string]any{"b": 1})

	var actualValues map[string]any = gfg.Raw()
	actualValues["a"].(map[string]any)["b"] = 99
	actualValues["new"] = true

	assertValue(t, gfg, 1, true, "a", "b")
	assertValue(t, gfg, nil, false, "new")
}

// TestRawSliceReturnsCopy covers the slice depth: a slice obtained from Raw
// keeps no backing array in common with the container.
func TestRawSliceReturnsCopy(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"list"}, []string{"a"})

	var raw map[string]any = gfg.Raw()
	raw["list"].([]string)[0] = "b"

	assertValue(t, gfg, []string{"a"}, true, "list")
}

// TestRawNeverReturnsNil covers the nil normalization on the empty state.
func TestRawNeverReturnsNil(t *testing.T) {
	var gfg Gonfig

	var raw map[string]any = gfg.Raw()

	if raw == nil {
		t.Fatal("Raw() = nil, want an empty map")
	}
}

// DeepCopyMap — recursive copy; exported, so it has its own pack.

// TestDeepCopyMapNil covers the nil input: nil returns nil.
func TestDeepCopyMapNil(t *testing.T) {
	var copied map[string]any = DeepCopyMap(nil)

	if copied != nil {
		t.Errorf("DeepCopyMap(nil) = %v, want nil", copied)
	}
}

// TestDeepCopyMapEmptyMap covers the empty input: a non-nil empty map.
func TestDeepCopyMapEmptyMap(t *testing.T) {
	var copied map[string]any = DeepCopyMap(map[string]any{})

	if copied == nil {
		t.Fatal("DeepCopyMap(empty) = nil, want a non-nil empty map")
	}
	if len(copied) != 0 {
		t.Errorf("DeepCopyMap(empty) = %v, want an empty map", copied)
	}
}

// TestDeepCopyMapCopiesNestedMaps asserts the recursive copy: nested map
// levels are independent.
func TestDeepCopyMapCopiesNestedMaps(t *testing.T) {
	var src map[string]any = map[string]any{"a": map[string]any{"b": map[string]any{"c": 1}}}
	var copied map[string]any = DeepCopyMap(src)

	src["a"].(map[string]any)["b"].(map[string]any)["c"] = 99

	if got := copied["a"].(map[string]any)["b"].(map[string]any)["c"]; got != 1 {
		t.Errorf("copied nested value = %v, want 1 (nested maps must be copied)", got)
	}
}

// TestDeepCopyMapCopiesSlices asserts the slice depth: the copied slice keeps
// its own backing array.
func TestDeepCopyMapCopiesSlices(t *testing.T) {
	var src map[string]any = map[string]any{"list": []string{"a"}}
	var copied map[string]any = DeepCopyMap(src)

	src["list"].([]string)[0] = "b"

	if got := copied["list"].([]string)[0]; got != "a" {
		t.Errorf("copied slice = %q, want an independent copy %q", got, "a")
	}
}

// TestDeepCopyMapCopiesTypedMaps asserts the typed-map depth: maps other than
// map[string]any are copied too.
func TestDeepCopyMapCopiesTypedMaps(t *testing.T) {
	var src map[string]any = map[string]any{"cfg": map[string]string{"a": "1"}}
	var copied map[string]any = DeepCopyMap(src)

	src["cfg"].(map[string]string)["a"] = "2"

	if got := copied["cfg"].(map[string]string)["a"]; got != "1" {
		t.Errorf("copied typed map = %q, want an independent copy %q", got, "1")
	}
}

// TestDeepCopyMapCopiesPointers asserts the pointer depth: the copy points at
// a fresh value.
func TestDeepCopyMapCopiesPointers(t *testing.T) {
	var srcValue int = 1
	var src map[string]any = map[string]any{"ptr": &srcValue}
	var copied map[string]any = DeepCopyMap(src)

	*src["ptr"].(*int) = 2

	if got := *copied["ptr"].(*int); got != 1 {
		t.Errorf("copied pointer value = %d, want an independent copy 1", got)
	}
}

// TestDeepCopyMapCopiesStructFields covers the struct branch: exported
// reference fields are copied with the struct value.
func TestDeepCopyMapCopiesStructFields(t *testing.T) {
	type payload struct {
		Values []string
	}

	var src map[string]any = map[string]any{"cfg": payload{Values: []string{"a"}}}
	var copied map[string]any = DeepCopyMap(src)

	src["cfg"].(payload).Values[0] = "b"

	if got := copied["cfg"].(payload).Values[0]; got != "a" {
		t.Errorf("copied struct field = %q, want an independent copy %q", got, "a")
	}
}

// TestDeepCopyMapCopiesUnexportedStructFields covers a struct type with
// unexported fields: the copy must not panic on the fields it cannot set.
func TestDeepCopyMapCopiesUnexportedStructFields(t *testing.T) {
	var moment time.Time = time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	var src map[string]any = map[string]any{"moment": moment}
	var copied map[string]any = DeepCopyMap(src)

	if got := copied["moment"].(time.Time); !got.Equal(moment) {
		t.Errorf("copied time = %v, want %v", got, moment)
	}
}

// Fluent chaining — one test for the receiver-return postconditions of Set,
// MergeMap and Delete.
func TestFluentChaining(t *testing.T) {
	var gfg *Gonfig = New().Set([]string{"a"}, 1).MergeMap(map[string]any{"b": 2}).Delete("a")

	assertRaw(t, gfg, map[string]any{"b": 2})
}

// Concurrency — every exported method is guarded by sync.RWMutex; Merge takes
// the source read lock before the receiver write lock, so the lock order
// cannot deadlock.

// TestConcurrentDistinctPaths covers write serialization: all writers land.
func TestConcurrentDistinctPaths(t *testing.T) {
	const writers int = 100

	var gfg *Gonfig = New()
	var waitGroup sync.WaitGroup = sync.WaitGroup{}
	waitGroup.Add(writers)

	for index := 0; index < writers; index++ {
		go func(path string) {
			defer waitGroup.Done()
			gfg.Set([]string{path}, path)
		}(fmt.Sprintf("key%d", index))
	}

	waitGroup.Wait()

	var want map[string]any = make(map[string]any, writers)
	for index := 0; index < writers; index++ {
		var path string = fmt.Sprintf("key%d", index)
		want[path] = path
	}

	assertRaw(t, gfg, want)
}

// TestConcurrentRawSnapshot covers the snapshot-under-lock invariant.
func TestConcurrentRawSnapshot(t *testing.T) {
	const writers int = 50
	const readers int = 50

	var gfg *Gonfig = New()
	var waitGroup sync.WaitGroup = sync.WaitGroup{}
	waitGroup.Add(writers + readers)

	for index := 0; index < writers; index++ {
		go func(path string) {
			defer waitGroup.Done()
			gfg.Set([]string{path}, path)
		}(fmt.Sprintf("key%d", index))
	}

	for index := 0; index < readers; index++ {
		go func() {
			defer waitGroup.Done()
			_ = gfg.Raw()
		}()
	}

	waitGroup.Wait()

	var want map[string]any = make(map[string]any, writers)
	for index := 0; index < writers; index++ {
		var path string = fmt.Sprintf("key%d", index)
		want[path] = path
	}

	assertRaw(t, gfg, want)
}

// TestConcurrentGetAndSet covers the read/write exclusion under the race detector.
func TestConcurrentGetAndSet(t *testing.T) {
	const writers int = 50
	const readers int = 50

	var gfg *Gonfig = New()
	var waitGroup sync.WaitGroup = sync.WaitGroup{}
	waitGroup.Add(writers + readers)

	for index := 0; index < writers; index++ {
		go func(path string) {
			defer waitGroup.Done()
			gfg.Set([]string{path}, path)
		}(fmt.Sprintf("key%d", index))
	}

	for index := 0; index < readers; index++ {
		go func(path string) {
			defer waitGroup.Done()
			gfg.Get(path)
		}(fmt.Sprintf("key%d", index))
	}

	waitGroup.Wait()

	var want map[string]any = make(map[string]any, writers)
	for index := 0; index < writers; index++ {
		var path string = fmt.Sprintf("key%d", index)
		want[path] = path
	}

	assertRaw(t, gfg, want)
}

// TestConcurrentGetAndDelete covers the read/delete exclusion; every key is
// deleted exactly once, so the final state is deterministic.
func TestConcurrentGetAndDelete(t *testing.T) {
	const writers int = 50
	const readers int = 50

	var gfg *Gonfig = New()
	var waitGroup sync.WaitGroup = sync.WaitGroup{}
	waitGroup.Add(writers + readers)

	for index := 0; index < writers; index++ {
		var path string = fmt.Sprintf("key%d", index)
		gfg.Set([]string{path}, path)
	}

	for index := 0; index < writers; index++ {
		go func(path string) {
			defer waitGroup.Done()
			gfg.Delete(path)
		}(fmt.Sprintf("key%d", index))
	}

	for index := 0; index < readers; index++ {
		go func(path string) {
			defer waitGroup.Done()
			gfg.Get(path)
		}(fmt.Sprintf("key%d", index))
	}

	waitGroup.Wait()

	assertRaw(t, gfg, map[string]any{})
}

// TestConcurrentMutualMergeDoesNotDeadlock pins Merge's lock order: the source
// is copied under its read lock and released before the receiver write lock, so
// two containers can merge into each other without a cycle.
func TestConcurrentMutualMergeDoesNotDeadlock(t *testing.T) {
	var first *Gonfig = New().Set([]string{"owner"}, "first")
	var second *Gonfig = New().Set([]string{"owner"}, "second")

	done := make(chan struct{}, 2)
	go func() {
		first.Merge(second)
		done <- struct{}{}
	}()
	go func() {
		second.Merge(first)
		done <- struct{}{}
	}()

	for index := 0; index < 2; index++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("mutual Merge deadlocked")
		}
	}

	if _, ok := first.Get("owner"); !ok {
		t.Error("first lost its owner key")
	}
	if _, ok := second.Get("owner"); !ok {
		t.Error("second lost its owner key")
	}
}

// assertRaw is the single check of a container's whole observable state.
func assertRaw(t *testing.T, gfg *Gonfig, want map[string]any) {
	t.Helper()

	if gfg == nil {
		t.Fatal("assertRaw: nil Gonfig, want initialized container")
	}

	var got map[string]any = gfg.Raw()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Raw() = %v, want %v", got, want)
	}
}

// assertValue is the single check of one Get lookup.
func assertValue(t *testing.T, gfg *Gonfig, wantValue any, wantFound bool, parts ...string) {
	t.Helper()

	if gfg == nil {
		t.Fatal("assertValue: nil Gonfig, want initialized container")
	}

	var gotValue any
	var gotFound bool
	gotValue, gotFound = gfg.Get(parts...)

	if gotFound != wantFound {
		t.Fatalf("Get(%v) found = %v, want %v", parts, gotFound, wantFound)
	}
	if !wantFound {
		return
	}

	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("Get(%v) = %v, want %v", parts, gotValue, wantValue)
	}
}

// assertContainerContract is the one place that checks the shared constructor
// postconditions: no error, non-nil container, exact contents.
func assertContainerContract(t *testing.T, gfg *Gonfig, err error, want map[string]any) {
	t.Helper()

	if err != nil {
		t.Fatalf("constructor error = %v, want nil", err)
	}
	if gfg == nil {
		t.Fatal("constructor returned nil container, want non-nil")
	}

	assertRaw(t, gfg, want)
}

// assertConstructorError covers the constructor failure contract: the error is
// the expected one and no partially built container escapes.
func assertConstructorError(t *testing.T, gfg *Gonfig, err error, wantErr error) {
	t.Helper()

	if !errors.Is(err, wantErr) {
		t.Errorf("constructor error = %v, want %v", err, wantErr)
	}
	if gfg != nil {
		t.Errorf("constructor = %v, want nil container on error", gfg)
	}
}
