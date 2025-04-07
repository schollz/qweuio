package player

import (
	"asdfgh/src/step"

	log "github.com/schollz/logger"
)

type Player interface {
	NoteOn(ch int, note int, velocity int) error
	NoteOff(ch int, note int) error
	Close() error
}

func Play(p Player, step step.Step) error {
	// Implement the logic for playing a step
	log.Tracef("Playing step: %v", step)
	return nil
}
