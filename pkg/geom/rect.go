package geom

import (
	"image"
	"math"
)

type Rectangle struct {
	Min, Max Vec2
}

func RectangleBySize(pos Vec2, size Vec2) Rectangle {
	return Rectangle{
		Min: pos,
		Max: size.Add(pos),
	}
}

func (r Rectangle) Empty() bool {
	return r.Min[0] >= r.Max[0] || r.Min[1] >= r.Max[1]
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min[0] < s.Max[0] && s.Min[0] < r.Max[0] &&
		r.Min[1] < s.Max[1] && s.Min[1] < r.Max[1]
}

func (r Rectangle) Sub(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.Sub(v),
		Max: r.Max.Sub(v),
	}
}

func (r Rectangle) Add(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.Add(v),
		Max: r.Max.Add(v),
	}
}

func (r Rectangle) Center() Vec2 {
	return r.Max.Add(r.Min).Mul(0.5)
}

func (r Rectangle) Size() Vec2 {
	return r.Max.Sub(r.Min)
}

// HadamardProduct returns element-wise product of two rectangles
func (r Rectangle) HadamardProduct(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.HadamardProduct(v),
		Max: r.Max.HadamardProduct(v),
	}
}

// HadamardDivide returns element-wise division of two rectangles
func (r Rectangle) HadamardDivide(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.HadamardDevide(v),
		Max: r.Max.HadamardDevide(v),
	}
}

func (r Rectangle) Round() Rectangle {
	return Rectangle{
		Min: r.Min.Round(),
		Max: r.Max.Round(),
	}
}

func FromRectangle(r image.Rectangle) Rectangle {
	return Rectangle{
		Min: FromPoint(r.Min),
		Max: FromPoint(r.Max),
	}
}

func Rect(x, y, w, h float64) Rectangle {
	return Rectangle{
		Min: Vec2{x, y},
		Max: Vec2{x + w, y + h},
	}
}

func (r Rectangle) Normalize() Rectangle {
	return Rectangle{
		Min: MinVec2(r.Min, r.Max),
		Max: MaxVec2(r.Min, r.Max),
	}
}

func (r Rectangle) Mul(k float64) Rectangle {
	return Rectangle{
		Min: r.Min.Mul(k),
		Max: r.Max.Mul(k),
	}
}

// IntersectsCircle checks if the rectangle intersects with a circle centered at center with given radius
func (r Rectangle) IntersectsCircle(center Vec2, radius float64) bool {
	// Find the closest point on the rectangle to the circle center
	closestX := math.Max(r.Min[0], math.Min(center[0], r.Max[0]))
	closestY := math.Max(r.Min[1], math.Min(center[1], r.Max[1]))

	// Calculate distance between closest point and circle center
	distanceX := center[0] - closestX
	distanceY := center[1] - closestY

	// If distance is less than radius, they intersect
	return (distanceX*distanceX + distanceY*distanceY) < (radius * radius)
}
