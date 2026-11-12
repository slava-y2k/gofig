---
name: Go code
description: Use when creating or editing go code in compilable go files.
---


# Code style rules
- Use explicit type declarations where it is possible.
- Naming: `MixedCaps`, no underscores.
- Avoid package/name stutter at call sites. Constructors are `NewWithXXX` (`NewWithMap`, `NewWithInt`,
  `NewWith<SmthElse>`).
- Specify types explicitly for variables, except in obvious cases where the type can be inferred from the line.
- Use 2-3 characters for the receiver name.
- Don't use single characters for arguments and variable names; use more meaningful names.
- Errors: return `error` as the last result. Wrap with `fmt.Errorf("...: %w", err)`.
  Use sentinel errors (`var ErrUnsupportedData = errors.New(...)`) for expected
  conditions; never `panic` for recoverable failures.


## Code blocks formatting
Declare types and assignments in one block and their usage in another block, separating them
with a blank line. The comments in the code snippet are for you; don't put them into the edited code.
```
// --- declaration and assignment block ---
var firstValues map[string]any
var err error
firstValues, err = loadValues(path)
if err != nil {
    t.Fatalf("loadValues() error = %v, want nil", err)
}

// --- actual usage block ---
firstValues["a"].(map[string]any)["b"] = 99
firstValues["new"] = true
```


## Method overloading

Go has no method overloading, so this limitation is overcome by adding a suffix to the method name;
the suffix should be 1-2 words long. In the patterns below, XXX is the overload suffix that mirrors
the argument signature.
- For constructors, the typical signature is `NewWith<XXX>`.
- For methods, use the pattern `MethodName<XXX>`.
