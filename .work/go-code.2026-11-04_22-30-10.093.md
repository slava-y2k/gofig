# Go Code Documentation Run: `gofig.go`

- Run timestamp: `2026-11-04_22:30:10.093`
- Skill: `@go-code` (documentation part), `llm-debug:true`
- Result: **re-validated, no source changes required** — `gofig.go` already carries the documentation produced by the previous run and still matches the skill rules, the updated `AGENTS.md`, and `docs/code/gofig.md`.
- Inputs used:
  - `.opencode/skills/go-code/SKILL.md` (incl. *Special cases*: emulated overloads)
  - `AGENTS.md` (`docs/code` mirrors the code structure and is additional generation context)
  - `docs/code/gofig.md` — `New<WithXXX>` are overloaded methods
  - `docs/model.md`, `docs/code-layout.md` — model information
  - `gofig.go` (documented file), `provider.go`, `gofig_util.go` (context)

---

## 1. Documentation Context

### 1.1 Method instructions itself

- **`package gofig`** — stdlib-only (`errors`, `sync`) core package of the configuration library.
- **`ErrNilProvider`** — sentinel error for a nil `Provider` source, returned by `NewWithProvider`.
- **`Gofig`** — `{ rwMutex sync.RWMutex; kvMap map[string]any }`; the managed state and its concurrency guard.
- **`New`** — allocates a `Gofig` with an empty initialized map.
- **`NewWithMap`** — `New()` → `DeepCopyMap` into `mapCopy` → `mergeMaps`; always returns a nil error.
- **`NewWithProvider`** — nil check (`ErrNilProvider`), `provider.Read()`, error propagation, delegate to `NewWithMap`.
- **`NewWithGofig`** — nil source → empty `Gofig`; otherwise `other.Raw()` → `NewWithMap`.
- **`Get`** — RLock; nil map → `(nil, false)`; zero parts → whole map; otherwise `getPath`.
- **`Set`** — Lock; `ensureMapInitialized`; zero parts + map → deep-copied root replacement; otherwise `setPath`; returns receiver.
- **`Delete`** — Lock; zero parts → fresh empty map; otherwise `deletePath` on a non-nil map; returns receiver.
- **`Merge`** — deep-copied snapshot of `other` under its RLock, then Lock, `ensureMapInitialized`, `mergeMaps`; returns receiver.
- **`MergeMap`** — Lock; `ensureMapInitialized`; `mergeMaps` with a deep copy of the input; returns receiver.
- **`Raw`** — RLock; `DeepCopyMap`; non-nil result guaranteed.
- **`ensureMapInitialized`** — lazy `kvMap` allocation; returns receiver.
- **`setPath`** — walks parts, allocates missing intermediate maps, deep-copies a map leaf.
- **`deletePath`** — walks parts, silently stops on missing/non-map intermediate, deletes the leaf.
- **`getPath`** — walks parts, returns the leaf value and found flag.
- **`mergeMaps`** — recursive deep merge; map pairs merge recursively, other values overwrite `dst`.

### 1.2 Internally invoked methods (stack depth ≤ 3)

1. `New` → `make(map[string]any)`.
2. `NewWithMap` → `New`, `DeepCopyMap`, `mergeMaps`; depth 2: recursive `mergeMaps`.
3. `NewWithProvider` → `provider.Read()`, `NewWithMap`; depth 2: `New`/`DeepCopyMap`/`mergeMaps`; depth 3: recursive `mergeMaps`.
4. `NewWithGofig` → `other.Raw()`, `NewWithMap`; depth 2: `Raw` → `DeepCopyMap`; `NewWithMap` → `DeepCopyMap`/`mergeMaps`; depth 3: recursive `mergeMaps`.
5. `Get` → `getPath`.
6. `Set` → `ensureMapInitialized`, `DeepCopyMap` (root replace), `setPath`; depth 2: `setPath` → `DeepCopyMap` (map leaf).
7. `Delete` → `deletePath`.
8. `Merge` → `other.rwMutex.RLock` + `DeepCopyMap`, `ensureMapInitialized`, `mergeMaps`; depth 2: recursive `mergeMaps`.
9. `MergeMap` → `ensureMapInitialized`, `DeepCopyMap`, `mergeMaps`; depth 2: recursive `mergeMaps`.
10. `Raw` → `DeepCopyMap`.
11. `ensureMapInitialized` → `make(map[string]any)`.
12. `setPath` → `make(map[string]any)`, `DeepCopyMap`.
13. `deletePath` → `delete`.
14. `getPath` → map lookup.
15. `mergeMaps` → recursive `mergeMaps`.

### 1.3 Model information

- Configuration is a hierarchical `map[string]any` key-value collection (`docs/model.md`).
- `Gofig` is the central manager: construct from sources (`NewWithXXX`), merge, mutate, read.
- **Boundary safety** (`docs/model.md` "Gofig implementation details"): inbound collections are copied internally; outbound methods return copies. This is the abstraction behind `DeepCopyMap` at every boundary.
- Public methods are safe for concurrent use (`docs/code-layout.md` Conventions) — internal `sync.RWMutex`.
- Mutators return `*Gofig` for a fluent interface.
- `Provider` is a one-method source: `Read() (map[string]any, error)`.

### 1.4 Additional context: `docs/code`

`docs/code/gofig.md` states: **`New<WithXXX>` are overloaded methods.** Go lacks overloads, so the family `New`, `NewWithMap`, `NewWithProvider`, `NewWithGofig` emulates one overloaded constructor keyed by the source argument.

### 1.5 Special case applied (emulated overloads)

Per the skill's *Special cases*: overloaded methods use a **single description varying only with the provided args**. The constructor family shares one template — "returns a new Gofig loaded with a deep copy of ⟨source⟩":

- `New` — no source → "returns a new empty Gofig".
- `NewWithMap` — source: the given key-value pairs.
- `NewWithProvider` — source: the key-value pairs read from the given Provider.
- `NewWithGofig` — source: the given Gofig's configuration.

Variant-specific detail (nil handling via `ErrNilProvider`) remains documented on the sentinel error, not duplicated per constructor.

---

## 2. Structured Context

```text
Package: gofig
│  (Level-wide: thread-safe hierarchical key-value configuration container;
│   stdlib-only core. Boundary safety is the package-level invariant — every
│   inbound collection is deep-copied, every outbound map is a deep copy.
│   Overloaded constructor family New<WithXXX> shares one description
│   varying only with the provided source argument.)
│
├── Sentinel error
│   └── ErrNilProvider (returned by NewWithProvider for a nil Provider)
│
├── Struct: Gofig
│   │  (Entity-level: central key-value manager; public methods are
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
│   │   ├── Set [..](..)              (path-parts write, creates intermediate maps)
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

## 3. Entity Documentation (current state of `gofig.go`)

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
- Documentation coverage: every item in `gofig.go` (package clause, sentinel error, struct, all constructors, accessors, mutators, internal helpers) carries a doc comment.

## 5. Observation (outside the documentation task)

`Get()` called with no path parts returns `gfg.kvMap` itself, not a deep copy, while `docs/model.md` states that outbound methods return a copy. This is a code-level inconsistency, not a documentation defect; the current doc comment stays neutral about it. Suggested follow-up: return a deep copy (or an error) for the zero-parts case, and record the finding in `docs/inconsistencies.md` via `/model-consistency-check`.
