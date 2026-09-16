# Neo COBOL Language Specification

Status: Draft

## 1. Purpose

Neo COBOL is a modern programming language that preserves the readability-oriented philosophy and recognizable syntax of COBOL while extending the language with contemporary type-system, OOP, callable, tooling, and compiler features.

Neo COBOL is intended to feel like a newer evolution of COBOL, particularly COBOL 2002, rather than an unrelated language with COBOL-inspired keywords.

## 2. Design requirements

Neo COBOL should:

1. Read naturally as English where practical.
2. Preserve established COBOL syntax by default.
3. Prefer explicit intent over symbolic terseness.
4. Prefer COBOL-like or full English terminology when multiple spellings have the same meaning.
5. Modernize runtime semantics and type behavior without needlessly replacing COBOL-style source syntax.
6. Reuse and modernize existing COBOL constructs when they already express the same concept.
7. Support deterministic parsing and tooling.
8. Be suitable for static analysis and automated review.
9. Support down-compilation to COBOL for a defined compatibility subset.
10. Support a native compiler pipeline.
11. Support transpilation to Bitlang.
12. Support a C backend suitable for GCC and Clang.
13. Extend the COBOL 2002 object model with modern features without discarding COBOL's character.

## 3. Modernization principle

Neo COBOL separates source-language style from runtime modernization.

The default rule is:

```text
Syntax and naming   -> COBOL-like / English-readable
Runtime semantics   -> modern and deterministic
Machine type widths -> modern fixed-width model
```

For example, `INTEGER` is preferred over an abbreviated `INT`, while its language-level size is a deterministic 32-bit signed integer.

When existing COBOL syntax already expresses the intended concept, Neo COBOL extends it instead of adding redundant modern-looking syntax. Examples include `OCCURS` for arrays/tables, level-88 condition names for named states, `PIC` for representation constraints, and `MOVE ... TO ...` for assignment.

## 4. Readability acceptance criterion

A proposed syntax feature should normally be rejected or redesigned when it makes valid source code substantially harder to read as an English description of program intent without providing a compelling technical benefit.

This does not require grammatically perfect English. Compiler determinism, compatibility, unambiguous grammar, and implementation safety take precedence where natural phrasing would introduce ambiguity.

## 5. Compatibility model

Compatibility with historical COBOL is defined explicitly rather than assumed.

The specification distinguishes:

- features that lower directly to COBOL,
- features that lower through generated support code,
- features that are native Neo COBOL features and cannot be represented faithfully on a selected COBOL target.

When a Neo shorthand conflicts with an established COBOL interpretation, the COBOL interpretation takes precedence unless Neo COBOL explicitly defines otherwise.

## 6. Compilation targets

Planned targets:

- COBOL
- C
- Bitlang
- Native Neo COBOL compiler pipeline

## 7. Language specification files

The language specification is intentionally split by concern.

- [`GRAMMAR.md`](GRAMMAR.md) — grammar entry point, common principles, and specification map.
- [`SOURCE_FORMAT.md`](SOURCE_FORMAT.md) — free-format source model.
- [`LEXICAL.md`](LEXICAL.md) — case, identifiers, comments, literals, and periods.
- [`PROGRAM_STRUCTURE.md`](PROGRAM_STRUCTURE.md) — divisions, sections, inference, and normalization.
- [`STATEMENTS.md`](STATEMENTS.md) — executable statements, conditions, control flow, and error phrases.
- [`TYPES.md`](TYPES.md) — built-in types, fixed numeric semantics, Unicode strings, references, aggregates, and function types.
- [`PIC.md`](PIC.md) — integration of `TYPE` with traditional `PIC` representation constraints.
- [`VARIABLES.md`](VARIABLES.md) — explicit declarations plus `VAR` and `LET` inference bindings.
- [`DATA_MODEL.md`](DATA_MODEL.md) — records, `STRUCT`, `OCCURS`, dynamic arrays, and level-88 condition names.
- [`CONVERSIONS.md`](CONVERSIONS.md) — numeric promotion, explicit conversion, reference conversion, and nullability.
- [`OOP.md`](OOP.md) — classes, inheritance, interfaces, namespaces, visibility, modifiers, creation, and destruction.
- [`FUNCTIONS.md`](FUNCTIONS.md) — named/anonymous functions, first-class callables, and closures.

