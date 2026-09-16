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

## Compiler bootstrap

The initial compiler is implemented in **Go** and lowers through a backend-independent **Neo IR (NIR)**. The first backend emits **C11**, which is then compiled by GCC, Clang, or another compatible C compiler.

```text
Neo COBOL -> Lexer -> Parser -> AST -> Sema -> NIR -> C11 -> native executable
```

The executable bootstrap currently supports literal `DISPLAY`, elementary level-01/77 data declarations, `TYPE`, an initial `PIC` subset, `VALUE`, symbol resolution, `MOVE`, and variable-backed `DISPLAY`.

```text
neoc check source.ncob
neoc emit-c source.ncob
neoc build source.ncob
neoc run source.ncob
neoc version
```

Try it with:

```sh
go run ./cmd/neoc run examples/hello.ncob
go run ./cmd/neoc run examples/variables.ncob
```

Compiler architecture and current bootstrap limitations are documented in `docs/compiler/ARCHITECTURE.md`.

## Planned toolchain

Neo COBOL source may target multiple backends:

```text
Neo COBOL / NIR
  ├─> COBOL
  ├─> C (GCC / Clang)
  ├─> Bitlang
  └─> future native / VM backends
```

The exact compatibility level and lowering rules are defined incrementally in `docs/language/` and `docs/compiler/`.

## Project status

Early implementation and active language-design phase. The compiler has an executable end-to-end path, while most language features remain under specification and implementation.

See:

- `docs/language/SPECIFICATION.md` — language specification draft
- `docs/language/GRAMMAR.md` — grammar and EBNF draft
- `docs/compiler/ARCHITECTURE.md` — compiler architecture and bootstrap scope
- `ROADMAP.md` — implementation roadmap
- GitHub Issues — individual design and implementation tasks

## License

MIT License. See `LICENSE`.
