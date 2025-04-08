package tli

import (
	"asdfgh/src/constants"
	"asdfgh/src/player"
	"time"

	log "github.com/schollz/logger"
)

func (t *TLI) Stop() {
	if t.isPlaying {
		t.stopChan <- true
	}
}

func (t *TLI) Play() (err error) {
	if t.isPlaying {
		return
	}
	t.isPlaying = true
	t.lastSeconds = -1

	// create a millisecond ticker
	ticker := time.NewTicker(1 * time.Millisecond)
	t.startTime = time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				t.Emit()

			case <-t.stopChan:
				log.Tracef("Stopping playback")
				t.isPlaying = false
				ticker.Stop()
				// stop all players
				for _, p := range t.Players {
					if err := p.Close(); err != nil {
						log.Errorf("Error closing player: %s", err)
					}
				}
			}
		}
	}()

	return
}

func (t *TLI) Emit() (err error) {
	ct := time.Since(t.startTime).Seconds()

	// ctCurrent := ct
	// ctLast := t.lastSeconds
	// loop through all components
	for i, component := range t.Components {
		cts := [2]float64{t.lastSeconds, ct}
		for cts[1] > component.ChainDuration {
			cts[0] -= component.ChainDuration
			cts[1] -= component.ChainDuration
		}
		// log.Tracef("%s [%2.3f,%2.3f]", component.Type, cts[0], cts[1])
		// find any step that is between ctCurrent and ctLast
		for j, step := range component.ChainSteps {
			if step.TimeStart <= cts[1] && step.TimeStart > cts[0] {
				log.Tracef("[%2.3f] [%s] %s (%d)", step.TimeStart, component.Type, step.Original, step.Iteration)
				if component.Type == constants.MODIFIER_NOTE {
					for _, p := range t.Players {
						ops := player.Options{
							Channel:   5,
							Velocity:  t.Velocity,
							Gate:      t.Gate,
							Transpose: 0,
						}
						if err := player.Play(p, step, ops); err != nil {
							log.Errorf("Error playing step: %s", err)
						}
					}
				} else if component.Type == constants.MODIFIER_VELOCITY && len(step.Velocity) > 0 {
					t.Velocity = step.Velocity[step.Iteration%len(step.Velocity)]
				} else if component.Type == constants.MODIFIER_TRANSPOSE && len(step.Transpose) > 0 {
					t.Transpose = step.Transpose[step.Iteration%len(step.Transpose)]
				} else if component.Type == constants.MODIFIER_PROBABILITY && len(step.Probability) > 0 {
					t.Probability = int(step.Probability[step.Iteration%len(step.Probability)])
				} else if component.Type == constants.MODIFIER_GATE && len(step.Gate) > 0 {
					t.Gate = step.Gate[step.Iteration%len(step.Gate)]
				}
				t.Components[i].ChainSteps[j].Iteration++
			}
		}
	}

	t.lastSeconds = ct
	return
}
