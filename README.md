# Gonfig

Gonfig (Go + config) is a micromodule extensible config library to deal with configuration in Go.


# Design Brief

Library is implemented with micro modular approach to reduce coupling and to use only the modules it needs:
- `Gonfig` is the core module for working with KeyValues and defines the `Provider` interface for KeyValues.
- `Providers` modules - each data provider is a separate module and uses only the dependencies it needs.

for detail design please refer to ...

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



