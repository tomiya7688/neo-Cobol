package compiler

import (
	"strings"
	"testing"
)

func TestRecordHierarchyAndQualifiedReferences(t *testing.T) {
	source := `
01 CUSTOMER.
    05 ID TYPE INTEGER VALUE 7.
    05 NAME TYPE STRING VALUE "KADOKA".
    05 ADDRESS.
        10 ZIP TYPE INTEGER VALUE 12345.
01 SUPPLIER.
    05 ID TYPE INTEGER VALUE 9.

DISPLAY NAME OF CUSTOMER " " ZIP OF ADDRESS OF CUSTOMER.
MOVE 8 TO ID OF CUSTOMER.
IF ID OF CUSTOMER IS EQUAL TO 8
    DISPLAY "UPDATED"
END-IF.
`
	generated, err := EmitC(source)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(generated, "struct {") < 3 {
		t.Fatalf("expected nested record structs:\n%s", generated)
	}
	for _, want := range []string{"= 7;", "= 12345;", "= 9;", "UPDATED"} {
		if !strings.Contains(generated, want) {
			t.Fatalf("generated C missing %q:\n%s", want, generated)
		}
	}
}

func TestUniqueBareFieldReference(t *testing.T) {
	_, err := EmitC(`
01 CUSTOMER.
    05 NAME TYPE STRING VALUE "KADOKA".
DISPLAY NAME.
`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAmbiguousBareFieldRequiresQualification(t *testing.T) {
	_, err := EmitC(`
01 CUSTOMER.
    05 ID TYPE INTEGER VALUE 1.
01 SUPPLIER.
    05 ID TYPE INTEGER VALUE 2.
DISPLAY ID.
`)
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("got error %v", err)
	}
}

func TestIncorrectQualificationRejected(t *testing.T) {
	_, err := EmitC(`
01 CUSTOMER.
    05 ID TYPE INTEGER VALUE 1.
DISPLAY ID OF SUPPLIER.
`)
	if err == nil || !strings.Contains(err.Error(), "qualification") {
		t.Fatalf("got error %v", err)
	}
}

func TestNestedRecordFieldUseBeforeInitialization(t *testing.T) {
	_, err := EmitC(`
01 CUSTOMER.
    05 ADDRESS.
        10 ZIP TYPE INTEGER.
DISPLAY ZIP OF ADDRESS OF CUSTOMER.
`)
	if err == nil || !strings.Contains(err.Error(), "before initialization") {
		t.Fatalf("got error %v", err)
	}
}

func TestRecordFieldDefiniteAssignmentAcrossIf(t *testing.T) {
	_, err := EmitC(`
01 CUSTOMER.
    05 ID TYPE INTEGER.
IF TRUE
    MOVE 1 TO ID OF CUSTOMER
ELSE
    MOVE 2 TO ID OF CUSTOMER
END-IF.
DISPLAY ID OF CUSTOMER.
`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGroupCannotHaveScalarDescriptor(t *testing.T) {
	_, err := EmitC(`
01 CUSTOMER TYPE STRING.
    05 ID TYPE INTEGER VALUE 1.
`)
	if err == nil || !strings.Contains(err.Error(), "group item") {
		t.Fatalf("got error %v", err)
	}
}

func TestEmptyGroupRejected(t *testing.T) {
	_, err := EmitC(`01 CUSTOMER.`)
	if err == nil || !strings.Contains(err.Error(), "requires TYPE or PIC") {
		t.Fatalf("got error %v", err)
	}
}
