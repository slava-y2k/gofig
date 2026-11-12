
# Design Idea
Main doc that describes the main idea from the technical perspective.

# Domain terms
- **key value(s)** - refers to the configuration key-value pair or pairs.
- **Gonfig** - refers to the `Gonfig` struct and its functions.
- **Provider** - an abstraction of a KeyValue collection for Gonfig.


# Gonfig module
Gonfig provides bound safety. All methods that accepts a key-value collection creates an internal copy, 
to protect changes by passed reference. 
Same approach is applied to the outbound methods - every method makes a copy for a reference type, 
to protect internal changes for outgoing values.


## Collection methods
- `NewXXX` - inbound constructor functions, creates the Gonfig struct and fills it with the initial or provided
  **key value(s)**. 
- `MergeXXX` - inbound instance functions (functions with receiver) that accepts the value for current `gonfig` struct.

Main difference `NewXXX` vs `MergeXXX` is than new creates and initializes `Gonfig` struct, 
whereas `MergeXXX` applies provided key values to the existing `Gonfig` struct.

- `Raw()` - returns the copy of the internal map 


## Key methods
all these methods accepts key as an argument, on empty key nill will be retuned  

- `Get` - returns key values by provided key, empty .
- `Set` - set value by provided key.
- `Delete` - set value by provided key.


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
- Should provider interface return the error, will empty map will be enough



## Code layout


