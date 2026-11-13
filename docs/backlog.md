# Backlog

- implement "gonfig tag"
- Prepare the provider list for the implementation.
- In Get method bool return is not required

# Code items 
Major	contradiction	gonfig.go:97-111 vs docs/design.md:34 — the model describes Delete as "set value by provided key"; the implementation removes the leaf (gonfig.go:107, deletePath at gonfig.go:191-200).
Minor	undocumented-code	gonfig.go:9-10 — ErrNilProvider is public and returned on a nil provider (gonfig.go:35-37); the model defines no error value.
Minor	undocumented-code	gonfig.go:12-22 — the usable zero value and the empty New() constructor are not in the model; only generic NewXXX constructors are mentioned (docs/design.md:19).
Minor	undocumented-code	gonfig.go:25-57 — concrete constructors and their semantics (nil input to empty, provider error propagation, snapshot of another container) are undefined; the model has only the NewXXX placeholder.
Minor	undocumented-code	gonfig.go:61 — Get returns (any, bool); the found flag is not part of the documented contract (docs/design.md:32).
Minor	undocumented-code	gonfig.go:77-142 — fluent returns (*Gonfig) and the deep-merge rule ("later sources win", recursively) are not documented (docs/design.md:21-24 only says merge "collect[s] key-values").
Minor	undocumented-code	gonfig.go:151-162 — Raw never returns nil; the model says only "returns a deep copy of the internal map" (docs/design.md:72).
Minor	undocumented-code	gonfig_util.go:4 — exported DeepCopyMap is not defined by the model.
Minor	undocumented-code	provider/keyvalue/keyvalue.go:8-34 — KeyValue, New, WithDeepCopy, and Read have no model contract; the docs only require providers to be separate modules (docs/design.md:49).


# model gaps 
Model gaps · not findings by themselves
gap No method signatures or return types for Gonfig (docs/design.md:19-38).
gap No definition of the zero value, errors, or helper functions.

# Implemented

- [IC-015] "deep copy" wording replaced with "map-level copy"; nested `map[string]any` levels are copied,
  slices and typed maps are shared. No behavior change, no test change.
  Superseded | by `[Get-copy]` below — the direction returned to a full deep-copy contract, so this
  map-level wording is history, not the current state.
  Applied | `gonfig.go:13` — copy-contract copies inbound collections and Raw results at map level;
  `gonfig.go:140-142` — Raw returns a map-level copy, slices and typed maps stay shared, result is never nil;
  `docs/design.md:9-13` — copy-contract term reworded, **map-level copy** term added to the glossary;
  `docs/design.md:74` — Raw returns a map-level copy of the internal map;
  `docs/code.md:11` — Raw() returns a map-level copy.

- [Get-copy] `Get` now returns a full deep copy instead of a live reference, so a caller that keeps a result
  cannot mutate the container's internal map through it.
  Applied | `gonfig_util.go:8-138` — `DeepCopyMap` deep-copies maps, slices, arrays, pointers, and structs via
  reflection (`deepCopyValue`/`deepCopyReflect`); un-copyable kinds pass through unchanged;
  `gonfig.go:59-76` — `Get` returns `DeepCopyMap(kvMap)` with no parts and `deepCopyValue(value)` on a path hit;
  `gonfig.go:175-199` — `setPath` deep-copies any leaf, not just `map[string]any`;
  `gonfig.go:13-14, :149-150` — `Gonfig` and `Raw` doc comments state the deep-copy contract;
  `gonfig_test.go` — `TestGetWholeConfigReturnsCopy`, `TestGetMapValueReturnsCopy`, `TestGetSliceReturnsCopy`,
  `TestGetTypedMapReturnsCopy` and the `DeepCopyMap` slice/typed-map/pointer/struct cases replace the
  live-reference and sharing pins.

- "empty key → nil" wording fixed: the access-methods prose now says a call with no key parts targets the
  root — `Get()` returns the whole map with `true`, `Set` replaces the root (maps only), and `Delete`
  clears it; `Set`/`Delete` return the `*Gonfig` receiver, never nil.
  Applied | `docs/design.md:78-80` — access-methods prose rewritten; `docs/design.md:87-88` — the
  `Set`/`Delete` bullets now state that they return the `*Gonfig` receiver.

- "overloaded methods" wording fixed: the `NewWithXXX` family is now described as an emulation of method
  overloading — Go has no overloading, so each variant is a distinct, separately named function.
  Applied | `docs/code.md:3` — the `NewWithXXX` line now reads "an emulation of method overloading";
  `docs/design.md:67` — aligned to the same wording.
