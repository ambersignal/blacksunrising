package state

import (
	"github.com/ambersignal/blacksunrising/pkg/geom"
)

// Group represents a collection of ships with shared behavior
type Group struct {
	Ships     []*Ship   // Ships belonging to this group
	Target    geom.Vec2 // Target position for the group
	HasTarget bool      // Whether the group has a target set
}

// NewGroup creates a new group
func NewGroup() *Group {
	return &Group{
		Ships: make([]*Ship, 0),
	}
}

// AddShip adds a ship to the group
func (g *Group) AddShip(ship *Ship) {
	// Check if ship is already in the group
	for _, s := range g.Ships {
		if s == ship {
			return // Already in group
		}
	}
	g.Ships = append(g.Ships, ship)
}

// RemoveShip removes a ship from the group
func (g *Group) RemoveShip(ship *Ship) {
	for i, s := range g.Ships {
		if s == ship {
			// Remove the ship by slicing
			g.Ships = append(g.Ships[:i], g.Ships[i+1:]...)
			return
		}
	}
}

// Contains checks if a ship is in the group
func (g *Group) Contains(ship *Ship) bool {
	for _, s := range g.Ships {
		if s == ship {
			return true
		}
	}
	return false
}

// IsEmpty checks if the group has no ships
func (g *Group) IsEmpty() bool {
	return len(g.Ships) == 0
}

// Size returns the number of ships in the group
func (g *Group) Size() int {
	return len(g.Ships)
}

// Clear removes all ships from the group
func (g *Group) Clear() {
	g.Ships = g.Ships[:0]
}
