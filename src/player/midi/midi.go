package midi

import (
	"asdfgh/src/midiconnector"

	log "github.com/schollz/logger"
)

type Player struct {
	Name    string
	Device  *midiconnector.Device
	opened  bool
	channel uint8
}

func New(name string, channel int) (p *Player, err error) {
	p0 := Player{Name: name, channel: uint8(channel)}
	p0.Device, err = midiconnector.New(name)
	if err != nil {
		log.Errorf("Error opening device: %s", err)
		return
	} else {
		p = &p0
		err = p.Device.Open()
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
	}
	return
}
