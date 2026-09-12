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
11. Neo COBOL extends the COBOL 2002 object model with modern access control, inheritance modifiers, static members, and namespaces.

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

### 2.5 Literals

Neo COBOL keeps ordinary COBOL literal forms and adds explicit modern literals where useful.

Single-quoted and double-quoted strings are accepted. Signed and unsigned decimal integers and decimals are accepted. Scientific notation uses `E` or `e`. Binary, octal, and hexadecimal integers use `0b`, `0o`, and `0x` prefixes. `TRUE`, `FALSE`, and `NULL` are built-in literals.

## 3. Program structure

Traditional COBOL divisions remain valid:

```cobol
IDENTIFICATION DIVISION.
DATA DIVISION.
PROCEDURE DIVISION.
```

Neo COBOL permits inferable divisions to be omitted.

## 4. Statements

### 4.1 MOVE

`MOVE` is the canonical assignment statement.

### 4.2 DISPLAY

COBOL-style `DISPLAY` is retained.

### 4.3 CREATE

Neo COBOL uses an English-oriented object creation statement.

```cobol
CREATE CUSTOMER AS CUSTOMER-OBJECT.
CREATE CUSTOMER USING NAME AGE AS CUSTOMER-OBJECT.
```

```ebnf
create-statement = "CREATE", class-identifier,
                   [ "USING", expression, { expression } ],
                   "AS", identifier,
                   sentence-terminator ;
```

### 4.4 DESTROY

```cobol
DESTROY CUSTOMER-OBJECT.
```

```ebnf
destroy-statement = "DESTROY", identifier, sentence-terminator ;
```

## 5. Conditions

Conditions use COBOL's English-oriented forms such as `IS EQUAL TO`, `IS GREATER THAN`, `AND`, `OR`, and `NOT`.

## 6. IF statement

COBOL-style `IF ... ELSE ... END-IF` is retained. `END IF` may be accepted as Neo shorthand and normalized to `END-IF`.

## 7. PERFORM

COBOL-style `PERFORM` is retained.

## 8. Error handling

Neo COBOL retains COBOL-style error phrases such as `ON ERROR`, `NOT ON ERROR`, `ON EXCEPTION`, and `NOT ON EXCEPTION`.

## 9. Classes, namespaces, and modern OOP

Neo COBOL keeps the COBOL 2002 object-oriented model and extends it with modern OOP facilities.

### 9.1 Class types

A class name also introduces a reference type.

```cobol
01 CUSTOMER-OBJECT TYPE CUSTOMER.
```

A class-typed variable may refer to a live instance of the declared class, an assignment-compatible derived instance, or `NULL`.

### 9.2 Access modifiers

The following access modifiers are part of Neo COBOL:

```text
PUBLIC
PRIVATE
PROTECTED
```

They may be applied to classes, methods, and members where meaningful.

```cobol
PUBLIC CLASS CUSTOMER.

PRIVATE 01 CUSTOMER-ID PIC X(20).
PUBLIC  01 DISPLAY-NAME PIC X(40).

PUBLIC METHOD GET-NAME
    MOVE DISPLAY-NAME TO RESULT
END METHOD.

PRIVATE METHOD CALCULATE-INTERNAL
    ...
END METHOD.

END CLASS CUSTOMER.
```

Section-oriented grouping may also be supported for a more COBOL-like style:

```cobol
PUBLIC SECTION.
    METHOD GET-NAME
    END METHOD.

PRIVATE SECTION.
    METHOD CALCULATE-INTERNAL
    END METHOD.
```

The exact interaction between section visibility and per-member modifiers remains to be finalized.

### 9.3 STATIC

`STATIC` declares a member or method that belongs to the class rather than to a specific instance.

```cobol
PUBLIC STATIC 01 MAX-CUSTOMERS PIC 9(5).

PUBLIC STATIC METHOD CREATE-DEFAULT
    ...
END METHOD.
```

### 9.4 Inheritance

Neo COBOL uses `INHERITS` for class inheritance.

```cobol
PUBLIC CLASS EMPLOYEE INHERITS PERSON.
    ...
END CLASS EMPLOYEE.
```

A derived-class reference is assignment-compatible with an accessible base-class reference subject to the type-system rules.

### 9.5 ABSTRACT

`ABSTRACT` may be applied to a class or method.

```cobol
PUBLIC ABSTRACT CLASS SHAPE.

PUBLIC ABSTRACT METHOD CALCULATE-AREA
    RETURNING AREA
END METHOD.
```

