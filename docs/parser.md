# Configuration Model

Mermaid class diagram of the gofig configuration model.

Value-returning methods (`Read`, `Parse`, `Get`) also return a second result —
an error, and for `Get` a found flag — omitted from the signatures below to
keep the diagram readable.

```mermaid
---
title: gofig — Configuration Model
---
classDiagram
    direction LR

    namespace gofig {
        class Provider {
            <<interface>>
            +Read() map[string]any
        }

        class Gofig {
            -kvMap map[string]any
            +New() Gofig
            +Load(provider Provider) error
            +Merge(other Gofig) Gofig
            +MergeMap(keyValues map[string]any) Gofig
            +Set(parts []string, value any) Gofig
            +Delete(parts ...string) Gofig
            +Get(parts ...string) any
            +Raw() map[string]any
        }
    }

    namespace provider_keyvalue {
        class KeyValue {
            -keyValues map[string]any
            -deepCopy bool
            +WithDeepCopy(isEnabled bool) KeyValue
            +Read() map[string]any
        }
    }

    class Parser {
        <<contract>>
        +Parse(keyValues map[string]any) map[string]any
    }

    Provider <|.. KeyValue : implements
    Gofig ..> Provider : reads in Load()

    note for Parser "A Parser only returns the key-value pairs (map[string]any). Defined in the model; implementations will be standalone modules under parser/."
    note for Gofig "Key-value pair manager: fluent interface — Merge, MergeMap, Set, Delete return the Gofig itself."
```
