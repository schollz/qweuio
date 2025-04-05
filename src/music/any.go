package music

func Parse(anything string, midiNear int) (notes []Note, err error) {
	notes, err = ParseChord(anything, midiNear)
	if err == nil {
		return
	}
	return ParseNote(anything, midiNear)
}
