package keyvalue

import (
	"reflect"
	"testing"

	"github.com/slava-y2k/gonfig"
)

// Compile-time check that *KeyValue implements gonfig.Provider.
var _ gonfig.Provider = (*KeyValue)(nil)

func TestReadReturnsMapAsIs(t *testing.T) {
	testValues := map[string]any{"key": "value"}
	var keyValue *KeyValue = New(testValues)

	var actualValues map[string]any = keyValue.Read()

	if !reflect.DeepEqual(actualValues, testValues) {
		t.Errorf("Read() = %v, want %v", actualValues, testValues)
	}
}

func TestNewNilMapServedAsEmpty(t *testing.T) {
	var keyValue *KeyValue = New(nil)

	var emptyMap map[string]any = keyValue.Read()

	if emptyMap == nil {
		t.Fatal("Read() = nil, want an empty map (a nil map is normalized)")
	}
	if len(emptyMap) != 0 {
		t.Errorf("Read() = %v, want an empty map", emptyMap)
	}
}

func TestReadWithoutCopyAliasesSource(t *testing.T) {
	testData := map[string]any{"a": 1}
	var keyValue *KeyValue = New(testData)

	var actualValues map[string]any = keyValue.Read()

	actualValues["b"] = 2
	if testData["b"] != 2 {
		t.Error("Read() must return the served map itself when copy is off")
	}
}

func TestWithDeepCopyReturnsIndependentCopies(t *testing.T) {
	sourceValues := map[string]any{"a": map[string]any{"b": 1}}
	var keyValue *KeyValue = New(sourceValues).WithDeepCopy(true)

	var firstValues map[string]any = keyValue.Read()

	firstValues["a"].(map[string]any)["b"] = 99
	firstValues["new"] = true

	var secondValues map[string]any = keyValue.Read()

	if got := secondValues["a"].(map[string]any)["b"]; got != 1 {
		t.Errorf("second Read() = %v, want 1 (every Read must return a fresh copy)", got)
	}
	if _, ok := secondValues["new"]; ok {
		t.Error("second Read() contains a key set on the first copy")
	}
	if got := sourceValues["a"].(map[string]any)["b"]; got != 1 {
		t.Errorf("source map changed to %v, want untouched", got)
	}
}

func TestCreateKeyValueProvider(t *testing.T) {
	var keyValue *KeyValue = New(map[string]any{"key": "value"})

	var actualValues map[string]any = keyValue.Read()

	if actualValues["key"] != "value" {
		t.Errorf("Read()[key] = %v, want %q", actualValues["key"], "value")
	}
}

// TestWithDeepCopyToggleOffRestoresAliasing pins both toggle directions: after
// switching copying off, Read serves the stored map again.
func TestWithDeepCopyToggleOffRestoresAliasing(t *testing.T) {
	sourceValues := map[string]any{"a": 1}
	var keyValue *KeyValue = New(sourceValues).WithDeepCopy(true)

	var copied map[string]any = keyValue.Read()

	copied["b"] = 2
	if _, ok := sourceValues["b"]; ok {
		t.Fatal("copying on must not reach the source map")
	}

	keyValue.WithDeepCopy(false)

	var served map[string]any = keyValue.Read()

	served["c"] = 3
	if sourceValues["c"] != 3 {
		t.Error("Read() must alias the served map again after copying is disabled")
	}
}

// TestWithDeepCopyReturnsReceiver pins the fluent return value.
func TestWithDeepCopyReturnsReceiver(t *testing.T) {
	var keyValue *KeyValue = New(map[string]any{})

	var got *KeyValue = keyValue.WithDeepCopy(true)

	if got != keyValue {
		t.Errorf("WithDeepCopy() = %p, want receiver %p", got, keyValue)
	}
}
