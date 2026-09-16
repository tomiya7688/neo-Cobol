package lexer

import (
	"testing"

	"github.com/tomiya7688/neo-Cobol/internal/token"
)

func TestLexDisplayAndComment(t *testing.T) {
	tokens, err := Lex("display \"HELLO\". *> comment\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 4 {
		t.Fatalf("got %d tokens, want 4", len(tokens))
	}
	if !tokens[0].IsWord("DISPLAY") {
		t.Fatalf("first token = %#v", tokens[0])
	}
	if tokens[1].Kind != token.String || tokens[1].Lexeme != "HELLO" {
		t.Fatalf("string token = %#v", tokens[1])
	}
	if tokens[2].Kind != token.Period {
		t.Fatalf("period token = %#v", tokens[2])
	}
}

func TestLexDoubledQuote(t *testing.T) {
	tokens, err := Lex(`DISPLAY "He said ""HELLO"".".`)
	if err != nil {
		t.Fatal(err)
	}
	if got := tokens[1].Lexeme; got != `He said "HELLO".` {
		t.Fatalf("got %q", got)
	}
}
