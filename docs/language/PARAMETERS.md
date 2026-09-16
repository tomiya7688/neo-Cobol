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

## 5. OPTIONAL

`OPTIONAL` marks a parameter that may be omitted by the caller.

```cobol
METHOD WRITE-REPORT
    USING REPORT
          OPTIONAL FORMATTER
END METHOD.
```

Omitted optional arguments receive their declared default when one exists; otherwise they use the type-defined absent/default behavior specified for that callable declaration.

Required parameters must precede optional parameters unless a later named-argument rule explicitly permits another ordering.

## 6. RETURNING

A callable has at most one direct `RETURNING` result.

```cobol
FUNCTION CALCULATE-TOTAL
    USING ORDER
    RETURNING TOTAL
END FUNCTION.
```

Multiple logical results should use a `STRUCT`, record, class, or explicit reference parameters rather than introducing a second incompatible return mechanism.

## 7. Common use across callables

The same rules apply to:

- `FUNCTION`,
- `METHOD`,
- callable values,
- `CREATE ... USING ... AS ...`,
- compatible external `CALL` forms.

Constructors therefore reuse `USING` parameter semantics rather than defining a separate constructor-argument system.

## 8. Compatibility

Traditional COBOL parameter forms remain valid where applicable. Neo COBOL normalizes compatible traditional and object-oriented forms to one internal parameter model.

## 9. Open items

- Named-argument syntax, if added
- Exact interaction of `OPTIONAL` with explicit default values
- External ABI details for COBOL/C interop
