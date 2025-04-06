package music

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
)

type Notes struct {
	Original string
	Note     []Note
}

type Note struct {
	MidiValue  int      `json:"midi_value,omitempty"`
	NameSharp  string   `json:"name_sharp,omitempty"`
	Frequency  float64  `json:"frequency,omitempty"`
	NamesOther []string `json:"names_other,omitempty"`
}

func (n Note) Add(interval int) (result Note) {
	result = Note{MidiValue: n.MidiValue + interval, NameSharp: n.NameSharp}
	for _, d := range noteDB {
		if d.MidiValue == n.MidiValue+interval {
			result = Note{MidiValue: d.MidiValue, NameSharp: d.NameSharp}
			break
		}
	}
	return
}

func findMaxPrefix(a string, b string) string {
	i := 0
	for i < len(a) && i < len(b) {
		if a[i] != b[i] {
			break
		}
		i++
	}
	return a[:i]
}

func exactMatch(n string) (note Note, ok bool) {
	for _, m := range noteDB {
		for _, noteFullName := range append(m.NamesOther, m.NameSharp) {
			if n == noteFullName {
				return Note{MidiValue: m.MidiValue, NameSharp: m.NameSharp}, true
			}
		}
	}
	return
}

func ParseNote(midiString string, midiNear int) (notes []Note, err error) {
	log.Tracef("ParseNote(%s, %d)", midiString, midiNear)
	// can be a single midi note like "c" in which case we need to find the closest note to midiNear
	// or can be a single note like "c4" in which case we want an exact match
	// or can be a sequence of notes like "c4eg" in which case we want to need to split them
	midiString = strings.ToLower(midiString)

	// check if split if it has multiple of any letter [a-g]
	noteStrings := []string{}
	lastAdded := 0
	for i := 1; i < len(midiString); i++ {
		if midiString[i] >= 'a' && midiString[i] <= 'g' {
			noteStrings = append(noteStrings, midiString[lastAdded:i])
			lastAdded = i
		}
	}
	if lastAdded != len(midiString) {
		noteStrings = append(noteStrings, midiString[lastAdded:])
	}
	// log.Debugf("%s' -> %v", midiString, noteStrings)

	// convert '#' to 's'
	for i, n := range noteStrings {
		noteStrings[i] = strings.Replace(n, "#", "s", -1)
	}

	// convert '♭' to 'b'
	for i, n := range noteStrings {
		noteStrings[i] = strings.Replace(n, "♭", "b", -1)
	}

	notes = make([]Note, len(noteStrings))
	for i, n := range noteStrings {
		if note, ok := exactMatch(n); ok {
			log.Tracef("found exact match %s %d", n, note.MidiValue)
			notes[i] = Note{MidiValue: note.MidiValue, NameSharp: note.NameSharp}
			midiNear = note.MidiValue
		} else {
			// find closes to midiNear
			newNote := Note{MidiValue: 300, NameSharp: ""}
			closestDistance := math.Inf(1)
			for _, m := range noteDB {
				for octave := -1; octave <= 8; octave++ {
					for _, noteFullName := range append(m.NamesOther, m.NameSharp) {
						noteName := findMaxPrefix(n, noteFullName)
						if noteName != "" && (noteName == noteFullName || (noteName+strconv.Itoa(octave)) == noteFullName) {
							if math.Abs(float64(m.MidiValue-midiNear)) < closestDistance {
								closestDistance = math.Abs(float64(m.MidiValue - midiNear))
								log.Tracef("found %s %d", noteFullName, m.MidiValue)
								newNote = Note{MidiValue: m.MidiValue, NameSharp: m.NameSharp}
							}
						}
					}
				}
			}
			if newNote.MidiValue != 300 {
				log.Tracef("found %s %d", newNote.NameSharp, newNote.MidiValue)
				notes[i] = newNote
				notes[i].NameSharp = newNote.NameSharp
				midiNear = newNote.MidiValue
			} else {
				err = fmt.Errorf("parsemidi could not parse %s", n)
				return
			}

		}
	}

	log.Tracef("%+v -> %v", midiString, notes)
	return

}
