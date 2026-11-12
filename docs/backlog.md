# Backlog

- implement "gonfig tag"
- Prepare the provider list for the implementation.
- In Get method bool return is not required
- [Get-copy] Make `Get` return a full deep copy instead of a live reference, so callers cannot mutate the
  container's internal map through the returned map. Requires a full deep copy (maps and reference types such
  as slices and typed maps) — the current `DeepCopyMap` copies only nested `map[string]any` levels and shares
  slices/typed maps. Also update `Get`'s doc comments and the copy-contract wording in `gonfig.go`.
- [Get-copy] Test the deep copy: mutating a `Get` result (whole map, map-valued key, slice/typed-map value)
  must not change the container. Update `TestGetWholeConfigReturnsLiveMap` and
  `TestGetMapValueReturnsReference`, which currently pin the live-reference behavior.
- [IC-015] Raw "deep copy" wording ignores shared slices and typed maps
  Severity| Minor, Kind| contradiction, Status| implemented (see "# Implemented"),
  Model says | gonfig.go:140 — "Raw returns a deep copy per the copy-contract"; docs/design.md:72 — "returns a deep copy of the internal map"; docs/code.md:11 — "Raw() deep-copies its result".
  Reality | DeepCopyMap copies map[string]any levels only and shares every other value (gonfig_util.go:3-19); a slice obtained from Raw still points at the container's backing array (gonfig_test.go:678, TestRawSliceSharesBackingArray), and typed maps stay shared (gonfig_test.go:748, TestDeepCopyMapSharesTypedMap).
  Prepared change | Reword gonfig.go:140, docs/design.md:72 and docs/code.md:11 to "map-level copy": nested map[string]any levels are copied, slices and typed maps are shared.

# Code items 
Critical	contradiction	gonfig.go:59-72 vs docs/design.md:14-15 — the model promises "every method makes a copy for a reference type" on outbound values; Get returns live state: the whole map when no parts are given (gonfig.go:68-69) and map values by reference (gonfig.go:72, documented at gonfig.go:59-60). Tests pin this behavior (gonfig_test.go:140, :161).
Major	contradiction	gonfig_util.go:3-16 vs docs/design.md:12-13 — the model promises an internal copy for any accepted key-value collection; DeepCopyMap copies only map[string]any levels and shares other values (gonfig_util.go:15), so slices and typed maps stay aliased.
Major	contradiction	gonfig.go:97-111 vs docs/design.md:34 — the model describes Delete as "set value by provided key"; the implementation removes the leaf (gonfig.go:107, deletePath at gonfig.go:191-200).
Major	contradiction	gonfig.go:61-111 vs docs/design.md:30 — the model says "on empty key nill will be retuned"; the code returns *Gonfig from Set/Delete and replaces the root/clears it (gonfig.go:82-87, :101-104), and Get() returns the root map with true (gonfig.go:68-69).
Major	contradiction	gonfig.go:20-57 vs docs/code.md:3 — the model says New<WithXXX> are "overloaded methods"; Go has no overloading. The code exposes distinct functions: New (gonfig.go:20), NewWithMap (:25), NewWithProvider (:34), NewWithGonfig (:51).
Minor	undocumented-code	gonfig.go:9-10 — ErrNilProvider is public and returned on a nil provider (gonfig.go:35-37); the model defines no error value.
Minor	undocumented-code	gonfig.go:12-22 — the usable zero value and the empty New() constructor are not in the model; only generic NewXXX constructors are mentioned (docs/design.md:19).
Minor	undocumented-code	gonfig.go:25-57 — concrete constructors and their semantics (nil input to empty, provider error propagation, snapshot of another container) are undefined; the model has only the NewXXX placeholder.
Minor	undocumented-code	gonfig.go:61 — Get returns (any, bool); the found flag is not part of the documented contract (docs/design.md:32).
Minor	undocumented-code	gonfig.go:77-142 — fluent returns (*Gonfig) and the deep-merge rule ("later sources win", recursively) are not documented (docs/design.md:21-24 only says merge "collect[s] key-values").
Minor	undocumented-code	gonfig.go:144-155 — Raw never returns nil and copies only map levels; the model says only "returns the copy of the internal map" (docs/design.md:26).
Minor	undocumented-code	gonfig_util.go:4 — exported DeepCopyMap is not defined by the model.
Minor	undocumented-code	provider/keyvalue/keyvalue.go:8-34 — KeyValue, New, WithDeepCopy, and Read have no model contract; the docs only require providers to be separate modules (docs/design.md:49).


# model gaps 
Model gaps · not findings by themselves
gap No method signatures or return types for Gonfig (docs/design.md:19-38).
gap No definition of the zero value, errors, or helper functions.

# Implemented

- [IC-015] "deep copy" wording replaced with "map-level copy"; nested `map[string]any` levels are copied,
  slices and typed maps are shared. No behavior change, no test change.
  Applied | `gonfig.go:13` — copy-contract copies inbound collections and Raw results at map level;
  `gonfig.go:140-142` — Raw returns a map-level copy, slices and typed maps stay shared, result is never nil;
  `docs/design.md:9-13` — copy-contract term reworded, **map-level copy** term added to the glossary;
  `docs/design.md:74` — Raw returns a map-level copy of the internal map;
  `docs/code.md:11` — Raw() returns a map-level copy.
  Not covered | The inbound "deep-copied" claim in `docs/code.md:11` stays as is; it belongs to the separate
  Major code item about the copy depth of `DeepCopyMap`.
