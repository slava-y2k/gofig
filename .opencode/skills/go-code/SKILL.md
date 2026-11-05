---
name: Go code 
description: Use when creating or editing .go code files.
---



## Code style rules
- use explicit type declaration where it is possible.
- Naming: `MixedCaps`, no underscores.
- Avoid package/name stutter at call sites; prefer `gofig.New()` over `gofig.NewGofig()`.
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


## Documentation


<llm-debug> file for save is go-code.YYYY-MM-DD_HH:mm:ss.sss.md</llm-debug>

Documentation is the text description which should explain why and what 
for this entity created. What is the primary idea and describe the abstraction is used.     


### Documentation context
Form a documentation context from the following sources:
- method instructions itself
- other methods which are invoked internally, keep stack depth equal to 3 methods
- model information this. The information how it actually should work

<llm-debug>wriite the documentation context in a separate section </llm-debug>


### Documentation structure
Follow the code structure in crating document context: package -> interfaces/structs -> methods/constants
This means that approaches which are applicable to the whole entity should be specified on the entity level. 

<llm-debug>wriite the structured context </llm-debug>


### Documentation writing
For each code item: package class struct create a documentation depending on the
cardinality of the underlying encapsulated logic in 1-3 sentences depending on the code footprint.

<llm-debug>for each entity write the documentation with the figured out structure</llm-debug>

### Special cases

There is no method overloading in go, so to emulate overload a special suffixes are added.
For overloaded methods use single description which vary only with provided args 
