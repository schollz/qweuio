package music

import (
	"fmt"
	"testing"
)

func TestNote(t *testing.T) {
	n := Note{MidiValue: 60, NameSharp: "C"}
	n2 := n.Add(1)
	if n2.MidiValue != 61 {
		t.Fatalf("expected %d, got %d", 61, n2.MidiValue)
	}
	fmt.Println(n2)

}
