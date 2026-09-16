# Neo COBOL Variables and Bindings

Status: Draft

## 1. Explicit typed declarations

Traditional level-number declarations remain valid and may use Neo COBOL `TYPE` names.

```cobol
01 AGE TYPE INTEGER.
01 NAME TYPE STRING.
```

## 2. VAR

`VAR` declares a local binding whose type is inferred from its initializer. The binding is mutable and may later be assigned another value compatible with the inferred type.

```cobol
VAR COUNT VALUE 10.
VAR NAME VALUE "KADOKA".
```

Conceptually, the compiler infers `COUNT` as `INTEGER` and `NAME` as `STRING`.

`VAR` itself is not a runtime type.

## 3. LET

`LET` declares a local binding whose type is inferred from its initializer. The binding cannot be reassigned after initialization.

```cobol
LET MAX-COUNT VALUE 100.
LET TITLE VALUE "REPORT".
```

`LET` is the canonical Neo COBOL constant-binding mechanism for local values. A separate `CONST` keyword is not introduced because it would duplicate the same core meaning.

The value referred to by a `LET` binding may still contain mutable referenced state; `LET` guarantees that the binding itself is not reassigned.

`LET` itself is not a runtime type.

## 4. Inference rule

Both forms require an initializer.

```ebnf
var-declaration = "VAR", identifier, "VALUE", expression, sentence-terminator ;
let-declaration = "LET", identifier, "VALUE", expression, sentence-terminator ;
```

Type inference must produce one deterministic logical type. Ambiguous inference is a compile-time error.

Integer literals infer to `INTEGER` when the value fits 32 bits, otherwise `LONG` when it fits 64 bits. Decimal literals with a fractional component infer to `DECIMAL` unless an explicit floating-point suffix/form is later specified. String literals infer to `STRING`; `TRUE` and `FALSE` infer to `BOOLEAN`.

## 5. Assignment

A `VAR` binding may be updated using normal Neo COBOL assignment rules.

```cobol
MOVE 20 TO COUNT.
```

Reassignment of a `LET` binding is a compile-time error.

## 6. Scope

`VAR` and `LET` are block-scoped declarations.

An inner scope may shadow an outer binding only when the new declaration is explicit. The compiler should warn by default when shadowing could make code misleading.

Class/member/global data continues to use explicit data-description declarations unless a later specification explicitly permits inferred member declarations.

## 7. Definite assignment

A binding must be definitely initialized before use.

Because `VAR` and `LET` require an initializer, they satisfy this rule at declaration. Traditional declarations without an explicit `VALUE` follow the default-initialization rules of their declared type/data representation.

## 8. Closures

A closure may capture `VAR` or `LET` bindings.

- captured `LET` bindings cannot be rebound,
- captured `VAR` bindings preserve normal mutability,
- the runtime must preserve captured storage for as long as the closure can access it.

## 9. Open items

- Whether inferred declarations are ever allowed for class/static members
- Destructuring syntax, if later needed
- Detailed lifetime optimization rules for captured locals
