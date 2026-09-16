# Neo COBOL PIC and Representation Constraints

Status: Draft

This document defines how Neo COBOL integrates `TYPE` declarations with traditional COBOL `PIC` data descriptions.

## 1. Core rule

`TYPE` defines the logical type and its default representation.

`PIC` refines a compatible representation, such as digit count, character width, sign layout, decimal position, and COBOL-compatible field layout.

```text
TYPE = logical type and default representation
PIC  = compatible representation constraint
```

## 2. TYPE-only declarations

```cobol
01 AGE TYPE INTEGER.
01 NAME TYPE STRING.
01 PRICE TYPE DECIMAL.
01 IS-ACTIVE TYPE BOOLEAN.
```

## 3. TYPE with PIC

```cobol
01 AGE TYPE INTEGER PIC 9(3).
01 NAME TYPE STRING PIC X(20).
01 PRICE TYPE DECIMAL PIC S9(7)V99.
```

The logical type remains unchanged; `PIC` adds a field representation constraint.

```cobol
01 SMALL-VALUE TYPE INTEGER PIC 9(3).
01 LARGE-VALUE TYPE INTEGER PIC 9(8).
```

Both are `INTEGER` values even though their field layouts differ.

## 4. PIC-only compatibility syntax

Traditional COBOL declarations remain valid:

```cobol
01 AGE PIC 9(3).
01 NAME PIC X(20).
01 PRICE PIC S9(7)V99.
```

When `TYPE` is omitted, the compiler infers a compatible Neo COBOL logical type and keeps the original `PIC` constraint.

## 5. Intended combinations

```cobol
TYPE INTEGER PIC 9(5)
TYPE DECIMAL PIC S9(7)V99
TYPE STRING PIC X(20)
```

A `PIC` that is incompatible with the declared logical type is invalid.

Class, interface, and function reference types do not normally use `PIC`.

## 6. BOOLEAN

`BOOLEAN` has a language-defined default representation. A target-specific compatible layout may be defined later without changing its Boolean semantics.

## 7. Normalization

```text
PIC-only declaration
    -> infer compatible TYPE
    -> preserve PIC constraint
    -> continue with one unified type model
```

## 8. Open items

- Exact PIC-to-TYPE inference table
- Default representation of built-in types
- `USAGE` interaction
- `DECIMAL` precision defaults
- `STRING` representation details
- Edited PIC/display forms
- Backend mappings
