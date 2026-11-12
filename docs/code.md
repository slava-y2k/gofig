# Methods

`NewWithXXX` — constructor variants; the `XXX` suffix only emulates overloading, which Go does not support.

For locks and thread safety, specify thread safety for XXX if it keeps the explanation compact.

- `Set` — copies on addition and if key found  
- `Delete` — deletes the value by the provided key.

# Copy-contract
Every inbound collection is deep-copied before it is stored, and `Raw()` deep-copies its result;
`Get` is the only method that returns live values, so a returned map stays read-only for the caller.
A method must keep its documented side of the copy-contract when its arguments or results 
are maps or any other struct that can be mutated outside.
