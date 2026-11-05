# Inconsistencies

The registry of divergences between the design model (`docs/`) and the
implementation: what the model says, what the code or the other docs do
instead, and the prepared change that resolves it. The
`/model-consistency-check` command maintains this document — the reviewer
agent prepares the entries, the session that ran the check appends them.

## Rules

- IDs are sequential (`IC-001`, `IC-002`, …) and never reused, even after
  the entry is removed.
- Resolved entries are removed once the fix or the resolution is agreed.
- A repeated finding updates the existing open entry instead of adding a
  duplicate.

## Entry format

```md
### [IC-###] Short title

- Severity: Critical | Major | Minor
- Kind: idea-inconsistency | undocumented-code | unimplemented-model | contradiction
- Status: open
- Model says: docs/…:line — quote or reference
- Reality: file…:line — quote, diagram, or verification result
- Prepared change: the concrete edit — to the docs, the code, or both —
  that resolves the divergence
```

## Findings

### [IC-003] Parser contract referenced but not defined in model.md

- Severity: Major
- Kind: idea-inconsistency
- Status: open
- Model says: docs/code-layout.md:130 — "The `Parser` contract (`Parse(m map[string]any) (map[string]any, error)`) is defined in [model.md](model.md)".
- Reality: docs/model.md contains no `Parse`/`Parser` contract (grep for "Parse" returns no matches).
- Prepared change: Add the `Parser` contract (`Parse(m map[string]any) (map[string]any, error)`) to model.md, the stated source of truth; otherwise remove the "is defined in model.md" claim from code-layout.md:130 and state the contract is not yet defined.

### [IC-004] "gonfig tag" entity not implemented and not marked planned

- Severity: Major
- Kind: unimplemented-model
- Status: open
- Model says: docs/model.md:7 — "gonfig tag - marks the struct field" listed as a main entity.
- Reality: No struct-tag/reflection binding logic exists in any source file (gonfig.go, provider.go, gonfig_util.go, provider/keyvalue/*); docs/code-layout.md:140-145 (Planned) does not list struct-tag support.
- Prepared change: Either implement struct-tag binding (tag parsing + mapping into Gonfig), or move it into code-layout.md's Planned section and reword model.md:7 to present the gonfig tag as a future entity rather than a current main entity.

### [IC-005] Get signature in the model diagram omits the found result

- Severity: Minor
- Kind: contradiction
- Status: open
- Model says: docs/model.md:48 — `+Get(parts ...string) any`.
- Reality: gonfig.go:72 — `func (gfg *Gonfig) Get(parts ...string) (any, bool)`; tests depend on the bool (gonfig_test.go:19-22,30-35).
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
- Reality: `docs/parser.md` is absent from disk (`ls docs` shows only code-layout.md, inconsistencies.md, model.md) and its deletion is committed (git status no longer lists it).
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
- Model says: docs/code-layout.md:125 — the `Gonfig` method table defines no `New` constructor and no `ErrNilProvider`; model.md lists neither.
- Reality: gonfig.go:26 `func New() *Gonfig` and gonfig.go:10 `var ErrNilProvider` are exported — `New` backs `NewFromMap` (gonfig.go:53) — and used by tests (gonfig_test.go:10,335-346), with no model definition; the example no longer uses `New` (example/basic/main.go:22 uses `NewFromProvider`).
- Prepared change: Document `New()` and `ErrNilProvider` in code-layout.md (Core abstractions table and/or Conventions) or in model.md.

### [IC-010] Gonfig "writing" responsibility has no API and no Planned entry

- Severity: Major
- Kind: unimplemented-model
- Status: open
- Model says: docs/model.md:22 — Gonfig serves "writing - write to different type of sources" as one of its two responsibilities (docs/model.md:20-22).
- Reality: No write-to-source API exists in any Go file (grep for "Write" returns no matches); `Set` (gonfig.go:89) writes into the in-memory configuration, not to a source. docs/code-layout.md:125 (Core abstractions) and docs/code-layout.md:140-145 (Planned) list no writing path.
- Prepared change: Either mark the writing responsibility as planned — add a write-to-source item to code-layout.md's Planned section and reword model.md:22 as future work — or drop the writing bullet from model.md:22 until a design exists.

### [IC-011] Key normalization for nested non-string keys specified but unimplemented

- Severity: Major
- Kind: unimplemented-model
- Status: open
- Model says: docs/model.md:11-13 — a value "may be also a nested map with type `key:any, value:any`, in this case keys will be normalized to a string type".
- Reality: No normalization code exists: `DeepCopyMap` (gonfig_util.go:5-19), `mergeMaps` (gonfig.go:224-234), and `setPath`/`getPath`/`deletePath` (gonfig.go:176-220) type-assert `map[string]any` only, so a nested map with non-string keys is stored as an opaque value, never normalized. Not listed in docs/code-layout.md:140-145 (Planned).
- Prepared change: Either implement key normalization (convert non-string nested-map keys to strings when loading and merging) or reword model.md:12-13 to present normalization as planned and add it to code-layout.md's Planned section.

### [IC-012] Diagram shows constructor results as "error, Gonfig" instead of "(*Gonfig, error)"

- Severity: Minor
- Kind: idea-inconsistency
- Status: open
- Model says: docs/model.md:44 — `+NewFromXXX()$ error, Gonfig` (error first); docs/code-layout.md:125 — the constructors "return `(*Gonfig, error)`".
- Reality: The two model documents disagree on the constructor result order; the implementation follows code-layout.md (gonfig.go:33 `func NewFromProvider(provider Provider) (*Gonfig, error)`; also gonfig.go:52, gonfig.go:61).
- Prepared change: Update the model.md diagram to `+NewFromXXX()$ (Gonfig, error)`, matching code-layout.md:125 and the implementation.
