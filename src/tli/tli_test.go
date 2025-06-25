package tli

import (
	"encoding/json"
	"os"
	"runtime/pprof"
	"testing"
)

var tli1 = `
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

func BenchmarkParseTLI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(tli1)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkCopyTLI(b *testing.B) {
	t1, err := Parse(tli1)
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
	t2, err := Parse(tli1)
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < b.N; i++ {
		err = t1.Copy(t2)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestParseTLI(t *testing.T) {

	// 	tli = `
	// midi op-z
	// bpm 480
	// gate 0.1
	// velocity 30

	// +# a

	// # a
	// a4,e5 c e g4,g6,f7,d2
	// c4,d4,c4 e g b,a

	// `
	parsed, err := Parse(tli1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// write parsed to a file out.json
	b, _ := json.MarshalIndent(parsed, "", "  ")
	f, _ := os.Create("out.json")
	defer f.Close()
	f.Write(b)

	// // play it
	// parsed.Play()
	// time.Sleep(12 * time.Second)

	// // stop
	// parsed.Stop()
	// time.Sleep(1 * time.Second)

}

func TestParseProfile(t *testing.T) {
	// go tool pprof -http=":8080" cpu.prof
	f, err := os.Create("cpu.prof")
	if err != nil {
		t.Fatalf("could not create CPU profile: %v", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		t.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	// Run the function enough times to get meaningful data
	for i := 0; i < 50000; i++ {
		Parse(tli1)
	}
}
