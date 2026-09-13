# Neo COBOL Functions and Callables

Status: Draft

Neo COBOL keeps COBOL 2002-style named functions and adds first-class function values, anonymous functions, and closures.

## 1. Named functions

Named functions use explicit English-oriented clauses.

```cobol
FUNCTION ADD-NUMBERS USING A B RETURNING RESULT
    COMPUTE RESULT = A + B
END FUNCTION.
```

Traditional `FUNCTION-ID.` forms remain valid where supported.

## 2. Anonymous functions

Neo COBOL uses an unnamed `FUNCTION ... END FUNCTION` form instead of a symbolic lambda operator.

```cobol
FUNCTION USING X
    RETURN X + 1
END FUNCTION
```

The absence of a function name distinguishes an anonymous function expression from a named function definition.

Preliminary grammar:

```ebnf
anonymous-function = "FUNCTION",
                     [ "USING", identifier, { identifier } ],
                     [ "RETURNING", type-reference ],
                     statement-list,
                     "END", "FUNCTION" ;
```

Anonymous functions are expressions.

## 3. Function values

Function values may be:

- stored in variables,
- passed as arguments,
- returned from functions,
- used anywhere a compatible function type is expected.

Example:

```cobol
MOVE FUNCTION USING X
         RETURN X + 1
     END FUNCTION
TO INCREMENT.
```

Function type syntax is defined in `TYPES.md`.

## 4. Closures

Anonymous functions may capture variables from an enclosing lexical scope.

```cobol
01 LIMIT PIC 9(3) VALUE 100.

MOVE FUNCTION USING VALUE
         RETURN VALUE IS GREATER THAN LIMIT
     END FUNCTION
TO IS-LARGE.
```

`LIMIT` is captured by the closure.

The runtime must keep required captured state alive for as long as the closure remains reachable.

## 5. Compatibility and lowering

Anonymous functions and closures are Neo COBOL extensions.

When targeting traditional COBOL, the compiler may lower them to generated named functions/procedures plus compiler-managed environment data where faithful translation is possible.

## 6. Open items

- Exact capture-by-value/reference rules
- Mutation of captured variables
- Function overload interaction
- Function signature variance
- Callable invocation syntax in all contexts
- Lowering limits for COBOL targets
