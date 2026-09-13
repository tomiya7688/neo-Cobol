# Neo COBOL Lexical Specification

Status: Draft

This document defines lexical rules already decided for Neo COBOL.

## 1. Case sensitivity

Keywords and identifiers are case-insensitive.

```cobol
MOVE VALUE TO RESULT.
move value to result.
Move Value To Result.
```

These forms are equivalent.

## 2. Identifiers

Neo COBOL identifiers preserve COBOL-style hyphenated names and additionally permit underscores.

Rules:

- An identifier must contain at least one letter.
- Letters, digits, `-`, and `_` are permitted.
- An identifier must not begin with `-` or `_`.
- An identifier must not end with `-` or `_`.
- An identifier must not consist only of digits.
- Name comparison is case-insensitive.

Valid examples:

```text
CUSTOMER-NAME
customer_name
CUSTOMER-01
GET_CUSTOMER-NAME
A1
```

Invalid examples:

```text
-CUSTOMER
_CUSTOMER
CUSTOMER-
CUSTOMER_
12345
```

Preliminary EBNF:

```ebnf
identifier = identifier-start,
             { identifier-middle },
             identifier-end
           | letter ;

identifier-start = letter | digit ;
identifier-middle = letter | digit | "-" | "_" ;
identifier-end = letter | digit ;
```

The complete identifier must additionally contain at least one letter.

## 3. Comments

Neo COBOL uses the free-format COBOL comment marker `*>`.

`*>` begins a comment that continues to the end of the physical source line.

```cobol
*> Full-line comment.
MOVE SOURCE-VALUE TO TARGET-VALUE. *> Inline comment.
```

`*>` inside a character-string literal is not a comment marker.

```ebnf
comment = "*>", { any-character-except-line-terminator } ;
```

`//` and `/* ... */` are not canonical Neo COBOL comment syntax.

## 4. Character-string literals

Both single and double quotation marks are accepted.

```cobol
DISPLAY "HELLO".
DISPLAY 'HELLO'.
```

A matching quote inside the literal is represented by doubling it.

```cobol
DISPLAY "He said ""HELLO"".".
DISPLAY 'It''s fine.'.
```

## 5. Numeric literals

Decimal integer and decimal literals may be signed or unsigned.

```text
0
123
-123
+123
12.34
-0.5
+10.0
```

Scientific notation uses `E` or `e`.

```text
1E3
1.25E6
-2.5e-4
+6E+10
```

Neo COBOL also supports explicit based integer notation:

```text
0b1010   *> binary
0o755    *> octal
0xFF     *> hexadecimal
```

The radix prefix is case-insensitive.

## 6. Boolean and null literals

The following built-in literals are defined:

```text
TRUE
FALSE
NULL
```

They are case-insensitive keywords.

`TRUE` and `FALSE` have Boolean type. The type-system rules for `NULL` are specified in `TYPES.md`.

## 7. Sentence terminator

A period (`.`) is the canonical sentence terminator.

```ebnf
sentence-terminator = [ "." ] ;
```

A period may be omitted only when the parser can determine the statement boundary unambiguously from block structure, line structure, end of file, or another explicit syntactic boundary.

Omitting a period must never silently change the meaning of a valid program.

## 8. Open items

- Reserved-word set
- Maximum identifier length
- Unicode identifier policy
- Exact period-omission boundary rules
- Digit separators
- Locale-dependent numeric conventions
- Legacy COBOL literal edge cases