Compiler-specific language tooling is documented separately under `docs/compiler/`, including [`../compiler/FORMATTER.md`](../compiler/FORMATTER.md).

## 8. Current language direction

The following points are established at draft level:

- Free-format source code.
- Case-insensitive keywords and identifiers.
- COBOL-style `*>` comments.
- COBOL-style `MOVE ... TO ...` assignment.
- Periods remain canonical sentence terminators but may be omitted when unambiguous.
- Traditional divisions/sections remain valid and may be inferred when omitted.
- COBOL English-style conditions are preferred over symbolic comparisons.
- COBOL 2002-style classes/functions remain recognized foundations.
- Modern OOP additions include `PUBLIC`, `PRIVATE`, `PROTECTED`, `STATIC`, `ABSTRACT`, `OVERRIDE`, `SEALED`, `INHERITS`, `IMPLEMENTS`, interfaces, and namespaces.
- Class inheritance is initially single inheritance; multiple interfaces may be implemented.
- Built-in scalar types use full-word names such as `INTEGER`, `DECIMAL`, `STRING`, and `BOOLEAN`.
- `BYTE`, `INTEGER`, `LONG`, `FLOAT`, `DOUBLE`, and `DECIMAL` have deterministic modern numeric semantics.
- `STRING` is variable-length Unicode text by default; UTF-8 is the canonical external/tooling encoding.
- Reference-like values are non-null by default; nullable references are explicitly marked `NULLABLE`.
- Explicit conversion uses `CONVERT ... TO ...`.
- Safe widening conversions may be implicit; narrowing or potentially lossy conversions are explicit.
- Traditional `OCCURS` is the canonical fixed-array/table mechanism; `OCCURS DYNAMIC` is the runtime-sized form.
- Traditional level-88 condition names remain the canonical named-state mechanism and may receive modern exhaustiveness diagnostics.
- `STRUCT` is a reusable aggregate value type integrated with traditional record concepts.
- `VAR` is inferred and mutable; `LET` is inferred and non-reassignable and is the local constant-binding mechanism.
- `TRUE`, `FALSE`, and `NULL` are built-in literals.
- Binary/octal/hexadecimal and scientific numeric notation are supported.
- Anonymous functions use `FUNCTION ... END FUNCTION` rather than symbolic lambda syntax.
- Closures and first-class function values are part of the language.
- Function types use the same `FUNCTION / USING / RETURNING` vocabulary as callable definitions.

## 9. Normalization philosophy

Neo COBOL source should normalize mechanically into a smaller canonical AST wherever practical.

Examples include:

- omitted divisions becoming explicit internal division nodes,
- `END IF` and `END-IF` becoming the same construct,
- short class/method forms becoming canonical class/method AST nodes,
- PIC-only declarations becoming logical `TYPE` plus representation constraints,
- traditional records and reusable `STRUCT` forms sharing a common aggregate model where semantics match,
- modern features lowering to generated support structures when targeting older COBOL dialects.

The compiler should reject ambiguous shorthand rather than guess.

## 10. Items intentionally not frozen yet

The following require design decisions that cannot be derived mechanically from the current principles:

- generic/parameterized type syntax,
- complete dynamic-array mutation API,
- exact runtime overflow policy,
- optional closed-set syntax for exhaustive level-88 groups,
- detailed error object/exception propagation model,
- namespace import/use syntax,
- ABI and foreign-function interoperability rules.

## 11. Versioning

The specification will use explicit language revisions once the first grammar and core type system stabilize. Until then, these documents remain working drafts.
