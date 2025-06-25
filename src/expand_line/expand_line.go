package expand_line

import (
	"museq/src/expand_multiply"
	"museq/src/step"
	"strings"
	"unicode"

	log "github.com/schollz/logger"
)

func ExpandLine(s string) (steps step.Steps, err error) {
	steps = step.Steps{}
	// First expand multiplication
	sExpanded := expand_multiply.ExpandMultiplication(s)
	log.Tracef("%s -> %s", s, sExpanded)

	// Now parse and distribute
	result := ParseAndDistribute(sExpanded)

	startTime := 0.0
	for _, tokenWeight := range result {
		log.Tracef("Token: %s, Weight: %f", tokenWeight.Value, tokenWeight.Weight)
		// Create a new step for each token
		newStep := step.Step{
			TimeStart: startTime,
			Original:  tokenWeight.Value,
		}
		// Add the new step to the steps
		steps.Add(newStep)
		startTime += tokenWeight.Weight
	}

	log.Tracef("Parsed steps: %v", steps)

	return
}

type Node interface{}

type Token struct {
	Value string
}

type Group struct {
	Children []Node
}

type TokenWeight struct {
	Value  string
	Weight float64
}

func ParseAndDistribute(input string) []TokenWeight {
	tokens := tokenize(input)
	parsed, _ := parseGroup(tokens, 0)
	var result []TokenWeight
	evaluate(parsed, 1.0, &result)
	return result
}

// Tokenize splits input into bracket-aware tokens
func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	for _, ch := range input {
		switch {
		case unicode.IsSpace(ch):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case ch == '[' || ch == ']':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

// parseGroup recursively builds nested token structures
func parseGroup(tokens []string, index int) (Group, int) {
	var group Group
	for index < len(tokens) {
		token := tokens[index]
		switch token {
		case "[":
			subgroup, newIndex := parseGroup(tokens, index+1)
			group.Children = append(group.Children, subgroup)
			index = newIndex
		case "]":
			return group, index + 1
		default:
			group.Children = append(group.Children, Token{Value: token})
			index++
		}
	}
	return group, index
}

// evaluate recursively assigns weights to each token instance
func evaluate(group Group, weight float64, result *[]TokenWeight) {
	count := float64(len(group.Children))
	if count == 0 {
		return
	}
	subWeight := weight / count
	for _, node := range group.Children {
		switch v := node.(type) {
		case Token:
			*result = append(*result, TokenWeight{Value: v.Value, Weight: subWeight})
		case Group:
			evaluate(v, subWeight, result)
		}
	}
}

// func main() {
// 	example := "[[a b] d] c d [d d [d d] d]"
// 	result := ParseAndDistribute(example)
// 	for _, tw := range result {
// 		fmt.Printf("%s: %.3f\n", tw.Value, tw.Weight)
// 	}
// }
