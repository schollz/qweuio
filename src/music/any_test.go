package music

import (
	"testing"
)

func TestParseAny(t *testing.T) {
	tests := []struct {
		midiString string
		midiNear   int
		expected   []Note
	}{
		{"C/G", 71, []Note{
			{MidiValue: 67, NameSharp: "g4"},
			{MidiValue: 72, NameSharp: "c5"},
			{MidiValue: 76, NameSharp: "e5"},
		}},
		{"c6", 20, []Note{{MidiValue: 84, NameSharp: "c6"}}},
		{"c#", 60, []Note{{MidiValue: 61, NameSharp: "c#4"}}},
	}
	for _, test := range tests {
		notes, err := Parse(test.midiString, test.midiNear)
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
