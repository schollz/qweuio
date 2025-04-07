package tli

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParseTLI(t *testing.T) {
	tli := `
+# [first_part second_part] * 2 second_part

# first_part
Cmaj@u2d4

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
	fmt.Printf("Parsed TLI: %s\n", parsed)
	// write parsed to a file out.json
	b, _ := json.MarshalIndent(parsed, "", "  ")
	f, _ := os.Create("out.json")
	defer f.Close()
	f.Write(b)
}
