package music

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	log "github.com/schollz/logger"
)

// Memoization cache for ParseChord results
type chordCacheKey struct {
	chordString string
	midiNear    int
}

var (
	chordCache = make(map[chordCacheKey][]Note)
	chordCacheMutex sync.RWMutex
)

func ParseChord(chordString string, midiNear int) (result []Note, err error) {
	// Check cache first
	cacheKey := chordCacheKey{chordString: chordString, midiNear: midiNear}
	chordCacheMutex.RLock()
	if cachedResult, exists := chordCache[cacheKey]; exists {
		chordCacheMutex.RUnlock()
		return cachedResult, nil
	}
	chordCacheMutex.RUnlock()
	
	// Compute result
	result, err = parseChordUncached(chordString, midiNear)
	
	// Cache the result if successful
	if err == nil {
		chordCacheMutex.Lock()
		chordCache[cacheKey] = result
		chordCacheMutex.Unlock()
	}
	
	return result, err
}

func parseChordUncached(chordString string, midiNear int) (result []Note, err error) {
	chordStringOriginal := chordString
	chordMatch := ""
	_ = chordMatch

	log.Tracef("chordString: %s", chordStringOriginal)

	octave := 4
	if midiNear != 0 {
		octave = midiNear/12 - 1
	}
	if strings.Contains(chordString, ";") {
		chordSplit := strings.Split(chordString, ";")
		chordString = chordSplit[0]
		if len(chordSplit) > 1 {
			octave, err = strconv.Atoi(chordSplit[1])
			if err != nil {
				return
			}
		}
	}
	log.Tracef("octave: %d", octave)
	log.Tracef("chordString: %s", chordString)

	transposeNote := ""
	if strings.Contains(chordString, "/") {
		chordSplit := strings.Split(chordString, "/")
		chordString = chordSplit[0]
		if len(chordSplit) > 1 {
			transposeNote = strings.ToLower(chordSplit[1])
		}
	}
	log.Tracef("transposeNote: %s", transposeNote)

	// find the root note name - optimized to avoid repeated ToLower calls
	noteMatch := ""
	transposeNoteMatch := ""
	chordRest := ""
	chordStringLower := strings.ToLower(chordString)
	for _, n := range notesAll {
		if transposeNote != "" && len(n) > len(transposeNoteMatch) {
			if n == transposeNote {
				transposeNoteMatch = n
			}
		}
		if len(n) > len(noteMatch) {
			// check if has prefix - original case first
			if strings.HasPrefix(chordString, n) {
				noteMatch = n
				chordRest = chordString[len(n):]
			} else if strings.HasPrefix(chordStringLower, n) {
				noteMatch = n
				chordRest = chordString[len(n):]
			}
		}
	}
	if noteMatch == "" {
		err = fmt.Errorf("no chord found")
	}
	log.Tracef("noteMatch: %s", noteMatch)
	log.Tracef("chordRest: %s", chordRest)

	// convert to canonical sharp scale using lookup map - O(1) access
	// e.g. Fb -> E, Gs -> G#
	if canonical, exists := noteConversionMap[noteMatch]; exists {
		noteMatch = canonical
	}
	if transposeNoteMatch != "" {
		if canonical, exists := noteConversionMap[transposeNoteMatch]; exists {
			transposeNoteMatch = canonical
		}
	}
	log.Tracef("noteMatch: %s", noteMatch)
	log.Tracef("transposeNoteMatch: %s", transposeNoteMatch)

	// find longest matching chord pattern - use lookup map for O(1) access
	chordMatch = "" // (no chord match is major chord)
	chordIntervals := "1P 3M 5P"
	chordRestLower := strings.ToLower(chordRest)
	if intervals, exists := chordPatternMap[chordRestLower]; exists {
		chordMatch = chordRest
		chordIntervals = intervals
	}
	log.Tracef("chordMatch for %s: %s", chordRest, chordMatch)
	log.Tracef("chordIntervals: %s", chordIntervals)

	// find location of root
	rootPosition := 0
	for i, n := range notesScaleSharp {
		if n == noteMatch {
			rootPosition = i
			break
		}
	}
	log.Tracef("rootPosition: %d", rootPosition)

	/** lua code
		-- find notes from intervals
	  whole_note_semitones={0,2,4,5,7,9,11,12}
	  notes_in_chord={}
	  for interval in string.gmatch(chord_intervals,"%S+") do
	    -- get major note position
	    major_note_position=(string.match(interval,"%d+")-1)%7+1
	    -- find semitones from root
	    semitones=whole_note_semitones[major_note_position]
	    -- adjust semitones based on interval
	    if string.match(interval,"m") then
	      semitones=semitones-1
	    elseif string.match(interval,"A") then
	      semitones=semitones+1
	    end
	    if self.debug then
	      print("interval: "..interval)
	      print("major_note_position: "..major_note_position)
	      print("semitones: "..semitones)
	      print("root_position+semitones: "..root_position+semitones)
	    end
	    -- get note in scale from root position
	    note_in_chord=self.notes_scale_sharp[root_position+semitones]
	    table.insert(notes_in_chord,note_in_chord)
	  end
	  **/

	// go code
	// find notes from intervals
	wholeNoteSemitones := []int{0, 2, 4, 5, 7, 9, 11, 12}
	notesInChord := []string{}
	for _, interval := range strings.Fields(chordIntervals) {
		// get major note position
		majorNotePosition, _ := strconv.Atoi(strings.TrimRight(interval, "mMAP"))
		majorNotePosition = ((majorNotePosition - 1) % 7) + 1
		// find semitones from root
		semitones := wholeNoteSemitones[majorNotePosition]
		// adjust semitones based on interval
		if strings.Contains(interval, "A") || strings.Contains(interval, "M") {
			semitones = semitones + 1
		}
		log.Trace("-------------------------")
		log.Tracef("interval: %s", interval)
		log.Tracef("majorNotePosition: %d", majorNotePosition)
		log.Tracef("semitones: %d", semitones)
		log.Tracef("rootPosition+semitones: %d", rootPosition+semitones)
		// get note in scale from root position
		noteInChord := notesScaleSharp[(rootPosition+semitones)-2+12]
		notesInChord = append(notesInChord, noteInChord)
	}
	log.Tracef("notesInChord: %v", notesInChord)
	log.Tracef("notesScaleSharp: %v", notesScaleSharp)

	// if tranposition, rotate until new root
	if transposeNoteMatch != "" {
		foundNote := false
		for i := 0; i < len(notesInChord); i++ {
			if notesInChord[0] == transposeNoteMatch {
				foundNote = true
				break
			}
			notesInChord = append(notesInChord[1:], notesInChord[0])
		}
		if !foundNote {
			notesInChord = append([]string{transposeNoteMatch}, notesInChord...)
		} else {
			log.Tracef("transposeNoteMatch: %s", transposeNoteMatch)
		}
	}
	log.Tracef("notesInChord: %v", notesInChord)

	// go code
	// convert to midi - use lookup map for O(1) note access
	midiNotesInChord := []int{}
	lastNote := 0
	
	// Use lookup map for faster note finding
	for i, n := range notesInChord {
		// Try different octaves in order
		for oct := octave; oct <= octave+3; oct++ {
			targetName := n + strconv.Itoa(oct)
			if note, exists := noteByNameMap[targetName]; exists && note.MidiValue > lastNote {
				lastNote = note.MidiValue
				midiNotesInChord = append(midiNotesInChord, note.MidiValue)
				notesInChord[i] = note.NameSharp
				break
			}
		}
	}
	log.Tracef("midiNotesInChord: %v", midiNotesInChord)
	log.Tracef("notesInChord: %v", notesInChord)

	result = make([]Note, len(midiNotesInChord))
	for i, m := range midiNotesInChord {
		result[i] = Note{MidiValue: m, NameSharp: notesInChord[i]}
	}
	log.Tracef("result: %v", result)

	return
}
