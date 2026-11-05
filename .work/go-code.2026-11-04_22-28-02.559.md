# Go Code Documentation Run: `gofig.go`

- Run timestamp: `2026-11-04_22:28:02.559`
- Skill: `@go-code` (documentation part), `llm-debug:true`
- Inputs used:
  - `.opencode/skills/go-code/SKILL.md` (including the new *Special cases* section on emulated overloads)
  - `AGENTS.md` (updated: `docs/code` repeats the code directory structure and is additional generation context)
  - `docs/code/gofig.md` — additional context: `New<WithXXX>` are overloaded methods
  - `docs/model.md`, `docs/code-layout.md` — model information
  - `gofig.go`, `provider.go`, `gofig_util.go` — source under documentation

---

## 1. Documentation Context

The documentation context is formed from three sources per the skill — the method instructions themselves, internally invoked methods (stack depth up to 3), and model information — plus the additional `docs/code` context mandated by `AGENTS.md`.

### 1.1 Method instructions itself

- **`package gofig`** — root package of the core module; stdlib-only (`errors`, `sync`).
- **`ErrNilProvider`** — sentinel error (`errors.New("gofig: nil provider")`); guard value returned when a constructor receives a nil `Provider`.
- **`Gofig`** — struct `{ rwMutex sync.RWMutex; kvMap map[string]any }`; encapsulates the configuration state and its concurrency guard.
- **`New`** — allocates a `Gofig` with an initialized empty map.
- **`NewWithMap`** — builds a `Gofig` from a map: `New()`, `DeepCopyMap(keyValues)` into `mapCopy`, then `mergeMaps` into `kvMap`; always returns `nil` error.
- **`NewWithProvider`** — nil-checks the provider (`ErrNilProvider`), calls `provider.Read()`, propagates the read error, delegates to `NewWithMap`.
- **`NewWithGofig`** — nil source yields an empty `Gofig`; otherwise snapshots the source via `other.Raw()` and delegates to `NewWithMap`.
- **`Get`** — RLock; nil map → `(nil, false)`; empty parts → whole map; otherwise `getPath`.
- **`Set`** — Lock; `ensureMapInitialized`; empty parts + map value → replaces root with a deep copy; otherwise `setPath`; returns receiver.
- **`Delete`** — Lock; empty parts → fresh empty map; otherwise `deletePath` on a non-nil map; returns receiver.
- **`Merge`** — snapshots `other` under its RLock (`DeepCopyMap`), then Lock on receiver, `ensureMapInitialized`, `mergeMaps`; returns receiver.
- **`MergeMap`** — Lock; `ensureMapInitialized`; `mergeMaps` with `DeepCopyMap(keyValues)`; returns receiver.
- **`Raw`** — RLock; `DeepCopyMap(kvMap)`; guarantees a non-nil result.
- **`ensureMapInitialized`** — lazy allocation of `kvMap`; returns receiver.
- **`setPath`** — walks `parts`, allocates missing intermediate `map[string]any` nodes, deep-copies a map value at the leaf.
- **`deletePath`** — walks `parts`, stops silently at a missing/non-map intermediate, deletes the leaf key.
- **`getPath`** — walks `parts` through nested maps, returns the leaf value and found flag.
- **`mergeMaps`** — recursive deep merge: nested map pairs merge recursively, all other values overwrite `dst`.

### 1.2 Internally invoked methods (stack depth ≤ 3)

