package pattern

import (
	"testing"
)

func TestParseNotePattern(t *testing.T) {
	patternString := `
# a_notes
_ _ a b
_ Cmaj@u2d4
[[a a] a] a`
	pattern, err := Parse(patternString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pattern
}

func TestParseNotePatternDouble(t *testing.T) {
	patternString := `
# a_notes
_ _ a4,a5 b
`
	pattern, err := Parse(patternString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pattern
}

func TestParseVelocityPattern(t *testing.T) {
	patternString := `
! a_velocity
30 40
60 80 90 90`
	pattern, err := Parse(patternString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pattern
}
