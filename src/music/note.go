package music

type Note struct {
	MidiValue  int
	NameSharp  string
	Frequency  float64
	NamesOther []string
	IsRest     bool
	IsLegato   bool
}

func (n Note) Add(interval int) (result Note) {
	result = Note{MidiValue: n.MidiValue + interval, NameSharp: n.NameSharp}
	for _, d := range noteDB {
		if d.MidiValue == n.MidiValue+interval {
			result = Note{MidiValue: d.MidiValue, NameSharp: d.NameSharp}
			break
		}
	}
	return
}
