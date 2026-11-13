package gonfig

import "reflect"

// DeepCopyMap returns a copy that shares no mutable state with the source.
func DeepCopyMap(keyValues map[string]any) map[string]any {
	if keyValues == nil {
		return nil
	}

	var copied map[string]any = make(map[string]any, len(keyValues))
	for key, value := range keyValues {
		copied[key] = deepCopyValue(value)
	}

	return copied
}

// deepCopyValue deep-copies the value behind an any; a nil any stays nil.
func deepCopyValue(value any) any {
	if value == nil {
		return nil
	}

	var reflected reflect.Value = reflect.ValueOf(value)
	var copied reflect.Value = deepCopyReflect(reflected)

	return copied.Interface()
}

// deepCopyReflect recurses into composite kinds; channels, functions, and
// unsafe pointers cannot be copied and come back as they are.
func deepCopyReflect(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface:
		return deepCopyInterfaceReflect(value)
	case reflect.Map:
		return deepCopyMapReflect(value)
	case reflect.Slice:
		return deepCopySliceReflect(value)
	case reflect.Array:
		return deepCopyArrayReflect(value)
	case reflect.Pointer:
		return deepCopyPointerReflect(value)
	case reflect.Struct:
		return deepCopyStructReflect(value)
	default:
		return value
	}
}

// deepCopyInterfaceReflect copies the dynamic value but keeps the interface type.
func deepCopyInterfaceReflect(value reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}

	var copiedValue reflect.Value = deepCopyReflect(value.Elem())
	var result reflect.Value = reflect.New(value.Type()).Elem()
	result.Set(copiedValue)

	return result
}

// deepCopyMapReflect copies values but keeps keys, so a pointer key still matches.
func deepCopyMapReflect(value reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}

	var copied reflect.Value = reflect.MakeMapWithSize(value.Type(), value.Len())
	var iter *reflect.MapIter = value.MapRange()
	for iter.Next() {
		copied.SetMapIndex(iter.Key(), deepCopyReflect(iter.Value()))
	}

	return copied
}

// deepCopySliceReflect allocates a fresh backing array; a typed nil slice stays nil.
func deepCopySliceReflect(value reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}

	var copied reflect.Value = reflect.MakeSlice(value.Type(), value.Len(), value.Len())
	for index := 0; index < value.Len(); index++ {
		copied.Index(index).Set(deepCopyReflect(value.Index(index)))
	}

	return copied
}

// deepCopyArrayReflect copies a value-type array element by element.
func deepCopyArrayReflect(value reflect.Value) reflect.Value {
	var copied reflect.Value = reflect.New(value.Type()).Elem()
	for index := 0; index < value.Len(); index++ {
		copied.Index(index).Set(deepCopyReflect(value.Index(index)))
	}

	return copied
}

// deepCopyPointerReflect points the copy at a fresh value; a typed nil pointer stays nil.
func deepCopyPointerReflect(value reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}

	var copied reflect.Value = reflect.New(value.Type().Elem())
	copied.Elem().Set(deepCopyReflect(value.Elem()))

	return copied
}

// deepCopyStructReflect copies exported fields; unexported fields keep the struct's own copy.
func deepCopyStructReflect(value reflect.Value) reflect.Value {
	var copied reflect.Value = reflect.New(value.Type()).Elem()
	copied.Set(value)
	for index := 0; index < value.NumField(); index++ {
		if !copied.Field(index).CanSet() {
			continue
		}
		copied.Field(index).Set(deepCopyReflect(value.Field(index)))
	}

	return copied
}
