# Neo COBOL Type System

Status: Draft

This document records logical type-system decisions already made. Representation constraints based on traditional COBOL `PIC` are specified separately in [PIC.md](PIC.md).

## 1. Type categories

Neo COBOL currently distinguishes at least:

- primitive/scalar values,
- traditional COBOL data descriptions and records,
- class reference types,
- interface reference types,
- function reference types,
- Boolean values,
- null references.

## 2. TYPE and representation

`TYPE` defines the logical type of a value and provides its language-defined default representation.

Traditional `PIC` syntax remains part of Neo COBOL, but it is integrated into the same type model rather than acting as an unrelated second type system.

A `PIC` clause refines or overrides representation constraints of a compatible `TYPE`.

```cobol
01 AGE TYPE INTEGER.
01 AGE-3-DIGITS TYPE INTEGER PIC 9(3).
```

A traditional PIC-only declaration remains valid and is normalized by inferring a compatible logical `TYPE`.

```cobol
01 LEGACY-AGE PIC 9(3).
```

See [PIC.md](PIC.md) for the detailed normalization and compatibility rules.

## 3. Boolean

Neo COBOL has a Boolean type with literals:

```text
TRUE
FALSE
```

Boolean values are intended for conditions, function results, and normal data use.

## 4. Null

`NULL` denotes the absence of a reference.

`NULL` is valid for reference-like types where nullability is permitted, including class, interface, and function references.

The language will define static diagnostics for invalid null use. Exact nullable/non-null syntax is not yet fixed.

## 5. Class types

A class declaration introduces a class reference type.

```cobol
01 CUSTOMER-OBJECT TYPE CUSTOMER.
```

A variable of class type stores an object reference rather than embedding the full object value.

It may refer to:

- an instance of the declared class,
- an assignment-compatible derived-class instance,
- `NULL` where nullable references are permitted.

Derived-to-base assignment is supported subject to accessibility and inheritance rules. Explicit down-cast syntax remains to be specified.

## 6. Interface types

An interface declaration introduces an interface reference type.

An interface reference may refer to any object whose class implements that interface, or to `NULL` where nullable references are permitted.

Interfaces do not imply multiple class inheritance.

## 7. Function types

Functions are first-class values.

An inline function-reference type may be written as:

```cobol
01 TRANSFORMER TYPE FUNCTION USING NUMBER RETURNING NUMBER.
01 PREDICATE   TYPE FUNCTION USING CUSTOMER RETURNING BOOLEAN.
```

The general direction is:

```ebnf
function-type-reference = "FUNCTION",
                          [ "USING", type-reference, { type-reference } ],
                          [ "RETURNING", type-reference ] ;
```

Reusable named function types may be declared:

```cobol
FUNCTION TYPE CUSTOMER-PREDICATE
    USING CUSTOMER
    RETURNING BOOLEAN.
```

and then referenced as a normal type:

```cobol
01 IS-VALID TYPE CUSTOMER-PREDICATE.
```

Named functions, anonymous functions, and closures may be assigned where their signatures are compatible with the destination function type.

A function reference may be `NULL` where nullable references are permitted.

## 8. Function compatibility

A function value must have a compatible parameter list and return type.

Exact variance, overload interaction, implicit conversions, and callable covariance/contravariance are not yet fixed.

## 9. Traditional COBOL data

Traditional COBOL level numbers, records, `PIC`, and related data-description concepts remain part of Neo COBOL.

The unified rule is:

```text
TYPE = logical type and default representation
PIC  = compatible representation constraint/override
```

PIC-only legacy declarations are normalized to an inferred logical type plus the original PIC constraint.

## 10. Open items

- Full primitive/scalar type set
- Default representation for each scalar type
- Nullable/non-null reference syntax and defaults
- Explicit casting syntax
- Numeric conversion rules
- Function variance and overload compatibility
- Generic/parameterized types
- COBOL/C/Bitlang representation mappings
