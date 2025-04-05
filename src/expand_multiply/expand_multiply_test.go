package expand_multiply

import (
	"fmt"
	"testing"
)

func TestExpandMultiply(t *testing.T) {
	// Test cases for ExpandMultiplication function
	tests := []struct {
		line     string
		expected string
	}{
		{"a b c", "a b c"},
		{"a*3 b c", "[a a a] b c"},
		{"a*3 b*2 c", "[a a a] [b b] c"},
		{"[a b] c", "[a b] c"},
		{"[a b]*3 c", "[[a b] [a b] [a b]] c"},
		{"[[a b] * 2]*2 c", "[[[[a b] [a b]]] [[[a b] [a b]]]] c"},
		{"a  *3 b c", "[a a a] b c"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.line), func(t *testing.T) {
			result := ExpandMultiplication(test.line)
			if result != test.expected {
				t.Fatalf("\n\t%s -->\n\t%v != %v", test.line, result, test.expected)
			}
		})
	}

	// Test with removeBrackets set to true
	tests = []struct {
		line     string
		expected string
	}{
		{"a b c", "a b c"},
		{"a*3 - c", "a a a - c"},
		{"a*3 b*2 c", "a a a b b c"},
		{"[a b] c", "a b c"},
		{"[a b]*3 c", "a b a b a b c"},
		{"[[a b b]*2 a]*2 c", "a b b a b b a a b b a b b a c"},
		{"a  *3 b c", "a a a b c"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.line), func(t *testing.T) {
			result := ExpandMultiplication(test.line, true)
			if result != test.expected {
				t.Fatalf("\n\t%s -->\n\t%v != %v", test.line, result, test.expected)
			}
		})
	}
}
