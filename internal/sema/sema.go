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
	ID      string
	Name    string
	Type    types.Type
	Mutable bool
}

type Info struct {
	Symbols   map[string]Symbol
	Globals   map[string]string
	Qualified map[string]string
}

func (i *Info) DataID(ref ast.DataReference) (string, bool) {
	if len(ref.Qualifiers) == 0 {
		id, ok := i.Globals[normalize(ref.Name)]
		return id, ok
	}
	id, ok := i.Qualified[ReferenceKey(ref)]
	return id, ok
}

type scope struct {
	parent  *scope
	symbols map[string]Symbol
}

type flow map[string]bool

type checker struct {
	info         *Info
	nextID       int
	dataNames    map[string][]string
	allDataNames map[string]bool
}

func Check(program *ast.Program) (*Info, error) {
	c := &checker{
		info: &Info{
			Symbols:   make(map[string]Symbol),
			Globals:   make(map[string]string),
			Qualified: make(map[string]string),
		},
		dataNames:    make(map[string][]string),
		allDataNames: make(map[string]bool),
	}
	state := make(flow)
	rootNames := make(map[string]bool)
	for _, decl := range program.Declarations {
		key := normalize(decl.Name)
		if rootNames[key] {
			return nil, fmt.Errorf("duplicate top-level data declaration %q", decl.Name)
		}
		rootNames[key] = true
		if err := c.registerData(decl, nil, state); err != nil {
			return nil, err
		}
	}
	for name, ids := range c.dataNames {
		if len(ids) == 1 {
			c.info.Globals[name] = ids[0]
		}
	}

	procedureScope := &scope{symbols: make(map[string]Symbol)}
	if err := c.validateInitializersByPath(program.Declarations, nil, procedureScope, state); err != nil {
		return nil, err
	}
	if _, err := c.checkStatements(procedureScope, state, program.Statements); err != nil {
		return nil, err
	}
	return c.info, nil
}

func (c *checker) registerData(decl ast.DataDeclaration, ancestors []string, state flow) error {
	key := normalize(decl.Name)
	c.allDataNames[key] = true
	path := append(append([]string{}, ancestors...), key)
	if decl.IsGroup() {
		for _, child := range decl.Children {
			if err := c.registerData(child, path, state); err != nil {
				return err
			}
		}
		return nil
	}

	typeInfo, err := resolveDeclarationType(decl)
	if err != nil {
		return fmt.Errorf("%s: %w", strings.Join(path, " OF "), err)
	}
	id := "d:" + strings.Join(path, "/")
	if _, exists := c.info.Symbols[id]; exists {
		return fmt.Errorf("duplicate data path for %q", decl.Name)
	}
	symbol := Symbol{ID: id, Name: decl.Name, Type: typeInfo, Mutable: true}
	c.info.Symbols[id] = symbol
	c.dataNames[key] = append(c.dataNames[key], id)
	state[id] = false

	if len(ancestors) > 0 {
		qualifiers := make([]string, 0, len(ancestors))
		for index := len(ancestors) - 1; index >= 0; index-- {
			qualifiers = append(qualifiers, ancestors[index])
		}
		ref := ast.DataReference{Name: decl.Name, Qualifiers: qualifiers}
		qkey := ReferenceKey(ref)
		if _, exists := c.info.Qualified[qkey]; exists {
			return fmt.Errorf("duplicate qualified data name %s", FormatReference(ref))
		}
		c.info.Qualified[qkey] = id
	}
	return nil
}

func (c *checker) validateInitializersByPath(decls []ast.DataDeclaration, ancestors []string, s *scope, state flow) error {
	for _, decl := range decls {
		key := normalize(decl.Name)
		path := append(append([]string{}, ancestors...), key)
		if decl.IsGroup() {
			if err := c.validateInitializersByPath(decl.Children, path, s, state); err != nil {
				return err
			}
			continue
		}
		id := "d:" + strings.Join(path, "/")
		symbol := c.info.Symbols[id]
		if decl.Initializer != nil {
			if err := c.checkAssignmentValue(decl.Initializer, symbol.Type, s, state); err != nil {
				return fmt.Errorf("%s VALUE: %w", strings.Join(path, " OF "), err)
			}
			state[id] = true
		}
	}
	return nil
}

