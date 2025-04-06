package expand_arpeggio

import (
	"asdfgh/src/music"
	"fmt"
	"testing"

	log "github.com/schollz/logger"
)

func TestExpandArpeggio(t *testing.T) {
	log.Tracef("ExpandArpeggio")
	// Test cases for ExpandArpeggio function
	tests := []struct {
		notestring string
		arpstring  string
		notes      []music.Note
	}{
		{"C", "u4d2", []music.Note{music.C4, music.E4, music.G4, music.C5, music.G4, music.E4}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.notestring), func(t *testing.T) {
			notes, err := ExpandArpeggio(test.notestring, test.arpstring)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(notes) != len(test.notes) {
				// print out notes
				for _, note := range notes {
					log.Tracef("note: %v", note.NameSharp)
				}
				t.Fatalf("expected %d notes, got %d", len(test.notes), len(notes))
			}
			for i, note := range notes {
				if note.MidiValue != test.notes[i].MidiValue {
					t.Fatalf("expected note %v, got %v", test.notes[i].NameSharp, note.NameSharp)
				}
			}
		})
	}
}
