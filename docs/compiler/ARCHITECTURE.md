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
  -> semantic analysis
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
internal/sema/        semantic analysis
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

## Bootstrap scope

The first executable slice supports:

- UTF-8 source input
- case-insensitive keywords
- COBOL free-format `*>` comments
- single- and double-quoted strings with doubled-quote escaping
- optional `IDENTIFICATION DIVISION.` and `PROCEDURE DIVISION.` headers
- `PROGRAM-ID.` metadata
- `DISPLAY` with string and numeric literal operands
- NIR lowering
- C11 emission
- `neoc check`, `emit-c`, `build`, `run`, and `version`

Unsupported syntax fails explicitly. It must never be silently discarded or compiled as a no-op.

## Types and PIC

Logical type identity and storage/display representation are separate compiler concepts. `INTEGER`, `DECIMAL`, `STRING`, `BOOLEAN`, `FLOAT`, `BYTE`, `LONG`, and `DOUBLE` are canonical logical built-ins. `PIC` is represented as a compatible representation constraint attached to a logical type, not as the type itself.

`VAR` and `LET` will be modeled as declaration/inference semantics rather than primitive type identities when their detailed grammar is finalized.

## Runtime boundary

The bootstrap C backend currently lowers `DISPLAY` through portable C stdio. `runtime/c/` establishes the boundary for behavior that later needs shared runtime support. Backend-specific mechanics must not redefine language semantics.

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

1. Data declarations and symbol tables.
2. Canonical Neo COBOL type checking and PIC validation.
3. `MOVE` and variable-backed `DISPLAY`.
4. Conditions and `IF`.
5. Functions, classes, interfaces, inheritance, and closures.
6. COBOL and Bitlang backends.
7. Expanded runtime and backend-equivalence/conformance tests.
