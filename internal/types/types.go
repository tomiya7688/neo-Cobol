package types

// Kind is the canonical logical type identity used by semantic analysis.
// PIC is intentionally modeled separately as a representation constraint.
type Kind uint8

const (
	Invalid Kind = iota
	Integer
	Decimal
	String
	Boolean
	Float
	Byte
	Long
	Double
)

func (k Kind) String() string {
	switch k {
	case Integer:
		return "INTEGER"
	case Decimal:
		return "DECIMAL"
	case String:
		return "STRING"
	case Boolean:
		return "BOOLEAN"
	case Float:
		return "FLOAT"
	case Byte:
		return "BYTE"
	case Long:
		return "LONG"
	case Double:
		return "DOUBLE"
	default:
		return "<invalid>"
	}
}

type Picture struct {
	Raw string
}

type Type struct {
	Kind    Kind
	Picture *Picture
}
