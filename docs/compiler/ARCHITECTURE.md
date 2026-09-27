# Neo COBOL Compiler Architecture

Status: Bootstrap implementation

## Implementation languages

- Compiler and tooling: **Go**
- Initial generated-code backend: **C11**
- Shared native runtime: **C11**
- Future backends may target COBOL, Bitlang, LLVM IR, WebAssembly, or other targets without changing front-end semantics.

## Pipeline

```text
Neo COBOL source
  -> lexer
  -> parser
  -> canonical AST
  -> semantic analysis / scoped symbol resolution
  -> Neo IR (NIR)
  -> backend
  -> C11 (initial backend)
  -> GCC / Clang / compatible C compiler
  -> native executable
```

NIR is a required boundary. Backends must not depend directly on parser-specific syntax choices.

## Current executable scope

The bootstrap compiler currently supports:

- UTF-8 source input and case-insensitive keywords/identifiers
- COBOL free-format `*>` comments
- elementary level-01 and level-77 declarations
- traditional level-01 records with nested level-02 through level-49 group/elementary items
- COBOL-style `field OF group` and nested `field OF group OF record` qualification
- ambiguity diagnostics for repeated unqualified field names
- canonical built-in `TYPE` names and the initial `PIC` subset
- `VALUE`, symbol resolution, definite initialization, and type/PIC validation
- mutable inferred `VAR` and non-reassignable inferred `LET`
- lexical local scopes and explicit shadowing of an outer binding
- `MOVE` and `DISPLAY`
- Boolean conditions
- canonical English comparisons: `IS EQUAL TO`, `IS NOT EQUAL TO`, `IS GREATER THAN`, `IS LESS THAN`, and the `OR EQUAL TO` forms
- logical `NOT`, `AND`, and `OR`, with `NOT` > `AND` > `OR` precedence
- parenthesized conditions
- `IF ... ELSE ... END-IF` and `END IF` shorthand
- branch-aware definite-assignment merging
- nested `IF` blocks
- NIR condition/block lowering and C11 emission
- record-tree lowering to nested anonymous C11 `struct` values with field-level initialization
- `neoc check`, `emit-c`, `build`, `run`, and `version`

Unsupported syntax fails explicitly. It must never be silently discarded or compiled as a no-op.

## Scope and flow model

Traditional level-01/77 data items and the current top-level procedure bindings share the implemented source-unit scope. Each `IF` branch creates a lexical child scope. A binding declared inside a branch is not visible after `END-IF` or in the sibling branch.

Definite assignment is flow-sensitive across `IF`. For an outer variable that was uninitialized before an `IF`, it is considered initialized after `END-IF` only when both the `THEN` and `ELSE` paths definitely initialize it. An `IF` without `ELSE` cannot make a previously uninitialized outer variable definitely initialized after the block.

## Condition model

Neo COBOL does not use integer truthiness in the implemented condition subset. A value used directly as an `IF` condition must have logical type `BOOLEAN`.

Equality permits equal logical types and compatible numeric types. Ordered comparison is currently limited to numeric types. String equality lowers through `strcmp` in the C11 backend; string ordering is intentionally rejected until language semantics are specified.

`DECIMAL` remains semantically recognized, but exact decimal128 C11 lowering is still intentionally unsupported rather than silently mapped to an imprecise C type.

## Data-model boundary

Traditional level-number records are now represented as a hierarchy rather than flattened variables. A group item has subordinate fields and no independent scalar type in the current executable slice. Elementary leaves carry the normal Neo COBOL logical type/PIC model.

Bare field lookup is permitted only when the field name is unique across visible traditional data items. Repeated names require `OF` qualification, and nested groups use the full outward chain such as `ZIP OF ADDRESS OF CUSTOMER`.

The C11 backend lowers each top-level record to an anonymous `struct`, including nested anonymous structs for subordinate groups. Individual field initializers remain field assignments so existing definite-initialization analysis stays field-sensitive.

Whole-record values/copies, level-88 conditions, `OCCURS`, `REDEFINES`, and layout-sensitive interoperability remain outside this slice.

## Next implementation slices

1. `OCCURS` fixed arrays and dynamic-sequence groundwork.
2. Level-88 condition names and additional traditional data-description entries.
3. Expanded `PERFORM` family and loops.
4. Functions, classes, interfaces, inheritance, and closures.
5. COBOL and Bitlang backends.
6. Expanded runtime and backend-equivalence/conformance tests.
