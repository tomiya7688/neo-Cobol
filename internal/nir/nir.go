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
	Records   []Record
	Ops       []Op
}

type Variable struct {
	Name    string
	Type    types.Type
	Initial *Value
	Mutable bool
}

type Record struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name     string
	Type     *types.Type
	Initial  *Value
	Children []Field
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

type If struct {
	Condition Condition
	Then      []Op
	Else      []Op
}

func (If) opNode() {}

type Condition interface{ conditionNode() }

type Truth struct{ Value Value }

func (Truth) conditionNode() {}

type ComparisonOperator uint8

const (
	CompareEqual ComparisonOperator = iota
	CompareNotEqual
	CompareGreater
	CompareLess
	CompareGreaterEqual
	CompareLessEqual
)

type Compare struct {
	Left  Value
	Op    ComparisonOperator
	Right Value
}

func (Compare) conditionNode() {}

type LogicalOperator uint8

const (
	LogicalAnd LogicalOperator = iota
	LogicalOr
)

type Logical struct {
	Left  Condition
	Op    LogicalOperator
	Right Condition
}

func (Logical) conditionNode() {}

type Not struct{ Inner Condition }

func (Not) conditionNode() {}

type binding struct {
	id  string
	typ types.Type
}

type lowerScope struct {
	parent   *lowerScope
	bindings map[string]binding
	info     *sema.Info
}

type lowerer struct{ nextID int }

func Lower(program *ast.Program, info *sema.Info) (*Program, error) {
	out := &Program{Name: program.Name}
	root := &lowerScope{bindings: make(map[string]binding), info: info}

	for _, decl := range program.Declarations {
		if decl.IsGroup() {
			record, err := lowerRecord(decl, nil, root)
			if err != nil {
				return nil, err
			}
			out.Records = append(out.Records, record)
			continue
		}
		id := "d:" + normalize(decl.Name)
		symbol, ok := info.Symbols[id]
		if !ok {
			return nil, fmt.Errorf("cannot lower unknown global %q", decl.Name)
		}
		variable := Variable{Name: id, Type: symbol.Type, Mutable: true}
		if decl.Initializer != nil {
			value, err := lowerExpression(decl.Initializer, root)
			if err != nil {
				return nil, err
			}
			variable.Initial = &value
		}
		out.Variables = append(out.Variables, variable)
	}

	l := &lowerer{}
	ops, err := l.lowerBlock(program.Statements, root)
	if err != nil {
		return nil, err
	}
	out.Ops = ops
	return out, nil
}

func lowerRecord(decl ast.DataDeclaration, ancestors []string, s *lowerScope) (Record, error) {
	path := append(append([]string{}, ancestors...), normalize(decl.Name))
	record := Record{Name: "r:" + strings.Join(path, "/")}
	for _, child := range decl.Children {
		field, err := lowerField(child, path, s)
		if err != nil {
			return Record{}, err
		}
		record.Fields = append(record.Fields, field)
	}
	return record, nil
}

func lowerField(decl ast.DataDeclaration, ancestors []string, s *lowerScope) (Field, error) {
	path := append(append([]string{}, ancestors...), normalize(decl.Name))
	if decl.IsGroup() {
		field := Field{Name: "g:" + strings.Join(path, "/")}
		for _, child := range decl.Children {
			nested, err := lowerField(child, path, s)
			if err != nil {
				return Field{}, err
			}
			field.Children = append(field.Children, nested)
		}
		return field, nil
	}

	id := "d:" + strings.Join(path, "/")
	symbol, ok := s.info.Symbols[id]
	if !ok {
		return Field{}, fmt.Errorf("cannot lower unknown data field %q", decl.Name)
	}
	typ := symbol.Type
	field := Field{Name: id, Type: &typ}
	if decl.Initializer != nil {
		value, err := lowerExpression(decl.Initializer, s)
		if err != nil {
			return Field{}, err
		}
		field.Initial = &value
	}
	return field, nil
}

