package music

func Parse(anything string, midiNear int) (notes []Note, err error) {
	// check if uppercase
	if anything[0] >= 'A' && anything[0] <= 'Z' {
		return ParseChord(anything, midiNear)
	}
	return ParseNote(anything, midiNear)
}
