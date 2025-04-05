package modifier

import (
	"asdfgh/src/constants"
	"regexp"
	"strings"
)

type ModifiedString struct {
	Original   string
	Unmodified string
	Modifiers  []string
}

func Split(s string) (result ModifiedString) {
	result.Original = s
	result.Unmodified = ""

	result.Modifiers = make([]string, 0)
	for i, part := range splitWithDelimiters(s) {
		if i == 0 && !strings.Contains(string(constants.MODIFIERS), part) {
			result.Unmodified = part
		} else {
			result.Modifiers = append(result.Modifiers, part)
		}
	}

	return
}

func Merge(s string, modifiers []string) (result string) {
	return s + strings.Join(modifiers, "")
}

func Remove(s string, modifier string) (result string) {
	sSplit := splitWithDelimiters(s)
	result = ""
	for i, part := range sSplit {
		if i > 0 && strings.Contains(part, modifier) {
			continue
		}
		result += part
	}
	return
}

func Pop(s string, modifier string) (resultWihoutModifier string, modifierPart string) {
	sSplit := splitWithDelimiters(s)
	resultWihoutModifier = ""
	for i, part := range sSplit {
		if i > 0 && strings.Contains(part, modifier) {
			modifierPart = part[1:]
			continue
		}
		resultWihoutModifier += part
	}
	return
}

// splitWithDelimiters splits the input string on the given delimiters,
// keeping the delimiters at the beginning of each resulting part.
func splitWithDelimiters(input string) []string {
	// Define the regex pattern to match delimiters: !, @, #, $
	pattern := regexp.MustCompile("[" + string(constants.MODIFIERS) + "]")
	// Split the string by capturing the delimiters.
	parts := pattern.Split(input, -1)
	delimiters := pattern.FindAllString(input, -1)

	// Reconstruct the parts, adding the delimiters to the front of the next part.
	var result []string
	for i, part := range parts {
		if i == 0 && part != "" {
			// Add the initial part if it doesn't start with a delimiter.
			result = append(result, part)
		} else if i > 0 {
			// Always prepend the delimiter to the part.
			result = append(result, delimiters[i-1]+part)
		}
	}

	return result
}
