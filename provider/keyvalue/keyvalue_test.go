package keyvalue

import (
	"reflect"
	"testing"

	"github.com/slava-y2k/gofig"
)

// Compile-time check that *KeyValue implements gofig.Provider.
var _ gofig.Provider = (*KeyValue)(nil)

func TestReadReturnsMapAsIs(t *testing.T) {
	testValues := map[string]any{"key": "value"}
	var keyValue *KeyValue = New(testValues)

	var actualValues map[string]any
	var err error
	actualValues, err = keyValue.Read()
	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

	if !reflect.DeepEqual(actualValues, testValues) {
		t.Errorf("Read() = %v, want %v", actualValues, testValues)
	}
}

func TestNewNilMapServedAsEmpty(t *testing.T) {
	var keyValue *KeyValue = New(nil)

	var emptyMap map[string]any
	var err error
	emptyMap, err = keyValue.Read()
	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

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

	var actualValues map[string]any
	var err error
	actualValues, err = keyValue.Read()
	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

	actualValues["b"] = 2
	if testData["b"] != 2 {
		t.Error("Read() must return the served map itself when copy is off")
	}
}

func TestWithDeepCopyReturnsIndependentCopies(t *testing.T) {
	sourceValues := map[string]any{"a": map[string]any{"b": 1}}
	var keyValue *KeyValue = New(sourceValues).WithDeepCopy(true)

	var firstValues map[string]any
	var err error
	firstValues, err = keyValue.Read()

	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

	firstValues["a"].(map[string]any)["b"] = 99
	firstValues["new"] = true

	var secondValues map[string]any
	secondValues, err = keyValue.Read()
	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

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

	var actualValues map[string]any
	var err error
	actualValues, err = keyValue.Read()
	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

	if actualValues["key"] != "value" {
		t.Errorf("Read()[key] = %v, want %q", actualValues["key"], "value")
	}
}
