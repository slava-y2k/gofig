# Gonfig

Gonfig (Go + config) is a micromodular, extensible configuration library for Go.


# Design Brief

The library uses a micromodular approach to reduce coupling and to include only the modules it needs:
- `Gonfig` is the core module for working with key-value pairs and defines the `Provider` interface for them.
- `Provider` modules — each provider module is a separate package and uses only the dependencies it needs.
- The copy-contract keeps the container isolated: collections passed in or read back through `Raw` are copied,
  while `Get` returns live values that the caller must not modify.

For detailed design, please refer to [docs/design.md](docs/design.md).

# Quick start

```go
package main

import (
	"fmt"
	"log"

	"github.com/slava-y2k/gonfig"
	"github.com/slava-y2k/gonfig/provider/keyvalue"
)

func main() {
	src := map[string]any{"app": map[string]any{"name": "gonfig"}}

	var g *gonfig.Gonfig
	var err error
	g, err = gonfig.NewWithProvider(keyvalue.New(src))
	if err != nil {
		log.Fatal(err)
	}

	name, ok := g.Get("app", "name")
	if !ok {
		log.Fatal("app.name is not set")
	}

	fmt.Println(name)
}
```

# Examples

A runnable example is in [example/basic/main.go](example/basic/main.go).

# See more

The full documentation index is in [docs/_index.md](docs/_index.md).


