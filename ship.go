package main

import (
	"math"
	"time"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	IsSelected bool      // Whether this ship is currently selected
}

// NewShip creates a new ship with the specified position, velocity, and image
func NewShip(pos geom.Vec2, vel geom.Vec2, img *ebiten.Image) *Ship {
	return &Ship{
		Pos:        pos,
		Image:      img,
		Vel:        vel,
		Accel:      geom.Vec2{0, 0},
		IsSelected: false,
	}
}

// LoadShipImage loads the ship image from file
func LoadShipImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("data/fighter.png")
	return img, err
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
	s.Pos = s.Pos.Add(s.Vel.Mul(deltaSeconds))

	// Reset acceleration for next frame
	s.Accel = geom.Vec2{0, 0}

	return nil
}

// Draw renders the ship
func (s *Ship) Draw(screen *ebiten.Image) {
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

	// Set the position (center the image on the ship's position)
	opts.GeoM.Translate(s.Pos[0], s.Pos[1])

	// Draw the image
	screen.DrawImage(s.Image, opts)

	// Draw selection indicator if selected
	if s.IsSelected {
		// Draw a circle around the ship
		radius := float32(math.Max(float64(width), float64(height)))/2 + 10
		vector.StrokeCircle(screen, float32(s.Pos[0]+centerX), float32(s.Pos[1]+centerY),
			radius, 1, selectionColor, false)
	}
}
