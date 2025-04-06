package step

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"asdfgh/src/constants"
	"asdfgh/src/music"

	log "github.com/schollz/logger"
)

type Notes struct {
	Note []music.Note
}

type Step struct {
	Original    string       `json:"original,omitempty"`
	Iteration   int          `json:"iteration,omitempty"`
	NotesString []string     `json:"notes_string,omitempty"`
	Notes       []music.Note `json:"notes,omitempty"`
	Velocity    []int        `json:"velocity,omitempty"`
	Transpose   []int        `json:"transpose,omitempty"`
	Probability []int        `json:"probability,omitempty"`
	Arpeggio    []string     `json:"arpeggio,omitempty"`
	Gate        []float64    `json:"gate,omitempty"`
	TimeStart   float64      `json:"time_start,omitempty"`
	Duration    float64      `json:"duration,omitempty"`
}

func (s Step) String() string {
	return fmt.Sprintf("%s[%2.2f,%2.2f]", s.Original, s.TimeStart, s.Duration)
}

var regexpParseModifiers = regexp.MustCompile("[" + string(constants.MODIFIERS) + "]")

func (s *Step) Parse() (err error) {
	// uses Orgnal string to parse the step
	// Split the string by capturing the delimiters.
	parts := regexpParseModifiers.Split(s.Original, -1)
	delimiters := regexpParseModifiers.FindAllString(s.Original, -1)
	log.Tracef("parts: %v", parts)
	log.Tracef("delimiters: %v", delimiters)
	for i, part := range parts {
		if i == 0 {
			// First part is the note
			s.NotesString = strings.Split(part, ",")
		} else {
			// Subsequent parts are modifiers
			switch delimiters[i-1] {
			case string(constants.MODIFIER_ARPEGGIO):
				s.Arpeggio = strings.Split(part, ",")
			}
		}
	}
	b, _ := json.Marshal(s)
	log.Tracef("step: %s", b)
	return
}

type Steps struct {
	Step  []Step
	Total float64
}

func (s Steps) Count() int {
	return len(s.Step)
}

func (s Steps) String() string {
	var result string
	for _, step := range s.Step {
		result += step.String() + " "
	}
	result += fmt.Sprintf("total: %2.2f", s.Total)
	return result
}

func (s *Steps) Add(step ...Step) {
	for i := range step {
		s.Step = append(s.Step, step[i])
	}
}

func (s *Steps) CalculateStart() {
	// Calculate start points based on the number of steps in the line
	count := float64(len(s.Step))
	for i := range s.Step {
		s.Step[i].TimeStart = float64(i) / count
	}
}

func (s *Steps) CalculateEnd() {
	// Recalculate start points
	s.Total = -1.0
	for i := range s.Step {
		if s.Step[i].TimeStart == 0 {
			s.Total += 1.0
		}
		s.Step[i].TimeStart += s.Total
	}
	s.Total += 1.0

	// Calculate end/durations
	for i := 0; i < len(s.Step); i++ {
		for j_ := 1; j_ < len(s.Step); j_++ {
			j := (j_ + i) % len(s.Step)
			if string(s.Step[j].Original[0]) == constants.HOLD {
				continue
			}
			startTime := s.Step[i].TimeStart
			endTime := s.Step[j].TimeStart
			if endTime < startTime {
				endTime += s.Total
			}
			s.Step[i].Duration = endTime - startTime
			break
		}
	}
}

func (s *Steps) ClearRests() {
	// Clear rests
	var stepsNew Steps
	for i := range s.Step {
		if string(s.Step[i].Original[0]) == constants.REST ||
			string(s.Step[i].Original[0]) == constants.HOLD {
			continue
		}
		stepsNew.Add(s.Step[i])
	}
	s.Step = stepsNew.Step
}
