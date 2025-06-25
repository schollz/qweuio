package tli

import (
	"encoding/json"
	"strconv"
	"strings"

	"museq/src/constants"
	"museq/src/expand_multiply"
	"museq/src/pattern"
	"museq/src/player"
	"museq/src/player/midi"
	"museq/src/player/supercollider"

	log "github.com/schollz/logger"
)

func (t TLI) String() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

func Parse(tliString string) (tli TLI, err error) {
	chains := make(map[string]string)
	patterns := make(map[string]string)
	inPattern := false
	currentPattern := ""
	currentPatternName := ""
	ordering := []string{}
	orderingHas := make(map[string]bool)
	type MidiParsed struct {
		Name    string
		Channel int
	}

	midiParsed := MidiParsed{Channel: 1}

	// make sure channel is non-blocking
	tli.stopChan = make(chan bool, 1)

	// reset scale settings to ensure clean re-parsing
	tli.Scale = ""
	tli.ScaleRoot = ""

	for _, line := range strings.Split(tliString, "\n") {
		fields := strings.Fields(line)
		line = strings.Join(fields, " ")
		if line == "" || line[0] == '/' {
			continue
		}
		if strings.ToLower(fields[0]) == "midi" {
			midiName := fields[1]
			if midiParsed.Name != "" {
				var p player.Player
				p, err = midi.New(midiName, midiParsed.Channel)
				if err != nil {
					log.Warnf("Error creating midi player: %s", err)
					continue
				} else {
					log.Debugf("connected: %+v", p)
					tli.Players = append(tli.Players, p)
				}
			}
			midiParsed.Name = midiName
		} else if strings.ToLower(fields[0]) == "channel" {
			if len(fields) > 1 {
				if midiParsed.Channel, err = strconv.Atoi(fields[1]); err != nil {
					log.Warnf("Error parsing channel: %s", err)
				}
			} else {
				log.Warnf("No channel value provided")
			}
		} else if strings.ToLower(fields[0]) == "supercollider" {
			player, playerErr := supercollider.Parse(line)
			if playerErr != nil {
				log.Warnf("Error parsing SuperCollider player: %s", playerErr)
				continue
			} else {
				log.Debugf("connected: %+v", player)
				tli.Players = append(tli.Players, player)
			}
		} else if strings.ToLower(fields[0]) == "transpose" {
			if len(fields) > 1 {
				if tli.Transpose, err = strconv.ParseFloat(fields[1], 64); err != nil {
					log.Warnf("Error parsing transpose: %s", err)
				}
			} else {
				log.Warnf("No transpose value provided")
			}
		} else if strings.ToLower(fields[0]) == "bpm" {
			if len(fields) > 1 {
				if tli.BPM, err = strconv.ParseFloat(fields[1], 64); err != nil {
					log.Warnf("Error parsing BPM: %s", err)
				}
			} else {
				log.Warnf("No BPM value provided")
			}
		} else if strings.ToLower(fields[0]) == "probability" {
			if len(fields) > 1 {
				if tli.Probability, err = strconv.Atoi(fields[1]); err != nil {
					log.Warnf("Error parsing probability: %s", err)
				}
			} else {
				log.Warnf("No probability value provided")
			}
		} else if strings.ToLower(fields[0]) == "velocity" {
			if len(fields) > 1 {
				if tli.Velocity, err = strconv.Atoi(fields[1]); err != nil {
					log.Errorf("Error parsing velocity: %s", err)
				}
			} else {
				log.Warnf("No velocity value provided")
			}
		} else if strings.ToLower(fields[0]) == "gate" {
			if len(fields) > 1 {
				if tli.Gate, err = strconv.ParseFloat(fields[1], 64); err != nil {
					log.Errorf("Error parsing gate: %s", err)
				}
			} else {
				log.Warnf("No gate value provided")
			}
		} else if strings.ToLower(fields[0]) == "scale" {
			if len(fields) > 1 {
				tli.Scale = fields[1]
				if len(fields) > 2 {
					tli.ScaleRoot = fields[2]
				} else {
					tli.ScaleRoot = "c"  // default to C if no root specified
				}
			} else {
				log.Warnf("No scale value provided")
			}
		} else if string(line[0]) == constants.SYMBOL_CHAIN {
			chains[string(line[1])] = line[1:]
			if inPattern {
				patterns[currentPatternName] = strings.TrimSpace(currentPattern)
			}
			inPattern = false
		} else if strings.Contains(constants.MODIFIERS, string(line[0])) {
			if _, ok := orderingHas[string(line[0])]; !ok {
				orderingHas[string(line[0])] = true
				ordering = append(ordering, string(line[0]))
			}
			if inPattern {
				patterns[currentPatternName] = strings.TrimSpace(currentPattern)
			}
			inPattern = true
			currentPattern = ""
			currentPatternName = line
		} else if inPattern {
			currentPattern += line + "\n"
		}
	}
	if inPattern {
		patterns[currentPatternName] = strings.TrimSpace(currentPattern)
	}
	if midiParsed.Name != "" {
		var p player.Player
		p, err = midi.New(midiParsed.Name, midiParsed.Channel)
		if err != nil {
			log.Warnf("Error creating midi player: %s", err)
			return
		} else {
			log.Debugf("connected: %+v", p)
			tli.Players = append(tli.Players, p)
		}
	}

	log.Tracef("Parsed chains: %v", chains)
	log.Tracef("Parsed patterns: %v", patterns)
	log.Tracef("Parsed ordering: %v", ordering)

	for _, t := range ordering {
		// find the chain
		var chainString string
		var ok bool
		if chainString, ok = chains[t]; !ok {
			// use the first pattern string
			// if no chain is found
			if len(patterns) > 0 {
				for k := range patterns {
					chainString = k
					break
				}
				chains[t] = chainString
			} else {
				log.Warnf("No chain found for %s", t)
				continue
			}
		}
		// find all the patterns for this chain
		chainString = strings.TrimSpace(chainString[1:])
		// expand the chain string
		chainString = expand_multiply.ExpandMultiplication(chainString, true)
		component := Component{
			Type:  t,
			Chain: strings.Fields(chainString),
		}

		// find all the patterns for this chain
		component.Patterns = make(map[string]pattern.Pattern)
		for patternName, patternString := range patterns {
			if strings.HasPrefix(patternName, t) {
				var parsedPattern pattern.Pattern
				parsedPattern, err = pattern.Parse(patternName + "\n" + patternString)
				if err != nil {
					log.Warnf("Error parsing pattern: %s", err)
				} else {
					component.Patterns[parsedPattern.Name] = parsedPattern
				}
			}
		}

		// validate that all elements in the chain have patterns
		isValid := true
		for _, chainElement := range component.Chain {
			if _, ok := component.Patterns[chainElement]; !ok {
				log.Warnf("Pattern %s not found for chain element %s", chainElement, t)
				isValid = false
				break
			}
		}

		if !isValid {
			log.Warnf("Chain %s is not valid", t)
			continue
		}

		// render chains
		total := 0.0
		for _, chainElement := range component.Chain {
			if pattern, ok := component.Patterns[chainElement]; ok {
				for _, step := range pattern.Steps.Step {
					step.TimeStart += total
					step.PatternName = pattern.Name
					component.ChainSteps = append(component.ChainSteps, step)
				}
				total += pattern.Steps.Total
			}
		}
		component.ChainDuration = total
		tli.Components = append(tli.Components, component)
	}

	if tli.BPM <= 0.0 {
		tli.BPM = 120.0
	}

	// multiple all the duration and start times by the BPM
	// with 4 beats per measure
	for i := range tli.Components {
		tli.Components[i].ChainDuration *= 60.0 / tli.BPM * 4.0
		for j := range tli.Components[i].ChainSteps {
			tli.Components[i].ChainSteps[j].Duration *= 60.0 / tli.BPM * 4.0
			tli.Components[i].ChainSteps[j].TimeStart *= 60.0 / tli.BPM * 4.0
		}
	}

	if tli.Velocity <= 0 {
		tli.Velocity = 100
	}
	if tli.Probability <= 0 {
		tli.Probability = 100
	}
	if tli.Gate <= 0.0 {
		tli.Gate = 0.5
	}
	return
}

func (tli TLI) IsPlaying() bool {
	return tli.isPlaying
}
