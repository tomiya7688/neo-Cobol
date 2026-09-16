package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/tomiya7688/neo-Cobol/internal/token"
)

// Lex converts UTF-8 Neo COBOL source into tokens. The bootstrap lexer
// intentionally implements only the lexical forms needed by the first
// compiler slice while preserving the language's case-insensitive model.
func Lex(source string) ([]token.Token, error) {
	if !utf8.ValidString(source) {
		return nil, fmt.Errorf("source is not valid UTF-8")
	}

	var out []token.Token
	line, col := 1, 1

	for i := 0; i < len(source); {
		ch := source[i]

		switch ch {
		case ' ', '\t', '\r':
			i++
			col++
			continue
		case '\n':
			i++
			line++
			col = 1
			continue
		case '.':
			out = append(out, token.Token{Kind: token.Period, Lexeme: ".", Line: line, Column: col})
			i++
			col++
			continue
		case '\'', '"':
			startLine, startCol := line, col
			quote := ch
			i++
			col++
			var value strings.Builder
			closed := false
			for i < len(source) {
				current := source[i]
				if current == quote {
					if i+1 < len(source) && source[i+1] == quote {
						value.WriteByte(quote)
						i += 2
						col += 2
						continue
					}
					i++
					col++
					closed = true
					break
				}
				if current == '\n' {
					return nil, fmt.Errorf("%d:%d: unterminated string literal", startLine, startCol)
				}
				value.WriteByte(current)
				i++
				col++
			}
			if !closed {
				return nil, fmt.Errorf("%d:%d: unterminated string literal", startLine, startCol)
			}
			out = append(out, token.Token{Kind: token.String, Lexeme: value.String(), Line: startLine, Column: startCol})
			continue
		}

		if ch == '*' && i+1 < len(source) && source[i+1] == '>' {
			for i < len(source) && source[i] != '\n' {
				i++
				col++
			}
			continue
		}

		if isIdentifierStart(ch) {
			start, startCol := i, col
			hasLetter := isLetter(ch)
			i++
			col++
			for i < len(source) && isIdentifierMiddle(source[i]) {
				hasLetter = hasLetter || isLetter(source[i])
				i++
				col++
			}
			lexeme := source[start:i]
			if !hasLetter || lexeme[len(lexeme)-1] == '-' || lexeme[len(lexeme)-1] == '_' {
				return nil, fmt.Errorf("%d:%d: invalid identifier %q", line, startCol, lexeme)
			}
			out = append(out, token.Token{Kind: token.Identifier, Lexeme: lexeme, Line: line, Column: startCol})
			continue
		}

		if isNumberStart(source, i) {
			start, startCol := i, col
			i, col = scanNumber(source, i, col)
			out = append(out, token.Token{Kind: token.Number, Lexeme: source[start:i], Line: line, Column: startCol})
			continue
		}

		if ch >= 0x80 {
			return nil, fmt.Errorf("%d:%d: non-ASCII identifiers are not defined yet", line, col)
		}
		return nil, fmt.Errorf("%d:%d: unexpected character %q", line, col, ch)
	}

	out = append(out, token.Token{Kind: token.EOF, Line: line, Column: col})
	return out, nil
}

func isLetter(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }

func isIdentifierStart(ch byte) bool { return isLetter(ch) }

func isIdentifierMiddle(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || ch == '-' || ch == '_'
}

func isNumberStart(source string, i int) bool {
	if isDigit(source[i]) {
		return true
	}
	return (source[i] == '+' || source[i] == '-') && i+1 < len(source) && isDigit(source[i+1])
}

func scanNumber(source string, i, col int) (int, int) {
	if source[i] == '+' || source[i] == '-' {
		i++
		col++
	}
	for i < len(source) && isDigit(source[i]) {
		i++
		col++
	}
	if i < len(source) && source[i] == '.' && i+1 < len(source) && isDigit(source[i+1]) {
		i++
		col++
		for i < len(source) && isDigit(source[i]) {
			i++
			col++
		}
	}
	if i < len(source) && (source[i] == 'e' || source[i] == 'E') {
		j := i + 1
		if j < len(source) && (source[j] == '+' || source[j] == '-') {
			j++
		}
		if j < len(source) && isDigit(source[j]) {
			for i < j {
				i++
				col++
			}
			for i < len(source) && isDigit(source[i]) {
				i++
				col++
			}
		}
	}
	return i, col
}
