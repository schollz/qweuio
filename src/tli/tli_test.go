package tli

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestParseTLI(t *testing.T) {
	tli := `
midi virtual 
bpm 240

+# [first_part second_part] * 2 second_part

# first_part
Cmaj@u2d4,u3d3

# second_part
d,d5 e
f g a

// velocity 
+! velocity_thing

! velocity_thing
30 30 90

// transposition thing
+$ ta tb

$ ta 
-0.1 1 2 5

$ tb
3 4


`

	tli = `
midi op-z
bpm 120

+# a

# a
Am;3@u11d8
C/G;4@u6d8u4
`
	parsed, err := Parse(tli)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// write parsed to a file out.json
	b, _ := json.MarshalIndent(parsed, "", "  ")
	f, _ := os.Create("out.json")
	defer f.Close()
	f.Write(b)

	// play it
	parsed.Play()
	time.Sleep(12 * time.Second)

	// stop
	parsed.Stop()
	time.Sleep(1 * time.Second)

}
