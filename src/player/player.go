package player

import (
	"fmt"
	"museq/src/music"
	"museq/src/step"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/schollz/logger"
)

const ARPEGGIO_UP = "u"
const ARPEGGIO_DOWN = "d"

const (
	Reset   = "\033[0m"
	Cyan    = "\033[36m"
	Magenta = "\033[35m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	Red     = "\033[31m"
)

type Player interface {
	NoteOn(note int, velocity int) error
	NoteOff(note int) error
	Close() error
}

type Options struct {
	Velocity    int
	Gate        float64
	Transpose   float64
	Probability float64
	Scale       string
	ScaleRoot   string
}

func Play(p Player, step step.Step, ops Options) (err error) {
	if len(step.NoteChoices) == 0 {
		return nil
	}
	// Implement the logic for playing a step
	log.Tracef("Playing step: %v with ops %v", step, ops)
	noteChoiceNum := step.Iteration % len(step.NoteChoices)

	noteList := step.NoteChoices[noteChoiceNum].NoteList
	if len(step.Arpeggio) == 0 {
		notesPlayed := make([]string, len(noteList))
		for i, note := range noteList {
			finalMidiValue := note.MidiValue + int(ops.Transpose)

			// Apply scale quantization if scale is specified
			if ops.Scale != "" {
				quantizedNote, err := music.QuantizeToScale(finalMidiValue, ops.Scale, ops.ScaleRoot)
				if err != nil {
					log.Warnf("Error quantizing note to scale: %s", err)
				} else {
					finalMidiValue = quantizedNote
				}
			}

			resultNote := note.Add(finalMidiValue - note.MidiValue)
			notesPlayed[i] = resultNote.NameSharp
			if err := p.NoteOn(finalMidiValue, ops.Velocity); err != nil {
				log.Errorf("Error playing note: %s", err)
			}
		}
		fmt.Printf("%s%-8.3f%s%s%-18s%s%s%-6s%s%s%-12s%s%s%-12s%s\n",
			Cyan, step.TimeStart, Reset,
			Magenta, p, Reset,
			Red, step.PatternName, Reset,
			Yellow, step.Original, Reset,
			Green, strings.Join(notesPlayed, ", "), Reset)

		go func() {
			// Wait for the duration of the step
			time.Sleep(time.Duration(int(step.Duration*1000000.0*ops.Gate)) * time.Microsecond)
			for _, note := range noteList {
				finalMidiValue := note.MidiValue + int(ops.Transpose)

				// Apply scale quantization if scale is specified and transpose is non-zero
				if ops.Scale != "" && ops.Transpose != 0 {
					quantizedNote, err := music.QuantizeToScale(finalMidiValue, ops.Scale, ops.ScaleRoot)
					if err != nil {
						log.Warnf("Error quantizing note to scale: %s", err)
					} else {
						finalMidiValue = quantizedNote
					}
				}

				if err := p.NoteOff(finalMidiValue); err != nil {
					log.Errorf("Error stopping note: %s", err)
				}
			}
			log.Tracef("Finished playing step: %v", step)
		}()
	} else {
		// select an arpeggio
		arpeggioString := step.Arpeggio[step.Iteration%len(step.Arpeggio)]
		log.Tracef("Playing arpeggio: %s", arpeggioString)
		var noteListArpeggio []music.Note
		notei := 0
		octave := 0

		re := regexp.MustCompile(`\d+`)
		for i, token := range tokenizeLetterNumberes(arpeggioString) {
			// get the number using regex
			match := re.FindString(token)
			if match == "" {
				err = fmt.Errorf("invalid arpeggio token: %s", token)
				return
			}
			var steps int
			steps, err = strconv.Atoi(match)
			if err != nil {
				return
			}
			log.Tracef("steps: %d, notei: %d", steps, notei)
			direction := 0
			if string(token[0]) == ARPEGGIO_UP {
				direction = 1
			} else if string(token[0]) == ARPEGGIO_DOWN {
				direction = -1
			} else {
				err = fmt.Errorf("invalid arpeggio direction: %s", token)
				return
			}

			log.Tracef("direction: %d", direction)
			for j := 0; j < steps; j++ {
				if !(i == 0 && j == 0) {
					notei += direction
				}
				for notei < 0 {
					notei += len(noteList)
					octave--
				}
				for notei >= len(noteList) {
					notei -= len(noteList)
					octave++
				}
				newNote := noteList[notei].Add(octave * 12)
				log.Tracef("notei: %d (%s)", notei, newNote.NameSharp)
				noteListArpeggio = append(noteListArpeggio, newNote)
			}
		}
		log.Tracef("noteListArpeggio: %v", noteListArpeggio)
		durationPerNote := step.Duration / float64(len(noteListArpeggio))
		for i, note := range noteListArpeggio {

			finalMidiValue := note.MidiValue + int(ops.Transpose)

			// Apply scale quantization if scale is specified
			if ops.Scale != "" {
				quantizedNote, err := music.QuantizeToScale(finalMidiValue, ops.Scale, ops.ScaleRoot)
				if err != nil {
					log.Warnf("Error quantizing note to scale: %s", err)
				} else {
					finalMidiValue = quantizedNote
				}
			}

			resultNote := note.Add(finalMidiValue - note.MidiValue)

			fmt.Printf("%s%-8.3f%s%s%-18s%s%s%-6s%s%s%-12s%s%s%-12s%s\n",
				Cyan, step.TimeStart+float64(i)*durationPerNote, Reset,
				Magenta, p, Reset,
				Red, step.PatternName, Reset,
				Yellow, step.Original, Reset,
				Green, resultNote.NameSharp, Reset)

			if err = p.NoteOn(finalMidiValue, ops.Velocity); err != nil {
				log.Errorf("Error playing note: %s", err)
			}
			go func(finalMidi int) {
				time.Sleep(time.Duration(int(durationPerNote*ops.Gate*1000000.0)) * time.Microsecond)
				if err = p.NoteOff(finalMidi); err != nil {
					log.Errorf("Error stopping note: %s", err)
				}
			}(finalMidiValue)
			time.Sleep(time.Duration(int(durationPerNote*1000000.0)) * time.Microsecond)
		}
	}

	return nil
}

func tokenizeLetterNumberes(s string) (tokens []string) {
	// take a string with letters and numbers and split it into tokens
	// where the letter determines the start of a new token
	for i := 0; i < len(s); i++ {
		if i == 0 {
			tokens = append(tokens, string(s[i]))
			continue
		}
		if s[i] >= '0' && s[i] <= '9' {
			tokens[len(tokens)-1] += string(s[i])
		} else {
			tokens = append(tokens, string(s[i]))
		}
	}
	return
}
