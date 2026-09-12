# Neo COBOL Grammar Specification

Status: Draft

This document defines the grammar direction of Neo COBOL. Neo COBOL is intentionally COBOL-like: existing COBOL concepts and spellings are preserved where practical, while redundant structural declarations may be omitted when the compiler can infer them unambiguously.

## 1. Core grammar principles

1. COBOL syntax is the baseline. Neo COBOL extensions should not replace familiar COBOL syntax without a strong reason.
2. Keywords are case-insensitive.
3. English-like expressions are preferred over symbolic alternatives.
4. Explicit COBOL block forms such as `IF ... END-IF` and `PERFORM ... END-PERFORM` are retained.
5. `MOVE ... TO ...` is the canonical assignment form.
6. A period (`.`) is the canonical sentence terminator, but may be omitted where the statement boundary is unambiguous.
7. Traditional `DIVISION` and `SECTION` declarations remain valid but may be omitted where their role can be inferred.
8. `CLASS`, `METHOD`, and `FUNCTION` facilities follow COBOL 2002 style as closely as practical, with optional Neo COBOL shorthand.
9. Error-handling phrases use COBOL-style forms such as `ON ERROR` and `ON EXCEPTION` rather than introducing unrelated exception syntax.
10. When Neo shorthand can be expanded mechanically into standard COBOL-like structure, that expansion defines its intended meaning.

## 2. Lexical rules

### 2.1 Case

Keywords and identifiers are case-insensitive for name resolution.

The following are equivalent:

```cobol
MOVE VALUE TO RESULT.
move value to result.
Move Value To Result.
```

A formatter may choose a canonical keyword style, but source spelling is not semantically significant.

### 2.2 Identifiers

Neo COBOL identifiers preserve conventional COBOL hyphenated naming and additionally allow underscores.

An identifier:

- must contain at least one letter,
- may contain letters, digits, hyphens (`-`), and underscores (`_`),
- must not begin with `-` or `_`,
- must not end with `-` or `_`,
- must not consist only of digits,
- is case-insensitive.

Examples of valid identifiers:

```text
CUSTOMER-NAME
customer_name
CUSTOMER-01
GET_CUSTOMER-NAME
A1
```

Examples of invalid identifiers:

```text
-CUSTOMER
_CUSTOMER
CUSTOMER-
CUSTOMER_
12345
```

Preliminary lexical form:

```ebnf
identifier = identifier-start,
             { identifier-middle },
             identifier-end
           | letter ;

identifier-start = letter | digit ;
identifier-middle = letter | digit | "-" | "_" ;
identifier-end = letter | digit ;
```

The additional semantic constraint applies that the complete identifier must contain at least one letter. Reserved words, maximum identifier length, and Unicode identifier policy remain to be specified.

### 2.3 Sentence terminator

```ebnf
sentence-terminator = [ "." ] ;
```

A period is preferred and always valid where COBOL permits it. Omission is permitted only when the parser can determine the end of a statement from block structure, line structure, end of file, or another unambiguous syntactic boundary.

A missing period must never change the meaning of an otherwise valid program.

### 2.4 Comments

Neo COBOL uses the standard free-format COBOL comment marker `*>`.

`*>` begins a comment and the comment continues to the end of the physical source line. It may appear at the beginning of a line or after source code.

```cobol
*> This is a full-line comment.
MOVE SOURCE-VALUE TO TARGET-VALUE. *> This is an inline comment.
```

Comments are lexical whitespace and have no semantic effect. The comment marker is not recognized while it appears inside a character-string literal.

```ebnf
comment = "*>", { any-character-except-line-terminator } ;
```

Neo COBOL does not introduce `//` or `/* ... */` as canonical comment syntax.

## 3. Program structure

Traditional COBOL divisions remain valid:

```cobol
IDENTIFICATION DIVISION.
DATA DIVISION.
PROCEDURE DIVISION.
```

Neo COBOL permits inferable divisions to be omitted.

For example:

```cobol
PROGRAM-ID. HELLO.

01 NAME PIC X(20) VALUE "KADOKA".

DISPLAY "HELLO " NAME.
```

is interpreted as if the omitted structure had been written explicitly:

```cobol
IDENTIFICATION DIVISION.
PROGRAM-ID. HELLO.

DATA DIVISION.
WORKING-STORAGE SECTION.
01 NAME PIC X(20) VALUE "KADOKA".

PROCEDURE DIVISION.
DISPLAY "HELLO " NAME.
```

The exact inference rules will be specified separately. Inference must be deterministic.

A preliminary top-level grammar is:

```ebnf
program = [ identification-division ],
          [ data-division ],
          [ procedure-division ],
          { class-definition | function-definition } ;
```

## 4. Statements

### 4.1 MOVE

`MOVE` is the canonical assignment statement.

```ebnf
move-statement = "MOVE", expression, "TO", identifier,
                 { identifier }, sentence-terminator ;
```

Example:

```cobol
MOVE CUSTOMER-NAME TO DISPLAY-NAME.
```

### 4.2 DISPLAY

```ebnf
display-statement = "DISPLAY", expression,
                    { expression }, sentence-terminator ;
```

Example:

```cobol
DISPLAY "HELLO WORLD".
```

## 5. Conditions

Conditions use COBOL's English-oriented forms.

```ebnf
condition = relation-condition
          | condition, "AND", condition
          | condition, "OR", condition
          | "NOT", condition ;

relation-condition = expression, relation-operator, expression ;

relation-operator = "IS", [ "NOT" ],
                    ( "EQUAL TO"
                    | "GREATER THAN"
                    | "LESS THAN"
                    | "GREATER THAN OR EQUAL TO"
                    | "LESS THAN OR EQUAL TO" ) ;
```

