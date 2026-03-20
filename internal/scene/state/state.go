package state

import "github.com/ambersignal/blacksunrising/pkg/geom"

// State holds all the game state information.
type State struct {
	Camera  geom.Rectangle // Camera viewpoint
	MiniMap geom.Rectangle

	WorldSize geom.Vec2

	Ships             []*Ship            // All ships in the game
	Asteroids         []*Asteroid        // All asteroids in the game (flattened from fields)
	AsteroidFields    []*AsteroidField   // All asteroid fields (for resource management)
	Groups            []*Group           // All groups of ships
	Selected          map[*Ship]struct{} // Set of currently selected ships
	CurrentGroupIndex int                // Index of the currently active group
	Planet            *Planet            // The planet in the game world
}

// NewState creates a new game state.
func NewState() *State {
	return &State{
		Ships:             make([]*Ship, 0),
		Asteroids:         make([]*Asteroid, 0),
		AsteroidFields:    make([]*AsteroidField, 0),
		Groups:            make([]*Group, 0),
		Selected:          make(map[*Ship]struct{}),
		CurrentGroupIndex: -1, // No group selected initially
		Planet:            nil,
	}
}

// AddShip adds a ship to the state.
func (s *State) AddShip(ship *Ship) {
	s.Ships = append(s.Ships, ship)
}

// AddAsteroid adds an asteroid to the state.
func (s *State) AddAsteroid(asteroid *Asteroid) {
	s.Asteroids = append(s.Asteroids, asteroid)
}

// AddAsteroidField adds an asteroid field to the state.
func (s *State) AddAsteroidField(field *AsteroidField) {
	s.AsteroidFields = append(s.AsteroidFields, field)
	// Also add all asteroids from the field to the flattened list
	s.Asteroids = append(s.Asteroids, field.Asteroids...)
}

// AddGroup adds a group to the state.
func (s *State) AddGroup(group *Group) {
	s.Groups = append(s.Groups, group)
}

// CleanupEmptyGroups removes groups that have no ships.
func (s *State) CleanupEmptyGroups() {
	// Iterate backwards to safely remove elements
	for i := len(s.Groups) - 1; i >= 0; i-- {
		if s.Groups[i].IsEmpty() {
			// Remove the group
			s.Groups = append(s.Groups[:i], s.Groups[i+1:]...)

			// Adjust current group index if needed
			if s.CurrentGroupIndex >= i && s.CurrentGroupIndex > 0 {
				s.CurrentGroupIndex--
			} else if s.CurrentGroupIndex >= len(s.Groups) {
				s.CurrentGroupIndex = len(s.Groups) - 1
			}
		}
	}

	// If no groups left, reset current group index
	if len(s.Groups) == 0 {
		s.CurrentGroupIndex = -1
	}
}

// GetGroupForShip returns the group that contains the ship, or nil if not in any group.
func (s *State) GetGroupForShip(ship *Ship) *Group {
	for _, group := range s.Groups {
		if group.Contains(ship) {
			return group
		}
	}
	return nil
}

// RemoveSelectedFromAllGroups removes selected ships from all existing groups.
func (s *State) RemoveSelectedFromAllGroups() {
	// For each selected ship, remove it from any group it might be in
	for ship := range s.Selected {
		for _, group := range s.Groups {
			group.RemoveShip(ship)
		}
	}
}
