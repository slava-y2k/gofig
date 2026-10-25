// Package maps provides helpers for working with key-value maps.
package maps

// DeepCopyMap returns a deep copy of m: nested maps are copied
// recursively, all other values are shared.
func DeepCopyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}

	out := make(map[string]any, len(m))
	for key, value := range m {
		if nestedMap, ok := value.(map[string]any); ok {
			out[key] = DeepCopyMap(nestedMap)
			continue
		}
		out[key] = value
	}

	return out
}
