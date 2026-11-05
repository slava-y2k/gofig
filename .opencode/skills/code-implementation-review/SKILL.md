 ---
name: code-implementation-review
description: Reviews the Go implementation against the model inventory — contract match, layout rules, behavior guarantees, tests, and build health. Use it as the final step of a model-consistency check, after the model inventory exists.
---

# Code Implementation Review

## Purpose

Verify the implementation is relevant to the model: 
it implements everything the model defines, nothing the model does not define, 
and it behaves the way the model promises. 
The input is the model inventory from the model-definition skill.

## Checks

1. **Contract match** — every entity in the inventory exists in the code with the exact contract: 
   type name, file, kind, and method signatures. Public surface beyond the inventory 
   is undocumented code; an inventory item without code is an unimplemented model.
2. **Layout match** — the repository follows the layout rules: module
   structure and import paths, `go.work` registration, one module per
   provider, `replace` directives for standalone builds.
3. **Semantics match** — behavior the model promises is behavior the code
   guarantees: deep copies at boundaries, nil normalization, concurrency
   safety where claimed, sentinel errors, fluent returns.
4. **Test match** — behavior the model promises is covered by tests, and
   the compile-time interface assertions the conventions require exist.
5. **Build health** — verify, from the repository root:

   ```sh
   go build github.com/slava-y2k/gonfig/...
   go vet github.com/slava-y2k/gonfig/...
   go test github.com/slava-y2k/gonfig/...
   ```

   Use the module-path pattern, never `./...` (directory-rooted patterns
   skip the nested workspace modules). Record the commands and results.

## Classification

| Kind | Meaning |
| --- | --- |
| `undocumented-code` | Public code surface the model does not define |
| `unimplemented-model` | Model item the code does not provide |
| `contradiction` | Both exist but disagree: signature, behavior, file, or kind |

## Output

```md
## Implementation-relevance findings
- Verdict: Faithful | Deviates
- Verification: <commands run and their results>
- <severity> — <kind> — file.go:N vs docs/…:M — what diverges and how
```

Severity: `Critical` — the documented contract is broken (wrong signature,
missing entity, failing build); `Major` — a real divergence, one side must
change; `Minor` — placement or naming drift that breaks a layout rule.

## Rules

- The review is exhaustive over the model, not a spot check: every
  inventory item is checked. Findings cite both sides — the `file:line` and
  the model reference.
- Run the verification commands; a verdict without them is invalid.
- Style is out of scope unless the model itself demands it — a convention
  the docs state is part of the model.
- Report model-versus-code facts, not code-quality opinions.
