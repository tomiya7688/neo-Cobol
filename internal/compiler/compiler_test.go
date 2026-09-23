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

func TestDataMoveDisplay(t *testing.T) {
	source := `
01 NAME TYPE STRING PIC X(20) VALUE "KADOKA".
01 COUNT TYPE INTEGER PIC 9(3) VALUE 7.
MOVE "MARU" TO NAME.
MOVE 42 TO COUNT.
DISPLAY NAME " " COUNT.
`
	generated, err := EmitC(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`const char * neo_v_0 = "KADOKA";`, `int32_t neo_v_1 = 7;`, `neo_v_0 = "MARU";`, `neo_v_1 = 42;`} {
		if !strings.Contains(generated, want) {
			t.Fatalf("generated C missing %q:\n%s", want, generated)
		}
	}
}

func TestPicOnlyInference(t *testing.T) {
	if _, err := EmitC(`01 NAME PIC X(8) VALUE "KADOKA". DISPLAY NAME.`); err != nil {
		t.Fatal(err)
	}
}

func TestUseBeforeInitialization(t *testing.T) {
	_, err := EmitC(`01 NAME TYPE STRING. DISPLAY NAME.`)
	if err == nil || !strings.Contains(err.Error(), "before initialization") {
		t.Fatalf("got error %v", err)
	}
}

func TestTypeMismatch(t *testing.T) {
	_, err := EmitC(`01 COUNT TYPE INTEGER VALUE 1. MOVE "NO" TO COUNT.`)
	if err == nil || !strings.Contains(err.Error(), "cannot assign STRING to INTEGER") {
		t.Fatalf("got error %v", err)
	}
}

func TestPicLiteralConstraint(t *testing.T) {
	_, err := EmitC(`01 CODE TYPE STRING PIC X(3) VALUE "TOOLONG".`)
	if err == nil || !strings.Contains(err.Error(), "exceeds PIC") {
		t.Fatalf("got error %v", err)
	}
}

func TestNestedLevelRejectedUntilRecordModelExists(t *testing.T) {
	_, err := EmitC(`05 CHILD TYPE INTEGER VALUE 1.`)
	if err == nil || !strings.Contains(err.Error(), "record hierarchy support") {
		t.Fatalf("got error %v", err)
	}
}

func TestBooleanVariable(t *testing.T) {
	generated, err := EmitC(`01 READY TYPE BOOLEAN VALUE TRUE. DISPLAY READY.`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(generated, `bool neo_v_0 = true;`) {
		t.Fatalf("generated C missing BOOLEAN declaration:\n%s", generated)
	}
}

func TestBasedIntegerFitsByte(t *testing.T) {
	generated, err := EmitC(`01 MASK TYPE BYTE VALUE 0xDE. DISPLAY MASK.`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(generated, `uint8_t neo_v_0 = 222;`) {
		t.Fatalf("generated C missing based integer lowering:\n%s", generated)
	}
}

func TestVARLETInferenceAndDeclarationOrder(t *testing.T) {
	source := `
01 SOURCE TYPE INTEGER VALUE 1.
MOVE 2 TO SOURCE.
VAR SNAPSHOT VALUE SOURCE.
LET LABEL VALUE "SNAPSHOT=".
DISPLAY LABEL SNAPSHOT.
MOVE 3 TO SNAPSHOT.
DISPLAY "VAR=" SNAPSHOT.
`
	generated, err := EmitC(source)
	if err != nil {
		t.Fatal(err)
	}
	moveSource := strings.Index(generated, "neo_v_0 = 2;")
	declareSnapshot := strings.Index(generated, "int32_t neo_v_1 = neo_v_0;")
	declareLabel := strings.Index(generated, `const char * neo_v_2 = "SNAPSHOT=";`)
	moveSnapshot := strings.Index(generated, "neo_v_1 = 3;")
	if moveSource < 0 || declareSnapshot < 0 || declareLabel < 0 || moveSnapshot < 0 {
		t.Fatalf("generated C missing binding lowering:\n%s", generated)
	}
	if !(moveSource < declareSnapshot && declareSnapshot < declareLabel && declareLabel < moveSnapshot) {
		t.Fatalf("binding declarations were reordered:\n%s", generated)
	}
}

func TestLETReassignmentRejected(t *testing.T) {
	_, err := EmitC(`LET LIMIT VALUE 10. MOVE 20 TO LIMIT.`)
	if err == nil || !strings.Contains(err.Error(), "LET binding") {
		t.Fatalf("got error %v", err)
	}
}

func TestBindingRequiresValue(t *testing.T) {
	_, err := EmitC(`VAR COUNT.`)
	if err == nil || !strings.Contains(err.Error(), "expected VALUE") {
		t.Fatalf("got error %v", err)
	}
}

func TestBindingDuplicateRejected(t *testing.T) {
	_, err := EmitC(`01 NAME TYPE STRING VALUE "A". VAR NAME VALUE "B".`)
	if err == nil || !strings.Contains(err.Error(), "duplicate declaration") {
		t.Fatalf("got error %v", err)
	}
}

func TestBindingInitializerUsesDefiniteAssignment(t *testing.T) {
	_, err := EmitC(`01 SOURCE TYPE INTEGER. VAR COPY VALUE SOURCE.`)
	if err == nil || !strings.Contains(err.Error(), "before initialization") {
		t.Fatalf("got error %v", err)
	}
}

func TestLongBindingInference(t *testing.T) {
	generated, err := EmitC(`VAR BIG VALUE 2147483648. DISPLAY BIG.`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(generated, `int64_t neo_v_0 = 2147483648;`) {
		t.Fatalf("generated C missing LONG inference:\n%s", generated)
	}
}

func TestStringAndBooleanBindingInference(t *testing.T) {
	generated, err := EmitC(`LET TITLE VALUE "REPORT". VAR READY VALUE TRUE. DISPLAY TITLE " " READY.`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(generated, `const char * neo_v_0 = "REPORT";`) || !strings.Contains(generated, `bool neo_v_1 = true;`) {
		t.Fatalf("generated C missing STRING/BOOLEAN inference:\n%s", generated)
	}
}
