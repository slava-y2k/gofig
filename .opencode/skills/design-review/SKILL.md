---
name: design-review
description: Reviews the design docs for internal consistency — cross-document contradictions, diagram-versus-prose mismatches, stale references, and glossary discipline. Use it as the second step of a model-consistency check, after the model inventory exists.
---

# Design Review

## Purpose

Check the model itself: is the idea coherent? The design docs are reviewed
against each other, never against the code. The input is the model
inventory from the model-definition skill.

## Checks

For every check, cite `docs/…:line` for both sides of every finding.

1. **Glossary discipline** — `docs/design.md` promises exactly one term per
   entity. Flag synonyms (two names for the same entity) and one name used
   for two different things.
2. **Cross-document agreement** — a contract quoted in several docs must be
   identical everywhere: same method names, same signatures, same behavior
   claims. Compare `design.md`, `code.md`, `_index.md`,
   `roadmap.md`, and `README.md` pairwise.
3. **Diagram versus prose** — every class, method, and relation in the
   Mermaid diagrams must match the prose. Flag members that appear in only
   one of the two.
4. **Stale references** — files, directories, types, or commands a doc
   mentions that do not exist in the repository.
5. **Planned versus done** — planned work must be marked as planned. Flag
   prose that presents planned features as existing, or existing features
   as planned.
6. **README honesty** — README claims (usage, commands, module paths) must
   match what the other docs say.
7. **Model gaps** — revisit the gaps recorded in the inventory: each one is
   either harmless wording or an idea inconsistency (a contradiction, a
   missing definition that other docs rely on).

## Output

```md
## Idea-consistency findings
- Verdict: Consistent | Inconsistent
- <severity> — docs/a.md:N vs docs/b.md:M — what contradicts what
```

Severity: `Critical` — a contradiction that misleads a reader into building
the wrong thing; `Major` — a real disagreement that must be fixed on one
side; `Minor` — naming drift and stale references.

## Rules

- Compare docs with docs only. Implementation mismatches belong to the
  code-implementation-review skill.
- Every finding names both documents and their lines. A "feels wrong"
  without a quote is not a finding.
- A model gap becomes a finding only when two docs actively disagree or one
  contradicts itself; open questions are listed, not flagged.
