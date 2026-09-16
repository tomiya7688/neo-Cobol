# Neo COBOL Object-Oriented Programming Specification

Status: Draft

Neo COBOL keeps the COBOL 2002 object-oriented model and extends it with modern OOP facilities while preserving English-readable syntax.

## 1. Classes

Traditional COBOL 2002-style forms remain valid where supported:

```cobol
CLASS-ID. PERSON.
...
END CLASS PERSON.
```

Neo COBOL also accepts the shorter form:

```cobol
PUBLIC CLASS PERSON.
    ...
END CLASS PERSON.
```

A class name introduces a class reference type. See `TYPES.md`.

## 2. Object creation and destruction

Object creation uses an English-oriented statement:

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

The target must be assignment-compatible with the created class type.

Explicit destruction/release uses:

```cobol
DESTROY CUSTOMER-OBJECT.
```

```ebnf
destroy-statement = "DESTROY", identifier, sentence-terminator ;
```

Exact memory reclamation semantics remain a runtime concern.

## 3. Access control

Neo COBOL defines:

```text
PUBLIC
PRIVATE
PROTECTED
```

These may be applied where meaningful to classes, methods, and members.

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

A COBOL-like section grouping form may also be supported:

```cobol
PUBLIC SECTION.
    METHOD GET-NAME
    END METHOD.

PRIVATE SECTION.
    METHOD CALCULATE-INTERNAL
    END METHOD.
```

The exact precedence between section-level and per-member visibility is still open.

## 4. Static members

`STATIC` marks a member or method that belongs to the class rather than an instance.

```cobol
PUBLIC STATIC 01 MAX-CUSTOMERS PIC 9(5).

PUBLIC STATIC METHOD CREATE-DEFAULT
    ...
END METHOD.
```

## 5. Class inheritance

Neo COBOL uses `INHERITS`.

```cobol
PUBLIC CLASS EMPLOYEE INHERITS PERSON.
    ...
END CLASS EMPLOYEE.
```

The initial design uses single class inheritance: a class may directly inherit from at most one base class.

Derived references are assignment-compatible with accessible base-class references according to `TYPES.md`.

## 6. Abstract classes and methods

`ABSTRACT` may modify a class or method.

```cobol
PUBLIC ABSTRACT CLASS SHAPE.

PUBLIC ABSTRACT METHOD CALCULATE-AREA
    RETURNING AREA
END METHOD.
```

An abstract class cannot be instantiated directly. A concrete derived class must implement inherited abstract methods unless another language rule explicitly permits otherwise.

## 7. Override

`OVERRIDE` explicitly marks an overriding method.

```cobol
PUBLIC OVERRIDE METHOD CALCULATE-AREA
    ...
END METHOD.
```

The compiler must diagnose an `OVERRIDE` that does not match an overridable inherited member.

## 8. Sealed classes and methods

`SEALED` prevents further inheritance or overriding depending on the declaration.

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

A sealed method cannot be overridden further.

`ABSTRACT SEALED CLASS` is invalid unless a future specification explicitly defines such a combination.

## 9. Interfaces

Interfaces exist as a narrow contract mechanism rather than as a replacement for classes.

```cobol
PUBLIC INTERFACE PRINTABLE.

PUBLIC METHOD PRINT
END METHOD.

END INTERFACE PRINTABLE.
```

An interface normally declares callable/member contracts and does not contain ordinary instance data fields.

A class implements interfaces using `IMPLEMENTS`:

```cobol
PUBLIC CLASS REPORT IMPLEMENTS PRINTABLE.

PUBLIC METHOD PRINT
    DISPLAY "REPORT"
END METHOD.

END CLASS REPORT.
```

Multiple interfaces are allowed:

```cobol
PUBLIC CLASS REPORT
    INHERITS DOCUMENT
    IMPLEMENTS PRINTABLE SERIALIZABLE AUDITABLE.

END CLASS REPORT.
```

This does not introduce multiple class inheritance.

## 10. Namespaces

Neo COBOL provides namespaces for classes, interfaces, functions, and other public declarations.

```cobol
NAMESPACE COMPANY.CUSTOMER

PUBLIC CLASS CUSTOMER.
    ...
END CLASS CUSTOMER.

END NAMESPACE.
```

Namespace components are case-insensitive identifiers. A `.` between valid namespace components is treated as a namespace separator rather than a sentence terminator.

```ebnf
qualified-name = identifier, { ".", identifier } ;
```

Namespaces are imported with `USE NAMESPACE`:

```cobol
USE NAMESPACE COMPANY.CUSTOMER.
USE NAMESPACE COMPANY.ACCOUNT.
```

`USE NAMESPACE` makes the public declarations of the target namespace available for unqualified name lookup within the current compilation scope.

Neo COBOL does not initially provide a separate single-symbol import syntax. Code may either use `USE NAMESPACE` or refer to a declaration by its fully qualified name.

```cobol
01 CURRENT-CUSTOMER TYPE COMPANY.CUSTOMER.CUSTOMER.
```

Name conflicts introduced by multiple imported namespaces are compile-time errors unless the use site gives a fully qualified name.

```ebnf
use-namespace = "USE", "NAMESPACE", qualified-name, sentence-terminator ;
```

## 11. Property convention

Neo COBOL does not introduce a separate `PROPERTY` declaration syntax in the initial language.

Property-like access is expressed with ordinary methods whose names use the `GET-` and `SET-` conventions.

```cobol
PUBLIC METHOD GET-NAME
    RETURNING RESULT
    MOVE NAME TO RESULT
END METHOD.

PUBLIC METHOD SET-NAME
    USING VALUE
    MOVE VALUE TO NAME
END METHOD.
```

A getter uses normal `RETURNING` semantics. A setter uses normal `USING` parameter semantics.

A `GET-*` method without a corresponding `SET-*` method represents read-only property-like access. A `SET-*` method without a corresponding `GET-*` method represents write-only property-like access.

These are ordinary methods for access control, inheritance, interface contracts, `ABSTRACT`, `OVERRIDE`, and `SEALED` behavior. The compiler does not give `GET-*` or `SET-*` methods separate runtime dispatch semantics.

Tools and IDEs may present a compatible `GET-X` / `SET-X` pair as a logical property named `X`, but that presentation is tooling metadata rather than a distinct language construct.

## 12. Preliminary modifier grammar

```ebnf
access-modifier = "PUBLIC" | "PRIVATE" | "PROTECTED" ;
class-modifier = access-modifier | "ABSTRACT" | "SEALED" ;
member-modifier = access-modifier | "STATIC" | "ABSTRACT" | "OVERRIDE" | "SEALED" ;
```

## 13. Preliminary class grammar

```ebnf
class-definition = { class-modifier }, "CLASS", identifier,
                   [ "INHERITS", qualified-name ],
                   [ implements-clause ],
                   sentence-terminator,
                   { class-member },
                   ( "END CLASS" | "END-CLASS" ),
                   [ identifier ], sentence-terminator ;

implements-clause = "IMPLEMENTS", qualified-name,
                    { qualified-name } ;
```

Traditional `CLASS-ID.` and `METHOD-ID.` forms remain compatible forms and normalize into the same internal representation where supported.

## 14. Open items

- Constructor/initialization semantics behind `CREATE ... USING ... AS ...`
- Object lifetime and exact `DESTROY` semantics
- Explicit casting syntax
- Detailed inheritance conversion rules
- Section-level visibility precedence
- Interface default methods/events policy
- Exact mapping to COBOL targets lacking equivalent modern OOP features
