package pattern

import (
	"testing"
)

func TestParsePattern(t *testing.T) {
	patternString := `
_ _ 0,1 c
_ Cmaj@u2d4
[[a a] a] a`
	pattern, err := Parse(patternString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pattern
}
