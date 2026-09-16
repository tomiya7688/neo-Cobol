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

```cobol
01 SCORES TYPE INTEGER OCCURS 10.
```

Traditional fixed-size `OCCURS` is a fixed-length array whose element type is the declared logical `TYPE`.

Array indexing is one-based by default, matching COBOL table conventions.

Variable-length arrays use an explicit dynamic form:

```cobol
01 ITEMS TYPE CUSTOMER-DATA OCCURS DYNAMIC.
```

A dynamic `OCCURS` value has runtime-managed length and capacity. Exact allocation and growth strategy are implementation details and must not change language semantics.

Traditional `OCCURS ... DEPENDING ON ...` remains supported for COBOL-compatible externally described layouts. It is not the preferred general-purpose dynamic collection mechanism.

## 3. Multidimensional data

Nested `OCCURS` declarations define multidimensional data, preserving COBOL's hierarchical data-description style.

```cobol
01 BOARD.
    05 ROW TYPE INTEGER OCCURS 8.
        10 CELL TYPE INTEGER OCCURS 8.
```

No separate bracket-based multidimensional type syntax is required.

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

## 6. No redundant collection or enum syntax

Neo COBOL does not initially add separate `ARRAY`, `LIST`, or `ENUM` declaration keywords when existing COBOL concepts can express the same intent clearly.

- arrays/tables use `OCCURS`,
- records use level-number groups or reusable `STRUCT`,
- named state predicates use level-88 condition names.

Additional collection abstractions may later exist as library types without changing the core data-description grammar.

## 7. Open items

- Exact syntax for resizing and appending to `OCCURS DYNAMIC`
- Bounds-checking failure behavior
- Closed-set marker for exhaustive level-88 groups, if needed
- Aggregate literal/constructor syntax
- ABI/layout rules for externally interoperable structures
