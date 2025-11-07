package geom

import (
	"math"
	"testing"
)

func TestVec2Angle(t *testing.T) {
	tests := []struct {
		name string
		vec  Vec2
	}{
		{
			name: "Unit vector pointing right",
			vec:  Vec2{1, 0},
		},
		{
			name: "Unit vector pointing up",
			vec:  Vec2{0, 1},
		},
		{
			name: "Unit vector pointing left",
			vec:  Vec2{-1, 0},
		},
		{
			name: "Unit vector pointing down",
			vec:  Vec2{0, -1},
		},
		{
			name: "Zero vector",
			vec:  Vec2{0, 0},
		},
		{
			name: "Diagonal vector",
			vec:  Vec2{1, 1},
		},
		{
			name: "Negative diagonal vector",
			vec:  Vec2{-1, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			result := tt.vec.Angle()

			// For zero vector, just check it doesn't crash
			if tt.vec.Length() == 0 {
				_ = result
				return
			}

			// Verify that the angle is normalized
			normalized := result.Normalize()
			const epsilon = 1e-10
			if math.Abs(float64(result-normalized)) > epsilon {
				t.Errorf("Angle %v is not normalized, should be %v", result, normalized)
			}
		})
	}
}

func TestVec2AngleSpecificValues(t *testing.T) {
	const epsilon = 1e-10

	// Test vector pointing right: {1, 0}
	// Expected angle: 0 radians (pointing right)
	vec1 := Vec2{1, 0}
	angle1 := vec1.Angle()
	expected1 := Angle(0)
	if math.Abs(float64(angle1-expected1)) > epsilon {
		t.Errorf("Angle({1,0}) = %v, want %v", angle1, expected1)
	}

	// Test vector pointing up: {0, 1}
	// Expected angle: π/2 radians (pointing up)
	vec2 := Vec2{0, 1}
	angle2 := vec2.Angle()
	expected2 := Angle(math.Pi / 2)
	if math.Abs(float64(angle2-expected2)) > epsilon {
		t.Errorf("Angle({0,1}) = %v, want %v", angle2, expected2)
	}

	// Test vector pointing left: {-1, 0}
	// Expected angle: π radians (pointing left)
	vec3 := Vec2{-1, 0}
	angle3 := vec3.Angle()
	expected3 := Angle(math.Pi)
	if math.Abs(float64(angle3-expected3)) > epsilon {
		t.Errorf("Angle({-1,0}) = %v, want %v", angle3, expected3)
	}

	// Test vector pointing down: {0, -1}
	// Expected angle: -π/2 radians (pointing down)
	vec4 := Vec2{0, -1}
	angle4 := vec4.Angle()
	expected4 := Angle(-math.Pi / 2)
	if math.Abs(float64(angle4-expected4)) > epsilon {
		t.Errorf("Angle({0,-1}) = %v, want %v", angle4, expected4)
	}

	// Test diagonal vector: {1, 1}
	// Expected angle: π/4 radians (45 degrees)
	vec5 := Vec2{1, 1}
	angle5 := vec5.Angle()
	expected5 := Angle(math.Pi / 4)
	if math.Abs(float64(angle5-expected5)) > epsilon {
		t.Errorf("Angle({1,1}) = %v, want %v", angle5, expected5)
	}
}

func TestVec2AngleProperties(t *testing.T) {
	// Test various properties of the angle method

	// Test that all angles are in the valid range after normalization
	vectors := []Vec2{
		{1, 0}, {0, 1}, {-1, 0}, {0, -1},
		{1, 1}, {-1, 1}, {-1, -1}, {1, -1},
		{3, 4}, {-2, 5}, {7, -3}, {-4, -6},
	}

	for _, vec := range vectors {
		if vec.Length() == 0 {
			continue // Skip zero vector
		}

		angle := vec.Angle()
		normalized := angle.Normalize()

		const epsilon = 1e-10
		if math.Abs(float64(angle-normalized)) > epsilon {
			t.Errorf("Angle %v for vector %v is not normalized", angle, vec)
		}

		// Check that angle is in valid range [-π, π]
		if float64(normalized) < -math.Pi || float64(normalized) > math.Pi {
			t.Errorf("Normalized angle %v is out of range [-π, π]", normalized)
		}
	}
}

func TestVec2AngleConsistency(t *testing.T) {
	// Test that angle calculation is consistent with trigonometric expectations

	tests := []struct {
		vec      Vec2
		expected Angle
	}{
		{Vec2{1, 0}, 0},
		{Vec2{0, 1}, math.Pi / 2},
		{Vec2{-1, 0}, math.Pi},
		{Vec2{0, -1}, -math.Pi / 2},
		{Vec2{1, 1}, math.Pi / 4},
		{Vec2{-1, 1}, 3 * math.Pi / 4},
		{Vec2{-1, -1}, -3 * math.Pi / 4},
		{Vec2{1, -1}, -math.Pi / 4},
	}

	const epsilon = 1e-10

	for _, tt := range tests {
		if tt.vec.Length() == 0 {
			continue
		}

		angle := tt.vec.Angle()
		normalized := angle.Normalize()

		if math.Abs(float64(normalized-tt.expected)) > epsilon {
			t.Errorf("Angle(%v) = %v, want %v", tt.vec, normalized, tt.expected)
		}
	}
}
