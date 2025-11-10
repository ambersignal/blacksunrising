package geom

import (
	"testing"
)

func TestRectangleIntersectsCircle(t *testing.T) {
	tests := []struct {
		name       string
		rect       Rectangle
		center     Vec2
		radius     float64
		intersects bool
	}{
		{
			name:       "Circle completely inside rectangle",
			rect:       Rect(0, 0, 10, 10),
			center:     Vec2{5, 5},
			radius:     2,
			intersects: true,
		},
		{
			name:       "Circle completely outside rectangle",
			rect:       Rect(0, 0, 10, 10),
			center:     Vec2{15, 15},
			radius:     2,
			intersects: false,
		},
		{
			name:       "Circle intersects rectangle edge",
			rect:       Rect(0, 0, 10, 10),
			center:     Vec2{5, 12},
			radius:     3,
			intersects: true,
		},
		{
			name:       "Circle intersects rectangle corner",
			rect:       Rect(0, 0, 10, 10),
			center:     Vec2{11, 11},
			radius:     2,
			intersects: true,
		},
		{
			name:       "Circle near but not touching rectangle corner",
			rect:       Rect(0, 0, 10, 10),
			center:     Vec2{13, 13},
			radius:     1,
			intersects: false,
		},
		{
			name:       "Large circle encompasses rectangle",
			rect:       Rect(2, 2, 2, 2),
			center:     Vec2{5, 5},
			radius:     10,
			intersects: true,
		},
		{
			name:       "Circle touches rectangle top edge",
			rect:       Rect(0, 0, 10, 10),
			center:     Vec2{5, -1},
			radius:     2,
			intersects: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.IntersectsCircle(tt.center, tt.radius)
			if result != tt.intersects {
				t.Errorf("IntersectsCircle(%v, %v, %v) = %v, want %v",
					tt.rect, tt.center, tt.radius, result, tt.intersects)
			}
		})
	}
}

func TestRectangleEmpty(t *testing.T) {
	tests := []struct {
		name  string
		rect  Rectangle
		empty bool
	}{
		{
			name:  "Valid rectangle",
			rect:  Rect(0, 0, 10, 10),
			empty: false,
		},
		{
			name:  "Zero width rectangle",
			rect:  Rect(0, 0, 0, 10),
			empty: true,
		},
		{
			name:  "Zero height rectangle",
			rect:  Rect(0, 0, 10, 0),
			empty: true,
		},
		{
			name:  "Negative width rectangle",
			rect:  Rect(5, 0, -2, 10),
			empty: true,
		},
		{
			name:  "Negative height rectangle",
			rect:  Rect(0, 5, 10, -2),
			empty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.Empty()
			if result != tt.empty {
				t.Errorf("Empty(%v) = %v, want %v", tt.rect, result, tt.empty)
			}
		})
	}
}

func TestRectangleOverlaps(t *testing.T) {
	tests := []struct {
		name     string
		rect1    Rectangle
		rect2    Rectangle
		overlaps bool
	}{
		{
			name:     "Identical rectangles",
			rect1:    Rect(0, 0, 10, 10),
			rect2:    Rect(0, 0, 10, 10),
			overlaps: true,
		},
		{
			name:     "Separate rectangles",
			rect1:    Rect(0, 0, 5, 5),
			rect2:    Rect(10, 10, 5, 5),
			overlaps: false,
		},
		{
			name:     "Touching rectangles",
			rect1:    Rect(0, 0, 5, 5),
			rect2:    Rect(5, 0, 5, 5),
			overlaps: false, // Touching edges don't count as overlapping
		},
		{
			name:     "Partially overlapping rectangles",
			rect1:    Rect(0, 0, 10, 10),
			rect2:    Rect(5, 5, 10, 10),
			overlaps: true,
		},
		{
			name:     "One rectangle inside another",
			rect1:    Rect(0, 0, 20, 20),
			rect2:    Rect(5, 5, 5, 5),
			overlaps: true,
		},
		{
			name:     "Empty rectangle doesn't overlap",
			rect1:    Rect(0, 0, 0, 0),
			rect2:    Rect(5, 5, 5, 5),
			overlaps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect1.Overlaps(tt.rect2)
			if result != tt.overlaps {
				t.Errorf("Overlaps(%v, %v) = %v, want %v", tt.rect1, tt.rect2, result, tt.overlaps)
			}
		})
	}
}

func TestRectangleOperations(t *testing.T) {
	rect := Rect(1, 2, 8, 6) // Min={1,2}, Max={9,8}

	// Test Center
	center := rect.Center()
	expectedCenter := Vec2{5, 5}
	if center != expectedCenter {
		t.Errorf("Center() = %v, want %v", center, expectedCenter)
	}

	// Test Size
	size := rect.Size()
	expectedSize := Vec2{8, 6}
	if size != expectedSize {
		t.Errorf("Size() = %v, want %v", size, expectedSize)
	}

	// Test Add
	offset := Vec2{2, 3}
	movedRect := rect.Add(offset)
	expectedMovedRect := Rectangle{Vec2{3, 5}, Vec2{11, 11}}
	if movedRect != expectedMovedRect {
		t.Errorf("Add(%v) = %v, want %v", offset, movedRect, expectedMovedRect)
	}

	// Test Sub
	subtractedRect := rect.Sub(offset)
	expectedSubtractedRect := Rectangle{Vec2{-1, -1}, Vec2{7, 5}}
	if subtractedRect != expectedSubtractedRect {
		t.Errorf("Sub(%v) = %v, want %v", offset, subtractedRect, expectedSubtractedRect)
	}
}
