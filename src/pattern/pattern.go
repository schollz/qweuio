package pattern

import (
	"asdfgh/src/expand_line"
	"asdfgh/src/step"
	"strings"

	log "github.com/schollz/logger"
)

type Pattern struct {
	Steps step.Steps
}

func Parse(s string) (p Pattern, err error) {
	p = Pattern{}
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var steps step.Steps
		steps, err = expand_line.ExpandLine(line)
		p.Steps.Add(steps.Step...)
	}
	p.Steps.Expand()
	log.Tracef("Pattern: %s", p.Steps)
	return
}