1. `New` → `make(map[string]any)`.
2. `NewWithMap` → `New`, `DeepCopyMap`, `mergeMaps`; depth 2: `mergeMaps` → `mergeMaps` (recursive).
3. `NewWithProvider` → `provider.Read()`, `NewWithMap`; depth 2: `New`/`DeepCopyMap`/`mergeMaps`; depth 3: recursive `mergeMaps`.
4. `NewWithGofig` → `other.Raw()`, `NewWithMap`; depth 2: `Raw` → `DeepCopyMap`; `NewWithMap` → `DeepCopyMap`/`mergeMaps`; depth 3: recursive `mergeMaps`.
5. `Get` → `getPath` (RLock/RUnlock at depth 0).
6. `Set` → `ensureMapInitialized`, `DeepCopyMap` (root replace), `setPath`; depth 2: `setPath` → `DeepCopyMap` (leaf map value).
7. `Delete` → `deletePath`.
8. `Merge` → `other.rwMutex.RLock` + `DeepCopyMap` (source snapshot), `ensureMapInitialized`, `mergeMaps`; depth 2: recursive `mergeMaps`.
9. `MergeMap` → `ensureMapInitialized`, `DeepCopyMap`, `mergeMaps`; depth 2: recursive `mergeMaps`.
10. `Raw` → `DeepCopyMap`.
11. `ensureMapInitialized` → `make(map[string]any)`.
12. `setPath` → `make(map[string]any)`, `DeepCopyMap`.
13. `deletePath` → `delete`.
14. `getPath` → map lookup.
15. `mergeMaps` → recursive `mergeMaps`.

### 1.3 Model information

- Configuration is a hierarchical `map[string]any` key-value collection (`docs/model.md`).
- `Gofig` is the central key-value manager: construct from sources (`NewWithXXX`), merge, mutate, read.
- **Boundary safety** (`docs/model.md` "Gofig implementation details"): every method accepting a key-value collection makes an internal copy; every outbound method returns a copy. This is the primary abstraction behind `DeepCopyMap` usage at every boundary.
- Public `Gofig` methods are safe for concurrent use (`docs/code-layout.md` Conventions) — realized by the internal `sync.RWMutex`.
- Mutators return `*Gofig` for a fluent interface (`docs/code-layout.md`).
- `Provider` is a one-method source interface: `Read() (map[string]any, error)`.

### 1.4 Additional context: `docs/code` (per updated `AGENTS.md`)

`docs/code/gofig.md` mirrors the code file and states: **`New<WithXXX>` are overloaded methods.** Go has no method overloading, so the family `New`, `NewWithMap`, `NewWithProvider`, `NewWithGofig` emulates one overloaded constructor with argument-shape suffixes.

### 1.5 Special case applied (skill: emulated overloads)

Per the skill's *Special cases*: for overloaded methods use a **single description which varies only with the provided args**. Therefore the constructor family now shares one description template — "returns a new Gofig loaded with a deep copy of ⟨source⟩" — where only the source argument differs:

- `New` — no source → "returns a new empty Gofig".
- `NewWithMap` — source: the given key-value pairs.
- `NewWithProvider` — source: the key-value pairs read from the given Provider.
- `NewWithGofig` — source: the given Gofig's configuration.

Variant-specific details (nil handling, `ErrNilProvider`) stay documented on `ErrNilProvider` and in the structured context rather than duplicated per constructor.

---

## 2. Structured Context

Structure follows the code: package → struct → methods/constants. Statements applicable to a whole entity are specified at that entity's level.

