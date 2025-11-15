package scene

import "github.com/ambersignal/blacksunrising/pkg/geom"

// State holds all the game state information
type State struct {
	Camera    geom.Rectangle // Camera viewpoint
	WorldSize geom.Vec2

	Ships             []*Ship            // All ships in the game
	Groups            []*Group           // All groups of ships
	Selected          map[*Ship]struct{} // Set of currently selected ships
	CurrentGroupIndex int                // Index of the currently active group
}

// NewState creates a new game state
func NewState() *State {
	return &State{
		Ships:             make([]*Ship, 0),
		Groups:            make([]*Group, 0),
		Selected:          make(map[*Ship]struct{}),
		CurrentGroupIndex: -1, // No group selected initially
	}
}

// AddShip adds a ship to the state
func (s *State) AddShip(ship *Ship) {
	s.Ships = append(s.Ships, ship)
}

// AddGroup adds a group to the state
func (s *State) AddGroup(group *Group) {
	s.Groups = append(s.Groups, group)
}
