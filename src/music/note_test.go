package music

import (
	"fmt"
	"testing"
)

func TestNote(t *testing.T) {
	n := Note{MidiValue: 60, NameSharp: "C"}
	n2 := n.Add(1)
	if n2.MidiValue != 61 {
		t.Fatalf("expected %d, got %d", 61, n2.MidiValue)
	}
	fmt.Println(n2)

}

func TestParseNote(t *testing.T) {
	// table driven tests
	tests := []struct {
		midiString string
		midiNear   int
		expected   []Note
	}{
		{"c", 71, []Note{{MidiValue: 72, NameSharp: "c5"}}},
		{"c6", 20, []Note{{MidiValue: 84, NameSharp: "c6"}}},
		{"c", 62, []Note{{MidiValue: 60, NameSharp: "c4"}}},
		{"d", 32, []Note{{MidiValue: 26, NameSharp: "d1"}}},
		{"f#3", 32, []Note{{MidiValue: 54, NameSharp: "f#3"}}},
		{"g7", 100, []Note{{MidiValue: 103, NameSharp: "g7"}}},
		{"gb", 100, []Note{{MidiValue: 103, NameSharp: "g7"}, {MidiValue: 107, NameSharp: "b7"}}},
		{"g♭c", 100, []Note{{MidiValue: 102, NameSharp: "f#7"}, {MidiValue: 96, NameSharp: "c7"}}},
		{"c4eg", 52, []Note{
			{MidiValue: 60, NameSharp: "c4"},
			{MidiValue: 64, NameSharp: "e4"},
			{MidiValue: 67, NameSharp: "g4"},
		}},
		{"ceg6", 52, []Note{
			{MidiValue: 48, NameSharp: "c3"},
			{MidiValue: 52, NameSharp: "e3"},
			{MidiValue: 91, NameSharp: "g6"},
		}},
	}
	for _, test := range tests {
		notes, err := ParseNote(test.midiString, test.midiNear)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if len(notes) != len(test.expected) {
			t.Errorf("expected %v, got %v", test.expected, notes)
			continue
		}
		for i, note := range notes {
			if note.MidiValue != test.expected[i].MidiValue || note.NameSharp != test.expected[i].NameSharp {
				t.Errorf("test: %s (%d), expected %v, got %v", test.midiString, test.midiNear, test.expected, notes)
				break
			}
		}
	}
}
