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
1A
123-CUSTOMER
```

Invalid examples:

```text
-CUSTOMER
_CUSTOMER
CUSTOMER-
CUSTOMER_
12345
```

EBNF:

```ebnf
identifier = letter
           | letter, { identifier-middle }, identifier-end
           | digit, { identifier-prefix-continuation }, letter
           | digit, { identifier-prefix-continuation }, letter,
             { identifier-middle }, identifier-end ;

identifier-prefix-continuation = digit | "-" | "_" ;
identifier-middle = letter | digit | "-" | "_" ;
identifier-end = letter | digit ;
```

The alternatives allow digit-leading identifiers only when they contain a letter, and require every multi-character identifier to end with a letter or digit. The grammar therefore rejects identifiers made only of digits or ending in a hyphen or underscore.

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

When a digit-leading source span could match both a numeric literal and an identifier, the lexer selects the longest complete match; a numeric literal wins a tie. For example, `1E3` is a number and `123ABC` is an identifier. A leading `+` or `-` belongs to a signed numeric literal and never starts an identifier.

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
period = "." ;
sentence-terminator = period | omitted-statement-boundary ;
```

An omitted terminator is a parser-context boundary, not a whitespace rule. It is permitted only after a complete executable statement or local `VAR`/`LET` binding whose grammar allows omission, and only when the next token is unambiguous:

- end of file after a complete statement;
- the first token of another statement recognized by the grammar;
- a block marker such as `ELSE`, `END-IF`, or `END-PERFORM` that is valid at that point.

A line break or indentation change alone never ends a statement. The statement keyword, not the line break, makes this boundary unambiguous:

```cobol
MOVE SOURCE TO TARGET
DISPLAY "DONE"
```

Conversely, a line break does not split a `DISPLAY` operand list:

```cobol
DISPLAY FIRST-VALUE
SECOND-VALUE.
```

This is one `DISPLAY` statement with two operands.

The period remains required after division and section headers, `PROGRAM-ID` declarations, and traditional level-number data descriptions. A local `VAR` or `LET` binding may omit its period where its grammar permits. If a statement boundary cannot be determined from the grammar, the source is invalid; the parser must not guess or silently change meaning.

## 8. Reserved words

Reserved words are case-insensitive and cannot be used as identifiers, regardless of capitalization. This initial set covers words used by language forms described in the current drafts, including forms whose detailed grammar remains open. Words mentioned only under Open items are not reserved until their syntax is specified.

- Program structure: `DATA`, `DIVISION`, `IDENTIFICATION`, `PROGRAM-ID`, `PROCEDURE`, `SECTION`, `WORKING-STORAGE`.
- Data declarations and types: `BOOLEAN`, `BYTE`, `DECIMAL`, `DEPENDING`, `DOUBLE`, `FLOAT`, `INTEGER`, `LONG`, `NULL`, `NULLABLE`, `OCCURS`, `PIC`, `SEQUENCE`, `STRUCT`, `STRING`, `TYPE`, `VALUE`.
- Statements and conditions: `ADD`, `AND`, `AS`, `CALL`, `COMPUTE`, `CONVERT`, `CREATE`, `DESTROY`, `DISPLAY`, `ELSE`, `END`, `END-CLASS`, `END-IF`, `END-PERFORM`, `EQUAL`, `ERROR`, `EXCEPTION`, `GREATER`, `IF`, `IS`, `LAST`, `LENGTH`, `LESS`, `MOVE`, `NOT`, `OF`, `ON`, `OR`, `PERFORM`, `REMOVE`, `RESIZE`, `THAN`, `TO`.
- Functions, parameters, classes, and namespaces: `ABSTRACT`, `BY`, `CLASS`, `CLASS-ID`, `END CLASS`, `END FUNCTION`, `END INTERFACE`, `END METHOD`, `END NAMESPACE`, `FALSE`, `FUNCTION`, `IMPLEMENTS`, `INHERITS`, `INITIALIZE`, `INTERFACE`, `LET`, `METHOD`, `METHOD-ID`, `NAMESPACE`, `OPTIONAL`, `OVERRIDE`, `PRIVATE`, `PROTECTED`, `PUBLIC`, `REFERENCE`, `RETURN`, `RETURNING`, `SEALED`, `STATIC`, `TRUE`, `USE`, `USING`, `VAR`.

Multiword forms such as `END IF`, `END PERFORM`, `END CLASS`, and `END FUNCTION` are sequences of reserved words. Where a hyphenated spelling is defined, it is one reserved token: `PROGRAM-ID`, `WORKING-STORAGE`, `CLASS-ID`, `METHOD-ID`, `END-CLASS`, `END-IF`, or `END-PERFORM`. Spaced and hyphenated block terminators accepted for the same construct normalize identically.

## 9. Open items

- Reserved words for syntax whose grammar is still open
- Maximum identifier length
- Unicode identifier policy
- Exact period-omission boundary rules
- Digit separators
- Locale-dependent numeric conventions
- Legacy COBOL literal edge cases
