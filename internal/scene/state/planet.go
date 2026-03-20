package state

import (
	"github.com/ambersignal/blacksunrising/pkg/geom"
)

// Planet represents a planet in the game.
type Planet struct {
	Pos geom.Vec2
}

// NewPlanet creates a new planet at the specified position.
func NewPlanet(pos geom.Vec2) *Planet {
	return &Planet{
		Pos: pos,
	}
}
