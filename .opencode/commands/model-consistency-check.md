---
description: 
  Check the design model (docs/) and the implementation for mutual consistency via the reviewer agent, 
  and record the findings in docs/inconsistencies.md
---

Run a model-consistency check: the design model in `docs/` versus the Go
implementation.

1. Launch the **reviewer** subagent as a foreground task. 
   The reviewer's own prompt defines the method — the model-definition, design-review, and
   code-implementation-review skills — and the report format. In the launch
   prompt include:
   - the scope: the user's arguments, or the whole project;
   - any specific concern the user raised;
   - the instruction to read `docs/inconsistencies.md` first, for the entry
     format and for the entries that are already open.
2. Present the reviewer's three-part report to the user as received: idea
   consistency, implementation relevance, prepared changes.
3. Apply the prepared changes to `docs/inconsistencies.md`:
   - read the document first, then append each prepared entry exactly as
     written, keeping the `IC-###` numbering;
   - never delete or rewrite an unresolved entry — if the reviewer updated
     an existing entry, apply the update in place;
   - leave the document's Rules and Entry format sections untouched.
4. Summarize for the user: the two verdicts, the findings by severity and
   kind, and which entries were appended or updated.

If the reviewer agent is not available, stop and say so — do not run the
check in this session.
