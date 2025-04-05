package modifier

import "testing"

func TestRemove(t *testing.T) {
	tests := []struct {
		input    string
		remove   string
		expected string
	}{
		{
			input:    "hello!@#$world",
			remove:   "!",
			expected: "hello@#$world",
		},
		{
			input:    "hello@hlkj!ok#$world",
			remove:   "!",
			expected: "hello@hlkj#$world",
		},
	}

	for _, test := range tests {
		result := Remove(test.input, test.remove)
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		input    string
		expected ModifiedString
	}{
		{
			input: "!@#$",
			expected: ModifiedString{
				Original:   "!@#$",
				Unmodified: "",
				Modifiers:  []string{"!", "@", "#", "$"},
			},
		},
		{
			input: "hello!@#$world",
			expected: ModifiedString{
				Original:   "hello!@#$world",
				Unmodified: "hello",
				Modifiers:  []string{"!", "@", "#", "$world"},
			},
		},
	}

	for _, test := range tests {
		result := Split(test.input)
		if result.Original != test.expected.Original {
			t.Errorf("expected Original %s, got %s", test.expected.Original, result.Original)
		}
		if result.Unmodified != test.expected.Unmodified {
			t.Errorf("expected Unmodified %s, got %s", test.expected.Unmodified, result.Unmodified)
		}
		if len(result.Modifiers) != len(test.expected.Modifiers) {
			t.Errorf("expected %d modifiers, got %d", len(test.expected.Modifiers), len(result.Modifiers))
		}
		for i, modifier := range result.Modifiers {
			if modifier != test.expected.Modifiers[i] {
				t.Errorf("expected modifier %s, got %s", test.expected.Modifiers[i], modifier)
			}
		}
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		input    string
		pop      string
		newInput string
		expected string
	}{
		{
			input:    "hello!ok#$world",
			pop:      "!",
			newInput: "hello#$world",
			expected: "ok",
		},
	}

	for _, test := range tests {
		newInput, modifier := Pop(test.input, test.pop)
		if newInput != test.newInput {
			t.Errorf("expected %s, got %s", test.newInput, newInput)
		}
		if modifier != test.expected {
			t.Errorf("expected %s, got %s", test.expected, modifier)
		}
	}
}
