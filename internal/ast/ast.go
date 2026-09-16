package ast

type Program struct {
	Name       string
	Statements []Statement
}

type Statement interface {
	statementNode()
}

type DisplayStatement struct {
	Values []Expression
}

func (DisplayStatement) statementNode() {}

type Expression interface {
	expressionNode()
}

type StringLiteral struct{ Value string }

func (StringLiteral) expressionNode() {}

type NumberLiteral struct{ Value string }

func (NumberLiteral) expressionNode() {}

type Identifier struct{ Name string }

func (Identifier) expressionNode() {}