An abstract class cannot be instantiated directly. An abstract method has no concrete implementation in the declaring class and must be implemented by a non-abstract derived class unless another applicable rule is specified.

### 9.6 OVERRIDE

`OVERRIDE` explicitly marks a method that overrides an inherited method.

```cobol
PUBLIC OVERRIDE METHOD CALCULATE-AREA
    ...
END METHOD.
```

The compiler should diagnose an `OVERRIDE` method that does not match an overridable inherited member.

### 9.7 SEALED

`SEALED` prevents further inheritance or overriding, depending on where it is applied.

```cobol
PUBLIC SEALED CLASS CUSTOMER.
    ...
END CLASS CUSTOMER.
```

A sealed class cannot be inherited.

```cobol
PUBLIC SEALED OVERRIDE METHOD VALIDATE
    ...
END METHOD.
```

A sealed method may not be overridden by further derived classes.

### 9.8 Modifier combinations

Preliminary modifier sets are:

```ebnf
access-modifier = "PUBLIC" | "PRIVATE" | "PROTECTED" ;
class-modifier = access-modifier | "ABSTRACT" | "SEALED" ;
member-modifier = access-modifier | "STATIC" | "ABSTRACT" | "OVERRIDE" | "SEALED" ;
```

Illegal combinations such as `ABSTRACT SEALED CLASS` must be rejected unless explicitly defined otherwise in the type-system specification.

### 9.9 Namespace

Neo COBOL provides namespaces for organizing classes, functions, and other public program elements.

```cobol
NAMESPACE COMPANY.CUSTOMER

PUBLIC CLASS CUSTOMER.
    ...
END CLASS CUSTOMER.

END NAMESPACE.
```

Namespace components are case-insensitive identifiers. A period inside a namespace-qualified name is treated as a namespace separator rather than a sentence terminator when it occurs between valid namespace components.

```ebnf
namespace-definition = "NAMESPACE", qualified-name, sentence-terminator,
                       { namespace-member },
                       "END", "NAMESPACE", sentence-terminator ;

qualified-name = identifier, { ".", identifier } ;
```

The exact import/use syntax for accessing members of another namespace remains to be specified.

### 9.10 Preliminary class grammar

```ebnf
class-definition = { class-modifier }, "CLASS", identifier,
                   [ "INHERITS", qualified-name ],
                   sentence-terminator,
                   { class-member },
                   ( "END CLASS" | "END-CLASS" ),
                   [ identifier ], sentence-terminator ;
```

Traditional COBOL 2002 forms such as `CLASS-ID.` and `METHOD-ID.` remain valid where supported; Neo forms normalize to the same internal representation.

## 10. Functions

Functions remain based on COBOL 2002 concepts and use explicit English-oriented parameter and return clauses.

## 11. Division inference

When divisions are omitted, the compiler conceptually inserts canonical COBOL structure before semantic analysis. Executable statements such as `MOVE`, `DISPLAY`, `IF`, `PERFORM`, `CALL`, `CREATE`, and `DESTROY` belong to the Procedure Division. Top-level `CLASS`, `FUNCTION`, and `NAMESPACE` constructs are recognized directly.

## 12. Normalization model

Neo shorthand is normalized before later compiler stages wherever practical.

Examples:

```text
CLASS PERSON        -> CLASS-ID. PERSON.
METHOD PRINT-NAME   -> METHOD-ID. PRINT-NAME.
END IF              -> END-IF
END PERFORM         -> END-PERFORM
```

Omitted divisions and sections are represented internally as explicit AST nodes after parsing.

## 13. Compatibility rule

Where an accepted Neo shorthand conflicts with an established COBOL interpretation, the COBOL interpretation takes precedence unless the Neo COBOL specification explicitly states otherwise.

Neo syntax should extend COBOL rather than silently redefine familiar COBOL syntax.

## 14. Open grammar items

The following remain to be specified:

- Exact period-omission boundary rules
- Full reserved-word set
- Numeric and data-description grammar
- `PIC` and modern type syntax
- Detailed class-type declaration rules
- Constructor/initialization semantics
- Object lifetime and destruction semantics
- Detailed inheritance and conversion rules
- Namespace import/use syntax
- Interfaces, if supported
- Generic or parameterized facilities, if any
- `COMPUTE`
- `CALL`
- `EVALUATE`
- Complete `PERFORM` forms
- File handling
- Interoperability syntax
