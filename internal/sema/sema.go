package sema

import (
	"fmt"

	"github.com/tomiya7688/neo-Cobol/internal/ast"
)

// Check validates the semantic subset currently implemented by the bootstrap
// compiler. Identifier operands are deliberately rejected until data
// declarations and name resolution are implemented.
func Check(program *ast.Program) error {
	for _, statement := range program.Statements {
		switch stmt := statement.(type) {
		case ast.DisplayStatement:
			for _, value := range stmt.Values {
				if ident, ok := value.(ast.Identifier); ok {
					return fmt.Errorf("identifier %q requires data declarations/name resolution, which are not implemented yet", ident.Name)
				}
			}
		default:
			return fmt.Errorf("unsupported statement node %T", statement)
		}
	}
	return nil
}
