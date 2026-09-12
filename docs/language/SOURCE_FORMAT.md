# Neo COBOL Source Format

Status: Draft

## 1. Source form

Neo COBOL uses **free-format source code** as part of the language specification.

Unlike historical fixed-format COBOL source forms, Neo COBOL does not assign semantic meaning to specific character columns.

Source indentation and horizontal alignment do not affect program meaning.

Example:

```cobol
IDENTIFICATION DIVISION.
PROGRAM-ID. HELLO.

PROCEDURE DIVISION.
    DISPLAY "HELLO WORLD".
```

The language parser must not require traditional COBOL column areas for valid Neo COBOL source.

## 2. Indentation

Indentation is non-semantic. It is used for readability only.

Poor or misleading indentation may be reported by compiler tooling as a formatting warning, but it does not make otherwise valid Neo COBOL source syntactically invalid.

Canonical formatting and indentation diagnostics are specified separately in `../compiler/FORMATTER.md`.

## 3. Legacy fixed-format COBOL

Legacy fixed-format COBOL is not the native Neo COBOL source format.

Support for importing, converting, or otherwise accepting legacy fixed-format COBOL may be added as a compatibility feature, but such support is separate from the core Neo COBOL grammar.
