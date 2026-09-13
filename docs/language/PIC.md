# Neo COBOL PIC and Representation Constraints

Status: Draft

This document defines how Neo COBOL integrates modern `TYPE` declarations with traditional COBOL `PIC` data descriptions.

## 1. Core rule

`TYPE` defines the logical type of a value.

`PIC` refines or overrides the representation constraints of that type, such as digit count, character width, sign layout, decimal position, and COBOL-compatible storage/format expectations.

In short:

```text
TYPE = logical value type and default representation
PIC  = representation constraint/override for that type
```

`PIC` does not create a separate unrelated type system.

## 2. TYPE-only declarations

A Neo COBOL type has a language-defined default representation.

```cobol
01 AGE TYPE INTEGER.
01 NAME TYPE STRING.
01 PRICE TYPE DECIMAL.
01 IS-ACTIVE TYPE BOOLEAN.
```

The exact default width, precision, storage class, and target-specific lowering rules for each built-in type are defined by the type-system and backend specifications.

## 3. TYPE with PIC

`PIC` may refine the representation of a compatible `TYPE`.

```cobol
01 AGE TYPE INTEGER PIC 9(3).
01 NAME TYPE STRING PIC X(20).
01 PRICE TYPE DECIMAL PIC S9(7)V99.
```

The logical type remains `INTEGER`, `STRING`, or `DECIMAL` respectively.

The `PIC` clause adds representation constraints.

For example:

```cobol
01 SMALL-VALUE TYPE INTEGER PIC 9(3).
01 LARGE-VALUE TYPE INTEGER PIC 9(8).
```

Both variables have logical type `INTEGER`, so ordinary integer assignment is type-compatible. However, moving a value that does not fit the destination `PIC` constraint must produce the appropriate compile-time diagnostic where provable, or a defined runtime handling behavior where it is not statically knowable.

## 4. PIC-only compatibility syntax

Traditional COBOL declarations remain valid.

```cobol
01 AGE PIC 9(3).
01 NAME PIC X(20).
01 PRICE PIC S9(7)V99.
```

When a declaration contains `PIC` but no explicit `TYPE`, the compiler infers the corresponding Neo COBOL logical type and normalizes the declaration internally to an equivalent `TYPE + PIC` form.

Conceptually:

```text
PIC-only source
    -> infer compatible logical TYPE
    -> attach original PIC as representation constraint
    -> continue semantic analysis using the unified type model
```

The exact inference table is specified separately and must be deterministic.

## 5. Compatibility constraints

Not every `PIC` form is valid for every logical type.

Examples of intended compatible combinations include:

```cobol
TYPE INTEGER  PIC 9(5)
TYPE DECIMAL  PIC S9(7)V99
TYPE STRING   PIC X(20)
```

A semantically incompatible combination must be rejected rather than silently reinterpreted.

For example, a string-oriented `PIC` must not silently convert an unrelated class reference or function reference into a string type.

## 6. Reference and modern types

Class, interface, function-reference, and other reference-oriented types do not normally use `PIC` because their logical value is a reference rather than a fixed COBOL data field.

```cobol
01 CUSTOMER-OBJECT TYPE CUSTOMER.
01 PREDICATE TYPE FUNCTION USING CUSTOMER RETURNING BOOLEAN.
```

If future interoperability features require representation constraints for references, they should use a dedicated interoperability/storage mechanism rather than overloading ordinary `PIC` semantics without an explicit specification.

## 7. Boolean

`BOOLEAN` is a logical Neo COBOL type.

```cobol
01 IS-ACTIVE TYPE BOOLEAN.
```

Its default representation is language-defined.

A future `PIC` mapping for Boolean may be defined for COBOL interoperability, but Boolean semantics remain Boolean even when represented using a target-compatible field.

## 8. Normalization model

The compiler should normalize declarations early so later compiler stages operate on one model.

For example:

```cobol
01 AGE PIC 9(3).
```

becomes conceptually:

```cobol
01 AGE TYPE <inferred-numeric-type> PIC 9(3).
```

while:

```cobol
01 AGE TYPE INTEGER PIC 9(3).
```

already contains both the logical type and representation constraint.

Backends should consume the normalized type plus representation information rather than reparsing source-level `PIC` syntax independently.

## 9. Diagnostics

The compiler should diagnose at least:

- a `PIC` incompatible with the declared `TYPE`,
- a statically known value that cannot fit the destination `PIC`,
- ambiguous PIC-only type inference,
- representation combinations that the selected backend cannot faithfully lower.

Backend limitations should not silently change program meaning.

## 10. Open items

- Exact PIC-to-TYPE inference table
- Default representation for each built-in `TYPE`
- `USAGE` interaction with `TYPE` and `PIC`
- Signed/unsigned integer policy
- Decimal precision defaults
- Fixed-length versus variable-length string defaults
- Runtime behavior for representation overflow/truncation
- Edited PIC forms and display formatting
- Backend-specific representation mappings
