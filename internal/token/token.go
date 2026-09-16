package token

type Kind uint8

const (
	EOF Kind = iota
	Identifier
	String
	Number
	Period
)

type Token struct {
	Kind   Kind
	Lexeme string
	Line   int
	Column int
}

func (t Token) IsWord(word string) bool {
	if t.Kind != Identifier {
		return false
	}
	if len(t.Lexeme) != len(word) {
		return false
	}
	for i := range t.Lexeme {
		a, b := t.Lexeme[i], word[i]
		if a >= 'a' && a <= 'z' {
			a -= 'a' - 'A'
		}
		if b >= 'a' && b <= 'z' {
			b -= 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}
