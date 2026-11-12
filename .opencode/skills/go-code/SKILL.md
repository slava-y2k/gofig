---
name: Go code
description: Use when creating or editing go code in compilable go files.
---


## Code style rules
- Use explicit type declarations where it is possible.
- Naming: `MixedCaps`, no underscores.
- Avoid package/name stutter at call sites; prefer `gonfig.New()` over `gonfig.NewGonfig()`.
- Specify types explicitly for variables, except in obvious cases where the type can be inferred from the line.
- Use 2-3 characters for the receiver name.
- Don't use single characters for arguments and variable names; use more meaningful names.
- Errors: return `error` as the last result. Wrap with `fmt.Errorf("...: %w", err)`.
  Use sentinel errors (`var ErrUnsupportedData = errors.New(...)`) for expected
  conditions; never `panic` for recoverable failures.


### Code blocks formatting
Declare types and assignments in one block and their usage in another block, separating the blocks
with a blank line. The comments in the code snippet are for you; don't put them into the edited code.
The declaration and usage blocks should be separated by an empty line.
```
// --- declaration and assignment block ---
var variable1 VariableType1
var err error
variable1, err = keyValue.Read()
if err != nil {
    t.Fatalf("Read() error = %v, want nil", err)
}

// --- actual usage block ---
firstValues["a"].(map[string]any)["b"] = 99
firstValues["new"] = true
```
