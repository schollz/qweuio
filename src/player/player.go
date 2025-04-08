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

type Options struct {
	Channel     int
	Velocity    int
	Gate        float64
	Transpose   float64
	Probability float64
}

func Play(p Player, step step.Step, ops Options) error {
	if len(step.NoteChoices) == 0 {
		return nil
	}
	// Implement the logic for playing a step
	log.Tracef("Playing step: %v", step)
	noteChoiceNum := step.Iteration % len(step.NoteChoices)
	// TODO if arpeggio, use the note list to generate the arpeggio
	for _, note := range step.NoteChoices[noteChoiceNum].NoteList {
		if err := p.NoteOn(ops.Channel, note.MidiValue+int(ops.Transpose), ops.Velocity); err != nil {
			log.Errorf("Error playing note: %s", err)
		}
	}
	go func() {
		// Wait for the duration of the step
		time.Sleep(time.Duration(int(step.Duration*1000000.0*ops.Gate)) * time.Microsecond)
		for _, note := range step.NoteChoices[noteChoiceNum].NoteList {
			if err := p.NoteOff(ops.Channel, note.MidiValue+int(ops.Transpose)); err != nil {
				log.Errorf("Error stopping note: %s", err)
			}
		}
		log.Tracef("Finished playing step: %v", step)
	}()

	return nil
}
