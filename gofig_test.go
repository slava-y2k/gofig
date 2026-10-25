package gofig

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewReturnsEmptyConfig(t *testing.T) {
	var gfg *Gofig = New()

	var actualValues map[string]any = gfg.Raw()

	if len(actualValues) != 0 {
		t.Errorf("Raw() = %v, want empty map", actualValues)
	}

	var ok bool
	_, ok = gfg.Get("missing")
	if ok {
		t.Error("Get(missing) reported found, want not found")
	}
}

func TestZeroValueIsUsable(t *testing.T) {
	var gfg Gofig
	gfg.Set([]string{"a"}, 1) // must not panic on nil map

	var got any
	var ok bool
	got, ok = gfg.Get("a")

	if !ok || got != 1 {
		t.Errorf("Get(a) = %v, %v; want 1, true", got, ok)
	}
}

func TestMergeDeepMergesNestedMaps(t *testing.T) {
	var gfg *Gofig = New()
	gfg.MergeMap(map[string]any{
		"app":  "gofig",
		"db":   map[string]any{"host": "localhost", "port": 5432},
		"tags": []string{"a", "b"},
	})
	gfg.MergeMap(map[string]any{
		"app": "gofig2",
		"db":  map[string]any{"port": 5433, "user": "admin"},
	})

	want := map[string]any{
		"app":  "gofig2",
		"db":   map[string]any{"host": "localhost", "port": 5433, "user": "admin"},
		"tags": []string{"a", "b"},
	}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestMergeDoesNotAliasSourceMap(t *testing.T) {
	var gfg *Gofig = New()
	src := map[string]any{"db": map[string]any{"host": "localhost"}}
	gfg.MergeMap(src)

	src["db"].(map[string]any)["host"] = "changed"

	var got any
	got, _ = gfg.Get("db", "host")

	if got != "localhost" {
		t.Errorf("Get(db.host) = %v, want localhost (source map must not alias)", got)
	}
}

func TestMergeMergesGofigConfig(t *testing.T) {
	var other *Gofig = New().Set([]string{"db", "host"}, "remote").Set([]string{"db", "port"}, 3306)
	var gfg *Gofig = New().Set([]string{"db"}, map[string]any{"host": "localhost", "user": "admin"})
	gfg.Merge(other)

	want := map[string]any{"db": map[string]any{"host": "remote", "port": 3306, "user": "admin"}}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestMergeDoesNotAliasSourceGofig(t *testing.T) {
	var other *Gofig = New().Set([]string{"a"}, map[string]any{"b": 1})
	var gfg *Gofig = New()
	gfg.Merge(other)

	other.Set([]string{"a", "b"}, 99)
	other.Set([]string{"new"}, true)

	var got any
	var ok bool
	got, ok = gfg.Get("a", "b")

	if !ok || got != 1 {
		t.Errorf("Get(a.b) = %v, %v; want 1, true (Merge must snapshot the source Gofig)", got, ok)
	}

	_, ok = gfg.Get("new")
	if ok {
		t.Error("Get(new) reported found after the source Gofig changed")
	}
}

func TestMergeNilGofigIsNoop(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a"}, 1)
	gfg.Merge(nil)

	want := map[string]any{"a": 1}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestSetCreatesNestedMaps(t *testing.T) {
	var gfg *Gofig = New()
	gfg.Set([]string{"a", "b", "c"}, 1)

	want := map[string]any{"a": map[string]any{"b": map[string]any{"c": 1}}}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestSetOverwritesNonMapOnPath(t *testing.T) {
	var gfg *Gofig = New()
	gfg.Set([]string{"a"}, "scalar")
	gfg.Set([]string{"a", "b"}, 2)

	var got any
	var ok bool
	got, ok = gfg.Get("a", "b")

	if !ok || got != 2 {
		t.Errorf("Get(a.b) = %v, %v; want 2, true", got, ok)
	}
}

func TestSetCopiesMapValue(t *testing.T) {
	var gfg *Gofig = New()
	value := map[string]any{"b": 1}
	gfg.Set([]string{"a"}, value)

	value["b"] = 99

	var got any
	got, _ = gfg.Get("a", "b")

	if got != 1 {
		t.Errorf("Get(a.b) = %v, want 1 (Set must copy map values)", got)
	}
}

func TestSetEmptyKeyReplacesConfig(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"old"}, 1)
	gfg.Set(nil, map[string]any{"new": 2})

	want := map[string]any{"new": 2}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestSetEmptyKeyWithNonMapIsNoop(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"old"}, 1)
	gfg.Set(nil, 42)

	want := map[string]any{"old": 1}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestGetDotPath(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a", "b", "c"}, "deep")

	var got any
	var ok bool
	got, ok = gfg.Get("a", "b", "c")

	if !ok || got != "deep" {
		t.Errorf("Get(a.b.c) = %v, %v; want deep, true", got, ok)
	}

	_, ok = gfg.Get("a", "b", "missing")
	if ok {
		t.Error("Get(a.b.missing) reported found, want not found")
	}

	_, ok = gfg.Get("x", "y")
	if ok {
		t.Error("Get(x.y) reported found, want not found")
	}
}

func TestDeleteRemovesNestedKey(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a", "b"}, 1).Set([]string{"a", "c"}, 2)
	gfg.Delete("a", "b")

	var got any
	var ok bool
	got, ok = gfg.Get("a", "c")

	if !ok || got != 2 {
		t.Errorf("Get(a.c) = %v, %v; want 2, true", got, ok)
	}

	_, ok = gfg.Get("a", "b")
	if ok {
		t.Error("Get(a.b) reported found after Delete, want not found")
	}
}

func TestDeleteMissingPathIsNoop(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a"}, 1)
	gfg.Delete("x", "y", "z") // must not panic

	var got any
	var ok bool
	got, ok = gfg.Get("a")

	if !ok || got != 1 {
		t.Errorf("Get(a) = %v, %v; want 1, true", got, ok)
	}
}

func TestDeleteEmptyKeyClearsConfig(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a"}, 1)
	gfg.Delete()

	var actualValues map[string]any = gfg.Raw()

	if len(actualValues) != 0 {
		t.Errorf("Raw() = %v, want empty map", actualValues)
	}
}

func TestRawReturnsCopy(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a"}, map[string]any{"b": 1})

	var actualValues map[string]any = gfg.Raw()
	actualValues["a"].(map[string]any)["b"] = 99
	actualValues["new"] = true

	var got any
	got, _ = gfg.Get("a", "b")

	if got != 1 {
		t.Errorf("Get(a.b) = %v, want 1 (Raw must return a copy)", got)
	}

	var ok bool
	_, ok = gfg.Get("new")
	if ok {
		t.Error("Get(new) reported found, want not found")
	}
}

func TestFluentChaining(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"a"}, 1).MergeMap(map[string]any{"b": 2}).Delete("a")

	want := map[string]any{"b": 2}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

// stubProvider is a Provider returning a fixed result, for Load tests.
type stubProvider struct {
	keyValues map[string]any
	err       error
}

func (sp stubProvider) Read() (map[string]any, error) { return sp.keyValues, sp.err }

func TestLoadMergesProviderOutput(t *testing.T) {
	var gfg *Gofig = New().Set([]string{"keep"}, true)
	provider := stubProvider{keyValues: map[string]any{"a": 1}}

	var err error
	err = gfg.Load(provider)
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	want := map[string]any{"keep": true, "a": 1}

	var actualValues map[string]any = gfg.Raw()

	if !reflect.DeepEqual(actualValues, want) {
		t.Errorf("Raw() = %v, want %v", actualValues, want)
	}
}

func TestLoadReturnsProviderError(t *testing.T) {
	var wantErr error = errors.New("boom")
	var gfg *Gofig = New()

	var err error
	err = gfg.Load(stubProvider{err: wantErr})

	if !errors.Is(err, wantErr) {
		t.Errorf("Load() error = %v, want %v", err, wantErr)
	}
}

func TestLoadNilProvider(t *testing.T) {
	var gfg *Gofig = New()

	var err error
	err = gfg.Load(nil)

	if err == nil {
		t.Error("Load(nil provider) error = nil, want error")
	}
}
