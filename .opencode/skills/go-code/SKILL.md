---
name: Go code 
description: Use when creating or editing go code in compilable go files.
---



## Code style rules
- use explicit type declaration where it is possible.
- Naming: `MixedCaps`, no underscores.
- Avoid package/name stutter at call sites; prefer `gonfig.New()` over `gonfig.NewGonfig()`.
- specify types explicitly for the variables, except obvious cases, when type can be inferred from the line 
- As the receiver use 2-3 chars.
- Don't use single char for arguments and variables naming, use more meaningful names
- Errors: return `error` as the last result. Wrap with `fmt.Errorf("...: %w", err)`.
  Use sentinel errors (`var ErrUnsupportedData = errors.New(...)`) for expected
  conditions; never `panic` for recoverable failures.

## Code logic tips



 

### Code blocks formatting
Declare types and assignment in one block and usage of them in another separating with space comments 
in the code snippet is for you shouldn't put them into editing code
Declaration and usage block should be separated by empty line
```	
// --- declaration and assignment block----
var variable1 VariablType1
var err error
variable1, err = keyValue.Read()
if err != nil {
    t.Fatalf("Read() error = %v, want nil", err)
}

//--- actual usage block ----
firstValues["a"].(map[string]any)["b"] = 99
firstValues["new"] = true
```
