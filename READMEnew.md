# Gonfig

Gonfig (Go + config) is an extensible micromodule library to deal with configuration in Go.

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
- **ConfigStruct** - a struct marked with the struct tag `gonfig` (planned; not implemented yet).
- **Gonfig** - refers to the `Gonfig` struct and its functions.
- **Provider** - an abstraction of a KeyValue collection for Gonfig.


## The Design

The library is multi-module to reduce coupling and to use only the modules it needs.

- `Gonfig` is the core module for working with KeyValues and defines the `Provider`
  interface for KeyValues.
- `Providers` modules - each data provider is a separate module and uses only the
  dependencies it needs.

See the model diagram in [docs/parser.md](docs/parser.md).


# Module details

## Gonfig (core) module

Supports add, remove, modify, and merge operations.
Defines the `Provider` interface.

### API

```go
func New() *Gonfig
func (gfg *Gonfig) Load(provider Provider) error
func (gfg *Gonfig) Get(parts ...string) (any, bool)
func (gfg *Gonfig) Set(parts []string, value any) *Gonfig
func (gfg *Gonfig) Delete(parts ...string) *Gonfig
func (gfg *Gonfig) Merge(other *Gonfig) *Gonfig
func (gfg *Gonfig) MergeMap(keyValues map[string]any) *Gonfig
func (gfg *Gonfig) Raw() map[string]any
func DeepCopyMap(m map[string]any) map[string]any
```

| Method | Behavior |
| ------ | -------- |
| `New` | Creates a `Gonfig` with an empty configuration. |
| `Load` | Reads a `Provider` and deep-merges its output into the configuration. Returns `ErrNilProvider` for a nil provider; otherwise returns the provider's error. |
| `Get` | Returns the value at the given path parts and whether it was found. Without parts, returns the whole configuration. |
| `Set` | Stores `value` at the given path parts, creating missing intermediate maps. Without parts, replaces the whole configuration when `value` is a `map[string]any` (a non-map value is ignored). |
| `Delete` | Removes the value at the given path parts; missing paths are ignored. Without parts, clears the whole configuration. |
| `Merge` | Deep-merges another `Gonfig` into this one. A nil `Gonfig` is ignored. |
| `MergeMap` | Deep-merges the given key-value pairs into the configuration. |
| `Raw` | Returns a deep copy of the whole configuration map, so callers cannot mutate the internal state. |
| `DeepCopyMap` | Returns a deep copy of a map: nested maps are copied recursively, all other values are shared. |

Nested maps are merged recursively on `Load`, `Merge`, and `MergeMap`; every other
value overwrites the existing one. `Get`, `Set`, and `Delete` address nested maps
through path parts: `Get` and `Delete` take the parts as a variadic list, `Set` as
a slice.

`Set` deep-copies map values, and `Load`/`Merge`/`MergeMap` snapshot their sources,
so the internal configuration never aliases caller-owned maps.

Methods that cannot fail (`Merge`, `MergeMap`, `Set`, `Delete`) return the `Gonfig`
itself, enabling a fluent interface. All public methods of `Gonfig` are safe for
concurrent use (guarded by a `sync.RWMutex`).

The zero value of `Gonfig` is usable.

## Provider interface

The main source of key-value pairs is the `Provider` abstraction,
with the single function:

```go
type Provider interface {
    Read() (map[string]any, error)
}
```

Any source can implement `Provider` to feed KeyValues into `Gonfig`.

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
- `KeyValue` implements `gonfig.Provider` (asserted at compile time in its tests).


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


## Gonfig internal storage

The internal storage of `Gonfig` is a nested `map[string]any`, guarded by a
`sync.RWMutex`.


## Quick start

```go
package main

import (
  "fmt"
  "log"

  "github.com/slava-y2k/gonfig"
  "github.com/slava-y2k/gonfig/provider/keyvalue"
)

func main() {
  src := map[string]any{
    "app": map[string]any{"name": "gonfig", "debug": true},
    "db":  map[string]any{"host": "localhost", "port": 5432},
  }
  g := gonfig.New()
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
| Core              | `github.com/slava-y2k/gonfig`                       |
| Key-value provider | `github.com/slava-y2k/gonfig/provider/keyvalue`    |
| Examples          | `github.com/slava-y2k/gonfig/example`               |

Build, test, and run everything from the repository root:

```sh
go build github.com/slava-y2k/gonfig/...
go test github.com/slava-y2k/gonfig/...
go vet github.com/slava-y2k/gonfig/...
go run ./example/basic
```


Gonfig is shorthand for Go and Config.
The library design is similar to the koanf library, with a cleaner design.