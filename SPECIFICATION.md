# Neo COBOL Language Specification

Status: Draft

## 1. Purpose

Neo COBOL is a modern programming language that preserves the readability-oriented philosophy of COBOL while providing a contemporary compiler architecture and multiple compilation targets.

## 2. Design requirements

Neo COBOL should:

1. Read naturally as English where practical.
2. Prefer explicit intent over symbolic terseness.
3. Support deterministic parsing and tooling.
4. Be suitable for static analysis and automated review.
5. Support down-compilation to COBOL for a defined compatibility subset.
6. Support a native compiler pipeline.
7. Support transpilation to Bitlang.
8. Support a C backend suitable for GCC and Clang.

## 3. Readability acceptance criterion

A proposed syntax feature should normally be rejected or redesigned when it makes valid source code substantially harder to read as an English description of program intent without providing a compelling technical benefit.

This criterion does not require grammatically perfect English. Compiler determinism, unambiguous grammar, and implementation safety take precedence when natural-language phrasing would introduce ambiguity.

## 4. Compatibility model

Compatibility with historical COBOL will be defined explicitly rather than assumed.

The specification will distinguish at least:

- Neo COBOL features that can be lowered directly to COBOL.
- Neo COBOL features that can be lowered through generated support code.
- Neo COBOL features that are native-only and cannot be represented faithfully in the chosen COBOL target.

## 5. Compilation targets

Planned targets:

- COBOL
- C
- Bitlang
- Native Neo COBOL compiler pipeline

## 6. Language areas to specify

The following sections are intentionally incomplete and will be expanded through design Issues:

- Lexical structure
- Grammar and EBNF
- Names and scopes
- Type system
- Variables and constants
- Expressions
- Control flow
- Procedures and functions
- Classes and modules
- Error handling
- I/O model
- Data description
- Interoperability
- COBOL lowering rules
- C lowering rules
- Bitlang lowering rules
- Diagnostics
- Undefined / implementation-defined behavior policy

## 7. Example direction

The canonical syntax is not yet frozen. Examples added before grammar stabilization are illustrative rather than normative.

## 8. Versioning

The specification will use explicit revisions once the first grammar is stabilized. Until then, this document is a working draft.
