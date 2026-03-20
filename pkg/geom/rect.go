package geom

import (
	"image"
	"math"
)

const (
	half = 2.0
)

// Rectangle represents a rectangle defined by its minimum and maximum corners.
// The rectangle is defined such that Min contains the lower bounds and Max contains the upper bounds.
type Rectangle struct {
	Min, Max Vec2
}

// RectangleBySize creates a rectangle from a position and size.
// The position defines the a center of the rectangle.
func RectangleBySize(pos Vec2, size Vec2) Rectangle {
	semiSize := size.Div(2)

	return Rectangle{
		Min: pos.Sub(semiSize),
		Max: pos.Add(semiSize),
	}
}

// Empty reports whether the rectangle has zero or negative area.
// A rectangle is considered empty if its width or height is less than or equal to zero.
func (r Rectangle) Empty() bool {
	return r.Min[0] >= r.Max[0] || r.Min[1] >= r.Max[1]
}

// Overlaps reports whether two rectangles intersect.
// Two rectangles overlap if they share any interior points.
// Empty rectangles never overlap.
func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min[0] < s.Max[0] && s.Min[0] < r.Max[0] &&
		r.Min[1] < s.Max[1] && s.Min[1] < r.Max[1]
}

// Sub returns a rectangle translated by subtracting the given vector from both corners.
// This effectively moves the rectangle in the opposite direction of the vector.
func (r Rectangle) Sub(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.Sub(v),
		Max: r.Max.Sub(v),
	}
}

// Add returns a rectangle translated by adding the given vector to both corners.
// This effectively moves the rectangle in the direction of the vector.
func (r Rectangle) Add(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.Add(v),
		Max: r.Max.Add(v),
	}
}

// Center returns the center point of the rectangle.
// The center is calculated as the midpoint between the minimum and maximum corners.
func (r Rectangle) Center() Vec2 {
	return r.Max.Add(r.Min).Div(half)
}

// Size returns the dimensions of the rectangle as a vector.
// The size is calculated as the difference between the maximum and minimum corners.
func (r Rectangle) Size() Vec2 {
	return r.Max.Sub(r.Min)
}

// HadamardProduct returns element-wise product of the rectangle's corners with a vector.
// Both the Min and Max corners are multiplied component-wise by the given vector.
func (r Rectangle) HadamardProduct(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.HadamardProduct(v),
		Max: r.Max.HadamardProduct(v),
	}
}

// HadamardDivide returns element-wise division of the rectangle's corners by a vector.
// Both the Min and Max corners are divided component-wise by the given vector.
// Note: Division by zero components will result in infinity or NaN values.
func (r Rectangle) HadamardDivide(v Vec2) Rectangle {
	return Rectangle{
		Min: r.Min.HadamardDevide(v),
		Max: r.Max.HadamardDevide(v),
	}
}

// Round returns a rectangle with both corners rounded to the nearest integer coordinates.
func (r Rectangle) Round() Rectangle {
	return Rectangle{
		Min: r.Min.Round(),
		Max: r.Max.Round(),
	}
}

// FromRectangle converts an image.Rectangle to a geom.Rectangle.
// The integer coordinates are converted to float64 values.
func FromRectangle(r image.Rectangle) Rectangle {
	return Rectangle{
		Min: FromPoint(r.Min),
		Max: FromPoint(r.Max),
	}
}

// Rect creates a rectangle given position (x,y) and size (w,h).
// The minimum corner is set to (x,y) and the maximum corner is set to (x+w, y+h).
func Rect(x, y, w, h float64) Rectangle {
	return Rectangle{
		Min: Vec2{x, y},
		Max: Vec2{x + w, y + h},
	}
}

// Normalize returns a rectangle with properly ordered corners.
// The minimum corner will contain the smaller values and the maximum corner will contain the larger values.
func (r Rectangle) Normalize() Rectangle {
	return Rectangle{
		Min: MinVec2(r.Min, r.Max),
		Max: MaxVec2(r.Min, r.Max),
	}
}

// Mul returns a rectangle with both corners scaled by the scalar k.
// Each coordinate of both corners is multiplied by k.
func (r Rectangle) Mul(k float64) Rectangle {
	return Rectangle{
		Min: r.Min.Mul(k),
		Max: r.Max.Mul(k),
	}
}

// IntersectsCircle checks if the rectangle intersects with a circle.
// The circle is defined by its center point and radius.
// Returns true if the rectangle and circle share any points.
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
