package midi

import (
	"fmt"
	"museq/src/midiconnector"

	log "github.com/schollz/logger"
)

type Player struct {
	Name         string
	nameOriginal string // original name, used for debugging
	Device       *midiconnector.Device
	opened       bool
	channel      uint8
}

func New(name string, channel int) (p *Player, err error) {

	p0 := Player{Name: fmt.Sprintf("midi-%s-%d", name, channel), channel: uint8(channel), nameOriginal: name}
	p0.Device, err = midiconnector.New(name)
	if err != nil {
		log.Errorf("Error opening device: %s", err)
		return
	} else {
		p = &p0
		err = p.Device.Open()
		p.opened = true
		log.Infof("opened device %+v", p.Device)
	}
	return
}

func (m Player) String() string {
	return m.Name
}

func (m *Player) Close() (err error) {
	log.Tracef("close")
	if m.opened {
		err = m.Device.Close()
		if err != nil {
			log.Errorf("Error closing device: %s", err)
		} else {
			log.Infof("closed device %+v", m.Device)
		}
	}
	return
}

func (m *Player) NoteOn(note int, velocity int) (err error) {
	log.Tracef("note_on  (%d,%d,%d)\n", m.channel, note, velocity)
	if m.opened {
		err = m.Device.NoteOn(m.channel, uint8(note), uint8(velocity))
	}
	return
}

func (m *Player) NoteOff(note int) (err error) {
	log.Tracef("note_off (%d,%d)", m.channel, note)
	if m.opened {
		err = m.Device.NoteOff(m.channel, uint8(note))
	} else {
		// Player was closed, but we still need to send note_off to avoid stuck notes
		// Create a temporary device to send the note_off message
		tempDevice, tempErr := midiconnector.New(m.nameOriginal)
		if tempErr == nil {
			if openErr := tempDevice.Open(); openErr == nil {
				err = tempDevice.NoteOff(m.channel, uint8(note))
				tempDevice.Close() // Clean up temporary device
			}
		}
		if tempErr != nil || err != nil {
			log.Warnf("Could not send note_off via temp device: %v", err)
		}
	}
	return
}
