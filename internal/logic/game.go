package logic

import (
	"time"

	"github.com/ambersignal/blacksunrising/internal/state"
)

// Game represents the core game logic structure
type Game struct {
	state *state.State
}

// Update updates the game logic based on the current time
func (g *Game) Update(shift time.Duration) error {
	// Calculate next positions, velocity and angles via steering algorithm

	return nil
}
