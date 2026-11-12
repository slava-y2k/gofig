# Roadmap

- implement "gonfig tag"
- Prepare the provider list for the implementation.
- In Get method bool return is not required
- [IC-015] Raw "deep copy" wording ignores shared slices and typed maps
  Severity| Minor, Kind| contradiction, Status| open, 
  Model says | gonfig.go:140 — "Raw returns a deep copy per the copy-contract"; docs/design.md:72 — "returns a deep copy of the internal map"; docs/code.md:11 — "Raw() deep-copies its result".
  Reality | DeepCopyMap copies map[string]any levels only and shares every other value (gonfig_util.go:3-19); a slice obtained from Raw still points at the container's backing array (gonfig_test.go:678, TestRawSliceSharesBackingArray), and typed maps stay shared (gonfig_test.go:748, TestDeepCopyMapSharesTypedMap).
  Prepared change | Reword gonfig.go:140, docs/design.md:72 and docs/code.md:11 to "map-level copy": nested map[string]any levels are copied, slices and typed maps are shared.
