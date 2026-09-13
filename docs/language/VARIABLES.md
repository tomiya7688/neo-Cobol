# Neo COBOL Variables and Bindings

Status: Draft

## 1. Explicit typed declarations

Traditional level-number declarations remain valid and may use Neo COBOL `TYPE` names.

```cobol
01 AGE TYPE INT.
01 NAME TYPE STR.
```

## 2. VAR

`VAR` declares a local binding whose type is inferred from its initializer. The binding is mutable and may later be assigned another value compatible with the inferred type.

```cobol
VAR COUNT VALUE 10.
VAR NAME VALUE "KADOKA".
```

Conceptually, the compiler infers `COUNT` as `INT` and `NAME` as `STR`.

`VAR` itself is not a runtime type.

## 3. LET

`LET` declares a local binding whose type is inferred from its initializer. The binding cannot be reassigned after initialization.

```cobol
LET MAX-COUNT VALUE 100.
LET TITLE VALUE "REPORT".
```

The value referred to by a `LET` binding may still have its own mutability rules; `LET` guarantees that the binding itself is not reassigned.

`LET` itself is not a runtime type.

## 4. Inference rule

Both forms require an initializer unless a later specification explicitly defines another inference source.

```ebnf
var-declaration = "VAR", identifier, "VALUE", expression, sentence-terminator ;
let-declaration = "LET", identifier, "VALUE", expression, sentence-terminator ;
```

Type inference must produce one deterministic logical type. Ambiguous inference is a compile-time error.

## 5. Assignment

A `VAR` binding may be updated using normal Neo COBOL assignment rules.

```cobol
MOVE 20 TO COUNT.
```

Reassignment of a `LET` binding is a compile-time error.

## 6. Scope

`VAR` and `LET` are intended primarily for local/block-scoped declarations. Exact class/member/global applicability remains to be specified separately.

## 7. Open items

- Definite-assignment rules
- Shadowing rules
- Global/member use of `VAR` and `LET`
- Destructuring, if added
- Interaction with captured closure variables
