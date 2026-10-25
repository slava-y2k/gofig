package gofig

// Provider represents a configuration provider that can read
// configuration data and return it as a map of key-value pairs.
type Provider interface {
	Read() (map[string]any, error)
}