func (c *checker) checkStatements(s *scope, state flow, statements []ast.Statement) (flow, error) {
	current := cloneFlow(state)
	for _, statement := range statements {
		switch stmt := statement.(type) {
		case ast.BindingDeclaration:
			key := normalize(stmt.Name)
			if _, exists := s.symbols[key]; exists {
				return nil, fmt.Errorf("duplicate declaration of %q in the same scope", stmt.Name)
			}
			if s.parent == nil && c.allDataNames[key] {
				return nil, fmt.Errorf("duplicate declaration of %q in the same scope", stmt.Name)
			}
			kind, err := c.expressionType(stmt.Initializer, s, current, true)
			if err != nil {
				return nil, fmt.Errorf("%s VALUE: %w", stmt.Name, err)
			}
			id := fmt.Sprintf("l:%d", c.nextID)
			c.nextID++
			symbol := Symbol{ID: id, Name: stmt.Name, Type: types.Type{Kind: kind}, Mutable: stmt.Mutable}
			s.symbols[key] = symbol
			c.info.Symbols[id] = symbol
			current[id] = true

		case ast.DisplayStatement:
			for _, value := range stmt.Values {
				if _, err := c.expressionType(value, s, current, true); err != nil {
					return nil, err
				}
			}

		case ast.MoveStatement:
			target, err := c.resolveReference(stmt.Target, s)
			if err != nil {
				return nil, fmt.Errorf("MOVE target %s: %w", FormatReference(stmt.Target), err)
			}
			if !target.Mutable {
				return nil, fmt.Errorf("MOVE target %q is a LET binding and cannot be reassigned", stmt.Target.Name)
			}
			if err := c.checkAssignmentValue(stmt.Source, target.Type, s, current); err != nil {
				return nil, fmt.Errorf("MOVE to %s: %w", FormatReference(stmt.Target), err)
			}
			current[target.ID] = true

		case ast.IfStatement:
			if err := c.checkCondition(stmt.Condition, s, current); err != nil {
				return nil, fmt.Errorf("IF condition: %w", err)
			}
			thenScope := &scope{parent: s, symbols: make(map[string]Symbol)}
			thenFlow, err := c.checkStatements(thenScope, current, stmt.Then)
			if err != nil {
				return nil, err
			}
			if len(stmt.Else) == 0 {
				continue
			}
			elseScope := &scope{parent: s, symbols: make(map[string]Symbol)}
			elseFlow, err := c.checkStatements(elseScope, current, stmt.Else)
			if err != nil {
				return nil, err
			}
			merged := cloneFlow(current)
			for id := range current {
				merged[id] = thenFlow[id] && elseFlow[id]
			}
			current = merged

		default:
			return nil, fmt.Errorf("unsupported statement node %T", statement)
		}
	}
	return current, nil
}

func (c *checker) checkCondition(condition ast.Condition, s *scope, state flow) error {
	switch cond := condition.(type) {
	case ast.ValueCondition:
		kind, err := c.expressionType(cond.Value, s, state, true)
		if err != nil {
			return err
		}
		if kind != types.Boolean {
			return fmt.Errorf("condition must be BOOLEAN, got %s", kind)
		}
		return nil

	case ast.ComparisonCondition:
		left, err := c.expressionType(cond.Left, s, state, true)
		if err != nil {
			return err
		}
		right, err := c.expressionType(cond.Right, s, state, true)
		if err != nil {
			return err
		}
		return checkComparisonTypes(cond.Op, left, right)

	case ast.LogicalCondition:
		if err := c.checkCondition(cond.Left, s, state); err != nil {
			return err
		}
		return c.checkCondition(cond.Right, s, state)

	case ast.NotCondition:
		return c.checkCondition(cond.Inner, s, state)

	default:
		return fmt.Errorf("unsupported condition node %T", condition)
	}
}

