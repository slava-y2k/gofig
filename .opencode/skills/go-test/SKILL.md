---
name: Go unit tests
description: Use when creating, modifying, removing unit tests for the go language
---

# Unit testing context
Figure out the methods that implements the main responsibility of the class.
To figure out such methods build a stacktrace and find the root underlying methods which are 
implement the business logic 


# Unit test approach 

For every method identify pre conditions, post conditions and type invariant, 
only cover what is not covered  by other tests in the same module 

On assertions follow the code structure, and methods, including unexported should be 
in the relation to method they cover.
E.g. if two different methods calls same method and the result of the upper methods contains 
the result of underlying, the underlying result should be covered only by one assert.
                                                    

# Unit test formating
don't explain the test suits 
