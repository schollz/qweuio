package tli

import (
	"asdfgh/src/pattern"
	"asdfgh/src/player"
	"asdfgh/src/step"
	"sync"
	"time"
)

type Component struct {
	Type          string                     `json:"type,omitempty"`
	Chain         []string                   `json:"chain,omitempty"`
	Patterns      map[string]pattern.Pattern `json:"patterns,omitempty"`
	ChainSteps    []step.Step                `json:"chain_steps,omitempty"`
	ChainDuration float64                    `json:"chain_duration,omitempty"`
}

type TLI struct {
	BPM        float64         `json:"bpm,omitempty"`
	Components []Component     `json:"components,omitempty"`
	Players    []player.Player `json:"players,omitempty"`
	// for realtime playback
	Probability int     `json:"probability,omitempty"`
	Velocity    int     `json:"velocity,omitempty"`
	Transpose   float64 `json:"transpose,omitempty"`
	Gate        float64 `json:"gate,omitempty"`
	// create a mutex
	mutex sync.Mutex
	// create a channel for stopping playback
	stopChan    chan bool
	isPlaying   bool
	lastSeconds float64
	startTime   time.Time
}
