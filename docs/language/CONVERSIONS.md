# Neo COBOL Conversions and Nullability

Status: Draft

## 1. General rule

Neo COBOL prefers safe implicit widening and explicit potentially lossy conversion.

The source syntax remains English-like and COBOL-oriented.

## 2. Numeric promotion

The default widening order is:

```text
BYTE -> INTEGER -> LONG
FLOAT -> DOUBLE
```

`DECIMAL` is an exact decimal type and does not implicitly convert to or from binary floating-point types when precision could change silently.

Rules:

- `BYTE` may widen implicitly to `INTEGER`, `LONG`, `FLOAT`, `DOUBLE`, or `DECIMAL` when the destination can represent the source value category safely.
- `INTEGER` may widen implicitly to `LONG`, `DOUBLE`, or `DECIMAL` where exactness is guaranteed by the language rules.
- `LONG` does not implicitly convert to `FLOAT` or `DOUBLE` when integer precision may be lost.
- `FLOAT` may widen implicitly to `DOUBLE`.
- narrowing conversions require an explicit conversion.
- decimal/binary floating conversion requires an explicit conversion unless the compiler can prove exact preservation.

Arithmetic between compatible numeric types promotes to the wider safe type. The compiler must diagnose ambiguous or potentially lossy implicit promotion rather than choosing silently.

## 3. Explicit conversion

The canonical explicit conversion form is:

```cobol
CONVERT VALUE TO INTEGER
CONVERT TOTAL TO DECIMAL
CONVERT OBJECT TO CUSTOMER
```

`CONVERT ... TO ...` is an expression and yields the converted value.

Example:

```cobol
MOVE CONVERT INPUT-VALUE TO INTEGER TO COUNT.
```

The keyword `AS` is not the canonical general cast operator because Neo COBOL already uses `AS` in other English-oriented constructs such as object creation.

## 4. Failed conversion

A conversion that can fail at runtime must expose that possibility through the normal Neo COBOL error-handling model. It must not silently reinterpret bits or return an unrelated value.

Exact integration with `ON ERROR` / `ON EXCEPTION` is specified with the statement and error model.

## 5. Reference conversions

Derived-class references convert implicitly to accessible base-class references and implemented interface references.

Down-casts and unrelated interface/class conversions require explicit `CONVERT ... TO ...` and runtime type validation when not statically provable.

`SEALED` and inheritance rules participate in compile-time conversion diagnostics.

## 6. Nullability

Reference-like types are non-null by default.

A nullable reference is declared explicitly with `NULLABLE`:

```cobol
01 CUSTOMER-OBJECT TYPE CUSTOMER NULLABLE.
01 HANDLER TYPE CUSTOMER-PREDICATE NULLABLE.
```

Only nullable references may contain `NULL`.

This rule applies to class, interface, and function reference types. Scalar values and `STRUCT` values are not nullable unless a later wrapper/optional type is explicitly introduced.

The compiler must reject use of a nullable reference as definitely non-null unless control-flow analysis proves it non-null or an explicit checked conversion/guard is used.

## 7. NULL tests

English-oriented comparisons are canonical:

```cobol
IF CUSTOMER-OBJECT IS NULL
    ...
END-IF.

IF CUSTOMER-OBJECT IS NOT NULL
    ...
END-IF.
```

## 8. Representation constraints

`PIC` constraints do not change logical conversion rules. A value may be type-compatible yet fail a destination representation constraint.

Where overflow/truncation is statically provable, the compiler must diagnose it. Otherwise the runtime behavior must be explicit and must not silently corrupt the logical value.

## 9. Open items

- Exact checked-conversion syntax if a non-throwing conversion form is added
- Detailed arithmetic result table for mixed `DECIMAL` and integer expressions
- Runtime overflow mode selection, if configurable
- User-defined conversions, if later supported
