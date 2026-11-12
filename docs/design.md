
# Design Idea
Main doc that describes the idea from a technical perspective.

# Domain terms
- **key-value pair(s)** — refers to a configuration key together with its value.
- **Gonfig** — the configuration container type; its methods are listed under Collection methods and Access methods by key, and its constructors are `NewWithXXX`.
- **Provider** — an abstraction of a collection of key-value pairs for Gonfig.
- **copy-contract** — the rule for every value that crosses the `Gonfig` boundary: inbound collections are
  deep-copied, `Raw()` returns a deep copy, and `Get` returns live values; a returned map is a mutable entity
  shared with the container, so the caller must treat it as read-only.

## Diagram

```mermaid
---
title: Gonfig library design
---
classDiagram
    direction LR

    namespace gonfig_pkg ["Gonfig package"] {
        class Gonfig["Gonfig"] {
            -kvMap map[string]any
            +NewWithXXX()$ error, Gonfig
            +MergeXXX(other Gonfig) Gonfig
            +Set(parts []string, value any) Gonfig
            +Delete(parts ...string) Gonfig
            +Get(parts ...string) any
            +Raw() map[string]any
        }

        class IProvider["Provider"] {
            <<interface>>
            +read() map[string, any]
        }
    }

    namespace Provider ["Provider package"] {
        class ConcreteProvider
    }

    Gonfig --> IProvider
    IProvider <|-- ConcreteProvider


```

# Gonfig module
Gonfig provides isolation through its **copy-contract**. All methods that accept a collection of key-value
pairs create an internal copy to protect against changes through the passed reference.
`Raw()` returns a copy for the same reason, whereas `Get` is the documented exception: it returns live
values, so a returned map is a mutable entity shared with the container and must not be modified.


## Collection methods
- `NewWithXXX` — inbound constructor functions that create the `Gonfig` struct and fill it with the initial or provided
  **key-value pair(s)**.
- `MergeXXX` — inbound instance functions (functions with a receiver) that accept the value for the current `Gonfig` struct.

The `XXX` suffix only emulates overloading; Go has no method overloading.

The main difference between `NewWithXXX` and `MergeXXX` is that `NewWithXXX` creates and initializes the `Gonfig` struct,
whereas `MergeXXX` applies the provided key-value pairs to the existing `Gonfig` struct.

- `Raw()` — returns a deep copy of the internal map, as the copy-contract requires.


## Access methods by key

If key is empty, nil or not specified, nill will be returned,
otherwise the specified value.

`Get` is the live-value exception of the copy-contract: a returned map stays shared with the container,
so the caller must treat it as read-only.

- `Get` — live value by the provided key; a returned map must not be modified
- `Set` — value which was set by the provided key
- `Delete` — value which were removed by the provided key 


## Gonfig internal storage
The internal storage of `Gonfig` is a nested `map[string]any`, guarded by a `sync.RWMutex`.


# Provider

The main source of key-value pairs is the `Provider` abstraction, with a single function:
```go
type Provider interface {
Read() map[string]any
}
```
Each provider module is a separate package.

# Questions to answer later
- Should the provider interface return an error, or will an empty map be enough?
