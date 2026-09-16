# Neo COBOL Exception Handling

Status: Draft

Neo COBOL modernizes exception handling while preserving COBOL-like vocabulary.

## 1. Raising exceptions

Use `RAISE` to signal an exception.

```cobol
IF AGE IS LESS THAN 0
    RAISE INVALID-AGE
END-IF.
```

Exception values are objects derived from the standard `EXCEPTION` base class.

```cobol
CLASS INVALID-AGE INHERITS EXCEPTION.
END CLASS INVALID-AGE.
```

## 2. Handling exceptions

Neo COBOL retains `ON EXCEPTION` as the canonical handler phrase.

```cobol
PERFORM
    CALL SERVICE
ON EXCEPTION USING ERROR
    DISPLAY ERROR
END-PERFORM.
```

`USING` binds the current exception object to a local name.

## 3. Normal completion

`NOT ON EXCEPTION` remains valid for normal-completion handling where the surrounding statement supports it.

## 4. Propagation

An exception not handled in the current scope propagates to the caller.

A handler may re-raise the current exception with:

```cobol
RAISE.
```

## 5. Cleanup

Cleanup that must execute on both normal and exceptional completion uses `ALWAYS`.

```cobol
PERFORM
    CALL SERVICE
ON EXCEPTION USING ERROR
    DISPLAY ERROR
ALWAYS
    CLOSE RESOURCE
END-PERFORM.
```

`ALWAYS` is the Neo COBOL equivalent of a finally-style cleanup clause.

## 6. Design rule

Neo COBOL does not use `TRY`, `CATCH`, `THROW`, or `FINALLY` as canonical syntax when COBOL-like wording can express the same semantics.

## 7. Open items

- Standard exception class hierarchy
- Exact interaction between statement-specific `ON ERROR` and object exceptions
- Stack trace/runtime diagnostic representation
