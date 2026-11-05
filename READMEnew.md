# Gofig

Gofig (Go + config) is an extensible micromodule library to deal with configuration in Go.

# The Design  

## Code

```go

config example

``` 

## Diagram



## Installation






## Terms

Terms are introduced so there is exactly one term for each entity/focus.

- **KeyValue(s)** - refers to the configuration key-value pair or pairs.
- **ConfigStruct** - a struct marked with the struct tag `gofig` (planned; not implemented yet).
- **Gofig** - refers to the `Gofig` struct and its functions.
- **Provider** - an abstraction of a KeyValue collection for Gofig.


## The Design

The library is multi-module to reduce coupling and to use only the modules it needs.

- `Gofig` is the core module for working with KeyValues and defines the `Provider`
  interface for KeyValues.
- `Providers` modules - each data provider is a separate module and uses only the
  dependencies it needs.

See the model diagram in [docs/parser.md](docs/parser.md).


# Module details

## Gofig (core) module

Supports add, remove, modify, and merge operations.
Defines the `Provider` interface.

### API

```go
func New() *Gofig
func (gfg *Gofig) Load(provider Provider) error
func (gfg *Gofig) Get(parts ...string) (any, bool)
func (gfg *Gofig) Set(parts []string, value any) *Gofig
func (gfg *Gofig) Delete(parts ...string) *Gofig
func (gfg *Gofig) Merge(other *Gofig) *Gofig
func (gfg *Gofig) MergeMap(keyValues map[string]any) *Gofig
func (gfg *Gofig) Raw() map[string]any
func DeepCopyMap(m map[string]any) map[string]any
```

| Method | Behavior |
| ------ | -------- |
| `New` | Creates a `Gofig` with an empty configuration. |
| `Load` | Reads a `Provider` and deep-merges its output into the configuration. Returns `ErrNilProvider` for a nil provider; otherwise returns the provider's error. |
| `Get` | Returns the value at the given path parts and whether it was found. Without parts, returns the whole configuration. |
| `Set` | Stores `value` at the given path parts, creating missing intermediate maps. Without parts, replaces the whole configuration when `value` is a `map[string]any` (a non-map value is ignored). |
| `Delete` | Removes the value at the given path parts; missing paths are ignored. Without parts, clears the whole configuration. |
| `Merge` | Deep-merges another `Gofig` into this one. A nil `Gofig` is ignored. |
| `MergeMap` | Deep-merges the given key-value pairs into the configuration. |
| `Raw` | Returns a deep copy of the whole configuration map, so callers cannot mutate the internal state. |
| `DeepCopyMap` | Returns a deep copy of a map: nested maps are copied recursively, all other values are shared. |

Nested maps are merged recursively on `Load`, `Merge`, and `MergeMap`; every other
value overwrites the existing one. `Get`, `Set`, and `Delete` address nested maps
through path parts: `Get` and `Delete` take the parts as a variadic list, `Set` as
a slice.

`Set` deep-copies map values, and `Load`/`Merge`/`MergeMap` snapshot their sources,
so the internal configuration never aliases caller-owned maps.

Methods that cannot fail (`Merge`, `MergeMap`, `Set`, `Delete`) return the `Gofig`
itself, enabling a fluent interface. All public methods of `Gofig` are safe for
concurrent use (guarded by a `sync.RWMutex`).

The zero value of `Gofig` is usable.

## Provider interface

The main source of key-value pairs is the `Provider` abstraction,
with the single function:

```go
type Provider interface {
    Read() (map[string]any, error)
}
```

Any source can implement `Provider` to feed KeyValues into `Gofig`.

### Key-value provider

The `keyvalue` provider serves an in-memory map of key-value pairs:

```go
func New(keyValues map[string]any) *KeyValue
func (kv *KeyValue) WithDeepCopy(isEnabled bool) *KeyValue
func (kv *KeyValue) Read() (map[string]any, error)
```

- `New` serves the given map; a nil map is served as an empty map.
- `WithDeepCopy` enables/disables deep copying. When enabled, every `Read` returns
  a fresh deep copy, so the served map is never mutated by the caller. When
  disabled, `Read` returns the served map itself.
- `KeyValue` implements `gofig.Provider` (asserted at compile time in its tests).


## Parser

A `Parser` is an auxiliary package that only returns the key-value pairs
(`map[string]any`). Its contract is a single method:

```go
Parse(m map[string]any) (map[string]any, error)
```

No parser module exists yet; the first parser modules are planned under
`parser/` (JSON, YAML, TOML).


## Code layout

Code is represented as a multi-module structure: each Provider and each Parser
is a separate module. The full rules live in [docs/code-layout.md](docs/code-layout.md).


## Gofig internal storage

The internal storage of `Gofig` is a nested `map[string]any`, guarded by a
`sync.RWMutex`.


## Quick start

```go
package main

import (
  "fmt"
  "log"

  "github.com/slava-y2k/gofig"
  "github.com/slava-y2k/gofig/provider/keyvalue"
)

func main() {
  src := map[string]any{
    "app": map[string]any{"name": "gofig", "debug": true},
    "db":  map[string]any{"host": "localhost", "port": 5432},
  }
  g := gofig.New()
  // Load from a key-value provider. WithDeepCopy makes Read return a copy
  // on every call, so the served map is never mutated.
  if err := g.MergeProvider(keyvalue.New(src).WithDeepCopy(true)); err != nil {
    log.Fatal(err)
  }

  // MergeMap overrides existing values; Set and Delete then modify keys fluently.
  g.MergeMap(map[string]any{"db": map[string]any{"port": 5433}}).
    Set([]string{"app", "mode"}, "dev").
    Delete("app", "debug")

  host, _ := g.Get("db", "host")
  port, _ := g.Get("db", "port")
  mode, _ := g.Get("app", "mode")

  fmt.Println("host:", host)
  fmt.Println("port:", port)
  fmt.Println("mode:", mode)
  fmt.Println("raw: ", g.Raw())
}
```

See [example/basic](example/basic/main.go) for the runnable version:
`go run ./example/basic`.


## Development

The repository is a multi-module Go workspace - the core library, one module
per provider (and later per parser), and runnable examples, independently
referenceable:

| Module            | Import path                                        |
| ----------------- | -------------------------------------------------- |
| Core              | `github.com/slava-y2k/gofig`                       |
| Key-value provider | `github.com/slava-y2k/gofig/provider/keyvalue`    |
| Examples          | `github.com/slava-y2k/gofig/example`               |

Build, test, and run everything from the repository root:

```sh
go build github.com/slava-y2k/gofig/...
go test github.com/slava-y2k/gofig/...
go vet github.com/slava-y2k/gofig/...
go run ./example/basic
```


Gofig is shorthand for Go and Config.
The library design is similar to the koanf library, with a cleaner design.