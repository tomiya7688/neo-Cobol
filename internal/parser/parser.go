package parser

import (
	"fmt"
	"github.com/tomiya7688/neo-Cobol/internal/ast"
	"github.com/tomiya7688/neo-Cobol/internal/token"
	"strconv"
	"strings"
)

type parser struct {
	tokens []token.Token
	pos    int
}

func Parse(tokens []token.Token) (*ast.Program, error) {
	p := &parser{tokens: tokens}
	program := &ast.Program{}
	seenExecutable := false
	for !p.atEnd() {
		switch {
		case p.current().IsWord("IDENTIFICATION"):
			if err := p.parseDivisionHeader("IDENTIFICATION"); err != nil {
				return nil, err
			}
		case p.current().IsWord("DATA"):
			if err := p.parseDivisionHeader("DATA"); err != nil {
				return nil, err
			}
		case p.current().IsWord("PROCEDURE"):
			if err := p.parseDivisionHeader("PROCEDURE"); err != nil {
				return nil, err
			}
		case p.current().IsWord("WORKING-STORAGE"):
			if err := p.parseSectionHeader("WORKING-STORAGE"); err != nil {
				return nil, err
			}
		case p.current().IsWord("PROGRAM-ID"):
			if err := p.parseProgramID(program); err != nil {
				return nil, err
			}
		case isLevelToken(p.current()):
			if seenExecutable {
				t := p.current()
				return nil, fmt.Errorf("%d:%d: traditional data declaration cannot appear after executable statements", t.Line, t.Column)
			}
			decl, err := p.parseDataDeclaration()
			if err != nil {
				return nil, err
			}
			program.Declarations = append(program.Declarations, decl)
		default:
			seenExecutable = true
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program, nil
}
func (p *parser) parseStatement() (ast.Statement, error) {
	switch {
	case p.current().IsWord("VAR"), p.current().IsWord("LET"):
		return p.parseBindingDeclaration()
	case p.current().IsWord("DISPLAY"):
		return p.parseDisplay()
	case p.current().IsWord("MOVE"):
		return p.parseMove()
	case p.current().IsWord("IF"):
		return p.parseIf()
	default:
		t := p.current()
		return nil, fmt.Errorf("%d:%d: unsupported statement starting at %q", t.Line, t.Column, t.Lexeme)
	}
}
func (p *parser) parseDivisionHeader(name string) error {
	p.advance()
	if !p.current().IsWord("DIVISION") {
		return p.expected("DIVISION after " + name)
	}
	p.advance()
	return p.requirePeriod(name + " DIVISION")
}
func (p *parser) parseSectionHeader(name string) error {
	p.advance()
	if !p.current().IsWord("SECTION") {
		return p.expected("SECTION after " + name)
	}
	p.advance()
	return p.requirePeriod(name + " SECTION")
}
func (p *parser) parseProgramID(program *ast.Program) error {
	p.advance()
	if p.current().Kind != token.Period {
		return p.expected("'.' after PROGRAM-ID")
	}
	p.advance()
	if p.current().Kind != token.Identifier {
		return p.expected("program name after PROGRAM-ID.")
	}
	program.Name = p.advance().Lexeme
	return p.requirePeriod("program name")
}
func (p *parser) parseDataDeclaration() (ast.DataDeclaration, error) {
	levelToken := p.advance()
	level, _ := strconv.Atoi(levelToken.Lexeme)
	decl := ast.DataDeclaration{Level: level}
	if p.current().Kind != token.Identifier {
		return decl, p.expected("data name after level number")
	}
	decl.Name = p.advance().Lexeme
	if p.current().IsWord("TYPE") {
		p.advance()
		if p.current().Kind != token.Identifier {
			return decl, p.expected("type name after TYPE")
		}
		decl.TypeName = p.advance().Lexeme
	}
	if p.current().IsWord("PIC") {
		p.advance()
		var picture strings.Builder
		for p.current().Kind != token.Period && !p.current().IsWord("VALUE") && p.current().Kind != token.EOF {
			t := p.current()
			switch t.Kind {
			case token.Identifier, token.Number, token.LParen, token.RParen:
				picture.WriteString(t.Lexeme)
				p.advance()
			default:
				return decl, fmt.Errorf("%d:%d: invalid PIC token %q", t.Line, t.Column, t.Lexeme)
			}
		}
		decl.Picture = picture.String()
		if decl.Picture == "" {
			return decl, p.expected("picture after PIC")
		}
	}
	if decl.TypeName == "" && decl.Picture == "" {
		return decl, fmt.Errorf("%d:%d: data declaration %q requires TYPE or PIC", levelToken.Line, levelToken.Column, decl.Name)
	}
	if p.current().IsWord("VALUE") {
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return decl, err
		}
		switch expr.(type) {
		case ast.StringLiteral, ast.NumberLiteral, ast.BooleanLiteral:
			decl.Initializer = expr
		default:
			return decl, fmt.Errorf("%d:%d: VALUE currently requires a literal", p.current().Line, p.current().Column)
		}
	}
	if err := p.requirePeriod("data declaration"); err != nil {
		return decl, err
	}
	return decl, nil
}
func (p *parser) parseBindingDeclaration() (ast.Statement, error) {
	keyword := p.advance()
	mutable := keyword.IsWord("VAR")
	if p.current().Kind != token.Identifier {
		return nil, p.expected("binding name after " + strings.ToUpper(keyword.Lexeme))
	}
	name := p.advance().Lexeme
	if !p.current().IsWord("VALUE") {
		return nil, p.expected("VALUE in " + strings.ToUpper(keyword.Lexeme) + " declaration")
	}
	p.advance()
	initializer, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	p.consumeOptionalPeriod()
	return ast.BindingDeclaration{Name: name, Mutable: mutable, Initializer: initializer}, nil
}
func (p *parser) parseDisplay() (ast.Statement, error) {
	p.advance()
	stmt := ast.DisplayStatement{}
	for p.current().Kind != token.Period && p.current().Kind != token.EOF && !p.atBlockBoundary() {
		if len(stmt.Values) > 0 && p.isStatementStart() {
			break
		}
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Values = append(stmt.Values, expr)
	}
	if len(stmt.Values) == 0 {
		return nil, p.expected("at least one DISPLAY operand")
	}
	p.consumeOptionalPeriod()
	return stmt, nil
}
func (p *parser) parseMove() (ast.Statement, error) {
	p.advance()
	source, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if !p.current().IsWord("TO") {
		return nil, p.expected("TO in MOVE statement")
	}
	p.advance()
	if p.current().Kind != token.Identifier {
		return nil, p.expected("target identifier after TO")
	}
	target := p.advance().Lexeme
	p.consumeOptionalPeriod()
	return ast.MoveStatement{Source: source, Target: target}, nil
}
func (p *parser) parseIf() (ast.Statement, error) {
	p.advance()
	condition, err := p.parseCondition()
	if err != nil {
		return nil, err
	}
	thenBlock, err := p.parseBlock(true)
	if err != nil {
		return nil, err
	}
	var elseBlock []ast.Statement
	if p.current().IsWord("ELSE") {
		p.advance()
		elseBlock, err = p.parseBlock(false)
		if err != nil {
			return nil, err
		}
	}
	if p.current().IsWord("END-IF") {
		p.advance()
	} else if p.current().IsWord("END") && p.peekWord(1, "IF") {
		p.advance()
		p.advance()
	} else {
		return nil, p.expected("END-IF")
	}
	p.consumeOptionalPeriod()
	return ast.IfStatement{Condition: condition, Then: thenBlock, Else: elseBlock}, nil
}
func (p *parser) parseBlock(stopAtElse bool) ([]ast.Statement, error) {
	var out []ast.Statement
	for !p.atEnd() {
		if p.current().IsWord("END-IF") || (p.current().IsWord("END") && p.peekWord(1, "IF")) {
			return out, nil
		}
		if stopAtElse && p.current().IsWord("ELSE") {
			return out, nil
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		out = append(out, stmt)
	}
	return nil, p.expected("END-IF")
}
func (p *parser) parseCondition() (ast.Condition, error) { return p.parseOrCondition() }
func (p *parser) parseOrCondition() (ast.Condition, error) {
	left, err := p.parseAndCondition()
	if err != nil {
		return nil, err
	}
	for p.current().IsWord("OR") {
		p.advance()
		right, err := p.parseAndCondition()
		if err != nil {
			return nil, err
		}
		left = ast.LogicalCondition{Left: left, Op: ast.LogicalOr, Right: right}
	}
	return left, nil
}
func (p *parser) parseAndCondition() (ast.Condition, error) {
	left, err := p.parseNotCondition()
	if err != nil {
		return nil, err
	}
	for p.current().IsWord("AND") {
		p.advance()
		right, err := p.parseNotCondition()
		if err != nil {
			return nil, err
		}
		left = ast.LogicalCondition{Left: left, Op: ast.LogicalAnd, Right: right}
	}
	return left, nil
}
func (p *parser) parseNotCondition() (ast.Condition, error) {
	if p.current().IsWord("NOT") {
		p.advance()
		inner, err := p.parseNotCondition()
		if err != nil {
			return nil, err
		}
		return ast.NotCondition{Inner: inner}, nil
	}
	return p.parsePredicateCondition()
}
func (p *parser) parsePredicateCondition() (ast.Condition, error) {
	if p.current().Kind == token.LParen {
		p.advance()
		inner, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		if p.current().Kind != token.RParen {
			return nil, p.expected("')' after condition")
		}
		p.advance()
		return inner, nil
	}
	left, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if !p.current().IsWord("IS") {
		return ast.ValueCondition{Value: left}, nil
	}
	op, err := p.parseComparisonOperator()
	if err != nil {
		return nil, err
	}
	right, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return ast.ComparisonCondition{Left: left, Op: op, Right: right}, nil
}
func (p *parser) parseComparisonOperator() (ast.ComparisonOperator, error) {
	p.advance()
	if p.current().IsWord("NOT") {
		p.advance()
		if !p.current().IsWord("EQUAL") {
			return 0, p.expected("EQUAL after IS NOT")
		}
		p.advance()
		if !p.current().IsWord("TO") {
			return 0, p.expected("TO after IS NOT EQUAL")
		}
		p.advance()
		return ast.CompareNotEqual, nil
	}
	if p.current().IsWord("EQUAL") {
		p.advance()
		if !p.current().IsWord("TO") {
			return 0, p.expected("TO after IS EQUAL")
		}
		p.advance()
		return ast.CompareEqual, nil
	}
	if p.current().IsWord("GREATER") || p.current().IsWord("LESS") {
		greater := p.current().IsWord("GREATER")
		p.advance()
		if !p.current().IsWord("THAN") {
			return 0, p.expected("THAN in comparison")
		}
		p.advance()
		if p.current().IsWord("OR") {
			p.advance()
			if !p.current().IsWord("EQUAL") {
				return 0, p.expected("EQUAL after OR")
			}
			p.advance()
			if !p.current().IsWord("TO") {
				return 0, p.expected("TO after OR EQUAL")
			}
			p.advance()
			if greater {
				return ast.CompareGreaterEqual, nil
			}
			return ast.CompareLessEqual, nil
		}
		if greater {
			return ast.CompareGreater, nil
		}
		return ast.CompareLess, nil
	}
	return 0, p.expected("comparison phrase after IS")
}
func (p *parser) parseExpression() (ast.Expression, error) {
	t := p.current()
	switch t.Kind {
	case token.String:
		p.advance()
		return ast.StringLiteral{Value: t.Lexeme}, nil
	case token.Number:
		p.advance()
		return ast.NumberLiteral{Value: t.Lexeme}, nil
	case token.Identifier:
		p.advance()
		if t.IsWord("TRUE") {
			return ast.BooleanLiteral{Value: true}, nil
		}
		if t.IsWord("FALSE") {
			return ast.BooleanLiteral{Value: false}, nil
		}
		return ast.Identifier{Name: t.Lexeme}, nil
	default:
		return nil, fmt.Errorf("%d:%d: expected expression, got %q", t.Line, t.Column, t.Lexeme)
	}
}
func isLevelToken(t token.Token) bool {
	if t.Kind != token.Number || strings.ContainsAny(t.Lexeme, "+-.eExXoObB") {
		return false
	}
	n, err := strconv.Atoi(t.Lexeme)
	if err != nil {
		return false
	}
	return (n >= 1 && n <= 49) || n == 77
}
func (p *parser) isStatementStart() bool {
	return p.current().IsWord("VAR") || p.current().IsWord("LET") || p.current().IsWord("DISPLAY") || p.current().IsWord("MOVE") || p.current().IsWord("IF")
}
func (p *parser) atBlockBoundary() bool {
	return p.current().IsWord("ELSE") || p.current().IsWord("END-IF") || (p.current().IsWord("END") && p.peekWord(1, "IF"))
}
func (p *parser) peekWord(offset int, word string) bool {
	idx := p.pos + offset
	return idx < len(p.tokens) && p.tokens[idx].IsWord(word)
}
func (p *parser) consumeOptionalPeriod() {
	if p.current().Kind == token.Period {
		p.advance()
	}
}
func (p *parser) requirePeriod(context string) error {
	if p.current().Kind != token.Period {
		return p.expected("'.' after " + context)
	}
	p.advance()
	return nil
}
func (p *parser) expected(what string) error {
	t := p.current()
	return fmt.Errorf("%d:%d: expected %s, got %q", t.Line, t.Column, what, t.Lexeme)
}
func (p *parser) current() token.Token { return p.tokens[p.pos] }
func (p *parser) atEnd() bool          { return p.current().Kind == token.EOF }
func (p *parser) advance() token.Token {
	t := p.current()
	if !p.atEnd() {
		p.pos++
	}
	return t
}
