package compiler

import (
	"strings"
	"testing"
)

func TestEmitCHelloWorld(t *testing.T) {
	generated, err := EmitC(`DISPLAY "HELLO WORLD".`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(generated, `fputs("HELLO WORLD", stdout);`) {
		t.Fatalf("generated C does not contain DISPLAY lowering:\n%s", generated)
	}
}

func TestIdentifierIsRejectedUntilNameResolutionExists(t *testing.T) {
	_, err := EmitC(`DISPLAY NAME.`)
	if err == nil || !strings.Contains(err.Error(), "name resolution") {
		t.Fatalf("got error %v", err)
	}
}
