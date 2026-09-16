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

## 2. OCCURS and arrays

Neo COBOL uses COBOL `OCCURS` as its canonical array/table syntax rather than introducing a second bracket-based array declaration language.

### 2.1 Fixed-length arrays

```cobol
01 SCORES TYPE INTEGER OCCURS 10.
```

Traditional fixed-size `OCCURS` is a fixed-length array whose element type is the declared logical `TYPE`.

Array indexing is one-based by default, matching COBOL table conventions.

### 2.2 Dynamic arrays

Neo COBOL provides a general-purpose variable-length array using `OCCURS DYNAMIC`.

```cobol
01 ITEMS TYPE CUSTOMER-DATA OCCURS DYNAMIC.
01 SCORES TYPE INTEGER OCCURS DYNAMIC.
```

A dynamic array has:

- a runtime length,
- implementation-managed capacity,
- one logical element type,
- one-based indexing,
- automatic storage growth when elements are appended.

Capacity management is not observable language semantics. An implementation may use any growth strategy as long as element order, values, and defined failure behavior are preserved.

A newly initialized dynamic array has length zero unless an initializer specifies elements.

### 2.3 Dynamic array operations

Appending an element uses an English-oriented statement:

```cobol
APPEND NEW-CUSTOMER TO ITEMS.
APPEND SCORE TO SCORES.
```

The appended value must be assignment-compatible with the array element type.

Removing an element uses `REMOVE ... FROM ...`:

```cobol
REMOVE LAST FROM ITEMS.
REMOVE ITEM-INDEX FROM ITEMS.
```

`LAST` removes the final element. An integer expression removes the element at that one-based index and shifts following elements toward the beginning.

Explicit resizing uses:

```cobol
RESIZE ITEMS TO 100.
```

When growing, newly created elements receive the element type's defined default value. When shrinking, elements beyond the new length are discarded.

The current element count is read with:

```cobol
LENGTH OF ITEMS
```

For example:

```cobol
IF LENGTH OF ITEMS IS GREATER THAN 0
    REMOVE LAST FROM ITEMS
END-IF.
```

Preliminary grammar:

```ebnf
dynamic-occurs-clause = "OCCURS", "DYNAMIC" ;
append-statement = "APPEND", expression, "TO", identifier, sentence-terminator ;
remove-statement = "REMOVE", ( "LAST" | expression ), "FROM", identifier, sentence-terminator ;
resize-statement = "RESIZE", identifier, "TO", expression, sentence-terminator ;
array-length-expression = "LENGTH", "OF", expression ;
```

Bounds checks are mandatory for dynamic-array indexing and removal. Invalid indexing must not silently access unrelated storage. The exact runtime error object/handling path is specified with the general runtime error model.

### 2.4 OCCURS DEPENDING ON

Traditional `OCCURS ... DEPENDING ON ...` remains supported for COBOL-compatible externally described layouts.

It is distinct from `OCCURS DYNAMIC`:

- `OCCURS ... DEPENDING ON ...` describes a COBOL-style variable logical extent tied to another data item,
- `OCCURS DYNAMIC` is a runtime-managed resizable array intended for normal Neo COBOL programming.

The two forms normalize to different representation requirements even though both may have runtime-varying element counts.

## 3. Multidimensional data

Nested `OCCURS` declarations define multidimensional data, preserving COBOL's hierarchical data-description style.

```cobol
01 BOARD.
    05 ROW TYPE INTEGER OCCURS 8.
        10 CELL TYPE INTEGER OCCURS 8.
```

Dynamic dimensions may also be nested where the element type permits it. No separate bracket-based multidimensional type syntax is required.

## 4. Level 88 conditions

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

## 5. Data initialization

`VALUE` remains the canonical data-description initializer keyword.

```cobol
01 RETRY-COUNT TYPE INTEGER VALUE 3.
01 ENABLED TYPE BOOLEAN VALUE TRUE.
```

For aggregate values, each field may define its own default. Aggregate construction syntax may additionally initialize fields explicitly when specified by the expression/object model.

## 6. Collection design rule

Neo COBOL does not introduce a separate bracket-style core array syntax when `OCCURS` can express the same concept clearly.

- fixed arrays use `OCCURS n`,
- resizable arrays use `OCCURS DYNAMIC`,
- externally described COBOL variable tables may use `OCCURS ... DEPENDING ON ...`,
- records use level-number groups or reusable `STRUCT`,
- named state predicates use level-88 condition names.

Higher-level collections such as maps, sets, queues, and linked structures may be provided as library types rather than being confused with the core array model.

## 7. Open items

- Exact general runtime-error integration for bounds failures
- Aggregate literal/constructor syntax
- Closed-set marker for exhaustive level-88 groups, if needed
- ABI/layout rules for externally interoperable structures
