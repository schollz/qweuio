package tli

import (
	"museq/src/pattern"
	"museq/src/player"
	"museq/src/step"
)

// Copy copies the contents of t2 into t1
func (t1 *TLI) Copy(t2 TLI) (err error) {
	t1.mutex.Lock()
	defer t1.mutex.Unlock()

	t1.BPM = t2.BPM
	t1.Gate = t2.Gate
	t1.Velocity = t2.Velocity
	t1.Probability = t2.Probability
	t1.Transpose = t2.Transpose
	t1.Scale = t2.Scale
	t1.ScaleRoot = t2.ScaleRoot

	t1.Components = make([]Component, len(t2.Components))
	for i, c := range t2.Components {
		t1.Components[i] = c
		t1.Components[i].Patterns = make(map[string]pattern.Pattern)
		for k, v := range c.Patterns {
			t1.Components[i].Patterns[k] = v
		}
		t1.Components[i].ChainSteps = make([]step.Step, len(c.ChainSteps))
		for j, s := range c.ChainSteps {
			t1.Components[i].ChainSteps[j] = s
		}
		t1.Components[i].ChainDuration = c.ChainDuration
		t1.Components[i].Type = c.Type
		t1.Components[i].Chain = make([]string, len(c.Chain))
		for j, s := range c.Chain {
			t1.Components[i].Chain[j] = s
		}
	}

	t1.Players = make([]player.Player, len(t2.Players))
	for i, p := range t2.Players {
		t1.Players[i] = p
	}

	return
}
