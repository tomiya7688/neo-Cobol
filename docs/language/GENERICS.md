# Neo COBOL Generics

Status: Draft

Neo COBOL expresses generic type parameters with COBOL-like `USING TYPE` syntax.

## 1. Generic classes

```cobol
CLASS BOX USING TYPE T.
    01 VALUE TYPE T.
END CLASS BOX.
```

Multiple type parameters are written in order:

```cobol
CLASS PAIR USING TYPE T-KEY T-VALUE.
    01 KEY TYPE T-KEY.
    01 VALUE TYPE T-VALUE.
END CLASS PAIR.
```

## 2. Generic functions and methods

```cobol
FUNCTION IDENTITY USING TYPE T USING VALUE TYPE T RETURNING T
    RETURN VALUE
END FUNCTION.
```

Methods may use the same `USING TYPE` form.

## 3. Type application

Parameterized types use `OF` for English-readable application where practical.

```cobol
01 NUMBERS TYPE SEQUENCE OF INTEGER.
```

User-defined generic type application will follow the same English-readable direction rather than symbolic angle-bracket syntax.

## 4. Core collections

Only the resizable ordered collection `SEQUENCE OF T` is part of the language core.

Set, map, queue, stack, and similar collections belong to the standard library rather than the core grammar.

## 5. No symbolic generic syntax

Neo COBOL does not use C++/Java/C#-style `<T>` syntax as canonical grammar.

Generic syntax should remain readable as COBOL-like English.

## 6. Open items

- Constraint syntax for type parameters
- Exact user-defined generic type application syntax where `OF` would be ambiguous
- Variance rules for generic reference types
