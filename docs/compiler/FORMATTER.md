# Neo COBOL Compiler Formatter Specification

Status: Draft / Design memo

This document records the initial formatter and formatting-diagnostic policy for the Neo COBOL compiler toolchain.

## 1. Role

Neo COBOL is a free-format language. Indentation is not part of the language grammar and does not change program meaning.

However, readability is a core design goal of Neo COBOL. The compiler toolchain therefore provides formatting support and may emit non-fatal diagnostics when indentation or layout is sufficiently poor to make the program difficult to read.

Formatting diagnostics are style diagnostics, not syntax errors.

## 2. Initial policy

- Source code is parsed independently of indentation.
- Bad indentation must not change semantics.
- The formatter may rewrite source into the canonical project style.
- The compiler / formatter may warn about clearly inconsistent or misleading indentation.
- Formatting warnings must not prevent compilation by default.
- A future strict mode may promote selected formatting warnings to errors for CI or project policy, but this is not part of the language semantics.

## 3. Examples of formatting warnings

The exact thresholds are not yet fixed. Candidate cases include:

- Statements inside `IF`, `ELSE`, `PERFORM`, `METHOD`, `FUNCTION`, `CLASS`, and similar blocks are not indented relative to their containing block.
- Sibling statements use inconsistent indentation without a structural reason.
- A statement is indented as though it belonged to a different block than the parser determines.
- Closing forms such as `END-IF`, `END-PERFORM`, `END METHOD`, or `END CLASS` are visually aligned with the wrong scope.
- Excessive indentation or indentation oscillation substantially harms readability.

Example:

```cobol
IF CUSTOMER-VALID
DISPLAY "VALID"
        MOVE 1 TO RESULT
    END-IF.
```

This remains syntactically valid if the tokens form a valid program, but the formatter should warn and can normalize it to:

```cobol
IF CUSTOMER-VALID
    DISPLAY "VALID"
    MOVE 1 TO RESULT
END-IF.
```

## 4. Canonical formatting direction

Initial formatter direction:

- Preserve free-format source.
- Indent nested executable blocks consistently.
- Keep matching block terminators visually aligned with their opening construct.
- Preserve COBOL-oriented keyword spelling and structure.
- Avoid formatting rules that require semantic rewrites.
- Prefer deterministic output: formatting the same source repeatedly must be idempotent.

The exact indentation width, line wrapping rules, keyword casing policy, alignment policy, and blank-line rules remain to be specified.

## 5. Diagnostics

Formatter diagnostics should be short and actionable.

Example direction:

```text
warning[format-indent]: DISPLAY is not indented inside IF block
warning[format-scope]: END-IF indentation does not match its IF
```

Diagnostics should include source location information when available.

## 6. Separation from the language specification

The language grammar defines what source is valid.

The formatter defines how valid source should preferably be laid out for readability.

A poorly formatted program may therefore compile successfully while still producing formatter warnings.

## 7. Open items

- Canonical indentation width
- Tabs versus spaces policy
- Maximum line length
- Continuation-line indentation
- Formatting of `DATA DIVISION` entries and level numbers
- Alignment of `PIC`, `VALUE`, `USAGE`, and related clauses
- Keyword casing in formatter output
- Blank-line rules around divisions, sections, classes, methods, and functions
- Warning thresholds and suppressions
- `--format`, `--check-format`, and CI integration
