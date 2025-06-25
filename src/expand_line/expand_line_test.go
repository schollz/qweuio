package expand_line

import "testing"

func TestExpandLine(t *testing.T) {
	tests := []struct {
		line       string
		startTimes []float64
	}{
		{"a*3 b c", []float64{0, 1.0 / 3.0 / 3.0, 2.0 / 3.0 / 3.0, 1 / 3.0, 2.0 / 3.0}}, // "[a a a] b c"
		{"[[a a] a] a", []float64{0, 0.125, 0.25, 0.5}},
		{"[a a] a,a", []float64{0, 0.25, 0.5}},
	}

	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			result, err := ExpandLine(test.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			startTime := []float64{}
			for _, step := range result.Step {
				startTime = append(startTime, step.TimeStart)
			}
			if len(startTime) != len(test.startTimes) {
				t.Fatalf("expected %+v, got %+v", test.startTimes, startTime)
			}
			for i, startTime := range startTime {
				if startTime != test.startTimes[i] {
					t.Fatalf("expected %f, got %f", test.startTimes[i], startTime)
				}
			}
		})
	}
}
