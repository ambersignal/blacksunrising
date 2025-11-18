package geom

import (
	"math"
)

// Angle represents an angle in radians.
// Angles are typically in the range [-π, π] after normalization.
type Angle float64

// Normalize returns an equivalent angle in the range [-π, π].
// This ensures consistent representation of angles regardless of how they were calculated.
func (a Angle) Normalize() Angle {
	rad := math.Remainder(float64(a), 2*math.Pi)
	if rad <= -math.Pi {
		rad = math.Pi
	}

	return Angle(rad)
}

// Abs returns the absolute value of the angle.
// This is useful when the magnitude of the angle is needed without regard to direction.
func (a Angle) Abs() Angle {
	return Angle(math.Abs(float64(a)))
}
