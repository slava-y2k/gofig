# Gonfig Roadmap — Nearest Work

*Generated from `docs/backlog.md` (Nov 13 2026). This is a planning snapshot;
see `docs/_index.md` for the convention that `docs/output` is generated HTML.*

## 1. Summary

The backlog has three feature/decision items, two applied items, and a set of
model↔implementation consistency findings:

- **Feature / decision items** (3): implement `gonfig` tag, prepare the
  provider list, and drop the `bool` return from `Get`.
- **Applied** (2): `IC-015` "deep copy" wording and `[Get-copy]` — the full
  deep copy for `Get` — both done (see §2).
- **Consistency findings**: 0 Critical, 3 Major, 8 Minor.
- **Model gaps** (2): no method signatures in prose; zero value / errors / helpers undefined.

Key insight: the model specifies a full *deep-copy* contract — `Get` and
`Raw` both return deep copies, and every accepted or returned item is a deep copy. The code now
implements that contract: `Get` returns a full deep copy, and `Raw`/`DeepCopyMap` deep-copy maps,
slices, arrays, pointers, and structs. The remaining work is the `Get` signature decision and the
doc-only wording fixes.

## 2. Applied

> **IC-015** — the "deep copy" wording was first replaced with "map-level copy" across
> `gonfig.go`, `docs/design.md`, and `docs/code.md`. The direction then
> changed to a full *deep-copy* contract, so that rewording is now superseded: the docs again say
> "deep copy", but with an explicit definition ("if a value is a map or another reference type, a new copy
> is created for it"). No behavior change either way.

> **[Get-copy]** — `Get` now returns a full deep copy instead of a live reference, so a caller that keeps
> a result cannot mutate the container's internal map through it. `DeepCopyMap` deep-copies maps, slices,
> arrays, pointers, and structs via reflection; `Get` returns the deep copy for both the whole map and a
> path hit; `setPath` deep-copies any leaf. The live-reference and sharing tests were rewritten as copy
> tests (`TestGetWholeConfigReturnsCopy`, `TestGetMapValueReturnsCopy`, `TestGetSliceReturnsCopy`,
> `TestGetTypedMapReturnsCopy`, plus `DeepCopyMap` slice/typed-map/pointer/struct cases).

## 3. Nearest work

### 3.1 Decide the `Get` signature

**Backlog item:** "In Get method bool return is not required". The model prose
(`docs/design.md`) already describes value-only returns ("nil when missing"), while the diagram
and code return `(any, bool)`. Decide and make all sides agree.

- If dropping `bool`: update `gonfig.go`, `getPath`, the diagram,
  `gonfig_test.go` call sites, and the README quick-start.
- If keeping `bool`: reword the model prose to match the diagram instead.

*Not trivial — this touches the public API, tests, and examples.*

### 3.2 Remaining doc-only fixes

These are comment/doc-only edits. No behavior change, no test change, no API change.

| Sev | Ref | Doc location | Fix |
|-----|-----|--------------|-----|
| **Major** | `Delete` described as "set value" | `docs/design.md` | "Delete — removes the leaf at the key path (no-op on missing/non-map segments)". |
| **Major** | "empty key → nil" wording | `docs/design.md` | `Set`/`Delete` return `*Gonfig`, not nil; `Get()` with no parts returns the root map with `true`. Fix the prose. |
| **Major** | "overloaded methods" | `docs/code.md` | "constructor variants" — Go has no overloading (already correct at `docs/design.md`). |

### 3.3 Close the Minor undocumented-code gaps

Add short model entries (a sentence each) so the model covers what the code exposes. These are purely
additive documentation:

- `gonfig.go` — public `ErrNilProvider` sentinel and its nil-provider contract.
- `gonfig.go` — usable zero value and `New()`.
- `gonfig.go` — `NewWithMap`/`NewWithProvider`/`NewWithGonfig`
  semantics (nil→empty, error propagation, snapshot).
- `gonfig.go` — `Get` returns `(any, bool)` (see §3.1 before finalizing).
- `gonfig.go` — fluent `*Gonfig` returns and the deep-merge rule ("later sources
  win", recursive).
- `gonfig.go` — `Raw` never returns nil.
- `gonfig_util.go` — exported `DeepCopyMap`.
- `provider/keyvalue/keyvalue.go` — `KeyValue`, `New`,
  `WithDeepCopy`, `Read`.

### 3.4 Feature work (in order)

**Phase B — Prepare the provider list**

**Backlog item:** "Prepare the provider list for the implementation." Produce a short
ranked list of provider packages (e.g. `env`, `file/json`, `file/yaml`,
`flags`) with their dependencies and a first-implementation target. Then implement the first one
following the existing `provider/keyvalue` pattern (separate module, `Read() map[string]any`).

**Phase C — Implement the `gonfig` struct tag**

**Backlog item:** "implement gonfig tag". Map structs to configuration via a struct tag
(like `envconfig`/`mapstructure`). Largest item; defer until the `Get` signature,
the doc-only fixes, and the provider list are settled.

## 4. Suggested sequence

1. **Next:** settle the `Get` signature (§3.1).
2. **Then:** the remaining doc-only fixes (§3.2–3.3).
3. **Then:** Phase B — provider list + first provider.
4. **Last:** Phase C — `gonfig` tag.
