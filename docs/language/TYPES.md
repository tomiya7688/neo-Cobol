# Neo COBOL Type System

Status: Draft

Representation constraints are specified in [PIC.md](PIC.md). Variable declarations are specified in [VARIABLES.md](VARIABLES.md).

## 1. Core rule

`TYPE` defines the logical type and its default representation. `PIC` may refine a compatible representation.

When two spellings have the same meaning, Neo COBOL prefers the more COBOL-like, English-readable full word over an abbreviated alias.

Neo COBOL uses COBOL-like names and source syntax while giving built-in machine-oriented types modern, deterministic widths. Source-level semantics must not change merely because a backend C compiler, COBOL compiler, or target CPU uses a different native width.

## 2. Built-in scalar types

Canonical built-in scalar types are:

```text
INTEGER
DECIMAL
STRING
BOOLEAN
FLOAT
BYTE
LONG
DOUBLE
```

Short aliases such as `INT`, `DEC`, `STR`, and `BOOL` are not canonical type names.

### 2.1 Numeric widths

The canonical language-level widths are:

| Type | Meaning | Canonical width |
| --- | --- | ---: |
| `BYTE` | unsigned byte-sized integer/data unit | 8 bits |
| `INTEGER` | normal signed integer | 32 bits |
| `LONG` | wide signed integer | 64 bits |
| `FLOAT` | IEEE-754 binary floating point | 32 bits |
| `DOUBLE` | IEEE-754 binary floating point | 64 bits |
| `DECIMAL` | exact decimal numeric type | 128 bits |

Backends may use another physical representation internally only when observable Neo COBOL semantics are preserved.

`PIC` may further constrain the representable field layout without changing the logical type.

Examples:

```cobol
01 AGE TYPE INTEGER.
01 ACCOUNT-ID TYPE LONG.
01 FLAGS TYPE BYTE.
01 PRICE TYPE DECIMAL PIC S9(7)V99.
01 RATIO TYPE FLOAT.
01 PRECISE-RATIO TYPE DOUBLE.
```

### 2.2 STRING

`STRING` is the normal Neo COBOL text type.

```cobol
01 NAME TYPE STRING.
01 FIXED-NAME TYPE STRING PIC X(20).
```

The language-level string is not limited to legacy fixed-width COBOL character fields. `PIC X(n)` may impose a compatible field representation constraint where required.

The canonical character encoding and exact in-memory string representation remain to be specified separately.

### 2.3 BOOLEAN

`BOOLEAN` is the logical Boolean type and has the literals:

```text
TRUE
FALSE
```

Its language semantics are independent of the backend's physical Boolean representation.

```cobol
01 IS-ACTIVE TYPE BOOLEAN.
```

## 3. STRUCT

`STRUCT` defines a user-defined aggregate value type. It is for grouped data that does not require class identity or inheritance.

```cobol
STRUCT CUSTOMER-DATA
    01 ID TYPE INTEGER.
    01 NAME TYPE STRING.
END STRUCT.

01 CUSTOMER TYPE CUSTOMER-DATA.
```

Detailed layout and interaction with traditional COBOL records remain open.

## 4. Reference types and NULL

Class, interface, and function types are reference-like types. `NULL` may be used where nullability is permitted.

Scalar and `STRUCT` values are not implicitly nullable.

## 5. Function types

```cobol
01 TRANSFORMER TYPE FUNCTION USING INTEGER RETURNING INTEGER.
01 PREDICATE TYPE FUNCTION USING CUSTOMER RETURNING BOOLEAN.
```

Reusable named function types remain supported:

```cobol
FUNCTION TYPE CUSTOMER-PREDICATE
    USING CUSTOMER
    RETURNING BOOLEAN.
```

## 6. Type inference declarations

`VAR` and `LET` are declarations, not types.

- `VAR`: inferred type, mutable binding.
- `LET`: inferred type, non-reassignable binding.

See [VARIABLES.md](VARIABLES.md).

## 7. Traditional COBOL data

Traditional level numbers, records, and `PIC` remain valid. PIC-only declarations are normalized to an inferred logical type plus the original `PIC` constraint.

## 8. Design rule for modernization

Neo COBOL modernizes semantics without needlessly modernizing surface spelling.

As a default rule:

- prefer COBOL or English-like terminology and statement forms,
- retain established COBOL constructs when they remain useful,
- use modern fixed-width numeric semantics,
- make backend-dependent representation differences non-observable where practical,
- add modern features by extending COBOL concepts rather than replacing them with unrelated syntax.

## 9. Open items

- Signed/unsigned arithmetic details for `BYTE`
- Default `DECIMAL` precision and scale within the 128-bit representation
- Canonical `STRING` encoding and storage model
- Nullable/non-null syntax
- Casting and numeric promotion rules
- Detailed `STRUCT` layout
- Backend mappings
