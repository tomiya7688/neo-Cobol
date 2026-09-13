# Neo COBOL Statements

Status: Draft

This document records statement-level syntax already decided.

## 1. Assignment with MOVE

`MOVE ... TO ...` is the canonical assignment form.

```cobol
MOVE CUSTOMER-NAME TO DISPLAY-NAME.
```

Neo COBOL does not replace the COBOL assignment form with a symbolic `=` assignment syntax.

## 2. DISPLAY

COBOL-style `DISPLAY` is retained.

```cobol
DISPLAY "HELLO WORLD".
```

## 3. Conditions

Conditions prefer COBOL's English-oriented forms.

```cobol
IF AGE IS GREATER THAN OR EQUAL TO 18
    DISPLAY "ADULT"
ELSE
    DISPLAY "MINOR"
END-IF.
```

Canonical relational wording includes:

```text
IS EQUAL TO
IS NOT EQUAL TO
IS GREATER THAN
IS LESS THAN
IS GREATER THAN OR EQUAL TO
IS LESS THAN OR EQUAL TO
```

Logical operators include:

```text
AND
OR
NOT
```

Symbolic comparison operators are not part of the initial canonical grammar.

## 4. IF

COBOL-style explicit block termination is retained.

```cobol
IF CONDITION
    ...
ELSE
    ...
END-IF.
```

`END IF` may be accepted as Neo shorthand and normalized to `END-IF`.

## 5. PERFORM

COBOL-style `PERFORM` is retained.

```cobol
PERFORM
    ...
END-PERFORM.
```

`END PERFORM` may be accepted as Neo shorthand and normalized to `END-PERFORM`.

`PERFORM UNTIL`, `PERFORM VARYING`, and procedure-invocation forms remain to be specified in detail.

## 6. Error handling phrases

Neo COBOL retains COBOL-style error phrases instead of replacing them with unrelated exception syntax.

Initial forms include:

```text
ON ERROR
NOT ON ERROR
ON EXCEPTION
NOT ON EXCEPTION
```

Example direction:

```cobol
CALL "SERVICE"
    ON EXCEPTION
        DISPLAY "CALL FAILED"
END-CALL.
```

Exact attachment rules remain open.

## 7. Object lifecycle statements

Object creation and destruction are defined in `OOP.md` and are executable statements:

```cobol
CREATE CUSTOMER AS CUSTOMER-OBJECT.
CREATE CUSTOMER USING NAME AGE AS CUSTOMER-OBJECT.
DESTROY CUSTOMER-OBJECT.
```

## 8. Open items

- Complete `COMPUTE` grammar
- Complete `CALL` grammar
- `EVALUATE`
- Full `PERFORM` family
- `READ`, `WRITE`, `OPEN`, `CLOSE`
- Error/exception phrase attachment rules
- Detailed I/O and file handling
