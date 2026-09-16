package parser

import (
	"fmt"

	"github.com/tomiya7688/neo-Cobol/internal/ast"
	"github.com/tomiya7688/neo-Cobol/internal/token"
)

type parser struct {
	tokens []token.Token
	pos    int
}

// Parse implements the first executable language slice. It accepts optional
// IDENTIFICATION/PROCEDURE division headers, PROGRAM-ID, and DISPLAY.
func Parse(tokens []token.Token) (*ast.Program, error) {
	p := &parser{tokens: tokens}
	program := &ast.Program{}

	for !p.atEnd() {
		switch {
		case p.current().IsWord("IDENTIFICATION"):
			if err := p.parseDivisionHeader("IDENTIFICATION"); err != nil {
				return nil, err
			}
		case p.current().IsWord("PROCEDURE"):
			if err := p.parseDivisionHeader("PROCEDURE"); err != nil {
				return nil, err
			}
		case p.current().IsWord("PROGRAM-ID"):
			if err := p.parseProgramID(program); err != nil {
				return nil, err
			}
		case p.current().IsWord("DISPLAY"):
			stmt, err := p.parseDisplay()
			if err != nil {
				return nil, err
			}
			program.Statements = append(program.Statements, stmt)
		default:
			t := p.current()
			return nil, fmt.Errorf("%d:%d: unsupported syntax starting at %q", t.Line, t.Column, t.Lexeme)
		}
	}
	return program, nil
}

func (p *parser) parseDivisionHeader(name string) error {
	start := p.advance()
	if !p.current().IsWord("DIVISION") {
		return fmt.Errorf("%d:%d: expected DIVISION after %s", p.current().Line, p.current().Column, name)
	}
	p.advance()
	if p.current().Kind != token.Period {
		return fmt.Errorf("%d:%d: expected '.' after %s DIVISION", p.current().Line, p.current().Column, name)
	}
	p.advance()
	_ = start
	return nil
}

func (p *parser) parseProgramID(program *ast.Program) error {
	p.advance()
	if p.current().Kind != token.Period {
		return fmt.Errorf("%d:%d: expected '.' after PROGRAM-ID", p.current().Line, p.current().Column)
	}
	p.advance()
	if p.current().Kind != token.Identifier {
		return fmt.Errorf("%d:%d: expected program name after PROGRAM-ID.", p.current().Line, p.current().Column)
	}
	program.Name = p.advance().Lexeme
	if p.current().Kind != token.Period {
		return fmt.Errorf("%d:%d: expected '.' after program name", p.current().Line, p.current().Column)
	}
	p.advance()
	return nil
}

func (p *parser) parseDisplay() (ast.Statement, error) {
	p.advance()
	stmt := ast.DisplayStatement{}
	for p.current().Kind != token.Period && p.current().Kind != token.EOF {
		t := p.advance()
		switch t.Kind {
		case token.String:
			stmt.Values = append(stmt.Values, ast.StringLiteral{Value: t.Lexeme})
		case token.Number:
			stmt.Values = append(stmt.Values, ast.NumberLiteral{Value: t.Lexeme})
		case token.Identifier:
			stmt.Values = append(stmt.Values, ast.Identifier{Name: t.Lexeme})
		default:
			return nil, fmt.Errorf("%d:%d: invalid DISPLAY operand %q", t.Line, t.Column, t.Lexeme)
		}
	}
	if len(stmt.Values) == 0 {
		t := p.current()
		return nil, fmt.Errorf("%d:%d: DISPLAY requires at least one operand", t.Line, t.Column)
	}
	if p.current().Kind == token.Period {
		p.advance()
	}
	return stmt, nil
}

func (p *parser) current() token.Token { return p.tokens[p.pos] }
func (p *parser) atEnd() bool           { return p.current().Kind == token.EOF }
func (p *parser) advance() token.Token {
	t := p.current()
	if !p.atEnd() {
		p.pos++
	}
	return t
}
