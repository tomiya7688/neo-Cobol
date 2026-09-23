package ast

type Program struct {
	Name         string
	Declarations []DataDeclaration
	Statements   []Statement
}

type DataDeclaration struct {
	Level       int
	Name        string
	TypeName    string
	Picture     string
	Initializer Expression
}

type Statement interface{ statementNode() }

type BindingDeclaration struct {
	Name        string
	Mutable     bool
	Initializer Expression
}

func (BindingDeclaration) statementNode() {}

type DisplayStatement struct{ Values []Expression }

func (DisplayStatement) statementNode() {}

type MoveStatement struct {
	Source Expression
	Target string
}

func (MoveStatement) statementNode() {}

type Expression interface{ expressionNode() }

type StringLiteral struct{ Value string }

func (StringLiteral) expressionNode() {}

type NumberLiteral struct{ Value string }

func (NumberLiteral) expressionNode() {}

type BooleanLiteral struct{ Value bool }

func (BooleanLiteral) expressionNode() {}

type Identifier struct{ Name string }

func (Identifier) expressionNode() {}
