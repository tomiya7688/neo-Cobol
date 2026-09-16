# Neo COBOL Parameter Passing

Status: Draft

Neo COBOL extends traditional COBOL calling conventions to functions, methods, constructors, and other callable forms rather than introducing a separate modern calling syntax.

## 1. Core clauses

Canonical parameter and return clauses are:

```text
USING
RETURNING
BY REFERENCE
BY VALUE
OPTIONAL
VALUE
```

`USING` declares or supplies arguments. `RETURNING` declares the result value.

## 2. Default passing mode

When neither `BY REFERENCE` nor `BY VALUE` is written, the parameter is `BY REFERENCE`, preserving traditional COBOL behavior.

```cobol
METHOD UPDATE-CUSTOMER
    USING CUSTOMER NEW-NAME
    RETURNING RESULT
END METHOD.
```

is equivalent to:

```cobol
METHOD UPDATE-CUSTOMER
    USING BY REFERENCE CUSTOMER
          BY REFERENCE NEW-NAME
    RETURNING RESULT
END METHOD.
```

## 3. BY VALUE

`BY VALUE` passes a value copy according to the logical type's value semantics.

```cobol
FUNCTION ADD
    USING BY VALUE LEFT-VALUE
          BY VALUE RIGHT-VALUE
    RETURNING RESULT
END FUNCTION.
```

For reference-like logical types, `BY VALUE` copies the reference value rather than cloning the referenced object.

## 4. BY REFERENCE

`BY REFERENCE` passes access to the caller's storage or binding according to Neo COBOL reference-passing rules.

A callee may modify the referenced value when the parameter and target are writable.

Passing a non-reassignable `LET` binding by reference does not make the binding itself reassignable.

## 5. Positional arguments

Neo COBOL uses positional arguments as the canonical calling model.

Named arguments using a separate syntax such as `NAME = VALUE` are not part of the core language.

This keeps callable syntax aligned with traditional COBOL `USING` conventions and avoids introducing a second argument-binding model.

## 6. OPTIONAL and default values

`OPTIONAL` marks a parameter that may be omitted by the caller.

A default value may be specified using the existing COBOL-style `VALUE` vocabulary.

```cobol
METHOD FIND-CUSTOMER
    USING CUSTOMER-ID
          OPTIONAL LIMIT VALUE 100
    RETURNING RESULT
END METHOD.
```

When `LIMIT` is omitted, its value is `100`.

Required parameters must precede optional parameters.

If an `OPTIONAL` parameter has no explicit `VALUE`, the callable may observe its omitted state according to the parameter/type rules. The exact representation of that omitted state is an implementation detail and must not be confused with `NULL` unless the declared type is nullable.

## 7. RETURNING

A callable has at most one direct `RETURNING` result.

```cobol
FUNCTION CALCULATE-TOTAL
    USING ORDER
    RETURNING TOTAL
END FUNCTION.
```

Multiple logical results should use a `STRUCT`, record, class, or explicit reference parameters rather than introducing a second incompatible return mechanism.

## 8. Common use across callables

The same rules apply to:

- `FUNCTION`,
- `METHOD`,
- callable values,
- `CREATE ... USING ... AS ...`,
- compatible external `CALL` forms.

Constructors therefore reuse `USING` parameter semantics rather than defining a separate constructor-argument system.

## 9. External ABI

Neo COBOL source keeps the same `USING`, `RETURNING`, `BY VALUE`, and `BY REFERENCE` vocabulary when interoperating with external COBOL or C code.

Target-specific ABI details such as register usage, stack layout, symbol decoration, and native calling convention selection are backend concerns.

Backends must preserve the observable Neo COBOL parameter semantics and must not require source-level syntax to change merely because the selected target uses a different ABI.

## 10. Compatibility

Traditional COBOL parameter forms remain valid where applicable. Neo COBOL normalizes compatible traditional and object-oriented forms to one internal parameter model.

## 11. Open items

- Detailed omitted-state inspection syntax for `OPTIONAL` parameters without `VALUE`
- Backend-specific COBOL/C interop mapping details
