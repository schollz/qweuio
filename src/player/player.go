package player

type Player interface {
	NoteOn(ch int, note int, velocity int) error
	NoteOff(ch int, note int) error
	Close() error
}
