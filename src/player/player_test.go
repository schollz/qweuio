package player

import (
	"testing"

	"asdfgh/src/player/midi"

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
