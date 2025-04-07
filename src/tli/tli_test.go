package tli

import (
	"asdfgh/src/player"
	"encoding/json"
	"os"
	"testing"

	log "github.com/schollz/logger"
)

func TestParseTLI(t *testing.T) {
	tli := `
midi virtual 

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
	parsed, err := Parse(tli)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log.Tracef("parsed.Components[0].ChainSteps[0]: %+v", parsed.Components[0].ChainSteps[0])
	player.Play(parsed.Players[0], parsed.Components[0].ChainSteps[0])

	// write parsed to a file out.json
	b, _ := json.MarshalIndent(parsed, "", "  ")
	f, _ := os.Create("out.json")
	defer f.Close()
	f.Write(b)
}
