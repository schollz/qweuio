package step

import "testing"

func TestParseStep(t *testing.T) {
	tests := []struct {
		stepString string
		parsed     Step
	}{
		{"Cmaj@u2d4", Step{NotesString: []string{"Cmaj"}, Arpeggio: []string{"u2d4"}}},
	}

	for _, test := range tests {
		t.Run(test.stepString, func(t *testing.T) {
			parsed := Step{Original: test.stepString}
			err := parsed.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(parsed.NotesString) != len(test.parsed.NotesString) {
				t.Fatalf("expected %d notes, got %d", len(test.parsed.NotesString), len(parsed.NotesString))
			}
			if len(parsed.Arpeggio) != len(test.parsed.Arpeggio) {
				t.Fatalf("expected %d arpeggios, got %d", len(test.parsed.Arpeggio), len(parsed.Arpeggio))
			}
			if parsed.NotesString[0] != test.parsed.NotesString[0] {
				t.Fatalf("expected %s, got %s", test.parsed.NotesString[0], parsed.NotesString[0])
			}
			if parsed.Arpeggio[0] != test.parsed.Arpeggio[0] {
				t.Fatalf("expected %s, got %s", test.parsed.Arpeggio[0], parsed.Arpeggio[0])
			}
		})
	}

}
