package gonfig

// Provider is the pluggable source of configuration key-value pairs;
// the provider decides whether Read returns its own storage or a copy.
type Provider interface {
	Read() map[string]any
}
