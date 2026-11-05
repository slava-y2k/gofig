
# Model

There main entities are:
 * Provider - main supplier of key-value collection.
 * Gonfig - main struct with number of methods that provides a management functionality for the configuration
 * gonfig tag - marks the struct field


## Key-value collection
Key value collection is represented by `key:string, value:any`. 
**Note!** value may be also a nested map with type `key:any, value:any`, in this case keys 
will be normalized to a string type.

## Provider
An interface that has only one simple method `Read() (map[string]any, error)`. To return 
a copy or a reference form the map is the providers responsibility depending on its own. 


## Gonfig
Main key-value collection manager, which server for the following responsibilities for key-value collection:
* New<WithXXX> - where XXX is the source suffix - a constructor methods, which performs a deep copy 
for the provided key-value collection
* merge - collect key-values from different sources
* writing - write to different type of sources


## Gonfig implementation details
Gonfig provides bound safety. All methods that accepts a key-value collection creates an internal copy.
same approach is applied to the outbound methods - every method makes a copy.


# Architecture

Library is implemented with modular approach: 
* Every provider is a separate module, which is included only when it is needed.
* Core module has only a std lib dependency
  

# Diagram

```mermaid
---
title: Gonfig library design
---
classDiagram
    direction LR

    namespace Gonfig_Lib {
        class Gonfig["Gonfig"] {
            -kvMap map[string]any
            +NewWithXXX()$ error, Gonfig
            +MergeXXX(other Gonfig) Gonfig
            +Set(parts []string, value any) Gonfig
            +Delete(parts ...string) Gonfig
            +Get(parts ...string) any
            +Raw() map[string]any
        }

        class IProvider["Provider"] {
            <<interface>>
            +read() map[string, any]
        }
    }

    namespace Provider {
        class Provider1
        class Provider2
    }

    Gonfig --> IProvider
    IProvider <|-- Provider1
    IProvider <|-- Provider2
```
