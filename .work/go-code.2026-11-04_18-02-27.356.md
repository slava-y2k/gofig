# Go Code Documentation: `gofig.go`

This document records the documentation context, structured context, and final entity documentation for `gofig.go`, created in accordance with the `@go-code` skill instructions (`llm-debug:true`).

---

## 1. Documentation Context

The documentation context is derived from three primary sources:
1. **Method instructions itself** (direct code logic and signatures within `gofig.go`).
2. **Internally invoked methods** (call graph analysis up to a stack depth of 3).
3. **Model information** (design contracts, architectural guarantees, and boundary safety specifications from `docs/model.md` and `docs/code-layout.md`).

### 1.1 Method Instructions Analysis
- **`package gofig`**: Declares the root package for the core gofig configuration library.
- **`ErrNilProvider`**: Sentinel error initialized via `errors.New("gofig: nil provider")`. Guard condition triggered when `NewWithProvider` receives a nil `Provider`.
- **`Gofig`**: Struct holding `rwMutex sync.RWMutex` and `kvMap map[string]any`. Encapsulates internal configuration state and concurrency primitives.
- **`New`**: Allocates and returns a pointer to a `Gofig` with an empty initialized `map[string]any`.
- **`NewWithMap`**: Accepts `keyValues map[string]any`, constructs a new `Gofig` via `New()`, deep-copies input values, and merges them into `kvMap`.
- **`NewWithProvider`**: Accepts `provider Provider`, checks for nil provider (returning `ErrNilProvider`), invokes `provider.Read()`, and delegates to `NewWithMap`.
- **`NewWithGofig`**: Accepts `other *Gofig`, handles nil source by returning an empty `Gofig`, extracts a snapshot via `other.Raw()`, and delegates to `NewWithMap`.
- **`Get`**: Acquires a read lock on `rwMutex`, handles nil/empty map and zero-length path parts (returning the full map), and queries `getPath`.
- **`Set`**: Acquires a write lock on `rwMutex`, ensures `kvMap` is initialized, replaces the root map if `parts` is empty, or invokes `setPath`. Returns `*Gofig` for fluent chaining.
- **`Delete`**: Acquires a write lock on `rwMutex`, resets `kvMap` if `parts` is empty, or calls `deletePath`. Returns `*Gofig` for fluent chaining.
- **`Merge`**: Acquires a read lock on `other.rwMutex` to extract a deep-copied snapshot of `other.kvMap`, acquires a write lock on `gfg.rwMutex`, ensures map initialization, and invokes `mergeMaps`. Returns `*Gofig`.
- **`MergeMap`**: Acquires a write lock on `gfg.rwMutex`, ensures map initialization, deep-copies input `keyValues`, and invokes `mergeMaps`. Returns `*Gofig`.
- **`Raw`**: Acquires a read lock on `rwMutex`, creates a deep copy of `kvMap` using `DeepCopyMap`, and guarantees a non-nil return value.
- **`ensureMapInitialized`**: Checks if `kvMap` is nil; if so, allocates an empty `map[string]any`.
- **`setPath`**: Iterates path segments, allocates intermediate `map[string]any` nodes if absent or of incompatible type, deep-copies map values at the leaf node, and assigns.
- **`deletePath`**: Iterates path segments, terminates early if any intermediate node is absent or not a map, and deletes the leaf key from the parent map.
- **`getPath`**: Iterates path segments, traversing nested `map[string]any` nodes, and returns the leaf value along with a boolean indicating presence.
- **`mergeMaps`**: Iterates source keys, recursively descending when both destination and source entries are nested `map[string]any`, or overwriting destination values for scalar/leaf types.

### 1.2 Internally Invoked Methods (Call Stack Depth <= 3)
1. **`New()`**
   - Depth 1: `make(map[string]any)`
2. **`NewWithMap(keyValues)`**
   - Depth 1: `New()`, `DeepCopyMap(keyValues)`, `mergeMaps(gfg.kvMap, ...)`
   - Depth 2: `mergeMaps` calls `mergeMaps(dstMap, srcMap)` on nested map structures.
3. **`NewWithProvider(provider)`**
   - Depth 1: `provider.Read()`, `NewWithMap(keyValues)`
   - Depth 2: `NewWithMap` calls `New()`, `DeepCopyMap()`, and `mergeMaps()`.
   - Depth 3: `mergeMaps` recursively calls `mergeMaps()` for nested maps.
4. **`NewWithGofig(other)`**
   - Depth 1: `other.Raw()`, `NewWithMap(rawMap)`
   - Depth 2: `other.Raw()` calls `other.rwMutex.RLock()` and `DeepCopyMap()`; `NewWithMap` calls `New()`, `DeepCopyMap()`, and `mergeMaps()`.
   - Depth 3: `mergeMaps` recursively calls `mergeMaps()`.
