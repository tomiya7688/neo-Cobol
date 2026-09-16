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
  -> semantic analysis / symbol resolution
  -> Neo IR (NIR)
  -> backend
  -> C11 (initial backend)
  -> GCC / Clang / compatible C compiler
  -> native executable
```

NIR is a required boundary. Backends must not depend directly on parser-specific syntax choices.

## Repository layout

```text
cmd/neoc/             neoc CLI
internal/lexer/       source -> tokens
internal/parser/      tokens -> AST
internal/ast/         canonical syntax tree
internal/sema/        symbols, definite initialization, type/PIC validation
internal/types/       canonical logical types and representation metadata
internal/nir/         backend-independent Neo IR and lowering
internal/backend/c/   initial C11 backend
runtime/c/            shared C11 runtime boundary
docs/language/        language specification
docs/compiler/        implementation architecture and tooling docs
examples/             executable examples
tests/                cross-package fixtures / future conformance suites
tools/                repository-local development tools
```

Language specification and implementation details are intentionally separated. `docs/language/` defines what Neo COBOL means; compiler packages define one implementation of those rules.

## Current executable scope

The bootstrap compiler currently supports:

- UTF-8 source input
- case-insensitive keywords and identifiers
- COBOL free-format `*>` comments
- single- and double-quoted strings with doubled-quote escaping
- decimal and explicit `0b` / `0o` / `0x` integer literals
- optional `IDENTIFICATION DIVISION.`, `DATA DIVISION.`, and `PROCEDURE DIVISION.` headers
- optional `WORKING-STORAGE SECTION.` header
- `PROGRAM-ID.` metadata
- elementary level-01 and level-77 declarations
- canonical built-in `TYPE` names
- initial `PIC X(...)`, `9(...)`, `S9(...)`, and implied-decimal `V` forms
- PIC-only inference into `STRING`, `INTEGER`, `LONG`, or `DECIMAL`
- `VALUE` with string, numeric, or Boolean literals
- case-insensitive symbol resolution and duplicate-name diagnostics
- definite-initialization checks for the implemented statement flow
- static type compatibility checks and literal PIC-bound checks
- `MOVE value TO variable`
- `DISPLAY` with literals and declared variables
- NIR lowering
- C11 emission
- `neoc check`, `emit-c`, `build`, `run`, and `version`

Unsupported syntax fails explicitly. It must never be silently discarded or compiled as a no-op.

## Types and PIC

Logical type identity and storage/display representation are separate compiler concepts. `INTEGER`, `DECIMAL`, `STRING`, `BOOLEAN`, `FLOAT`, `BYTE`, `LONG`, and `DOUBLE` are canonical logical built-ins. `PIC` is represented as a compatible representation constraint attached to a logical type, not as the type itself.

The semantic layer can validate the initial `DECIMAL`/PIC model, but exact decimal128 C11 lowering is intentionally not faked: the C backend currently reports `DECIMAL` as unsupported. The initial native backend can lower `STRING`, `BOOLEAN`, `BYTE`, `INTEGER`, and `LONG`; `FLOAT`/`DOUBLE` backend work remains incomplete.

Fixed-width string and numeric PICs currently provide compile-time representation constraints for literal assignments. Full COBOL field padding/truncation and runtime PIC conversion semantics are not implemented yet.

`VAR` and `LET` remain declarations/inference features to implement; they are not primitive runtime types.

## Data-model boundary

Only elementary level-01 and level-77 items are executable in the current slice. Levels 02-49 are recognized as level numbers but rejected until record hierarchy semantics are implemented. This avoids incorrectly flattening nested COBOL records.

Traditional declarations must precede executable statements in the current source-unit implementation.

## Runtime boundary

The bootstrap C backend lowers implemented `DISPLAY` behavior through portable C stdio. `runtime/c/` establishes the boundary for behavior that later needs shared runtime support. Backend-specific mechanics must not redefine language semantics.

## CLI

```text
neoc check source.ncob
neoc emit-c [-o file.c] source.ncob
neoc build [-o executable] source.ncob
neoc run source.ncob
neoc version
```

`build` uses `$CC` when set; otherwise it searches `cc`, `clang`, then `gcc`.

## Next implementation slices

1. `VAR` / `LET` inferred local bindings and mutability checks.
2. Conditions and `IF`.
3. Record hierarchy and additional traditional data-description entries.
4. `OCCURS` and dynamic-array lowering.
5. Functions, classes, interfaces, inheritance, and closures.
6. COBOL and Bitlang backends.
7. Expanded runtime and backend-equivalence/conformance tests.
