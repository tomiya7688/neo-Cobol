# Neo COBOL Data Model

Status: Draft

Neo COBOL keeps COBOL data-description concepts where they already express the required idea and modernizes their semantics rather than introducing redundant syntax.

## 1. Records and STRUCT

Traditional level-number records remain valid:

```cobol
01 CUSTOMER.
    05 ID TYPE INTEGER.
    05 NAME TYPE STRING PIC X(40).
```

`STRUCT` provides a reusable named aggregate value type:

```cobol
STRUCT CUSTOMER-DATA
    01 ID TYPE INTEGER.
    01 NAME TYPE STRING.
END STRUCT.

01 CUSTOMER TYPE CUSTOMER-DATA.
```

A `STRUCT` is a value type. It has no class identity, inheritance, virtual dispatch, or object lifetime semantics.

Traditional record layouts and `STRUCT` share the same aggregate field model internally. A compiler may normalize both to the same aggregate AST representation where their semantics are equivalent.

## 2. Fixed-length arrays with OCCURS

Neo COBOL keeps COBOL `OCCURS` for fixed-size arrays/tables.

```cobol
01 SCORES TYPE INTEGER OCCURS 10.
```

An `OCCURS` declaration has a fixed element count as part of its data description. Array indexing is one-based by default, matching COBOL table conventions.

Traditional `OCCURS ... DEPENDING ON ...` remains supported for COBOL-compatible externally described layouts. It is treated as a layout-oriented COBOL compatibility feature, not as Neo COBOL's general-purpose resizable sequence type.

## 3. Variable-length sequences

Neo COBOL uses a separate `SEQUENCE OF` type for ordered, resizable collections.

```cobol
01 ITEMS TYPE SEQUENCE OF CUSTOMER-DATA.
01 SCORES TYPE SEQUENCE OF INTEGER.
```

A sequence:

- preserves insertion order,
- has a runtime length,
- may grow or shrink,
- uses one-based indexing by default,
- performs bounds checking,
- manages capacity internally.

The runtime allocation and growth strategy are implementation details and must not change observable language semantics.

### 3.1 Sequence operations

Appending an element:

```cobol
ADD NEW-CUSTOMER TO ITEMS.
```

Removing an element by position:

```cobol
REMOVE ITEM-INDEX FROM ITEMS.
```

Removing the final element:

```cobol
REMOVE LAST FROM ITEMS.
```

Changing the logical length explicitly:

```cobol
RESIZE ITEMS TO 100.
```

Reading the current number of elements:

```cobol
LENGTH OF ITEMS
```

The exact value used when `RESIZE` grows a sequence of value types is defined by each element type's default-initialization rule.

`SEQUENCE OF T` and `T OCCURS n` are intentionally different types. A fixed `OCCURS` field models a fixed data layout; a `SEQUENCE` models a runtime-resizable ordered collection.

## 4. Multidimensional fixed data

Nested `OCCURS` declarations define multidimensional fixed-layout data, preserving COBOL's hierarchical data-description style.

```cobol
01 BOARD.
    05 ROW OCCURS 8.
        10 CELL TYPE INTEGER OCCURS 8.
```

Nested sequences may be used when a resizable multidimensional structure is required:

```cobol
01 ROWS TYPE SEQUENCE OF SEQUENCE OF INTEGER.
```

## 5. Level 88 conditions

Neo COBOL retains level-88 condition names and treats them as the canonical COBOL-style representation for named states over an underlying value.

```cobol
01 ACCOUNT-STATUS TYPE STRING PIC X(10).
    88 STATUS-ACTIVE VALUE "ACTIVE".
    88 STATUS-LOCKED VALUE "LOCKED".
    88 STATUS-CLOSED VALUE "CLOSED".
```

The condition names are Boolean predicates over the parent value and may be used directly in conditions.

Neo COBOL extends tooling around level-88 conditions with exhaustiveness diagnostics where the compiler can prove that a closed set of values is intended. This provides enum-like safety without requiring a redundant `ENUM` type for the same use case.

A future explicit closed-set marker may be added if ordinary level-88 declarations are insufficient to distinguish open and closed value sets.

## 6. Data initialization

`VALUE` remains the canonical data-description initializer keyword.

```cobol
01 RETRY-COUNT TYPE INTEGER VALUE 3.
01 ENABLED TYPE BOOLEAN VALUE TRUE.
```

For aggregate values, each field may define its own default. Aggregate construction syntax may additionally initialize fields explicitly when specified by the expression/object model.

## 7. No redundant collection syntax

Neo COBOL distinguishes two core ordered collection forms by purpose:

- fixed-layout arrays/tables use `OCCURS`,
- runtime-resizable ordered collections use `SEQUENCE OF`.

This separation prevents `OCCURS` from gaining unrelated runtime collection semantics while preserving familiar COBOL data descriptions.

## 8. Open items

- Precise sequence insertion operation, if insertion at arbitrary positions is required
- Bounds-checking failure behavior
- Closed-set marker for exhaustive level-88 groups, if needed
- Aggregate literal/constructor syntax
- ABI/layout rules for externally interoperable structures
