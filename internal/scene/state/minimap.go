package state

import "github.com/ambersignal/blacksunrising/pkg/geom"

// MinimapScreenToWorld converts minimap screen coordinates to world coordinates.
func (s *State) MinimapScreenToWorld(pos geom.Vec2) geom.Vec2 {
	// Calculate the position within the minimap's inner area
	innerPos := pos.Sub(s.MiniMap.Min)

	// Normalize the position within the inner area (0-1 range)
	normalizedPos := innerPos.HadamardDevide(s.MiniMap.Size())

	// Clamp to valid range [0, 1]
	normalizedPos[0] = geom.Clamp(normalizedPos[0], 0, 1)
	normalizedPos[1] = geom.Clamp(normalizedPos[1], 0, 1)

	// Convert to world coordinates
	return normalizedPos.HadamardProduct(s.WorldSize)
}
