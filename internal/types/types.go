package types

import (
	"fmt"
	"strconv"
	"strings"
)

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

func ParseBuiltin(name string) (Kind, bool) {
	switch strings.ToUpper(name) {
	case "INTEGER":
		return Integer, true
	case "DECIMAL":
		return Decimal, true
	case "STRING":
		return String, true
	case "BOOLEAN":
		return Boolean, true
	case "FLOAT":
		return Float, true
	case "BYTE":
		return Byte, true
	case "LONG":
		return Long, true
	case "DOUBLE":
		return Double, true
	default:
		return Invalid, false
	}
}

type PictureKind uint8

const (
	PictureInvalid PictureKind = iota
	PictureText
	PictureNumeric
)

type Picture struct {
	Raw    string
	Kind   PictureKind
	Width  int
	Scale  int
	Signed bool
}

type Type struct {
	Kind    Kind
	Picture *Picture
}

func ParsePicture(raw string) (*Picture, error) {
	s := strings.ToUpper(strings.ReplaceAll(raw, " ", ""))
	if s == "" {
		return nil, fmt.Errorf("empty PIC")
	}
	if s[0] == 'X' {
		width, rest, err := parsePicRun(s, 'X')
		if err != nil || rest != "" {
			return nil, fmt.Errorf("unsupported text PIC %q", raw)
		}
		return &Picture{Raw: s, Kind: PictureText, Width: width}, nil
	}
	p := &Picture{Raw: s, Kind: PictureNumeric}
	if strings.HasPrefix(s, "S") {
		p.Signed = true
		s = s[1:]
	}
	integerDigits, rest, err := parsePicRun(s, '9')
	if err != nil {
		return nil, fmt.Errorf("unsupported numeric PIC %q", raw)
	}
	p.Width = integerDigits
	if rest != "" {
		if !strings.HasPrefix(rest, "V") {
			return nil, fmt.Errorf("unsupported numeric PIC %q", raw)
		}
		fractionDigits, tail, err := parsePicRun(rest[1:], '9')
		if err != nil || tail != "" {
			return nil, fmt.Errorf("unsupported numeric PIC %q", raw)
		}
		p.Width += fractionDigits
		p.Scale = fractionDigits
	}
	return p, nil
}

func parsePicRun(s string, symbol byte) (int, string, error) {
	if s == "" || s[0] != symbol {
		return 0, s, fmt.Errorf("expected %c", symbol)
	}
	count := 0
	for len(s) > 0 && s[0] == symbol {
		count++
		s = s[1:]
	}
	if len(s) > 0 && s[0] == '(' {
		if count != 1 {
			return 0, s, fmt.Errorf("mixed repeated and parenthesized PIC")
		}
		end := strings.IndexByte(s, ')')
		if end <= 1 {
			return 0, s, fmt.Errorf("invalid PIC repetition")
		}
		n, err := strconv.Atoi(s[1:end])
		if err != nil || n <= 0 {
			return 0, s, fmt.Errorf("invalid PIC repetition")
		}
		count = n
		s = s[end+1:]
	}
	return count, s, nil
}

func InferFromPicture(pic *Picture) (Kind, error) {
	switch pic.Kind {
	case PictureText:
		return String, nil
	case PictureNumeric:
		if pic.Scale > 0 {
			return Decimal, nil
		}
		if pic.Width <= 9 {
			return Integer, nil
		}
		if pic.Width <= 18 {
			return Long, nil
		}
		return Invalid, fmt.Errorf("numeric PIC %s exceeds initial INTEGER/LONG inference range", pic.Raw)
	default:
		return Invalid, fmt.Errorf("invalid PIC")
	}
}

func CompatiblePicture(kind Kind, pic *Picture) bool {
	if pic == nil {
		return true
	}
	switch kind {
	case String:
		return pic.Kind == PictureText
	case Byte, Integer, Long:
		return pic.Kind == PictureNumeric && pic.Scale == 0
	case Decimal:
		return pic.Kind == PictureNumeric
	default:
		return false
	}
}

func InferNumberLiteral(literal string) (Kind, error) {
	unsigned := strings.TrimPrefix(strings.TrimPrefix(literal, "+"), "-")
	isBased := strings.HasPrefix(unsigned, "0x") || strings.HasPrefix(unsigned, "0X") ||
		strings.HasPrefix(unsigned, "0b") || strings.HasPrefix(unsigned, "0B") ||
		strings.HasPrefix(unsigned, "0o") || strings.HasPrefix(unsigned, "0O")
	if !isBased && strings.ContainsAny(literal, ".eE") {
		return Decimal, nil
	}
	v, err := strconv.ParseInt(literal, 0, 64)
	if err != nil {
		return Invalid, fmt.Errorf("integer literal %q is outside LONG range", literal)
	}
	if v >= -2147483648 && v <= 2147483647 {
		return Integer, nil
	}
	return Long, nil
}

func Assignable(from, to Kind) bool {
	if from == to {
		return true
	}
	if from == Byte && (to == Integer || to == Long) {
		return true
	}
	if from == Integer && to == Long {
		return true
	}
	return false
}
