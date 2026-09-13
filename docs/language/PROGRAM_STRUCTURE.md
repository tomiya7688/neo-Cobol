# Neo COBOL Program Structure

Status: Draft

Neo COBOL keeps traditional COBOL program structure while allowing redundant declarations to be omitted when their meaning is unambiguous.

## 1. Traditional divisions

Traditional COBOL divisions remain valid:

```cobol
IDENTIFICATION DIVISION.
DATA DIVISION.
PROCEDURE DIVISION.
```

## 2. Optional divisions and sections

Neo COBOL permits inferable divisions and sections to be omitted.

For example:

```cobol
PROGRAM-ID. HELLO.

01 NAME PIC X(20) VALUE "KADOKA".

DISPLAY "HELLO " NAME.
```

is interpreted as if written conceptually as:

```cobol
IDENTIFICATION DIVISION.
PROGRAM-ID. HELLO.

DATA DIVISION.
WORKING-STORAGE SECTION.
01 NAME PIC X(20) VALUE "KADOKA".

PROCEDURE DIVISION.
DISPLAY "HELLO " NAME.
```

## 3. Initial inference rules

When explicit structure is omitted:

- `PROGRAM-ID` and related program metadata belong to the Identification Division.
- Level-number data declarations such as `01` and `05` belong to the Data Division.
- Executable statements such as `MOVE`, `DISPLAY`, `IF`, `PERFORM`, `CALL`, `CREATE`, and `DESTROY` belong to the Procedure Division.
- Top-level `CLASS`, `INTERFACE`, `FUNCTION`, and `NAMESPACE` constructs are recognized directly rather than inferred from arbitrary statement text.
- If a source fragment could validly belong to more than one implicit region, the compiler must reject it rather than guess.

Inference must be deterministic.

## 4. Normalization model

Neo COBOL shorthand is normalized into a canonical COBOL-like internal representation before later compiler stages wherever practical.

Examples:

```text
CLASS PERSON        -> CLASS-ID. PERSON.
METHOD PRINT-NAME   -> METHOD-ID. PRINT-NAME.
END IF              -> END-IF
END PERFORM         -> END-PERFORM
```

Omitted divisions and sections become explicit AST nodes after parsing.

This allows semantic analysis and backends to operate on one canonical representation rather than maintaining separate traditional and Neo forms.

## 5. Compatibility priority

Where a Neo shorthand conflicts with an established COBOL interpretation, the COBOL interpretation takes precedence unless the Neo COBOL specification explicitly states otherwise.

Neo syntax extends COBOL rather than silently redefining familiar COBOL syntax.

## 6. Source format

Neo COBOL source is free-format. See `SOURCE_FORMAT.md`.

Column positions and indentation are not semantic. Formatting diagnostics are tooling concerns described in `../compiler/FORMATTER.md`.

## 7. Open items

- Exact inference boundaries between sections
- Ambiguity diagnostics
- Multi-program source-unit rules
- Relationship between namespaces and source units
- Module/import structure
