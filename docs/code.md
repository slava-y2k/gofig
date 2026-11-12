# Methods

`NewWithXXX` — constructor variants; the `XXX` suffix only emulates overloading, which Go does not support.

For locks and thread safety, specify thread safety for XXX if it keeps the explanation compact.

- `Set` — copies on addition and if key found  
- `Delete` — deletes the value by the provided key.

# Copy-contract

**Copy-contract** — every accepted or returned item should be a deep copy, to prevent implicit
modifications.
**Deep copy** — if a value is a map or a reference type, a new copy is created for it.

`Get` and `Raw` return deep copies, so a returned map can be modified without affecting the container.

A method must keep its documented side of the copy-contract when its arguments or results are maps or any
other struct that can be mutated outside.
