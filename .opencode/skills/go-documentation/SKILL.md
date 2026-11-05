---
name: Go code documentation
description: Use when documenting go code in compilable go files
---


# Documentation

Code documentation (in this skill documentation refers to code documentation/docs, if there is no special meaning, e.g.
model documentation and etc) is an additional text explanation that describes what happens here to code.
Documentation should not repeat the code terms. 
Documentation should only explain why it was implemented in the following form or what is the tricky case
Prefer declarative form and statements.
```
struct is [explanation], method does [explanation] 
```
Better to use less explanation. So every document line should have a proof to be written otherwise it should be omitted.
Proof - means there is no any source except documentation where this information can be obtained neither code 
nor model documentation, test documentation and other project documentation.

Do not do unnecessary line breaks and keep line to 110 symbols length. 


<llm-debug> For each step save the found context into the file go-documentation.YYYY-MM-DD_HH-mm-ss.sss.md where
ss - seconds
sss - milliseconds</llm-debug>


# Documentation context
Form a documentation context from the following sources:
- method instructions itself
- other methods which are invoked internally, keep stack depth equal to 3 methods
- model information this. The information how it actually should work


# Documentation structure
Follow the code structure in crating document context: package -> interfaces/structs -> methods/constants
This means that approaches which are applicable to the whole entity should be specified on the entity level.

# Documentation preparation
For each code item: package class struct create a documentation depending on the
cardinality of the underlying encapsulated logic in 1-3 sentences depending on the code footprint.

# Documentation compaction
Find duplication context between the code and comment, and remove these text from the
documentation.
Always document structs and interfaces. If comment is omitted to public member, 
then write very short test from the context. 
<llm-debug>Specify the sources where this description can be obtained, code or model description.</llm-debug>

# Documentation writing
Write the documentation for items found in code 

# Special cases
There is no method overloading in go, so to emulate overload a special suffixes are added.
For overloaded methods use single description which vary only with provided args 
