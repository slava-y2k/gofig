---
name: Go unit tests
description: Use when creating, modifying, removing unit tests for the go language
---

# Unit testing context
Unit test is a test that checks the result of method.
Unit is a go file with the same receiver.

To create a find all the stacktrace with in the module, then identify the business 
logic footprint for each method.

For every method identify pre conditions, post conditions and type invariant, methods 
which make invariant incorrect for a while with in a method, or a stack trace should be checked.
Pre and post conditions maybe checked once or omitted

You have no need to every time test all the hidden (encapsulated methods) one test is enough.
So consider the example (Method1 calls method2, method2 calls method3 ) and Method21 calls method2, 
and it calls method3.

So it's enough to create only one pack of tests for Method1, test for Method21 should only cover what is not covered in test for Method1


