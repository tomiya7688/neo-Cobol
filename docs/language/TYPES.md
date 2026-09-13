# Neo COBOL Type System

Status: Draft

This document records type-system decisions already made. Detailed scalar/PIC mapping remains open.

## 1. Type categories

Neo COBOL currently distinguishes at least:

- primitive/scalar values,
- traditional COBOL data descriptions and records,
- class reference types,
- interface reference types,
- function reference types,
- Boolean values,
- null references.

## 2. Boolean

Neo COBOL has a Boolean type with literals:

```text
TRUE
FALSE
```

Boolean values are intended for conditions, function results, and normal data use.

## 3. Null

`NULL` denotes the absence of a reference.

`NULL` is valid for reference-like types where nullability is permitted, including class, interface, and function references.

The language will define static diagnostics for invalid null use. Exact nullable/non-null syntax is not yet fixed.

## 4. Class types

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

## 5. Interface types

An interface declaration introduces an interface reference type.

An interface reference may refer to any object whose class implements that interface, or to `NULL` where nullable references are permitted.

Interfaces do not imply multiple class inheritance.

## 6. Function types

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

## 7. Function compatibility

A function value must have a compatible parameter list and return type.

Exact variance, overload interaction, implicit conversions, and callable covariance/contravariance are not yet fixed.

## 8. Traditional COBOL data

Traditional COBOL data-description concepts, including level numbers and `PIC`, remain part of Neo COBOL.

Modern types extend rather than replace traditional COBOL data descriptions.

Detailed `PIC`, `USAGE`, numeric storage, string storage, and mapping rules will be specified separately.

## 9. Open items

- Full primitive/scalar type set
- `PIC` and modern scalar type interaction
- Nullable/non-null reference syntax and defaults
- Explicit casting syntax
- Numeric conversion rules
- Function variance and overload compatibility
- Generic/parameterized types
- COBOL/C/Bitlang representation mappings
