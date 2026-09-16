# Neo COBOL Type System

Status: Draft

Representation constraints are specified in [PIC.md](PIC.md). Variable declarations are specified in [VARIABLES.md](VARIABLES.md).

## 1. Core rule

`TYPE` defines the logical type and its default representation. `PIC` may refine a compatible representation.

When two spellings have the same meaning, Neo COBOL prefers the more COBOL-like, English-readable full word over an abbreviated alias.

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

- `INTEGER`: normal signed integer.
- `LONG`: wider signed integer.
- `BYTE`: byte-sized integer/data unit; intended default is 8-bit.
- `DECIMAL`: exact decimal/fixed-point numeric type.
- `FLOAT`: normal floating-point type; intended default is 32-bit.
- `DOUBLE`: wider floating-point type; intended default is 64-bit.
- `STRING`: normal string type.
- `BOOLEAN`: Boolean type with `TRUE` and `FALSE` literals.

Examples:

```cobol
01 AGE TYPE INTEGER.
01 PRICE TYPE DECIMAL PIC S9(7)V99.
01 NAME TYPE STRING PIC X(20).
01 IS-ACTIVE TYPE BOOLEAN.
```

Exact backend representations remain separately specified.

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

## 8. Open items

- Exact integer and byte representations
- Default `DECIMAL` precision/scale
- Default `STRING` representation
- Nullable/non-null syntax
- Casting and numeric promotion rules
- Detailed `STRUCT` layout
- Backend mappings
