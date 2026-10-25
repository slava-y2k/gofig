# Code Layout

gofig is a Go configuration library inspired by [koanf](https://github.com/knadh/koanf).
Configuration is represented as `map[string]any` key-value pairs. Sources feed
data in through **Providers**, and the **Gofig** container merges and manages
the result. **Parsers** (planned) will normalize key-value maps in standalone
modules.

## Modules

The repository is multi-module: each provider and each parser is
a separate module, so every piece can be referenced and versioned
independently while staying in one repository. `go.work` ties all modules
together for local development.

| Module           | go.mod                                                                  | Package         |
| ---------------- | ----------------------------------------------------------------------- | --------------- |
| Core             | `go.mod` — `github.com/slava-y2k/gofig`                                 | `gofig`         |
| Key-value provider | `provider/keyvalue/go.mod` — `github.com/slava-y2k/gofig/provider/keyvalue` | `keyvalue`   |
| Examples         | `example/go.mod` — `github.com/slava-y2k/gofig/example`                 | `main` per example |

Rules:

- Core lives at the repository root as `package gofig`: its module path,
  directory, and package name all agree.
- **Every provider is a separate module** — directory `provider/<name>/` with
  `go.mod` declaring `module github.com/slava-y2k/gofig/provider/<name>`.
- **Every parser is a separate module** — directory `parser/<name>/` with
  `go.mod` declaring `module github.com/slava-y2k/gofig/parser/<name>`.
- All module paths nest under `github.com/slava-y2k/gofig` — that is what
  makes the workspace-wide build pattern below cover every module.
- Providers depend on the core module to assert the real `Provider`
  interface (`var _ gofig.Provider`) — their `go.mod` requires the core with a
  `replace` to the local directory (see `provider/keyvalue/go.mod`).
- Importing one provider or parser never pulls in the others.
- Consumers reference each module independently:
  `go get github.com/slava-y2k/gofig` or
  `go get github.com/slava-y2k/gofig/provider/keyvalue`.
- **Versioning**: nested modules are tagged with a path prefix — core tags as
  `v1.2.3`, the key-value provider tags as `provider/keyvalue/v0.1.0`.
- Examples are not published; the module exists for the workspace build.

## Workspace

`go.work` at the repository root lists every module in its `use` block, so a
developer builds, tests, and runs the whole project from the root:

```sh
go build github.com/slava-y2k/gofig/...   # compile every module
go test github.com/slava-y2k/gofig/...    # test every module
go vet github.com/slava-y2k/gofig/...     # lint every module
go run ./example/basic                   # run an example
```

Why the module-path pattern and not `./...`:

- Directory-rooted patterns (`./...`, `./provider/...`) are scoped to a single
  module: walking a tree **skips nested modules rooted below the pattern**.
  `./...` from the root matches only the core; `./provider/...` matches
  nothing. Patterns rooted inside a module directory (`./example/...`) do
  work, but then every new module must be listed by hand.
- The module-path pattern `github.com/slava-y2k/gofig/...` matches packages in
  every workspace module whose import path has that prefix — which is every
  module, by the layout rules — so new providers and parsers are covered
  automatically.
- Avoid `go test all`: in workspace mode it includes the standard library's
  own test suites.

Operational notes:

- `go.work` is committed so every developer and CI gets the same build-all view.
- **Register every new module in `go.work`** (`go work use ./provider/<name>`);
  otherwise it is invisible to the workspace.
- Shared map helpers live in `gofig_util.go` of the core module
  (`gofig.DeepCopyMap`); the keyvalue provider imports the core instead of
  duplicating them.
- Inside the workspace the `replace` directives are dormant (workspace
  resolution wins); they keep each module buildable standalone
  (`GOWORK=off go test ./...` inside the module directory).
- Run examples with `go run`. Building inside the `example/` module tries to
  write a binary named after the package (`basic`), which collides with its
  source directory.

## Directory tree

```text
gofig/                        # repository root
├── AGENTS.md                 # agent instructions
├── README.md
├── go.mod                    # core module github.com/slava-y2k/gofig
├── gofig.go                  # package gofig — the Gofig container
├── gofig_test.go             # tests for the Gofig container
├── provider.go               # package gofig — Provider interface
├── gofig_util.go             # package gofig — DeepCopyMap and other shared helpers
├── go.work                   # workspace: use (., ./provider/keyvalue, ./example)
├── parser/                   # Parser modules (one module per parser, none yet)
├── provider/                 # Provider modules (one module per provider)
│   └── keyvalue/
│       ├── go.mod            # module github.com/slava-y2k/gofig/provider/keyvalue
│       ├── keyvalue.go       # package keyvalue — in-memory key-value Provider
│       └── keyvalue_test.go
├── example/                  # module github.com/slava-y2k/gofig/example (not published)
│   ├── go.mod
│   └── basic/
│       └── main.go           # runnable example: core + key-value provider
├── docs/                     # design documentation
│   ├── code-layout.md        # this file
│   ├── model.md              # data model notes
│   └── parser.md             # Mermaid class diagram of the configuration model
└── opencode.jsonc            # OpenCode config (skill permissions, commands)
```

The empty `parser/` directory is a local placeholder — git does not track empty
directories; the first parser module will bring it under version control.

## Core abstractions

```text
Provider ──Read()──► map[string]any ──Load()──► Gofig      (merge & manage)
```

| Type       | File                            | Kind      | Purpose |
| ---------- | ------------------------------- | --------- | ------- |
| `Gofig`    | `gofig.go`                      | struct    | Central key-value pair manager. `Load` reads a `Provider` and folds in its output; `Merge` takes another `Gofig`, `MergeMap` takes a map, and `Set`/`Delete` mutate keys (they return `*Gofig`, enabling a fluent interface); `Get` and `Raw` read. `Get` reads by path parts (`g.Get("db", "host")`), `Set` writes by path parts (`g.Set([]string{"db"}, v)`); `Delete` takes path parts as a variadic list (no parts clear the whole configuration). |
| `Provider` | `provider.go`                   | interface | A configuration source. `Read()` returns the configuration as key-value pairs. |
| `KeyValue` | `provider/keyvalue/keyvalue.go` | struct    | Provider serving an in-memory map of key-value pairs. Built with `keyvalue.New(m)`; a nil map is served as an empty map; `WithDeepCopy(true)` makes `Read` return a deep copy on every call. |

The `Parser` contract (`Parse(m map[string]any) (map[string]any, error)`) is
defined in [model.md](model.md); no implementation exists yet.

Conventions:

- Interfaces are small and consumer-shaped (one method each).
- All configuration data crosses boundaries as `(map[string]any, error)`.
- Compile-time interface checks: the keyvalue provider asserts the real
  `Provider` interface (`var _ gofig.Provider`) in its tests.
- Public methods of `Gofig` are safe for concurrent use (guarded by `sync.RWMutex`).

## Planned

- Parser modules under `parser/`, each its own module (JSON, YAML, TOML need a
  byte-source path in the model first).
- More provider modules under `provider/` (files, env vars, remote sources).
- More shared helpers in `gofig_util.go` (map merge, flatten) as the need arises.
