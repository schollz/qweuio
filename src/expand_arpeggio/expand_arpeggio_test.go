package expand_arpeggio

import (
	"fmt"
	"testing"

	log "github.com/schollz/logger"
)

func TestExpandArpeggio(t *testing.T) {
	log.Tracef("ExpandArpeggio")
	// Test cases for ExpandArpeggio function
	tests := []struct {
		line     string
		expected string
	}{
		{"C@u4d3$v#a", "c4$v#a e4$v#a g4$v#a c5$v#a g4$v#a e4$v#a c4$v#a"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.line), func(t *testing.T) {
			result, err := ExpandArpeggio(test.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Fatalf("\n\t%s -->\n\t%v != %v", test.line, result, test.expected)
			}
		})
	}
}
