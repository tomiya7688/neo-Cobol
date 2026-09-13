# Neo COBOL Grammar Specification

Status: Draft

This document is the grammar entry point for Neo COBOL. Detailed rules are split into focused specification files so that syntax, types, OOP, and runtime-related concerns do not accumulate in one document.

## 1. Core grammar principles

1. COBOL syntax is the baseline.
2. Existing COBOL spellings and structures are preserved unless there is a clear reason to extend them.
3. Keywords and identifiers are case-insensitive.
4. English-readable syntax is preferred over symbolic shorthand.
5. `MOVE ... TO ...` remains the canonical assignment form.
6. Explicit COBOL block forms such as `IF ... END-IF` and `PERFORM ... END-PERFORM` are retained.
7. A period (`.`) is the canonical sentence terminator, but may be omitted where the boundary is unambiguous.
8. `DIVISION` and `SECTION` declarations remain valid but may be omitted where their role can be inferred deterministically.
9. Neo shorthand should normalize mechanically into a canonical COBOL-like internal representation wherever practical.
10. If Neo shorthand conflicts with an established COBOL interpretation, the COBOL interpretation takes precedence unless the Neo COBOL specification explicitly says otherwise.
11. Neo COBOL extends COBOL 2002 with modern OOP, namespaces, first-class functions, anonymous functions, closures, Boolean values, and null/reference semantics.

## 2. Specification map

Detailed language rules are split as follows:

- [`SOURCE_FORMAT.md`](SOURCE_FORMAT.md) — free-format source rules.
- [`LEXICAL.md`](LEXICAL.md) — identifiers, comments, literals, case, and sentence termination.
- [`PROGRAM_STRUCTURE.md`](PROGRAM_STRUCTURE.md) — divisions, sections, inference, and normalization.
- [`STATEMENTS.md`](STATEMENTS.md) — `MOVE`, `DISPLAY`, conditions, `IF`, `PERFORM`, error phrases, and statement-level behavior.
- [`TYPES.md`](TYPES.md) — Boolean, null, class/interface/function reference types, and type-system direction.
- [`OOP.md`](OOP.md) — classes, access control, inheritance, interfaces, namespaces, modifiers, object creation, and destruction.
- [`FUNCTIONS.md`](FUNCTIONS.md) — named functions, anonymous functions, first-class function values, and closures.
- [`../compiler/FORMATTER.md`](../compiler/FORMATTER.md) — non-semantic formatting and indentation diagnostics.

The higher-level language design and target model are described in [`SPECIFICATION.md`](SPECIFICATION.md).

## 3. Preliminary top-level grammar

```ebnf
source-unit = { namespace-definition
              | class-definition
              | interface-definition
              | function-definition
              | program-definition } ;

program-definition = [ identification-division ],
                     [ data-division ],
                     [ procedure-division ] ;
```

Omitted program divisions are reconstructed according to `PROGRAM_STRUCTURE.md`.

## 4. Shared lexical references

```ebnf
qualified-name = identifier, { ".", identifier } ;

sentence-terminator = [ "." ] ;
```

Exact lexical details are normative in `LEXICAL.md`.

## 5. Shared block normalization

The following spellings may normalize to the same internal construct:

```text
END IF       -> END-IF
END PERFORM  -> END-PERFORM
CLASS NAME   -> canonical class-definition AST
METHOD NAME  -> canonical method-definition AST
```

Traditional COBOL forms remain valid where supported.

## 6. Design status

The grammar is still a working draft. A feature described as decided in the focused specification files may still require implementation details, backend lowering rules, or additional validation before the first stable language revision.

## 7. Major remaining grammar work

- Full reserved-word set
- Exact period-omission boundaries
- Full `PIC` and data-description grammar
- Numeric/storage type rules
- Complete `COMPUTE`
- Complete `CALL`
- `EVALUATE`
- Complete `PERFORM` family
- File handling (`READ`, `WRITE`, `OPEN`, `CLOSE`, etc.)
- Namespace import/use syntax
- Casting and conversion syntax
- Generic/parameterized facilities, if adopted
- Interoperability syntax
