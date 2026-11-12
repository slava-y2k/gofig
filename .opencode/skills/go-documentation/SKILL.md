---
name: Go code documentation
description: Use when documenting go code in compilable go files
---


# Documentation

Code documentation (in this skill, documentation refers to code documentation/docs unless stated otherwise, e.g.
model documentation) is an additional textual explanation that describes what happens in the code.
Documentation should not repeat terms from the code.
Documentation should only explain why it was implemented in the following form or what the tricky case is.
Prefer the declarative form and statements.
```
struct is [explanation], method does [explanation]
```
Use as little explanation as possible: every documentation line should have a proof to be written; otherwise it
should be omitted.
Proof means there is no source other than documentation where this information can be obtained — neither code
nor model documentation, test documentation, or other project documentation.

Do not add unnecessary line breaks and keep lines within 110 characters.

If a struct implements an interface entirely (the whole contract), specify the implemented interface in the
documentation at the struct level.
If the implementation is partial, specify the implemented methods in brackets.


<llm-debug> For each step save the found context into the file go-documentation.YYYY-MM-DD_HH-mm-ss.sss.md where
ss - seconds
sss - milliseconds</llm-debug>


# Documentation context
Form a documentation context from the following sources:
- the method instructions themselves
- other methods that are invoked internally, keeping the stack depth equal to 3 methods
- model information: how it actually should work


# Documentation structure
Follow the code structure when creating the documentation context: package -> interfaces/structs -> methods/constants.
This means that approaches that are applicable to the whole entity should be specified on the entity level.

# Documentation preparation
For each code item — package, class, or struct — create documentation of one to three sentences,
depending on the amount of underlying encapsulated logic and the code footprint.

# Documentation compaction
Find duplicated context between the code and the comment, and remove this text from the documentation.
Always document structs and interfaces. If a comment is omitted for a public member,
write a very short text from the context.
<llm-debug>Specify the sources where this description can be obtained, code or model description.</llm-debug>

# Documentation writing
Write the documentation for items found in code.

# Special cases
There is no method overloading in Go, so to emulate overloading, special suffixes are added.
For overloaded methods, use a single description that varies only by the provided arguments.
