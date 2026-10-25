# Gofig

Gofig is a simple library to deal with configuration in Go.
Gofig is shorthand for Go and Config.
It is a simplified and reworked library similar to koanf.

## Idea

The config library was implemented with the following principles:
- A flexible connection between the raw key-value pairs and config structs
- Composable key-value configuration with overriding and a defined order of application
- An extensible set of configuration sources


## Model

The model is visualized in [parser.puml](parser.puml) - a PlantUML file.


### Gofig

`Gofig` is a key-value pair manager responsible for the following:

- Key-value management: add, remove, and modify keys through its contract.
- Flexible key-value combination, using a fluent interface where appropriate.
- An extension point for working with key-value pairs.


### Provider package

The main source of key-value pairs is the `Provider`.
Any source can implement the `Provider` interface, which in general consists of a single method:
`Provider.Read() (map[string]any, error)`


### Parser package

A `Parser` is an auxiliary package which only returns the key-value pair.
Its contract is a single method:

`Parser.Parse(m map[string]any) (map[string]any, error)`



## Code layout

Code is represented as a multi-module structure: each Provider and each Parser
is a separate module. The full rules live in [code-layout.md](code-layout.md).


## Gofig internal storage
config is a nested map 