func checkComparisonTypes(op ast.ComparisonOperator, left, right types.Kind) error {
	if op == ast.CompareEqual || op == ast.CompareNotEqual {
		if left == right || (isNumeric(left) && isNumeric(right)) {
			return nil
		}
		return fmt.Errorf("cannot compare %s with %s", left, right)
	}
	if !isNumeric(left) || !isNumeric(right) {
		return fmt.Errorf("ordered comparison requires numeric operands, got %s and %s", left, right)
	}
	return nil
}

func isNumeric(kind types.Kind) bool {
	switch kind {
	case types.Byte, types.Integer, types.Long, types.Decimal, types.Float, types.Double:
		return true
	default:
		return false
	}
}

func (c *checker) resolveReference(ref ast.DataReference, s *scope) (Symbol, error) {
	if len(ref.Qualifiers) == 0 {
		if symbol, ok := s.lookup(ref.Name); ok {
			return symbol, nil
		}
		ids := c.dataNames[normalize(ref.Name)]
		switch len(ids) {
		case 0:
			return Symbol{}, fmt.Errorf("is not declared")
		case 1:
			return c.info.Symbols[ids[0]], nil
		default:
			return Symbol{}, fmt.Errorf("is ambiguous; qualify it with OF")
		}
	}

	id, ok := c.info.Qualified[ReferenceKey(ref)]
	if !ok {
		return Symbol{}, fmt.Errorf("is not declared with that qualification")
	}
	return c.info.Symbols[id], nil
}

func (c *checker) checkAssignmentValue(expr ast.Expression, target types.Type, s *scope, state flow) error {
	sourceKind, err := c.expressionType(expr, s, state, true)
	if err != nil {
		return err
	}
	if !types.Assignable(sourceKind, target.Kind) && !literalFitsTarget(expr, target.Kind) {
		return fmt.Errorf("cannot assign %s to %s", sourceKind, target.Kind)
	}
	return validatePictureLiteral(expr, target.Picture)
}

func (c *checker) expressionType(expr ast.Expression, s *scope, state flow, requireInitialized bool) (types.Kind, error) {
	switch value := expr.(type) {
	case ast.StringLiteral:
		return types.String, nil
	case ast.NumberLiteral:
		return types.InferNumberLiteral(value.Value)
	case ast.BooleanLiteral:
		return types.Boolean, nil
	case ast.DataReference:
		symbol, err := c.resolveReference(value, s)
		if err != nil {
			return types.Invalid, fmt.Errorf("%s %w", FormatReference(value), err)
		}
		if requireInitialized && !state[symbol.ID] {
			return types.Invalid, fmt.Errorf("%s is used before initialization", FormatReference(value))
		}
		return symbol.Type.Kind, nil
	default:
		return types.Invalid, fmt.Errorf("unsupported expression node %T", expr)
	}
}

func (s *scope) lookup(name string) (Symbol, bool) {
	key := normalize(name)
	for current := s; current != nil; current = current.parent {
		if symbol, ok := current.symbols[key]; ok {
			return symbol, true
		}
	}
	return Symbol{}, false
}

func cloneFlow(in flow) flow {
	out := make(flow, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
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

func ReferenceKey(ref ast.DataReference) string {
	parts := []string{normalize(ref.Name)}
	for _, qualifier := range ref.Qualifiers {
		parts = append(parts, normalize(qualifier))
	}
	return strings.Join(parts, "|")
}

func FormatReference(ref ast.DataReference) string {
	if len(ref.Qualifiers) == 0 {
		return ref.Name
	}
	var b strings.Builder
	b.WriteString(ref.Name)
	for _, qualifier := range ref.Qualifiers {
		b.WriteString(" OF ")
		b.WriteString(qualifier)
	}
	return b.String()
}

func normalize(name string) string { return strings.ToUpper(name) }
