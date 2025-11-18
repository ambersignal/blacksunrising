package geom

import (
	"image"
	"math"
	"math/rand/v2"
)

// Vec2 represents a 2D vector with X and Y components.
// The vector can be used to represent positions, directions, or offsets in 2D space.
type Vec2 [2]float64

// FromPoint converts an image.Point to a Vec2.
// The integer coordinates of the point are converted to float64 values.
func FromPoint(pt image.Point) Vec2 {
	return Vec2{
		float64(pt.X),
		float64(pt.Y),
	}
}

// RandVec2 creates a random vector with components in the range [0, 1).
// The resulting vector lies within the unit square from (0,0) to (1,1).
func RandVec2() Vec2 {
	return Vec2{
		rand.Float64(),
		rand.Float64(),
	}
}

// Mul returns a new vector scaled by the scalar k.
// Each component of the vector is multiplied by k.
func (v Vec2) Mul(k float64) Vec2 {
	return Vec2{
		v[0] * k,
		v[1] * k,
	}
}

// HadamardProduct returns the element-wise product of two vectors.
// The result vector has components (v[0]*u[0], v[1]*u[1]).
func (v Vec2) HadamardProduct(u Vec2) Vec2 {
	return Vec2{
		v[0] * u[0],
		v[1] * u[1],
	}
}

// HadamardDevide returns the element-wise division of two vectors.
// The result vector has components (v[0]/u[0], v[1]/u[1]).
// Note: Division by zero will result in infinity or NaN values.
func (v Vec2) HadamardDevide(u Vec2) Vec2 {
	return Vec2{
		v[0] / u[0],
		v[1] / u[1],
	}
}

// Add returns the sum of two vectors.
// The result vector has components (v[0]+u[0], v[1]+u[1]).
func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{
		v[0] + u[0],
		v[1] + u[1],
	}
}

// Sub returns the difference between two vectors.
// The result vector has components (v[0]-u[0], v[1]-u[1]).
func (v Vec2) Sub(u Vec2) Vec2 {
	return Vec2{
		v[0] - u[0],
		v[1] - u[1],
	}
}

// Normalize returns a unit vector in the same direction as v.
// If v is a zero vector, the result will be a zero vector.
func (v Vec2) Normalize() Vec2 {
	return v.Mul(1 / v.Length())
}

// Length returns the Euclidean length (magnitude) of the vector.
// Calculated as sqrt(v[0]^2 + v[1]^2).
func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

// Distance returns the Euclidean distance between two vectors.
// This is equivalent to (v.Sub(u)).Length().
func (v Vec2) Distance(u Vec2) float64 {
	return v.Sub(u).Length()
}

// Round returns a new vector with each component rounded to the nearest integer.
func (v Vec2) Round() Vec2 {
	return Vec2{math.Round(v[0]), math.Round(v[1])}
}

// ToImagePoint converts the vector to an image.Point.
// Each component is rounded to the nearest integer and then cast to int.
func (v Vec2) ToImagePoint() image.Point {
	rounded := v.Round()
	return image.Pt(int(rounded[0]), int(rounded[1]))
}

// Dot returns the dot product of two vectors.
// Calculated as v[0]*u[0] + v[1]*u[1].
func (v Vec2) Dot(u Vec2) float64 {
	return v[0]*u[0] + v[1]*u[1]
}

// Angle returns the angle of the vector relative to the positive X-axis.
// The angle is measured in radians, counter-clockwise, in the range [-π, π].
// For a zero vector, the angle is defined as 0.
func (v Vec2) Angle() Angle {
	// Handle zero vector case
	if v[0] == 0 && v[1] == 0 {
		return Angle(0)
	}

	// Use Atan2 to calculate angle from positive X-axis
	// Atan2 returns values in range [-π, π]
	angle := math.Atan2(v[1], v[0])
	return Angle(angle)
}

// Unpack returns the X and Y components of the vector as separate float64 values.
func (v Vec2) Unpack() (float64, float64) {
	return v[0], v[1]
}

// MinVec2 returns a vector with the minimum components from two vectors.
// The result vector has components (min(u[0], v[0]), min(u[1], v[1])).
func MinVec2(u, v Vec2) Vec2 {
	return Vec2{
		math.Min(u[0], v[0]),
		math.Min(u[1], v[1]),
	}
}

// MaxVec2 returns a vector with the maximum components from two vectors.
// The result vector has components (max(u[0], v[0]), max(u[1], v[1])).
func MaxVec2(u, v Vec2) Vec2 {
	return Vec2{
		math.Max(u[0], v[0]),
		math.Max(u[1], v[1]),
	}
}
