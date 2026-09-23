package nir

import (
	"fmt"
	"strings"

	"github.com/tomiya7688/neo-Cobol/internal/ast"
	"github.com/tomiya7688/neo-Cobol/internal/sema"
	"github.com/tomiya7688/neo-Cobol/internal/types"
)

type Program struct {
	Name      string
	Variables []Variable
	Ops       []Op
}

type Variable struct {
	Name    string
	Type    types.Type
	Initial *Value
	Mutable bool
}

type ValueKind uint8

const (
	LiteralValue ValueKind = iota
	VariableValue
)

type Value struct {
	Kind ValueKind
	Type types.Kind
	Text string
}

type Op interface{ opNode() }

type Declare struct {
	Name    string
	Type    types.Type
	Initial Value
	Mutable bool
}

func (Declare) opNode() {}

type Display struct{ Values []Value }

func (Display) opNode() {}

type Move struct {
	Source Value
	Target string
}

func (Move) opNode() {}

func Lower(program *ast.Program, info *sema.Info) (*Program, error) {
	out := &Program{Name: program.Name}
	for _, decl := range program.Declarations {
		key := normalize(decl.Name)
		symbol := info.Symbols[key]
		variable := Variable{Name: key, Type: symbol.Type, Mutable: true}
		if decl.Initializer != nil {
			value, err := lowerExpression(decl.Initializer, info)
			if err != nil {
				return nil, err
			}
			variable.Initial = &value
		}
		out.Variables = append(out.Variables, variable)
	}
	for _, statement := range program.Statements {
		switch stmt := statement.(type) {
		case ast.BindingDeclaration:
			key := normalize(stmt.Name)
			symbol, ok := info.Symbols[key]
			if !ok {
				return nil, fmt.Errorf("cannot lower unknown binding %q", stmt.Name)
			}
			value, err := lowerExpression(stmt.Initializer, info)
			if err != nil {
				return nil, err
			}
			out.Ops = append(out.Ops, Declare{Name: key, Type: symbol.Type, Initial: value, Mutable: stmt.Mutable})
		case ast.DisplayStatement:
			op := Display{}
			for _, expression := range stmt.Values {
				value, err := lowerExpression(expression, info)
				if err != nil {
					return nil, err
				}
				op.Values = append(op.Values, value)
			}
			out.Ops = append(out.Ops, op)
		case ast.MoveStatement:
			value, err := lowerExpression(stmt.Source, info)
			if err != nil {
				return nil, err
			}
			out.Ops = append(out.Ops, Move{Source: value, Target: normalize(stmt.Target)})
		default:
			return nil, fmt.Errorf("cannot lower statement node %T", statement)
		}
	}
	return out, nil
}

func lowerExpression(expr ast.Expression, info *sema.Info) (Value, error) {
	switch value := expr.(type) {
	case ast.StringLiteral:
		return Value{Kind: LiteralValue, Type: types.String, Text: value.Value}, nil
	case ast.BooleanLiteral:
		text := "FALSE"
		if value.Value {
			text = "TRUE"
		}
		return Value{Kind: LiteralValue, Type: types.Boolean, Text: text}, nil
	case ast.NumberLiteral:
		kind, err := types.InferNumberLiteral(value.Value)
		if err != nil {
			return Value{}, err
		}
		return Value{Kind: LiteralValue, Type: kind, Text: value.Value}, nil
	case ast.Identifier:
		key := normalize(value.Name)
		symbol, ok := info.Symbols[key]
		if !ok {
			return Value{}, fmt.Errorf("cannot lower unknown identifier %q", value.Name)
		}
		return Value{Kind: VariableValue, Type: symbol.Type.Kind, Text: key}, nil
	default:
		return Value{}, fmt.Errorf("cannot lower expression node %T", expr)
	}
}

func normalize(name string) string { return strings.ToUpper(name) }