5. **`Get(parts...)`**
   - Depth 1: `gfg.rwMutex.RLock()`, `getPath(gfg.kvMap, parts)`, `gfg.rwMutex.RUnlock()`
6. **`Set(parts, value)`**
   - Depth 1: `gfg.rwMutex.Lock()`, `gfg.ensureMapInitialized()`, `DeepCopyMap(value)` (if `len(parts) == 0`), `setPath(gfg.kvMap, parts, value)`, `gfg.rwMutex.Unlock()`
   - Depth 2: `setPath` calls `DeepCopyMap(valueMap)` when leaf value is `map[string]any`.
7. **`Delete(parts...)`**
   - Depth 1: `gfg.rwMutex.Lock()`, `deletePath(gfg.kvMap, parts)`, `gfg.rwMutex.Unlock()`
8. **`Merge(other)`**
   - Depth 1: `other.rwMutex.RLock()`, `DeepCopyMap(other.kvMap)`, `other.rwMutex.RUnlock()`, `gfg.rwMutex.Lock()`, `gfg.ensureMapInitialized()`, `mergeMaps(gfg.kvMap, otherMap)`, `gfg.rwMutex.Unlock()`
   - Depth 2: `mergeMaps` recursively calls `mergeMaps()`.
9. **`MergeMap(keyValues)`**
   - Depth 1: `gfg.rwMutex.Lock()`, `gfg.ensureMapInitialized()`, `DeepCopyMap(keyValues)`, `mergeMaps(gfg.kvMap, ...)`, `gfg.rwMutex.Unlock()`
   - Depth 2: `mergeMaps` recursively calls `mergeMaps()`.
10. **`Raw()`**
    - Depth 1: `gfg.rwMutex.RLock()`, `DeepCopyMap(gfg.kvMap)`, `gfg.rwMutex.RUnlock()`
11. **`ensureMapInitialized()`**
    - Depth 1: `make(map[string]any)` (if nil)
12. **`setPath(kvMap, parts, value)`**
    - Depth 1: `make(map[string]any)` (for missing intermediate nodes), `DeepCopyMap(valueMap)` (at leaf)
13. **`deletePath(kvMap, parts)`**
    - Depth 1: `delete(kvMap, key)`
14. **`getPath(kvMap, parts)`**
    - Depth 1: map lookup
15. **`mergeMaps(dst, src)`**
    - Depth 1: recursive call `mergeMaps(dstMap, srcMap)` for nested map values.

### 1.3 Model Information
- **Core Abstraction**: `Gofig` is the central key-value pair manager. Configuration is modeled hierarchically as `map[string]any`.
- **Boundary Safety**: As detailed in `docs/model.md` under Gofig implementation details, Gofig strictly provides bound safety. All methods accepting an inbound key-value collection (`NewWithMap`, `Set`, `MergeMap`, etc.) create an internal deep copy. All outbound operations (`Raw`) return a deep copy. Internal state is never exposed as an aliased reference to callers.
- **Concurrency Safety**: Public `Gofig` methods must be safe for concurrent use across goroutines, guarded by an encapsulated `sync.RWMutex`.
- **Fluent Interface**: Mutating methods (`Set`, `Delete`, `Merge`, `MergeMap`) return the `*Gofig` receiver to enable method chaining.
- **Extensibility & Decoupling**: External data sources supply configuration by implementing the `Provider` interface (`Read() (map[string]any, error)`). The core `gofig` package depends solely on the Go standard library.

---

## 2. Structured Context

The structured context follows the code architecture (`package -> interfaces/structs -> methods/constants`), elevating common abstractions, boundary safety semantics, and concurrency invariants to the upper tiers.

