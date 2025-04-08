package player

import (
	"testing"
	"time"

	"asdfgh/src/constants"
	"asdfgh/src/player/midi"
	"asdfgh/src/step"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestPlayer(t *testing.T) {
	// sudo modprobe snd-virmidi midi_devs=1
	var p Player
	var err error

	p, err = midi.New("Raw MIDI")
	assert.Nil(t, err)
	assert.Nil(t, p.NoteOn(1, 60, 100))
	assert.Nil(t, p.NoteOff(1, 60))
	assert.Nil(t, p.NoteOn(1, 62, 100))
	assert.Nil(t, p.NoteOn(1, 64, 100))
	assert.Nil(t, p.Close())
}

func TestPlaying(t *testing.T) {
	var p Player
	var err error

	p, err = midi.New("OP-Z")
	if err != nil {
		return
	}

	stepString := "Cmaj@u2d4!120,30"
	step := step.Step{Original: stepString}
	_, err = step.Parse(constants.MODIFIER_NOTE, 60)
	step.Duration = 0.5
	log.Tracef("Parsed step: %+v", step)
	err = Play(p, step, Options{
		Channel:     5,
		Velocity:    70,
		Transpose:   1,
		Probability: 0.5,
		Gate:        0.5,
	})
	assert.Nil(t, err)

	time.Sleep(3 * time.Second)
	assert.Nil(t, p.Close())
}
