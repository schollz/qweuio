package music

import (
	"testing"
)

func TestParseChord(t *testing.T) {
	// table driven tests
	tests := []struct {
		chordString string
		midiNear    int
		expected    []Note
	}{
		{"Dm7/A;3", 60, []Note{
			{MidiValue: 57, NameSharp: "a3"},
			{MidiValue: 60, NameSharp: "c4"},
			{MidiValue: 62, NameSharp: "d4"},
			{MidiValue: 65, NameSharp: "f4"},
		}},
		{"Cm", 32, []Note{
			{MidiValue: 24, NameSharp: "c1"},
			{MidiValue: 27, NameSharp: "d#1"},
			{MidiValue: 31, NameSharp: "g1"},
		}},
		{"Amaj7", 70, []Note{
			{MidiValue: 69, NameSharp: "a4"},
			{MidiValue: 73, NameSharp: "c#5"},
			{MidiValue: 76, NameSharp: "e5"},
			{MidiValue: 80, NameSharp: "g#5"},
		}},
		{"G", 70, []Note{
			{MidiValue: 67, NameSharp: "g4"},
			{MidiValue: 71, NameSharp: "b4"},
			{MidiValue: 74, NameSharp: "d5"},
		}},
		{"G7", 70, []Note{
			{MidiValue: 67, NameSharp: "g4"},
			{MidiValue: 71, NameSharp: "b4"},
			{MidiValue: 74, NameSharp: "d5"},
			{MidiValue: 77, NameSharp: "f5"},
		}},
		{"Gmaj7", 70, []Note{
			{MidiValue: 67, NameSharp: "g4"},
			{MidiValue: 71, NameSharp: "b4"},
			{MidiValue: 74, NameSharp: "d5"},
			{MidiValue: 78, NameSharp: "f#5"},
		}},
		{"Gmaj7/F#", 70, []Note{
			{MidiValue: 66, NameSharp: "f#4"},
			{MidiValue: 67, NameSharp: "g4"},
			{MidiValue: 71, NameSharp: "b4"},
			{MidiValue: 74, NameSharp: "d5"},
		}},
		{"Cm", 70, []Note{
			{MidiValue: 60, NameSharp: "c4"},
			{MidiValue: 63, NameSharp: "d#4"},
			{MidiValue: 67, NameSharp: "g4"},
		}},
	}

	for _, test := range tests {
		midiNotes, err := ParseChord(test.chordString, test.midiNear)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if len(midiNotes) != len(test.expected) {
			t.Errorf("'%s': Expected %d notes, got %d", test.chordString, len(test.expected), len(midiNotes))
		}
		for i, note := range midiNotes {
			if note.MidiValue != test.expected[i].MidiValue {
				t.Errorf("'%s': Expected %d, got %d", test.chordString, test.expected[i].MidiValue, note.MidiValue)
			}
			if note.NameSharp != test.expected[i].NameSharp {
				t.Errorf("'%s': Expected %s, got %s", test.chordString, test.expected[i].NameSharp, note.NameSharp)
			}
		}
	}

}
