# Inconsistencies

The registry of divergences between the design model (`docs/`) and the
implementation: what the model says, what the code or the other docs do
instead, and the prepared change that resolves it. The
`/model-consistency-check` command maintains this document — the reviewer
agent prepares the entries, the session that ran the check appends them.

## Rules

- IDs are sequential (`IC-001`, `IC-002`, …) and never reused.
- Entries are never deleted. When a divergence is resolved, set its status
  to `resolved` and note how.
- A repeated finding updates the existing open entry instead of adding a
  duplicate.

## Entry format

```md
### [IC-###] Short title

- Severity: Critical | Major | Minor
- Kind: idea-inconsistency | undocumented-code | unimplemented-model | contradiction
- Status: open | resolved
- Model says: docs/…:line — quote or reference
- Reality: file…:line — quote, diagram, or verification result
- Prepared change: the concrete edit — to the docs, the code, or both —
  that resolves the divergence
```

## Findings

### [IC-001] Load method broken; workspace does not build

- Severity: Critical
- Kind: contradiction
- Status: open
- Model says: docs/code-layout.md:125 — "`Load` reads a `Provider` and folds in its output" (a `Gofig` method); callers use it as `g.Load(provider)` returning `error` (gofig_test.go:305,324,335,395; example/basic/main.go:21).
- Reality: gofig.go:30 — `func Load(provider Provider) (error, *Gofig)` is a free function whose body uses the undefined identifier `gfg` (gofig.go:42-45) and returns one value where two are declared (gofig.go:32,39,47). Verified: `go build/vet/test github.com/slava-y2k/gofig/...` all fail with "undefined: gfg" and "not enough return values"; every module reports `[build failed]`.
- Prepared change: In gofig.go replace the free function with the method `func (gfg *Gofig) Load(provider Provider) error`, returning `ErrNilProvider` when the provider is nil, the provider's error when `Read()` fails, and `nil` after merging `DeepCopyMap(keyValues)` into `gfg.kvMap` under the write lock. Also update the model.md:44 diagram to `+Load(provider) error`.

### [IC-002] Model names the core struct "Gonfig" instead of "Gofig"

- Severity: Major
- Kind: idea-inconsistency
- Status: open
- Model says: docs/model.md:36 (title "Gonfig library design"), :41 (`namespace Gonfig_Lib`), :42 (`class Gonfig["Gonfig"]`).
- Reality: docs/model.md:6, docs/code-layout.md:125 and README.md:4 all call the type "Gofig". One entity, two names.
- Prepared change: In model.md, rename the diagram to `Gofig` (`class Gofig["Gofig"]`, `namespace Gofig_Lib`, title "Gofig library design").

### [IC-003] Parser contract referenced but not defined in model.md

- Severity: Major
- Kind: idea-inconsistency
- Status: open
- Model says: docs/code-layout.md:130 — "The `Parser` contract (`Parse(m map[string]any) (map[string]any, error)`) is defined in [model.md](model.md)".
- Reality: docs/model.md contains no `Parse`/`Parser` contract (grep for "Parse" matches only model.md:7 and model.md:44).
- Prepared change: Add the `Parser` contract (`Parse(m map[string]any) (map[string]any, error)`) to model.md, the stated source of truth; otherwise remove the "is defined in model.md" claim from code-layout.md:130 and state the contract is not yet defined.

### [IC-004] "gofig tag" entity not implemented and not marked planned

- Severity: Major
- Kind: unimplemented-model
- Status: open
- Model says: docs/model.md:7 — "gofig tag - marks the struct field" listed as a main entity.
- Reality: No struct-tag/reflection binding logic exists in any source file (gofig.go, provider.go, gofig_util.go, provider/keyvalue/*); docs/code-layout.md:140-145 (Planned) does not list struct-tag support.
- Prepared change: Either implement struct-tag binding (tag parsing + mapping into Gofig), or move it into code-layout.md's Planned section and reword model.md:7 to present the gofig tag as a future entity rather than a current main entity.

### [IC-005] Get signature in the model diagram omits the found result

- Severity: Minor
- Kind: contradiction
- Status: open
- Model says: docs/model.md:48 — `+Get(parts ...string) any`.
- Reality: gofig.go:53 — `func (gfg *Gofig) Get(parts ...string) (any, bool)`; tests depend on the bool (gofig_test.go:19-22,30-35).
- Prepared change: Update the model.md diagram to `+Get(parts ...string) (any, bool)`.

### [IC-006] Provider method shown as lowercase `read()` in the diagram

- Severity: Minor
- Kind: idea-inconsistency
- Status: open
- Model says: docs/model.md:54 — `+read() map[string, any]` (lowercase `read`, no error result, malformed type).
- Reality: docs/model.md:16 and docs/code-layout.md:120 both state `Read() (map[string]any, error)`.
- Prepared change: Update the model.md diagram to `+Read() (map[string]any, error)`.

### [IC-007] docs/parser.md is a stale reference

- Severity: Minor
- Kind: idea-inconsistency
- Status: open
- Model says: docs/code-layout.md:110 — lists `parser.md` under `docs/` in the directory tree.
- Reality: `docs/parser.md` was deleted (git status shows `D docs/parser.md`) and is absent from disk (`find . -name parser.md` returns nothing).
- Prepared change: Remove the `parser.md` line from the directory tree in code-layout.md.

### [IC-008] parser/ directory listed as present but does not exist

- Severity: Minor
- Kind: idea-inconsistency
- Status: open
- Model says: docs/code-layout.md:96 — directory tree shows `parser/`; :114-115 — "The empty `parser/` directory is a local placeholder".
- Reality: `parser/` does not exist on disk (`ls parser` → "No such file or directory") and is not tracked in git.
- Prepared change: Reword code-layout.md:96 and :114-115 to state that parser modules will live under `parser/` once the first parser module is added, and drop the claim that the directory currently exists.

### [IC-009] New() and ErrNilProvider are undocumented public surface

- Severity: Minor
- Kind: undocumented-code
- Status: open
- Model says: docs/code-layout.md:125 — the `Gofig` method table defines no `New` constructor and no `ErrNilProvider`; model.md lists neither.
- Reality: gofig.go:25 `func New() *Gofig` and gofig.go:9 `var ErrNilProvider` are exported and used by the example (example/basic/main.go:18) and tests (gofig_test.go:10,331-340), with no model definition.
- Prepared change: Document `New()` and `ErrNilProvider` in code-layout.md (Core abstractions table and/or Conventions) or in model.md.
