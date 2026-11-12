# Gonfig

Gonfig (Go + config) is a micromodular, extensible configuration library for Go.


# Design Brief

The library uses a micromodular approach to reduce coupling and to include only the modules it needs:
- `Gonfig` is the core module for working with key-value pairs and defines the `Provider` interface for them.
- `Provider` modules — each provider module is a separate package and uses only the dependencies it needs.

For detailed design, please refer to ...

# Quick start

```go
package main

import (
  "fmt"
  "log"

  "github.com/slava-y2k/gonfig"
)

// TODO
```

# Examples

# See more


