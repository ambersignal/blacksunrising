package geom

import (
	"image"
	"math"
)

type Vec2 [2]float64

func FromPoint(pt image.Point) Vec2 {
	return Vec2{
		float64(pt.X),
		float64(pt.Y),
	}
}

func (v Vec2) Mul(k float64) Vec2 {
	return Vec2{
		v[0] * k,
		v[1] * k,
	}
}

func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{
		v[0] + u[0],
		v[1] + u[1],
	}
}

func (v Vec2) Sub(u Vec2) Vec2 {
	return Vec2{
		v[0] - u[0],
		v[1] - u[1],
	}
}

func (v Vec2) Normalize() Vec2 {
	return v.Mul(1 / v.Length())
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec2) Distance(u Vec2) float64 {
	return v.Sub(u).Length()
}

func (v Vec2) Round() Vec2 {
	return Vec2{math.Round(v[0]), math.Round(v[1])}
}

func (v Vec2) ToImagePoint() image.Point {
	rounded := v.Round()

	return image.Pt(int(rounded[0]), int(rounded[1]))
}

func (v Vec2) Dot(u Vec2) float64 {
	return v[0]*u[0] + v[1]*u[1]
}

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

func (v Vec2) Unpack() (float64, float64) {
	return v[0], v[1]
}

func MinVec2(u, v Vec2) Vec2 {
	return Vec2{
		math.Min(u[0], v[0]),
		math.Min(u[1], v[1]),
	}
}

func MaxVec2(u, v Vec2) Vec2 {
	return Vec2{
		math.Max(u[0], v[0]),
		math.Max(u[1], v[1]),
	}
}
