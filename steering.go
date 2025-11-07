package main

import (
	"github.com/ambersignal/blacksunrising/pkg/geom"
)

var (
	MaxSeparation = 50.0
	MinSeparation = 0.0
)

// Alignment calculates alignment force for a ship based on nearby ships
func (g *Game) Alignment(ship *Ship) geom.Vec2 {
	sum := geom.Vec2{0, 0}
	count := 0

	// Look at all other ships
	for _, other := range g.Ships {
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

// Separation calculates separation force for a ship based on nearby ships
func (g *Game) Separation(ship *Ship) geom.Vec2 {
	var neighbors float64
	sum := geom.Vec2{0, 0}

	// Look at all other ships
	for _, other := range g.Ships {
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
				(1 - SmoothStep(MinSeparation, MaxSeparation, distance))
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

// Cohesion calculates cohesion force for a ship based on nearby ships
func (g *Game) Cohesion(ship *Ship) geom.Vec2 {
	sum := geom.Vec2{0, 0}
	count := 0

	// Look at all other ships
	for _, other := range g.Ships {
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
func (g *Game) Seek(ship *Ship, target geom.Vec2) geom.Vec2 {
	// Calculate desired velocity toward target
	desired := target.Sub(ship.Pos)

	// If we're very close to the target, stop moving
	if desired.Length() < 5.0 {
		// Slow down gradually
		return ship.Vel.Mul(-1.0) // Counteract current velocity
	}

	// Normalize and scale to desired velocity
	if desired.Length() > 0 {
		desired = desired.Normalize().Mul(50) // Desired velocity
	}

	// Reynolds steering: subtract current velocity to get steering force
	steer := desired.Sub(ship.Vel)
	return steer
}

func SmoothStep(edge0, edge1, x float64) float64 {
	t := Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)

	return t * t * (3 - 2*t)
}

func Clamp(x, low, high float64) float64 {
	if x < low {
		return low
	}

	if x > high {
		return high
	}

	return x
}
