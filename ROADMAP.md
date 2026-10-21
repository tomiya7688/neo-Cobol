# Neo COBOL Roadmap

Neo COBOL is currently in the language-design phase. The roadmap intentionally establishes semantics before optimizing compiler implementation.

## Phase 0 — Project foundation

- [x] Initialize repository documentation.
- [x] Add MIT license.
- [ ] Confirm language philosophy and design principles.
- [ ] Define contribution and coding conventions when implementation begins.

## Phase 1 — Language core

- [ ] Define lexical structure and EBNF grammar.
- [ ] Define type system.
- [ ] Define names, scopes, modules, and classes.
- [ ] Define core statements and expressions.
- [ ] Define diagnostics and invalid-program behavior.

## Phase 2 — Lowering and backend architecture

- [ ] Define COBOL down-compilation subset and lowering rules.
- [ ] Define C backend for GCC / Clang.
- [ ] Define Bitlang transpilation model.
- [ ] Establish a compiler intermediate representation if needed.

## Phase 3 — Minimum viable compiler

- [ ] Implement lexer and parser.
- [ ] Implement semantic analysis.
- [ ] Compile a minimal Hello World program.
- [ ] Add useful, concise diagnostics.

## Phase 4 — Quality infrastructure

- [ ] Add automated test suite.
- [ ] Add conformance tests for language rules.
- [ ] Add backend equivalence tests where practical.
- [ ] Add sample programs.

## Phase 5 — Expansion

- [ ] Expand COBOL compatibility coverage.
- [ ] Expand standard library / runtime facilities.
- [ ] Improve optimization and generated-code quality.
- [ ] Stabilize language specification and versioning.

## Guiding priority

Readability and semantic clarity come before syntactic cleverness. A feature that makes programs significantly harder to understand as English-like descriptions should require a strong technical justification.
