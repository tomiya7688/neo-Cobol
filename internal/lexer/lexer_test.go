package lexer

import (
	"github.com/tomiya7688/neo-Cobol/internal/token"
	"testing"
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
}

func TestLexPicture(t *testing.T) {
	tokens, err := Lex(`01 NAME PIC X(20).`)
	if err != nil {
		t.Fatal(err)
	}
	var got []token.Kind
	for _, tok := range tokens {
		got = append(got, tok.Kind)
	}
	if len(got) < 8 || got[4] != token.LParen || got[6] != token.RParen {
		t.Fatalf("tokens = %#v", tokens)
	}
}

func TestLexBasedInteger(t *testing.T) {
	tokens, err := Lex(`DISPLAY 0xFF.`)
	if err != nil {
		t.Fatal(err)
	}
	if tokens[1].Kind != token.Number || tokens[1].Lexeme != "0xFF" {
		t.Fatalf("number = %#v", tokens[1])
	}
}
