package step

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"asdfgh/src/constants"
	"asdfgh/src/music"

	log "github.com/schollz/logger"
)

type Notes struct {
	Note []music.Note
}

type Step struct {
	Original    string        `json:"original,omitempty"`
	Iteration   int           `json:"iteration,omitempty"`
	Notes       []music.Notes `json:"notes,omitempty"`
	Velocity    []int         `json:"velocity,omitempty"`
	Transpose   []int         `json:"transpose,omitempty"`
	Probability []int         `json:"probability,omitempty"`
	Arpeggio    []string      `json:"arpeggio,omitempty"`
	Gate        []float64     `json:"gate,omitempty"`
	TimeStart   float64       `json:"time_start,omitempty"`
	Duration    float64       `json:"duration,omitempty"`
}

func (s Step) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

var regexpParseModifiers = regexp.MustCompile("[" + string(constants.MODIFIERS) + "]")

func splitStringToInts(s string) (result []int) {
	result = make([]int, 0)
	for _, v := range strings.Split(s, ",") {
		if v == "" {
			continue
		}
		if i, err := strconv.Atoi(v); err == nil {
			result = append(result, i)
		} else {
			log.Errorf("Error parsing string to int: %s", err)
		}
	}
	return
}

func splitStringToFloats(s string) (result []float64) {
	result = make([]float64, 0)
	for _, v := range strings.Split(s, ",") {
		if v == "" {
			continue
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			result = append(result, f)
		} else {
			log.Errorf("Error parsing string to float: %s", err)
		}
	}
	return
}

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
			s.Notes = make([]music.Notes, 0)
			lastMidi := 60
			for _, noteString := range strings.Split(part, ",") {
				if noteString == "" {
					continue
				}
				noteObj, err := music.Parse(noteString, lastMidi)
				if err != nil {
					log.Errorf("Error parsing note: %s", err)
					continue
				}
				s.Notes = append(s.Notes, music.Notes{
					Original: noteString,
					Note:     noteObj,
				})
				lastMidi = noteObj[len(noteObj)-1].MidiValue
			}
		} else {
			// Subsequent parts are modifiers
			switch delimiters[i-1] {
			case string(constants.MODIFIER_ARPEGGIO):
				s.Arpeggio = strings.Split(part, ",")
			case string(constants.MODIFIER_VELOCITY):
				s.Velocity = splitStringToInts(part)
			case string(constants.MODIFIER_TRANSPOSE):
				s.Transpose = splitStringToInts(part)
			case string(constants.MODIFIER_PROBABILITY):
				s.Probability = splitStringToInts(part)
			case string(constants.MODIFIER_GATE):
				s.Gate = splitStringToFloats(part)
			default:
				log.Errorf("Unknown modifier: %s", delimiters[i-1])
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
		result += step.String() + "\n"
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

func (s *Steps) Parse() {
	// Parse each step
	for i := range s.Step {
		if err := s.Step[i].Parse(); err != nil {
			log.Errorf("Error parsing step: %s", err)
		}
	}
}

func (s *Steps) Expand() {
	s.CalculateEnd()
	s.ClearRests()
	s.Parse()
}
