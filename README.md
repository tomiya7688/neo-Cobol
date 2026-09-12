# Neo COBOL

**Neo COBOL** is a modern programming language inspired by the original philosophy of COBOL.

## Goals

- Programs should be readable as English.
- Can be down-compiled to COBOL.
- Has its own compiler.
- Can transpile to Bitlang.
- Keep the intent of business logic explicit and easy to review.

## Core design principle

Neo COBOL source code should read as naturally as practical as an English description of the program's intent.

Natural readability is not only a style preference; it is a language-design acceptance criterion. New syntax should be evaluated not only for brevity and implementation cost, but also for whether a reader can understand its intent without mentally translating dense symbolic notation.

## Planned toolchain

Neo COBOL source may target multiple backends:

```text
Neo COBOL
  ├─> COBOL
  ├─> C (GCC / Clang)
  ├─> Bitlang
  └─> Native Neo COBOL compiler pipeline
```

The exact compatibility level and lowering rules will be defined in `SPECIFICATION.md`.

## Project status

Early design phase. Syntax, type system, module model, compatibility rules, and compiler architecture are still under active specification.

See:

- `SPECIFICATION.md` — language specification draft
- `ROADMAP.md` — implementation roadmap
- GitHub Issues — individual design and implementation tasks

## License

MIT License. See `LICENSE`.
