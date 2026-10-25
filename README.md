# gofig

Go configuration library. Configuration is represented as key-value pairs
(`map[string]any`): providers feed them in, and the `Gofig` container
merges and manages them.

## Development

The repository is a multi-module Go workspace — the core library, one module
per provider (and later per parser), and runnable examples, all in one repo
but independently referenceable:

| Module       | Import path                                |
| ------------ | ------------------------------------------ |
| Core         | `github.com/slava-y2k/gofig`               |
| Key-value provider | `github.com/slava-y2k/gofig/provider/keyvalue` |
| Examples     | `github.com/slava-y2k/gofig/example`       |

Build, test, and run everything from the repository root:

```sh
go build github.com/slava-y2k/gofig/...
go test github.com/slava-y2k/gofig/...
go run ./example/basic
```

See [docs/code-layout.md](docs/code-layout.md) for the layout rules and
[docs/model.md](docs/model.md) for the data model.
