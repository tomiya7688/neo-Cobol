package nir

import (
	"fmt"

	"github.com/tomiya7688/neo-Cobol/internal/ast"
)

type Program struct {
	Name string
	Ops  []Op
}

type Op interface{ opNode() }

type Display struct{ Values []string }

func (Display) opNode() {}

func Lower(program *ast.Program) (*Program, error) {
	out := &Program{Name: program.Name}
	for _, statement := range program.Statements {
		switch stmt := statement.(type) {
		case ast.DisplayStatement:
			op := Display{}
			for _, value := range stmt.Values {
				switch v := value.(type) {
				case ast.StringLiteral:
					op.Values = append(op.Values, v.Value)
				case ast.NumberLiteral:
					op.Values = append(op.Values, v.Value)
				default:
					return nil, fmt.Errorf("cannot lower expression node %T", value)
				}
			}
			out.Ops = append(out.Ops, op)
		default:
			return nil, fmt.Errorf("cannot lower statement node %T", statement)
		}
	}
	return out, nil
}
