
# Design Idea
Main doc that describes the idea from a technical perspective.

# Domain terms
- **key value(s)** — refers to the configuration key-value pair or pairs.
- **Gonfig** — refers to the `Gonfig` struct and its functions.
- **Provider** — an abstraction of a KeyValue collection for Gonfig.


# Gonfig module
Gonfig provides bound safety. All methods that accept a key-value collection create an internal copy
to protect against changes through the passed reference.
The same approach is applied to outbound methods: every method makes a copy for a reference type
to protect internal state from changes to outgoing values.


## Collection methods
- `NewXXX` — inbound constructor functions that create the `Gonfig` struct and fill it with the initial or provided
  **key value(s)**.
- `MergeXXX` — inbound instance functions (functions with a receiver) that accept the value for the current `Gonfig` struct.

The main difference between `NewXXX` and `MergeXXX` is that `NewXXX` creates and initializes the `Gonfig` struct,
whereas `MergeXXX` applies the provided key values to the existing `Gonfig` struct.

- `Raw()` — returns a copy of the internal map.


## Key methods
All these methods accept a key as an argument; on an empty key, nil will be returned.

- `Get` — returns the key values for the provided key; empty if the key is missing.
- `Set` — sets a value by the provided key.
- `Delete` — deletes the value by the provided key.


## Gonfig internal storage
The internal storage of `Gonfig` is a nested `map[string]any`, guarded by a `sync.RWMutex`.


# Provider

The main source of key-value pairs is the `Provider` abstraction, with the single function:
```go
type Provider interface {
    Read() (map[string]any, error)
}
```
Each provider should be implemented as a separate module.


# Questions to answer later
- Should the provider interface return an error, or will an empty map be enough?



## Code layout
