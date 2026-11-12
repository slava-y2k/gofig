
# Design Idea
Main doc that describes the idea from a technical perspective.

# Domain terms
- **key-value pair(s)** — refers to a configuration key together with its value.
- **Gonfig** — refers to the `Gonfig` struct and its functions.
- **Provider** — an abstraction of a collection of key-value pairs for Gonfig.


# Gonfig module
Gonfig provides bound safety. All methods that accept a collection of key-value pairs create an internal copy
to protect against changes through the passed reference.
The same approach is applied to outbound methods: every method makes a copy for reference types
to protect internal state from changes to outgoing values.


## Collection methods
- `NewWithXXX` — inbound constructor functions that create the `Gonfig` struct and fill it with the initial or provided
  **key-value pair(s)**.
- `MergeXXX` — inbound instance functions (functions with a receiver) that accept the value for the current `Gonfig` struct.

The main difference between `NewWithXXX` and `MergeXXX` is that `NewWithXXX` creates and initializes the `Gonfig` struct,
whereas `MergeXXX` applies the provided key-value pairs to the existing `Gonfig` struct.

- `Raw()` — returns a copy of the internal map.


## Key methods
All these methods accept a key as an argument; for an empty key, nil will be returned.

- `Get` — returns the value for the provided key; nil if the key is missing.
- `Set` — sets a value by the provided key, returns the set key, nil if the key is missing.
- `Delete` — deletes the value by the provided key, nil if the key is missing.


## Gonfig internal storage
The internal storage of `Gonfig` is a nested `map[string]any`, guarded by a `sync.RWMutex`.


# Provider

The main source of key-value pairs is the `Provider` abstraction, with a single function:
```go
type Provider interface {
Read() map[string]any
}
```
Each provider should be implemented as a separate module.


# Questions to answer later
- Should the provider interface return an error, or will an empty map be enough?
