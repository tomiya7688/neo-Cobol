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

Only elementary level-01 and level-77 items are executable in the current slice. Levels 02-49 are recognized but rejected until record hierarchy semantics are implemented, avoiding accidental flattening of COBOL records.

## Next implementation slices

1. Record hierarchy and additional traditional data-description entries.
2. `OCCURS` and dynamic-array lowering.
3. Expanded `PERFORM` family and loops.
4. Functions, classes, interfaces, inheritance, and closures.
5. COBOL and Bitlang backends.
6. Expanded runtime and backend-equivalence/conformance tests.
