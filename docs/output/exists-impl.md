# `Exists` — implementation and model consistency

*Generated for the `Gonfig.Exists` addition (Nov 13 2026). Companion to `docs/design.md` and `docs/code.md`.*

## 1. What was added

`Gonfig.Exists(parts ...string) bool` in `gonfig.go:77-92`:

- Reuses the same `getPath` lookup as `Get`, but returns only the found flag.
- `Exists("a", "b")` is `true` exactly when `Get("a", "b")` would report `ok == true`, including a
  stored `nil` value.
- `Exists()` with no parts mirrors `Get()`: `true` when the root map is present (non-nil), `false` for
  the zero value.
- Reads under `RLock`, so it is safe alongside writers.
- Performs no deep copy: a `bool` is a scalar and cannot alias container state.

## 2. Model consistency

| Model source | Statement | Verdict |
|--------------|-----------|---------|
| `docs/design.md:31` diagram | `+Exists(parts ...string) bool` | **Consistent** — matches the code signature at `gonfig.go:77`. |
| `docs/design.md:9-11` copy-contract | every accepted/returned item is a deep copy; a deep copy applies to maps and other reference types | **Consistent** — a `bool` is not a reference type, so the rule does not apply. |
| `docs/design.md:81` access contract | a missing key is reported as `nil, false` by `Get` | **Consistent** — `Exists` reports the same condition as `false`. |
| `docs/design.md:78-82` root behavior | a call with no key parts targets the root | **Consistent** — `Exists()` reports whether the root map is present. |
| `docs/design.md:88-89` access bullet | `Exists` returns a `bool` and needs no copy | **Consistent** — added with this change. |
| `docs/design.md:94-95` storage | the internal map is guarded by `sync.RWMutex` | **Consistent** — `Exists` takes `RLock`. |
| `docs/code.md:18-19` copy-contract scope | the contract applies when results are maps or other mutable structs | **Consistent** — a `bool` is neither. |
| `docs/code.md:5` thread-safety guidance | document thread safety when it keeps the explanation compact | **Consistent** — the `RLock` is visible; the doc comment stays one line. |

## 3. Findings

- **No contradictions.** The implementation matches the model.
- **Model updated** to include `Exists`: the class diagram (`docs/design.md:31`) and the access-methods
  bullet (`docs/design.md:88-89`).
- **Deliberate semantic choice:** `Exists()` with no parts reports the presence of the root map — true
  for an empty but initialized container, false for the zero value — matching `Get()`'s found flag rather
  than "has any entries". This keeps `Exists` a pure predicate over `Get`.
- **No README change:** the README does not enumerate methods, so it needs no entry.

## 4. Test coverage (`gonfig_test.go`)

- `TestExistsZeroValueIsFalse` — nil root, both call forms.
- `TestExistsNoPartsReportsRoot` — initialized empty container.
- `TestExistsPresentAndMissingKeys` — found and absent paths.
- `TestExistsStoredNilIsTrue` — a stored `nil` key still exists.
- `TestExistsNonMapSegmentIsFalse` — a scalar blocks the path.

All tests pass under `go test -race`.

## 5. Notes

- `Exists` is intentionally copy-free: the copy-contract exists so callers cannot mutate container state
  through a returned reference. A `bool` carries no reference, so there is nothing to protect.
- The method is read-only; it never initializes or mutates the root.