```text
Package: gofig
│  (Level-wide: thread-safe hierarchical key-value configuration container,
│   stdlib-only core. Boundary safety is the package-level invariant: every
│   inbound collection is deep-copied, every outbound map is a deep copy.
│   Overloaded constructor family New<WithXXX> shares a single description
│   varying only with the provided source argument.)
│
├── Sentinel error
│   └── ErrNilProvider (returned by NewWithProvider for a nil Provider source)
│
├── Struct: Gofig
│   │  (Entity-level: central key-value manager; all public methods are
│   │   concurrency-safe via the internal RWMutex; mutators return the
│   │   receiver for fluent chaining.)
│   │
│   ├── Overloaded constructor family (single shared description)
│   │   ├── New               (source: none — empty configuration)
│   │   ├── NewWithMap        (source: key-value pairs map)
│   │   ├── NewWithProvider   (source: Provider.Read() result)
│   │   └── NewWithGofig      (source: another Gofig's configuration)
│   │
│   ├── Accessors
│   │   ├── Get               (path-parts read with found flag)
│   │   └── Raw               (full deep-copied snapshot)
│   │
│   ├── Mutators (fluent, return *Gofig)
│   │   ├── Set               (path-parts write, creates intermediate maps)
│   │   ├── Delete            (path-parts remove; no parts clears all)
│   │   ├── Merge             (deep merge of another Gofig)
│   │   └── MergeMap          (deep merge of a map)
│   │
│   └── Internal helpers
│       ├── ensureMapInitialized (lazy kvMap allocation)
│       ├── setPath              (walk + write leaf, deep-copy map leaf)
│       ├── deletePath           (walk + delete leaf, ignore missing)
│       ├── getPath              (walk + read leaf with found flag)
│       └── mergeMaps            (recursive deep merge, src overwrites dst)
```

---

## 3. Entity Documentation (as written into `gofig.go`)

### Package `gofig`
Provides a thread-safe, modular configuration container for managing hierarchical key-value pairs. It enforces boundary safety by deep-copying inbound and outbound data maps, isolating internal state.

### `ErrNilProvider`
Sentinel error returned by `NewWithProvider` when a nil Provider source is supplied. It signals that configuration data cannot be loaded from a non-existent provider.

### `Gofig`
The central configuration manager that stores and manages hierarchical key-value pairs. It provides thread-safe access guarded by an internal RWMutex and enforces boundary safety across all operations.

### Constructor family `New<WithXXX>` (single description, varies only with the source argument)
- `New` — returns a new empty Gofig.
- `NewWithMap` — returns a new Gofig loaded with a deep copy of the given key-value pairs.
- `NewWithProvider` — returns a new Gofig loaded with a deep copy of the key-value pairs read from the given Provider.
- `NewWithGofig` — returns a new Gofig loaded with a deep copy of the given Gofig's configuration.

### `(gfg *Gofig) Get`
Returns the value at the given path parts, reporting whether the key was found. It performs a read-locked traversal over the nested configuration maps.

### `(gfg *Gofig) Set`
Stores a value at the given path parts, creating missing intermediate maps as needed. It enforces boundary safety by deep-copying map values and returns the receiver to enable method chaining.

### `(gfg *Gofig) Delete`
Removes the configuration value at the given path parts. When called with no path parts, it clears the entire configuration map and returns the receiver for fluent chaining.

### `(gfg *Gofig) Merge`
Deep-merges the configuration of another Gofig into this one. It snapshots the source container under a read lock and merges nested maps recursively, returning the receiver for method chaining.

### `(gfg *Gofig) MergeMap`
Deep-merges the given key-value pairs into the configuration. It creates a deep copy of the input map to preserve boundary safety and returns the receiver for method chaining.

### `(gfg *Gofig) Raw`
Returns a complete deep copy of the whole configuration map. It guarantees outbound boundary safety so that caller mutations do not alter internal Gofig state.

### `(gfg *Gofig) ensureMapInitialized`
Allocates the internal configuration map if it is uninitialized and returns the receiver.

### `setPath`
Walks path segments and stores value at the leaf, creating missing intermediate maps. Map values are deep-copied to prevent aliasing.

### `deletePath`
Walks path segments and deletes the leaf key. Missing paths are ignored without error.

### `getPath`
Walks path segments through nested maps and returns the leaf value along with a boolean indicating whether the key was found.

### `mergeMaps`
Deep-merges src into dst: nested maps are merged recursively, while every other value overwrites dst.

---

## 4. Verification

- `go test -count=1 github.com/slava-y2k/gofig/...` — ok (core, keyvalue provider; example has no tests).
- `go vet github.com/slava-y2k/gofig/...` — clean.
