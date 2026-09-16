package compiler

import (
	"github.com/tomiya7688/neo-Cobol/internal/ast"
	cbackend "github.com/tomiya7688/neo-Cobol/internal/backend/c"
	"github.com/tomiya7688/neo-Cobol/internal/lexer"
	"github.com/tomiya7688/neo-Cobol/internal/nir"
	"github.com/tomiya7688/neo-Cobol/internal/parser"
	"github.com/tomiya7688/neo-Cobol/internal/sema"
)

func Check(source string) (*ast.Program, error) {
	tokens, err := lexer.Lex(source)
	if err != nil {
		return nil, err
	}
	program, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}
	if err := sema.Check(program); err != nil {
		return nil, err
	}
	return program, nil
}

func EmitC(source string) (string, error) {
	program, err := Check(source)
	if err != nil {
		return "", err
	}
	ir, err := nir.Lower(program)
	if err != nil {
		return "", err
	}
	return cbackend.Emit(ir)
}
