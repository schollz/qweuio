package tli

import (
	"encoding/json"
	"strings"

	"asdfgh/src/constants"
	"asdfgh/src/expand_multiply"
	"asdfgh/src/pattern"
	"asdfgh/src/player"
	"asdfgh/src/player/midi"
	"asdfgh/src/step"

	log "github.com/schollz/logger"
)

type Component struct {
	Type       string                     `json:"type,omitempty"`
	Chain      []string                   `json:"chain,omitempty"`
	Patterns   map[string]pattern.Pattern `json:"patterns,omitempty"`
	ChainSteps []step.Step                `json:"chain_steps,omitempty"`
}

type TLI struct {
	BPM         float64         `json:"bpm,omitempty"`
	Probability int             `json:"probability,omitempty"`
	Velocity    int             `json:"velocity,omitempty"`
	Transpose   float64         `json:"transpose,omitempty"`
	Gate        float64         `json:"gate,omitempty"`
	Components  []Component     `json:"components,omitempty"`
	Players     []player.Player `json:"players,omitempty"`
}

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

	for _, line := range strings.Split(tliString, "\n") {
		fields := strings.Fields(line)
		line = strings.Join(fields, " ")
		if line == "" || line[0] == '/' {
			continue
		}
		if strings.ToLower(fields[0]) == "midi" {
			midiName := strings.Join(fields[1:], " ")
			var p player.Player
			p, err = midi.New(midiName)
			if err != nil {
				log.Errorf("Error creating midi player: %s", err)
				continue
			} else {
				log.Infof("Connected to midi device: %s", p)
				tli.Players = append(tli.Players, p)
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

	log.Tracef("Parsed chains: %v", chains)
	log.Tracef("Parsed patterns: %v", patterns)
	log.Tracef("Parsed ordering: %v", ordering)

	for _, t := range ordering {
		// find the chain
		var chainString string
		var ok bool
		if chainString, ok = chains[t]; !ok {
			log.Warnf("Chain %s not found", t)
			continue
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
					component.ChainSteps = append(component.ChainSteps, step)
				}
				total += pattern.Steps.Total
			}
		}

		tli.Components = append(tli.Components, component)
	}

	return
}
