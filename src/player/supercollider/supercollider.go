package supercollider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hypebeast/go-osc/osc"
	log "github.com/schollz/logger"
)

type Player struct {
	Name         string
	client       *osc.Client
	opened       bool
	oscPort      int
	functionPath string
}

func Parse(line string) (p *Player, err error) {
	fields := make([]string, 0)
	for _, field := range strings.Fields(line) {
		if field != "" {
			fields = append(fields, field)
		}
	}

	if len(fields) < 2 {
		err = fmt.Errorf("invalid SuperCollider player line: %s", line)
		return
	}

	name := fields[0]
	oscPort := 57120
	if len(fields) > 2 {
		oscPort, err = strconv.Atoi(fields[2])
		if err != nil {
			log.Errorf("Error parsing OSC port: %s", err)
			return
		}
	}

	functionPath := "/note"
	if len(fields) > 1 {
		functionPath = fields[1]
	}

	log.Debugf("Parsed SuperCollider player: name=%s, oscPort=%d, functionPath=%s", name, oscPort, functionPath)
	p, err = New(name, oscPort, functionPath)
	return
}

func New(name string, oscPort int, functionPath string) (p *Player, err error) {
	if oscPort == 0 {
		oscPort = 57120
	}
	if functionPath == "" {
		functionPath = "/note"
	}

	p0 := Player{
		Name:         fmt.Sprintf("supercollider-%s-%d", name, oscPort),
		oscPort:      oscPort,
		functionPath: functionPath,
	}

	p0.client = osc.NewClient("localhost", oscPort)
	if p0.client == nil {
		err = fmt.Errorf("failed to create OSC client for port %d", oscPort)
		log.Errorf("Error creating OSC client: %s", err)
		return
	}

	p = &p0
	p.opened = true
	log.Infof("opened SuperCollider OSC client on port %d with function path %s", oscPort, functionPath)

	return
}

func (sc Player) String() string {
	return sc.Name
}

func (sc *Player) Close() (err error) {
	log.Tracef("close")
	if sc.opened {
		sc.opened = false
		log.Infof("closed SuperCollider OSC client %s", sc.Name)
	}
	return
}

func (sc *Player) NoteOn(note int, velocity int) (err error) {
	log.Tracef("[%s] note_on (%d,%d)\n", sc.functionPath, note, velocity)
	if sc.opened {
		msg := osc.NewMessage(sc.functionPath + "/noteOn")
		msg.Append(int32(note))
		msg.Append(int32(velocity))
		err = sc.client.Send(msg)
		if err != nil {
			log.Errorf("Error sending OSC note_on message: %s", err)
		}
	}
	return
}

func (sc *Player) NoteOff(note int) (err error) {
	log.Tracef("[%s] note_off (%d)", sc.functionPath, note)
	if sc.opened {
		msg := osc.NewMessage(sc.functionPath + "/noteOff")
		msg.Append(int32(note))
		err = sc.client.Send(msg)
		if err != nil {
			log.Errorf("Error sending OSC note_off message: %s", err)
		}
	}
	return
}
