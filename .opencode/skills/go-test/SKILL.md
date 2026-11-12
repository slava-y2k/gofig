---
name: Go unit tests
description: Use when creating, modifying, removing unit tests for the go language
---

# Unit testing context
Figure out the methods that implement the main responsibility of the class.
To figure out such methods, build a stack trace and find the root underlying methods that
implement the business logic


# Unit test approach

For every method, identify preconditions, postconditions, and type invariants;
cover only what is not covered by other tests in the same module.

For assertions, follow the code structure; methods, including unexported ones, should be
in relation to the method they cover.
E.g. if two different methods call the same method and the result of the upper methods contains
the result of the underlying one, the underlying result should be covered by only one assertion.


# Unit test formatting
Don't explain the test suites.