Example:

```cobol
IF AGE IS GREATER THAN OR EQUAL TO 18
    DISPLAY "ADULT"
ELSE
    DISPLAY "MINOR"
END-IF.
```

Symbolic comparison operators are not part of the initial canonical grammar.

## 6. IF statement

```ebnf
if-statement = "IF", condition,
               statement-list,
               [ "ELSE", statement-list ],
               ( "END-IF" | "END IF" ),
               sentence-terminator ;
```

Both `END-IF` and the spaced Neo form `END IF` may be accepted. The compiler normalizes them to the same internal construct.

## 7. PERFORM

COBOL-style `PERFORM` is retained.

```ebnf
perform-statement = "PERFORM",
                    statement-list,
                    ( "END-PERFORM" | "END PERFORM" ),
                    sentence-terminator ;
```

Additional COBOL forms such as `PERFORM UNTIL`, `PERFORM VARYING`, and procedure invocation will be specified in later revisions.

## 8. Error handling

Neo COBOL retains COBOL-style error phrases.

Initial supported forms include:

```text
ON ERROR
NOT ON ERROR
ON EXCEPTION
NOT ON EXCEPTION
```

These phrases are attached only to statements for which the corresponding condition is meaningful.

Example direction:

```cobol
CALL "SERVICE"
    ON EXCEPTION
        DISPLAY "CALL FAILED"
END-CALL.
```

Exact attachment rules are still provisional.

## 9. Classes and methods

Neo COBOL includes object-oriented facilities based primarily on COBOL 2002.

Traditional forms remain valid where supported:

```cobol
CLASS-ID. PERSON.

OBJECT.

METHOD-ID. GET-NAME.
PROCEDURE DIVISION RETURNING RESULT.
    MOVE NAME TO RESULT.
END METHOD GET-NAME.

END OBJECT.
END CLASS PERSON.
```

Neo COBOL may also accept a shorter spelling that expands to the same structure:

```cobol
CLASS PERSON

METHOD GET-NAME RETURNING RESULT
    MOVE NAME TO RESULT
END METHOD

END CLASS
```

Preliminary grammar:

```ebnf
class-definition = class-header,
                   { data-description | method-definition },
                   class-end ;

class-header = ( "CLASS-ID.", identifier, sentence-terminator )
             | ( "CLASS", identifier, sentence-terminator ) ;

class-end = "END", [ "CLASS" ], [ identifier ], sentence-terminator
          | "END-CLASS", [ identifier ], sentence-terminator ;
```

The exact compatibility rules with standard COBOL object syntax remain to be finalized.

## 10. Functions

Functions are based on COBOL 2002 concepts and use explicit English-oriented parameter and return clauses.

Example Neo form:

```cobol
FUNCTION ADD-NUMBERS USING A B RETURNING RESULT
    COMPUTE RESULT = A + B
END FUNCTION.
```

Preliminary grammar:

```ebnf
function-definition = function-header,
                      statement-list,
                      function-end ;

function-header = ( "FUNCTION-ID.", identifier, sentence-terminator
                  | "FUNCTION", identifier )
                  [ "USING", identifier, { identifier } ]
                  [ "RETURNING", identifier ] ;

function-end = "END", "FUNCTION", [ identifier ], sentence-terminator
             | "END-FUNCTION", [ identifier ], sentence-terminator ;
```

The final relationship between `FUNCTION-ID`, intrinsic functions, user-defined functions, and methods remains to be specified.

## 11. Division inference

When divisions are omitted, the compiler conceptually inserts canonical COBOL structure before semantic analysis.

Initial rules:

- Program metadata such as `PROGRAM-ID` belongs to the Identification Division.
- Level-number data declarations such as `01`, `05`, and related entries belong to the Data Division.
- Executable statements such as `MOVE`, `DISPLAY`, `IF`, `PERFORM`, and `CALL` belong to the Procedure Division.
- Top-level `CLASS`, `METHOD`, and `FUNCTION` constructs are recognized directly and are not inferred from arbitrary statement text.
- If a source fragment could validly belong to more than one implicit region, the compiler must reject it rather than guess.

## 12. Normalization model

Neo shorthand is normalized before later compiler stages wherever practical.

Examples:

```text
CLASS PERSON        -> CLASS-ID. PERSON.
METHOD PRINT-NAME   -> METHOD-ID. PRINT-NAME.
END IF              -> END-IF
END PERFORM         -> END-PERFORM
```

Similarly, omitted divisions and sections are represented internally as explicit AST nodes after parsing.

This allows later compiler stages and backends to operate on one canonical representation rather than maintaining separate Neo and traditional COBOL semantics.

## 13. Compatibility rule

Where an accepted Neo shorthand conflicts with an established COBOL interpretation, the COBOL interpretation takes precedence unless the Neo COBOL specification explicitly states otherwise.

Neo syntax should extend COBOL rather than silently redefine familiar COBOL syntax.

## 14. Open grammar items

The following remain to be specified:

- Exact period-omission boundary rules
- Full lexical grammar and reserved-word set
- Literal syntax
- Numeric and data-description grammar
- `PIC` and modern type syntax, if any
- `COMPUTE`
- `CALL`
- `EVALUATE`
- complete `PERFORM` forms
- `READ`, `WRITE`, `OPEN`, `CLOSE`, and file handling
- exception/error phrase attachment rules
- modules and imports
- visibility
- inheritance and interfaces, if supported
- generic or parameterized facilities, if any
- interoperability syntax

This document will be expanded incrementally as these decisions are made.
