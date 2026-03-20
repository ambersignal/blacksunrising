package scene

import (
	"github.com/ambersignal/blacksunrising/internal/scene/state"
	"github.com/ambersignal/blacksunrising/pkg/geom"
)

var (
	// FIXME(evgenii.omelchenko): should depends on ship size
	MaxSeparation = 64.0
	MinSeparation = 0.0
)

// Alignment calculates alignment force for a ship based on nearby ships in the same group
func (g *Scene) Alignment(ship *state.Ship) geom.Vec2 {
	// Find which group this ship belongs to
	group := g.state.GetGroupForShip(ship)

	// Only apply alignment if ship is in a group
	if group == nil {
		return geom.Vec2{0, 0}
	}

	// Need at least 2 ships in group for alignment to make sense
	if group.Size() < 2 {
		return geom.Vec2{0, 0}
	}

	sum := geom.Vec2{0, 0}
	count := 0

	// Look at ships in the same group only
	for _, other := range group.Ships {
		// Skip self
		if other == ship {
			continue
		}

		sum = sum.Add(other.Vel)
		count++
	}

	// If no neighbors, return zero vector
	if count == 0 {
		return geom.Vec2{0, 0}
	}

	// Average the velocities of neighbors
	average := sum.Mul(1.0 / float64(count))

	// Steer towards desired velocity
	if average.Length() > 0 {
		average = average.Normalize().Mul(50) // Desired velocity
	}

	// Reynolds steering: subtract current velocity to get steering force
	steer := average.Sub(ship.Vel)
	return steer
}

// Separation calculates separation force for a ship based on nearby ships in the same group
func (g *Scene) Separation(ship *state.Ship) geom.Vec2 {
	// Find which group this ship belongs to
	group := g.state.GetGroupForShip(ship)

	// Only apply separation if ship is in a group
	if group == nil {
		return geom.Vec2{0, 0}
	}

	// Need at least 2 ships in group for separation to make sense
	if group.Size() < 2 {
		return geom.Vec2{0, 0}
	}

	var neighbors float64
	sum := geom.Vec2{0, 0}

	// Look at ships in the same group only
	for _, other := range group.Ships {
		// Skip self
		if other == ship {
			continue
		}

		// Calculate distance to other ship
		distance := ship.Pos.Distance(other.Pos)

		// If too close, contribute to separation force
		if distance < MaxSeparation && distance > MinSeparation {
			pushForce := ship.Pos.Sub(other.Pos)
			separationStrength := pushForce.Length() *
				(1 - geom.SmoothStep(MinSeparation, MaxSeparation, distance))
			pushForce = pushForce.Normalize().Mul(pushForce.Length() * separationStrength)

			sum = sum.Add(pushForce)
			neighbors++
		}
	}

	// If no nearby ships, return zero vector
	if sum.Length() == 0 {
		return geom.Vec2{0, 0}
	}

	if neighbors != 0 {
		sum = sum.Mul(1 / neighbors)
	}

	// Reynolds steering: subtract current velocity to get steering force
	steer := sum.Sub(ship.Vel)
	return steer
}

// Cohesion calculates cohesion force for a ship based on nearby ships in the same group
func (g *Scene) Cohesion(ship *state.Ship) geom.Vec2 {
	// Find which group this ship belongs to
	group := g.state.GetGroupForShip(ship)

	// Only apply cohesion if ship is in a group
	if group == nil {
		return geom.Vec2{0, 0}
	}

	// Need at least 2 ships in group for cohesion to make sense
	if group.Size() < 2 {
		return geom.Vec2{0, 0}
	}

	sum := geom.Vec2{0, 0}
	count := 0

	// Look at ships in the same group only
	for _, other := range group.Ships {
		// Skip self
		if other == ship {
			continue
		}

		// Calculate distance to other ship
		distance := ship.Pos.Distance(other.Pos)

		// If too close then disable cohesion to avoid concurrency with the separation.
		if distance > MaxSeparation {
			sum = sum.Add(other.Pos)
			count++
		}
	}

	// If no neighbors, return zero vector
	if count == 0 {
		return geom.Vec2{0, 0}
	}

	// Average position of neighbors
	average := sum.Mul(1.0 / float64(count))

	// Desired velocity towards target
	desired := average.Sub(ship.Pos)

	// Reynolds steering: subtract current velocity to get steering force
	steer := desired.Sub(ship.Vel)
	return steer
}

// Seek calculates a steering force to move a ship toward a target position
func (g *Scene) Seek(ship *state.Ship, target geom.Vec2) geom.Vec2 {
	// Calculate desired velocity toward target
	desired := target.Sub(ship.Pos)

	// If we're very close to the target, stop moving
	if desired.Length() < 5.0 {
		// Slow down gradually
		return ship.Vel.Mul(-1.0) // Counteract current velocity
	}

	// Normalize and scale to desired velocity
	if desired.Length() > 0 {
		desired = desired.Normalize().Mul(ship.MaxSpeed) // Desired velocity
	}

	// Reynolds steering: subtract current velocity to get steering force
	steer := desired.Sub(ship.Vel)
	return steer
}