func (l *lowerer) lowerBlock(statements []ast.Statement, s *lowerScope) ([]Op, error) {
	var out []Op
	for _, statement := range statements {
		switch stmt := statement.(type) {
		case ast.BindingDeclaration:
			value, err := lowerExpression(stmt.Initializer, s)
			if err != nil {
				return nil, err
			}
			id := fmt.Sprintf("l:%d", l.nextID)
			l.nextID++
			typ := types.Type{Kind: value.Type}
			out = append(out, Declare{Name: id, Type: typ, Initial: value, Mutable: stmt.Mutable})
			s.bindings[normalize(stmt.Name)] = binding{id: id, typ: typ}

		case ast.DisplayStatement:
			op := Display{}
			for _, expr := range stmt.Values {
				value, err := lowerExpression(expr, s)
				if err != nil {
					return nil, err
				}
				op.Values = append(op.Values, value)
			}
			out = append(out, op)

		case ast.MoveStatement:
			target, err := resolveReference(stmt.Target, s)
			if err != nil {
				return nil, err
			}
			value, err := lowerExpression(stmt.Source, s)
			if err != nil {
				return nil, err
			}
			out = append(out, Move{Source: value, Target: target.id})

		case ast.IfStatement:
			condition, err := lowerCondition(stmt.Condition, s)
			if err != nil {
				return nil, err
			}
			thenScope := &lowerScope{parent: s, bindings: make(map[string]binding), info: s.info}
			thenOps, err := l.lowerBlock(stmt.Then, thenScope)
			if err != nil {
				return nil, err
			}
			elseScope := &lowerScope{parent: s, bindings: make(map[string]binding), info: s.info}
			elseOps, err := l.lowerBlock(stmt.Else, elseScope)
			if err != nil {
				return nil, err
			}
			out = append(out, If{Condition: condition, Then: thenOps, Else: elseOps})

		default:
			return nil, fmt.Errorf("cannot lower statement node %T", statement)
		}
	}
	return out, nil
}

func lowerCondition(condition ast.Condition, s *lowerScope) (Condition, error) {
	switch cond := condition.(type) {
	case ast.ValueCondition:
		value, err := lowerExpression(cond.Value, s)
		if err != nil {
			return nil, err
		}
		return Truth{Value: value}, nil
	case ast.ComparisonCondition:
		left, err := lowerExpression(cond.Left, s)
		if err != nil {
			return nil, err
		}
		right, err := lowerExpression(cond.Right, s)
		if err != nil {
			return nil, err
		}
		return Compare{Left: left, Op: mapCompare(cond.Op), Right: right}, nil
	case ast.LogicalCondition:
		left, err := lowerCondition(cond.Left, s)
		if err != nil {
			return nil, err
		}
		right, err := lowerCondition(cond.Right, s)
		if err != nil {
			return nil, err
		}
		op := LogicalAnd
		if cond.Op == ast.LogicalOr {
			op = LogicalOr
		}
		return Logical{Left: left, Op: op, Right: right}, nil
	case ast.NotCondition:
		inner, err := lowerCondition(cond.Inner, s)
		if err != nil {
			return nil, err
		}
		return Not{Inner: inner}, nil
	default:
		return nil, fmt.Errorf("cannot lower condition node %T", condition)
	}
}

func mapCompare(op ast.ComparisonOperator) ComparisonOperator {
	switch op {
	case ast.CompareEqual:
		return CompareEqual
	case ast.CompareNotEqual:
		return CompareNotEqual
	case ast.CompareGreater:
		return CompareGreater
	case ast.CompareLess:
		return CompareLess
	case ast.CompareGreaterEqual:
		return CompareGreaterEqual
	case ast.CompareLessEqual:
		return CompareLessEqual
	default:
		return CompareEqual
	}
}

func lowerExpression(expr ast.Expression, s *lowerScope) (Value, error) {
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
	case ast.DataReference:
		resolved, err := resolveReference(value, s)
		if err != nil {
			return Value{}, err
		}
		return Value{Kind: VariableValue, Type: resolved.typ.Kind, Text: resolved.id}, nil
	default:
		return Value{}, fmt.Errorf("cannot lower expression node %T", expr)
	}
}

func resolveReference(ref ast.DataReference, s *lowerScope) (binding, error) {
	if len(ref.Qualifiers) == 0 {
		if local, ok := s.lookupLocal(ref.Name); ok {
			return local, nil
		}
	}
	id, ok := s.info.DataID(ref)
	if !ok {
		return binding{}, fmt.Errorf("cannot lower unresolved reference %s", sema.FormatReference(ref))
	}
	symbol := s.info.Symbols[id]
	return binding{id: id, typ: symbol.Type}, nil
}

func (s *lowerScope) lookupLocal(name string) (binding, bool) {
	key := normalize(name)
	for current := s; current != nil; current = current.parent {
		if b, ok := current.bindings[key]; ok {
			return b, true
		}
	}
	return binding{}, false
}

func normalize(name string) string { return strings.ToUpper(name) }
