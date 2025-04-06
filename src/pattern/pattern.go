package pattern

import (
	"asdfgh/src/step"
	"strings"

	log "github.com/schollz/logger"
)

type Pattern struct {
	Steps step.Steps
}

func Parse(s string) (p Pattern, err error) {
	p = Pattern{}
	for i, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var steps step.Steps
		for _, token := range strings.Fields(line) {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			steps.Add(step.Step{Original: token})
		}
		if steps.Count() == 0 {
			continue
		}
		steps.CalculateStart()
		log.Tracef("line %d, steps: %s", i, steps)
		p.Steps.Add(steps.Step...)
	}
	log.Tracef("steps: %s", p.Steps)
	p.Steps.CalculateEnd()
	log.Tracef("steps: %s", p.Steps)
	p.Steps.ClearRests()
	log.Tracef("steps: %s", p.Steps)
	return
}
