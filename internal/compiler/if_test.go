package compiler

import (
	"strings"
	"testing"
)

func TestIfComparisonAndLogical(t *testing.T) {
	source := `01 AGE TYPE INTEGER VALUE 20.
01 READY TYPE BOOLEAN VALUE TRUE.
IF AGE IS GREATER THAN OR EQUAL TO 18 AND NOT FALSE
    LET LABEL VALUE "ADULT"
    DISPLAY LABEL
ELSE
    DISPLAY "MINOR"
END-IF.
DISPLAY "DONE".`
	generated, err := EmitC(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"if (", ">=", "&&", "!", "ADULT", "MINOR"} {
		if !strings.Contains(generated, want) {
			t.Fatalf("missing %q:\n%s", want, generated)
		}
	}
}
func TestNestedIfAndEndIfShorthand(t *testing.T) {
	source := `VAR X VALUE 2.
IF X IS GREATER THAN 1
    IF TRUE
        DISPLAY "YES"
    END IF
ELSE
    DISPLAY "NO"
END-IF.`
	if _, err := EmitC(source); err != nil {
		t.Fatal(err)
	}
}
func TestBlockShadowing(t *testing.T) {
	source := `LET LABEL VALUE "OUTER".
IF TRUE
    LET LABEL VALUE "INNER"
    DISPLAY LABEL
ELSE
    DISPLAY "NO"
END-IF.
DISPLAY LABEL.`
	generated, err := EmitC(source)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(generated, "const char *") != 2 {
		t.Fatalf("expected shadowed locals:\n%s", generated)
	}
}
func TestBlockLocalNotVisibleAfterIf(t *testing.T) {
	_, err := EmitC(`IF TRUE LET TEMP VALUE 1. END-IF. DISPLAY TEMP.`)
	if err == nil || !strings.Contains(err.Error(), "not declared") {
		t.Fatalf("got %v", err)
	}
}
func TestConditionRequiresBoolean(t *testing.T) {
	_, err := EmitC(`IF 1 DISPLAY "NO" END-IF.`)
	if err == nil || !strings.Contains(err.Error(), "condition must be BOOLEAN") {
		t.Fatalf("got %v", err)
	}
}
func TestComparisonTypeMismatch(t *testing.T) {
	_, err := EmitC(`IF "A" IS EQUAL TO 1 DISPLAY "NO" END-IF.`)
	if err == nil || !strings.Contains(err.Error(), "cannot compare STRING with INTEGER") {
		t.Fatalf("got %v", err)
	}
}
func TestOrderedStringComparisonRejected(t *testing.T) {
	_, err := EmitC(`IF "A" IS LESS THAN "B" DISPLAY "NO" END-IF.`)
	if err == nil || !strings.Contains(err.Error(), "ordered comparison requires numeric operands") {
		t.Fatalf("got %v", err)
	}
}
func TestDefiniteAssignmentAcrossBothBranches(t *testing.T) {
	source := `01 VALUE-A TYPE INTEGER.
IF TRUE
 MOVE 1 TO VALUE-A
ELSE
 MOVE 2 TO VALUE-A
END-IF.
DISPLAY VALUE-A.`
	if _, err := EmitC(source); err != nil {
		t.Fatal(err)
	}
}
func TestDefiniteAssignmentWithoutElseNotGuaranteed(t *testing.T) {
	_, err := EmitC(`01 VALUE-A TYPE INTEGER. IF TRUE MOVE 1 TO VALUE-A END-IF. DISPLAY VALUE-A.`)
	if err == nil || !strings.Contains(err.Error(), "before initialization") {
		t.Fatalf("got %v", err)
	}
}
func TestDuplicateSameScopeStillRejected(t *testing.T) {
	_, err := EmitC(`VAR NAME VALUE "A". LET NAME VALUE "B".`)
	if err == nil || !strings.Contains(err.Error(), "same scope") {
		t.Fatalf("got %v", err)
	}
}
func TestStringEqualityUsesStrcmp(t *testing.T) {
	generated, err := EmitC(`IF "A" IS NOT EQUAL TO "B" DISPLAY "YES" END-IF.`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(generated, `strcmp("A", "B") != 0`) {
		t.Fatalf("missing strcmp:\n%s", generated)
	}
}