```text
Package: gofig
│  (Common context: Thread-safe, stdlib-only hierarchical configuration container.
│   Enforces boundary safety via deep copying across all boundaries.)
│
├── Constants & Sentinel Errors
│   └── ErrNilProvider (Guard sentinel for uninitialized Provider inputs)
│
├── Struct: Gofig
│   │  (Abstraction: Thread-safe in-memory container encapsulating kvMap and rwMutex.
│   │   Provides bound safety, concurrent access, and a fluent mutation interface.)
│   │
│   ├── Constructors
│   │   ├── New               (Empty container allocation)
│   │   ├── NewWithMap        (Map ingestion with deep-copy boundary protection)
│   │   ├── NewWithProvider   (External Provider ingestion via Read())
│   │   └── NewWithGofig      (Isolated cloning of an existing Gofig container)
│   │
│   ├── Accessors
│   │   ├── Get               (Read-locked path traversal retrieving leaf values)
│   │   └── Raw               (Read-locked full configuration snapshot with deep copy)
│   │
│   ├── Mutators
│   │   ├── Set               (Write-locked path-based assignment with fluent return)
│   │   ├── Delete            (Write-locked path-based key deletion with fluent return)
│   │   ├── Merge             (Write-locked recursive merge of another Gofig with fluent return)
│   │   └── MergeMap          (Write-locked recursive merge of external map with fluent return)
│   │
│   └── Internal Helpers
│       ├── ensureMapInitialized (Lazy allocation of internal map)
│       ├── setPath              (Traverses/allocates intermediate path maps and writes leaf)
│       ├── deletePath           (Traverses path segments and removes leaf key)
│       ├── getPath              (Traverses path segments and reads leaf value)
│       └── mergeMaps            (Recursive map-combining logic preserving disjoint trees)
```

---

## 3. Entity Documentation

In accordance with the skill rules, each documentation entry consists of 1–3 sentences explaining why and what for the entity was created, its primary idea, and the abstraction used.

### Package `gofig`
Package gofig provides a thread-safe, modular configuration container for managing hierarchical key-value pairs. It enforces boundary safety by deep-copying inbound and outbound data maps, isolating internal state. The core package relies solely on the Go standard library.

### `ErrNilProvider`
`ErrNilProvider` is a sentinel error returned by `NewWithProvider` when a nil `Provider` source is supplied. It signals that configuration data cannot be loaded from a non-existent provider.

### `Gofig`
`Gofig` is the central configuration manager that stores and manages hierarchical key-value pairs. It provides thread-safe access guarded by an internal `RWMutex` and enforces boundary safety across all operations.

### `New()`
`New` constructs an empty `Gofig` instance initialized with an allocated internal map. It is intended for building configurations dynamically via in-memory mutations.

### `NewWithMap(keyValues map[string]any)`
`NewWithMap` creates a new `Gofig` populated with the provided key-value map. It performs a deep copy of the input map to ensure boundary safety and isolate caller data from internal modifications.

### `NewWithProvider(provider Provider)`
`NewWithProvider` creates a new `Gofig` loaded from the supplied `Provider` source. It reads the provider's key-value collection and returns an error if the provider is nil or if reading fails.

### `NewWithGofig(other *Gofig)`
`NewWithGofig` creates a new `Gofig` initialized from another `Gofig` instance. It takes a thread-safe deep copy snapshot of the source container to guarantee complete isolation between instances.

### `(gfg *Gofig) Get(parts ...string)`
`Get` returns the value at the given path parts, reporting whether the key was found. It performs a read-locked traversal over the nested configuration maps.

### `(gfg *Gofig) Set(parts []string, value any)`
`Set` stores a value at the given path parts, creating missing intermediate maps as needed. It enforces boundary safety by deep-copying map values and returns the receiver to enable method chaining.

### `(gfg *Gofig) Delete(parts ...string)`
`Delete` removes the configuration value at the given path parts. When called with no path parts, it clears the entire configuration map and returns the receiver for fluent chaining.

### `(gfg *Gofig) Merge(other *Gofig)`
`Merge` deep-merges the configuration of another `Gofig` into this one. It snapshots the source container under a read lock and merges nested maps recursively, returning the receiver for method chaining.

### `(gfg *Gofig) MergeMap(keyValues map[string]any)`
`MergeMap` deep-merges the given key-value pairs into the configuration. It creates a deep copy of the input map to preserve boundary safety and returns the receiver for method chaining.

### `(gfg *Gofig) Raw()`
`Raw` returns a complete deep copy of the whole configuration map. It guarantees outbound boundary safety so that caller mutations do not alter internal Gofig state.

### `(gfg *Gofig) ensureMapInitialized()`
`ensureMapInitialized` allocates the internal configuration map if it is uninitialized and returns the receiver.

### `setPath(kvMap map[string]any, parts []string, value any)`
`setPath` walks path segments and stores value at the leaf, creating missing intermediate maps. Map values are deep-copied to prevent aliasing.

### `deletePath(kvMap map[string]any, parts []string)`
`deletePath` walks path segments and deletes the leaf key. Missing paths are ignored without error.

### `getPath(kvMap map[string]any, parts []string)`
`getPath` walks path segments through nested maps and returns the leaf value along with a boolean indicating whether the key was found.

### `mergeMaps(dst, src map[string]any)`
`mergeMaps` deep-merges `src` into `dst`: nested maps are merged recursively, while every other value overwrites `dst`.
