package music

import (
	"fmt"
	"strings"
)

type Scale struct {
	Name      string
	Intervals []int
	Notes     []string
}

var scaleDefinitions = map[string]Scale{
	"major": {
		Name:      "Major",
		Intervals: []int{0, 2, 4, 5, 7, 9, 11},
		Notes:     []string{"C", "D", "E", "F", "G", "A", "B"},
	},
	"minor": {
		Name:      "Natural Minor",
		Intervals: []int{0, 2, 3, 5, 7, 8, 10},
		Notes:     []string{"C", "D", "Eb", "F", "G", "Ab", "Bb"},
	},
	"dorian": {
		Name:      "Dorian",
		Intervals: []int{0, 2, 3, 5, 7, 9, 10},
		Notes:     []string{"C", "D", "Eb", "F", "G", "A", "Bb"},
	},
	"phrygian": {
		Name:      "Phrygian",
		Intervals: []int{0, 1, 3, 5, 7, 8, 10},
		Notes:     []string{"C", "Db", "Eb", "F", "G", "Ab", "Bb"},
	},
	"lydian": {
		Name:      "Lydian",
		Intervals: []int{0, 2, 4, 6, 7, 9, 11},
		Notes:     []string{"C", "D", "E", "F#", "G", "A", "B"},
	},
	"mixolydian": {
		Name:      "Mixolydian",
		Intervals: []int{0, 2, 4, 5, 7, 9, 10},
		Notes:     []string{"C", "D", "E", "F", "G", "A", "Bb"},
	},
	"locrian": {
		Name:      "Locrian",
		Intervals: []int{0, 1, 3, 5, 6, 8, 10},
		Notes:     []string{"C", "Db", "Eb", "F", "Gb", "Ab", "Bb"},
	},
	"chromatic": {
		Name:      "Chromatic",
		Intervals: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		Notes:     []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"},
	},
	"pentatonic": {
		Name:      "Pentatonic Major",
		Intervals: []int{0, 2, 4, 7, 9},
		Notes:     []string{"C", "D", "E", "G", "A"},
	},
	"pentatonic_minor": {
		Name:      "Pentatonic Minor",
		Intervals: []int{0, 3, 5, 7, 10},
		Notes:     []string{"C", "Eb", "F", "G", "Bb"},
	},
	"blues": {
		Name:      "Blues",
		Intervals: []int{0, 3, 5, 6, 7, 10},
		Notes:     []string{"C", "Eb", "F", "F#", "G", "Bb"},
	},
	"harmonic_minor": {
		Name:      "Harmonic Minor",
		Intervals: []int{0, 2, 3, 5, 7, 8, 11},
		Notes:     []string{"C", "D", "Eb", "F", "G", "Ab", "B"},
	},
	"melodic_minor": {
		Name:      "Melodic Minor",
		Intervals: []int{0, 2, 3, 5, 7, 9, 11},
		Notes:     []string{"C", "D", "Eb", "F", "G", "A", "B"},
	},
}

func GetScale(scaleName string) (Scale, bool) {
	scaleName = strings.ToLower(scaleName)
	scale, exists := scaleDefinitions[scaleName]
	return scale, exists
}

func GetScaleNotes(scaleName string, rootNote string) ([]int, error) {
	scale, exists := GetScale(scaleName)
	if !exists {
		return nil, fmt.Errorf("scale '%s' not found", scaleName)
	}

	rootMidi := 60
	if rootNote != "" {
		notes, err := ParseNote(rootNote, 60)
		if err != nil {
			return nil, err
		}
		if len(notes) > 0 {
			rootMidi = notes[0].MidiValue
		}
	}

	scaleNotes := make([]int, len(scale.Intervals))
	for i, interval := range scale.Intervals {
		scaleNotes[i] = rootMidi + interval
	}
	
	return scaleNotes, nil
}

func QuantizeToScale(midiNote int, scaleName string, rootNote string) (int, error) {
	scaleNotes, err := GetScaleNotes(scaleName, rootNote)
	if err != nil {
		return midiNote, err
	}

	closestNote := midiNote
	minDistance := 127

	// Search enough octaves to cover the full MIDI range (0-127)
	// Calculate how many octaves we need to search both up and down
	maxOctaves := 11 // This covers MIDI 0-127 (about 10 octaves)
	
	for octave := -maxOctaves; octave <= maxOctaves; octave++ {
		for _, scaleNote := range scaleNotes {
			testNote := scaleNote + (octave * 12)
			// Only consider notes in valid MIDI range
			if testNote >= 0 && testNote <= 127 {
				distance := abs(midiNote - testNote)
				if distance < minDistance {
					minDistance = distance
					closestNote = testNote
				}
			}
		}
	}

	return closestNote, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ListScales() []string {
	scales := make([]string, 0, len(scaleDefinitions))
	for name := range scaleDefinitions {
		scales = append(scales, name)
	}
	return scales
}