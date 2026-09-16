# Neo COBOL Type System

Status: Draft

Representation constraints are specified in [PIC.md](PIC.md). Variable declarations are specified in [VARIABLES.md](VARIABLES.md). Data aggregates and arrays are specified in [DATA_MODEL.md](DATA_MODEL.md). Conversions and nullability are specified in [CONVERSIONS.md](CONVERSIONS.md).

## 1. Core rule

`TYPE` defines the logical type and its default representation. `PIC` may refine a compatible representation.

When two spellings have the same meaning, Neo COBOL prefers the more COBOL-like, English-readable full word over an abbreviated alias.

Neo COBOL uses COBOL-like names and source syntax while giving built-in machine-oriented types modern, deterministic semantics. Source-level behavior must not change merely because a backend C compiler, COBOL compiler, or target CPU uses a different native representation.

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

| Type | Meaning | Canonical representation |
| --- | --- | --- |
| `BYTE` | unsigned byte-sized integer/data unit | 8-bit |
| `INTEGER` | normal signed integer | 32-bit two's-complement semantics |
| `LONG` | wide signed integer | 64-bit two's-complement semantics |
| `FLOAT` | binary floating point | IEEE-754 binary32 |
| `DOUBLE` | wide binary floating point | IEEE-754 binary64 |
| `DECIMAL` | exact decimal floating/fixed-point numeric type | IEEE-754 decimal128 semantics |

`DECIMAL` therefore has up to 34 significant decimal digits under the canonical language model. Backends may lower it differently only when observable Neo COBOL semantics are preserved.

`PIC` may further constrain a field layout without changing the logical type.

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

`STRING` is the normal Neo COBOL Unicode text type.

Its logical content is a sequence of Unicode characters. UTF-8 is the canonical external/interchange encoding used by Neo COBOL tooling and native text I/O unless an interoperability boundary explicitly specifies another encoding.

The in-memory representation is implementation-defined and must not affect language semantics.

`STRING` is variable-length by default:

```cobol
01 NAME TYPE STRING.
```

`PIC X(n)` may impose a fixed-width COBOL-compatible field representation:

```cobol
01 FIXED-NAME TYPE STRING PIC X(20).
```

### 2.3 BOOLEAN

`BOOLEAN` has exactly two logical values:

```text
TRUE
FALSE
```

Its physical representation is backend-defined, while logical semantics are fixed.

```cobol
01 IS-ACTIVE TYPE BOOLEAN.
```

## 3. STRUCT

`STRUCT` defines a user-defined aggregate value type for grouped data that does not require object identity or inheritance.

```cobol
STRUCT CUSTOMER-DATA
    01 ID TYPE INTEGER.
    01 NAME TYPE STRING.
END STRUCT.

01 CUSTOMER TYPE CUSTOMER-DATA.
```

`STRUCT` values use value semantics. Copying a `STRUCT` copies its logical field values.

See [DATA_MODEL.md](DATA_MODEL.md) for the relationship with traditional level-number records.

## 4. Reference types and NULL

Class, interface, and function types are reference-like types.

Reference types are non-null by default. A reference that may hold `NULL` must be explicitly marked `NULLABLE`.

```cobol
01 CUSTOMER-OBJECT TYPE CUSTOMER NULLABLE.
```

Scalar and `STRUCT` values are not implicitly nullable.

See [CONVERSIONS.md](CONVERSIONS.md).

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

`VAR` and `LET` are declarations, not runtime types.

- `VAR`: inferred type, mutable binding.
- `LET`: inferred type, non-reassignable binding.

See [VARIABLES.md](VARIABLES.md).

## 7. Traditional COBOL data

Traditional level numbers, records, level-88 condition names, `OCCURS`, and `PIC` remain valid and are integrated into the same logical type/data model.

PIC-only declarations are normalized to an inferred logical type plus the original `PIC` constraint.

## 8. Modernization rule

Neo COBOL modernizes semantics without needlessly modernizing surface spelling.

As a default rule:

- prefer COBOL or English-like terminology and statement forms,
- retain established COBOL constructs when they remain useful,
- use modern deterministic numeric/string/reference semantics,
- make backend-dependent representation differences non-observable where practical,
- add modern features by extending COBOL concepts rather than replacing them with unrelated syntax.

## 9. Open items

- Exact runtime overflow handling mode
- User-defined generic/parameterized types
- ABI/backend mappings
