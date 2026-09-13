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
4. Support deterministic parsing and tooling.
5. Be suitable for static analysis and automated review.
6. Support down-compilation to COBOL for a defined compatibility subset.
7. Support a native compiler pipeline.
8. Support transpilation to Bitlang.
9. Support a C backend suitable for GCC and Clang.
10. Extend the COBOL 2002 object model with modern features without discarding COBOL's character.

## 3. Readability acceptance criterion

A proposed syntax feature should normally be rejected or redesigned when it makes valid source code substantially harder to read as an English description of program intent without providing a compelling technical benefit.

This does not require grammatically perfect English. Compiler determinism, compatibility, unambiguous grammar, and implementation safety take precedence where natural phrasing would introduce ambiguity.

## 4. Compatibility model

Compatibility with historical COBOL is defined explicitly rather than assumed.

The specification distinguishes:

- features that lower directly to COBOL,
- features that lower through generated support code,
- features that are native Neo COBOL features and cannot be represented faithfully on a selected COBOL target.

When a Neo shorthand conflicts with an established COBOL interpretation, the COBOL interpretation takes precedence unless Neo COBOL explicitly defines otherwise.

## 5. Compilation targets

Planned targets:

- COBOL
- C
- Bitlang
- Native Neo COBOL compiler pipeline

## 6. Language specification files

The language specification is intentionally split by concern.

- [`GRAMMAR.md`](GRAMMAR.md) — grammar entry point, common principles, and specification map.
- [`SOURCE_FORMAT.md`](SOURCE_FORMAT.md) — free-format source model.
- [`LEXICAL.md`](LEXICAL.md) — case, identifiers, comments, literals, and periods.
- [`PROGRAM_STRUCTURE.md`](PROGRAM_STRUCTURE.md) — divisions, sections, inference, and normalization.
- [`STATEMENTS.md`](STATEMENTS.md) — executable statements, conditions, control flow, and error phrases.
- [`TYPES.md`](TYPES.md) — current type-system decisions, including Boolean/null/reference/function types.
- [`OOP.md`](OOP.md) — classes, inheritance, interfaces, namespaces, visibility, modifiers, creation, and destruction.
- [`FUNCTIONS.md`](FUNCTIONS.md) — named/anonymous functions, first-class callables, and closures.

Compiler-specific language tooling is documented separately under `docs/compiler/`, including [`../compiler/FORMATTER.md`](../compiler/FORMATTER.md).

## 7. Current language direction

The following points are already established at draft level:

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
- Classes/interfaces/functions introduce reference-like types where applicable.
- `TRUE`, `FALSE`, and `NULL` are built-in literals.
- Binary/octal/hexadecimal and scientific numeric notation are supported.
- Anonymous functions use `FUNCTION ... END FUNCTION` rather than symbolic lambda syntax.
- Closures and first-class function values are part of the language direction.
- Function types use the same `FUNCTION / USING / RETURNING` vocabulary as callable definitions.

## 8. Normalization philosophy

Neo COBOL source should normalize mechanically into a smaller canonical AST wherever practical.

Examples include:

- omitted divisions becoming explicit internal division nodes,
- `END IF` and `END-IF` becoming the same construct,
- short class/method forms becoming canonical class/method AST nodes,
- modern features lowering to generated support structures when targeting older COBOL dialects.

The compiler should reject ambiguous shorthand rather than guess.

## 9. Versioning

The specification will use explicit language revisions once the first grammar and core type system stabilize. Until then, these documents remain working drafts.
