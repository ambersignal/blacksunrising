package state

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/ambersignal/blacksunrising/pkg/geom"
)

var (
	selectionColor = color.RGBA{30, 188, 115, 200}
)

// Global constants for ship movement limits
const (
	MaxSpeed         = 128.0       // Maximum speed in pixels per second
	MaxRotationSpeed = math.Pi / 4 // Maximum rotation speed in radians per second
)

// Ship represents a ship in the game
type Ship struct {
	Pos        geom.Vec2
	Image      *ebiten.Image
	Vel        geom.Vec2 // Velocity in pixels per second
	Accel      geom.Vec2 // Acceleration in pixels per second squared
	MaxSpeed   float64
	IsSelected bool // Whether this ship is currently selected
}

// NewShip creates a new ship with the specified position, velocity, and image
func NewShip(pos geom.Vec2, vel geom.Vec2, img *ebiten.Image) *Ship {
	return &Ship{
		Pos:      pos,
		Image:    img,
		Vel:      vel,
		Accel:    geom.Vec2{0, 0},
		MaxSpeed: MaxSpeed,

		IsSelected: false,
	}
}

// Update updates the ship's state
func (s *Ship) Update(elapsedTime time.Duration) error {
	// Convert elapsed time to seconds for calculations
	deltaSeconds := elapsedTime.Seconds()

	// Update velocity based on acceleration
	s.Vel = s.Vel.Add(s.Accel.Mul(deltaSeconds))

	// Enforce maximum speed
	if s.Vel.Length() > MaxSpeed {
		s.Vel = s.Vel.Normalize().Mul(MaxSpeed)
	}

	// Update position based on velocity
	s.Pos = s.Pos.Add(s.Vel.Mul(deltaSeconds)).Round()

	// Reset acceleration for next frame
	s.Accel = geom.Vec2{0, 0}

	return nil
}

// Draw renders the ship with an optional camera offset
func (s *Ship) Draw(screen *ebiten.Image, cameraOffset ...geom.Vec2) {
	if s.Image == nil {
		return
	}

	// Get image dimensions
	bounds := s.Image.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create options for drawing
	opts := &ebiten.DrawImageOptions{}

	// Calculate angle from velocity vector
	// Sprites are oriented in -pi/2 angle
	angle := s.Vel.Angle() + math.Pi/2

	// Apply rotation around the center of the image
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	opts.GeoM.Translate(-centerX, -centerY)
	opts.GeoM.Rotate(float64(angle))
	opts.GeoM.Translate(centerX, centerY)

	// Apply camera offset if provided
	drawPos := s.Pos
	if len(cameraOffset) > 0 {
		drawPos = s.Pos.Sub(cameraOffset[0])
	}

	// Set the position (center the image on the ship's position)
	opts.GeoM.Translate(drawPos.Round().Unpack())

	// Draw the image
	screen.DrawImage(s.Image, opts)

	// Draw selection indicator if selected
	if s.IsSelected {
		// Draw a circle around the ship
		vector.StrokeCircle(screen, float32(math.Round((drawPos[0] + centerX))),
			float32(math.Round(drawPos[1]+centerY)),
			float32(s.Radius()), 1, selectionColor, false)
	}
}

// Radius of the ship.
func (s *Ship) Radius() float64 {
	// Get ship size
	bounds := s.Image.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// Calculate the radius of the ship for circular collision detection
	return math.Max(width, height) / 2
}
