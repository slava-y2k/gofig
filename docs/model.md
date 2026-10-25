# Gofig

Gofig is shorthand for Go and Config.
Gofig is a library to deal with configuration in Go.

It has similar ideas to the koanf configuration library.

# Terms
Terms are introduced so there is exactly one term for each entity/focus.

- **KeyValue(s)** - refers to the configuration key-value pair or pairs.
- **ConfigStruct** - a struct marked with the struct tag `gofig`.
- **Gofig** - refers to the Gofig struct and its functions.
- **Provider** - an agnostic source of KeyValue(s) for Gofig.


## Model

Gofig is intended to handle multilevel configurations.

The model is visualized in [parser.md](parser.md) - a Mermaid class diagram.


### Gofig

`Gofig` is a key-value pair manager responsible for the following:

- Key-value management: add, remove, and modify keys through its contract.
- Flexible key-value combination, using a fluent interface where appropriate.
- An extension point for working with key-value pairs.


### Provider

The main source of key-value pairs is the `Provider`.
Any source can implement the `Provider` interface, which consists of a single method:
`Provider.Read() (map[string]any, error)`


### Parser

A `Parser` is an auxiliary package that only returns the key-value pairs
(`map[string]any`). Its contract is a single method:

`Parser.Parse(m map[string]any) (map[string]any, error)`



## Code layout

Code is represented as a multi-module structure: each Provider and each Parser
is a separate module. The full rules live in [code-layout.md](code-layout.md).


## Gofig internal storage
The internal storage of `Gofig` is a nested `map[string]any`.





