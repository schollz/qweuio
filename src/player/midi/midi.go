package midi

import (
	"asdfgh/src/midiconnector"

	log "github.com/schollz/logger"
)

type Player struct {
	Name   string
	Device *midiconnector.Device
	opened bool
}

func New(name string) (p *Player, err error) {
	p0 := Player{Name: name}
	p0.Device, err = midiconnector.New(name)
	p = &p0
	if err == nil {
		log.Infof("connected to %+v", p.Device)
	}
	err = p.Device.Open()
	if err != nil {
		log.Errorf("Error opening device: %s", err)
	} else {
		p.opened = true
	}
	return
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

func (m *Player) NoteOn(ch int, note int, velocity int) (err error) {
	log.Tracef("note_on  (%d,%d,%d)\n", ch, note, velocity)
	if m.opened {
		err = m.Device.NoteOn(uint8(ch), uint8(note), uint8(velocity))
	}
	return
}

func (m *Player) NoteOff(ch int, note int) (err error) {
	log.Tracef("note_off (%d,%d)", ch, note)
	if m.opened {
		err = m.Device.NoteOff(uint8(ch), uint8(note))
	}
	return
}
