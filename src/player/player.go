package player

import (
	"asdfgh/src/step"
	"time"

	log "github.com/schollz/logger"
)

type Player interface {
	NoteOn(ch int, note int, velocity int) error
	NoteOff(ch int, note int) error
	Close() error
}

func Play(p Player, channel int, step step.Step) error {
	if len(step.NoteChoices) == 0 {
		return nil
	}
	// Implement the logic for playing a step
	log.Tracef("Playing step: %v", step)
	noteChoiceNum := step.Iteration % len(step.NoteChoices)
	// TODO if arpeggio, use the note list to generate the arpeggio
	for _, note := range step.NoteChoices[noteChoiceNum].NoteList {
		if err := p.NoteOn(channel, note.MidiValue, 100); err != nil {
			log.Errorf("Error playing note: %s", err)
		}
	}
	go func() {
		// Wait for the duration of the step
		time.Sleep(time.Duration(int(step.Duration*1000000)) * time.Microsecond)
		for _, note := range step.NoteChoices[noteChoiceNum].NoteList {
			if err := p.NoteOff(channel, note.MidiValue); err != nil {
				log.Errorf("Error stopping note: %s", err)
			}
		}
		log.Tracef("Finished playing step: %v", step)
	}()

	return nil
}
