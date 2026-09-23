package sema

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/tomiya7688/neo-Cobol/internal/ast"
	"github.com/tomiya7688/neo-Cobol/internal/types"
)

type Symbol struct {
	Name    string
	Type    types.Type
	Mutable bool
}

type Info struct {
	Symbols map[string]Symbol
}

func Check(program *ast.Program) (*Info, error) {
	info := &Info{Symbols: make(map[string]Symbol)}
	initialized := make(map[string]bool)

	for _, decl := range program.Declarations {
		if decl.Level != 1 && decl.Level != 77 {
			return nil, fmt.Errorf("level %02d data item %q requires record hierarchy support, which is not implemented yet", decl.Level, decl.Name)
		}
		key := normalize(decl.Name)
		if _, exists := info.Symbols[key]; exists {
			return nil, fmt.Errorf("duplicate declaration of %q", decl.Name)
		}

		typeInfo, err := resolveDeclarationType(decl)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", decl.Name, err)
		}
		if decl.Initializer != nil {
			if err := checkAssignmentValue(decl.Initializer, typeInfo, info, initialized); err != nil {
				return nil, fmt.Errorf("%s VALUE: %w", decl.Name, err)
			}
			initialized[key] = true
		}
		info.Symbols[key] = Symbol{Name: decl.Name, Type: typeInfo, Mutable: true}
	}

	for _, statement := range program.Statements {
		switch stmt := statement.(type) {
		case ast.BindingDeclaration:
			key := normalize(stmt.Name)
			if _, exists := info.Symbols[key]; exists {
				return nil, fmt.Errorf("duplicate declaration of %q", stmt.Name)
			}
			kind, err := expressionType(stmt.Initializer, info, initialized, true)
			if err != nil {
				return nil, fmt.Errorf("%s VALUE: %w", stmt.Name, err)
			}
			info.Symbols[key] = Symbol{
				Name:    stmt.Name,
				Type:    types.Type{Kind: kind},
				Mutable: stmt.Mutable,
			}
			initialized[key] = true
		case ast.DisplayStatement:
			for _, value := range stmt.Values {
				if _, err := expressionType(value, info, initialized, true); err != nil {
					return nil, err
				}
			}
		case ast.MoveStatement:
			targetKey := normalize(stmt.Target)
			target, ok := info.Symbols[targetKey]
			if !ok {
				return nil, fmt.Errorf("MOVE target %q is not declared", stmt.Target)
			}
			if !target.Mutable {
				return nil, fmt.Errorf("MOVE target %q is a LET binding and cannot be reassigned", stmt.Target)
			}
			if err := checkAssignmentValue(stmt.Source, target.Type, info, initialized); err != nil {
				return nil, fmt.Errorf("MOVE to %s: %w", stmt.Target, err)
			}
			initialized[targetKey] = true
		default:
			return nil, fmt.Errorf("unsupported statement node %T", statement)
		}
	}
	return info, nil
}

func resolveDeclarationType(decl ast.DataDeclaration) (types.Type, error) {
	var result types.Type
	if decl.TypeName != "" {
		kind, ok := types.ParseBuiltin(decl.TypeName)
		if !ok {
			return result, fmt.Errorf("unknown type %q", decl.TypeName)
		}
		result.Kind = kind
	}
	if decl.Picture != "" {
		pic, err := types.ParsePicture(decl.Picture)
		if err != nil {
			return result, err
		}
		result.Picture = pic
		if result.Kind == types.Invalid {
			kind, err := types.InferFromPicture(pic)
			if err != nil {
				return result, err
			}
			result.Kind = kind
		}
		if !types.CompatiblePicture(result.Kind, pic) {
			return result, fmt.Errorf("PIC %s is incompatible with TYPE %s", pic.Raw, result.Kind)
		}
	}
	if result.Kind == types.Invalid {
		return result, fmt.Errorf("declaration has no logical type")
	}
	return result, nil
}

func checkAssignmentValue(expr ast.Expression, target types.Type, info *Info, initialized map[string]bool) error {
	sourceKind, err := expressionType(expr, info, initialized, true)
	if err != nil {
		return err
	}
	if !types.Assignable(sourceKind, target.Kind) && !literalFitsTarget(expr, target.Kind) {
		return fmt.Errorf("cannot assign %s to %s", sourceKind, target.Kind)
	}
	return validatePictureLiteral(expr, target.Picture)
}

func expressionType(expr ast.Expression, info *Info, initialized map[string]bool, requireInitialized bool) (types.Kind, error) {
	switch value := expr.(type) {
	case ast.StringLiteral:
		return types.String, nil
	case ast.NumberLiteral:
		return types.InferNumberLiteral(value.Value)
	case ast.BooleanLiteral:
		return types.Boolean, nil
	case ast.Identifier:
		key := normalize(value.Name)
		symbol, ok := info.Symbols[key]
		if !ok {
			return types.Invalid, fmt.Errorf("identifier %q is not declared", value.Name)
		}
		if requireInitialized && !initialized[key] {
			return types.Invalid, fmt.Errorf("identifier %q is used before initialization", value.Name)
		}
		return symbol.Type.Kind, nil
	default:
		return types.Invalid, fmt.Errorf("unsupported expression node %T", expr)
	}
}

func literalFitsTarget(expr ast.Expression, target types.Kind) bool {
	number, ok := expr.(ast.NumberLiteral)
	if !ok {
		return false
	}
	value, err := strconv.ParseInt(number.Value, 0, 64)
	if err != nil {
		return target == types.Decimal && strings.ContainsAny(number.Value, ".eE")
	}
	switch target {
	case types.Byte:
		return value >= 0 && value <= 255
	case types.Integer:
		return value >= -2147483648 && value <= 2147483647
	case types.Long:
		return true
	}
	return false
}

func validatePictureLiteral(expr ast.Expression, pic *types.Picture) error {
	if pic == nil {
		return nil
	}
	switch value := expr.(type) {
	case ast.StringLiteral:
		if pic.Kind == types.PictureText && utf8.RuneCountInString(value.Value) > pic.Width {
			return fmt.Errorf("string literal length exceeds PIC %s", pic.Raw)
		}
	case ast.NumberLiteral:
		if pic.Kind != types.PictureNumeric {
			return nil
		}
		text := value.Value
		if strings.HasPrefix(text, "-") {
			if !pic.Signed {
				return fmt.Errorf("negative literal is incompatible with unsigned PIC %s", pic.Raw)
			}
			text = text[1:]
		} else if strings.HasPrefix(text, "+") {
			text = text[1:]
		}
		if strings.ContainsAny(text, "xXoObBeE") {
			return nil
		}
		parts := strings.SplitN(text, ".", 2)
		integerDigits := len(strings.TrimLeft(parts[0], "0"))
		if integerDigits == 0 {
			integerDigits = 1
		}
		fractionDigits := 0
		if len(parts) == 2 {
			fractionDigits = len(parts[1])
		}
		if integerDigits > pic.Width-pic.Scale || fractionDigits > pic.Scale {
			return fmt.Errorf("numeric literal %s exceeds PIC %s", value.Value, pic.Raw)
		}
	}
	return nil
}

func normalize(name string) string { return strings.ToUpper(name) }
