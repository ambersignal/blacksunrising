package main

import (
	"time"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Ship represents a ship in the game
type Ship struct {
	Pos   geom.Vec2
	Image *ebiten.Image
	Vel   geom.Vec2 // Velocity in pixels per second
}

// NewShip creates a new ship with the specified position, velocity, and image
func NewShip(pos geom.Vec2, vel geom.Vec2, img *ebiten.Image) *Ship {
	return &Ship{
		Pos:   pos,
		Image: img,
		Vel:   vel,
	}
}

// LoadShipImage loads the ship image from file
func LoadShipImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("data/ship.png")
	return img, err
}

// Update updates the ship's state
func (s *Ship) Update(elapsedTime time.Duration) error {
	// Update position based on velocity
	s.Pos = s.Pos.Add(s.Vel.Mul(elapsedTime.Seconds()))

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
	angle := s.Vel.Angle()

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
}
