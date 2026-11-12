---
name: model-definition
description: Extracts the canonical model definition from the project's design docs into a structured inventory of terms, entities, contracts, and rules. Use it as the first step of a model-consistency check, or whenever the authoritative statement of what the project is supposed to be is needed.
---

# Model Definition

## Purpose

Build the model inventory: the single authoritative statement of what the
gonfig project is supposed to be, taken from the design docs. The inventory
is the reference against which both the docs themselves (design review) and
the implementation (code implementation review) are checked.

## Sources

Read them in this order; earlier sources win when they overlap:

| Source | Defines |
| ------ | ------ |
| `docs/design.md` | Terms (glossary), the data model, entity contracts, data flow |
| `docs/code.md` | Code conventions and API notes |
| `docs/_index.md` | Directory tree, documents index, generated-output policy |
| `docs/roadmap.md` | Planned work |
| `README.md` | Public self-description: what the library is, how it is used |
| `AGENTS.md` | Working rules for the repository |

## Process

1. Read every source in full.
2. Extract, with a `file:line` reference for every item:
   - **Glossary** — one term per entity. `docs/design.md` promises exactly one
     term per entity; note every synonym you meet.
   - **Entities and contracts** — for each type: name, file, kind (struct,
     interface, module), and its public contract (method names and exact
     signatures).
   - **Data flow** — how data moves between the entities.
   - **Layout rules** — module structure, import paths, workspace and
     versioning rules.
   - **Conventions** — interface size, data crossing boundaries as
     `map[string]any`, compile-time interface checks, concurrency guarantees.
   - **Planned items** — everything the docs explicitly mark as planned or
     not yet implemented.
3. Record **model gaps**: questions the docs leave open, ambiguous wording,
   and requirements stated in one doc but absent from the others. Gaps are
   input for the design review, not findings by themselves.

## Output

A single inventory block:

```md
## Model inventory
- Glossary: ... (docs/design.md:N)
- Entities: <name> — <kind>, <file>; contract: <signatures> (docs/…:N)
- Data flow: ... (docs/…:N)
- Layout rules: ... (docs/…:N)
- Conventions: ... (docs/…:N)
- Planned: ... (docs/…:N)
- Model gaps: ...
```

## Rules

- The docs are the source of truth for the model. Quote them; do not
  paraphrase contracts — copy signatures exactly.
- Distinguish what a doc asserts from what it plans: planned items go to
  the Planned section, not to the contracts.
- Do not consult the implementation while building the inventory. The
  inventory states the intent; comparing it with the code is the next
  skill's job.
