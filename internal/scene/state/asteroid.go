package state

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ambersignal/blacksunrising/pkg/animation"
	"github.com/ambersignal/blacksunrising/pkg/geom"
)

var (
	AsteroidColor = color.RGBA{139, 69, 19, 255} // Brown for minimap
)

const (
	AsteroidSpriteSize   = 24
	AsteroidSpriteCount  = 24
	AsteroidSpriteCols   = 6 // Spritesheet columns (6x4 grid = 24 sprites)
	AsteroidMinCycleTime = 1 * time.Second
	AsteroidMaxCycleTime = 5 * time.Second
)

// Asteroid represents an asteroid in the game.
type Asteroid struct {
	Pos       geom.Vec2
	Animation *animation.Animation
}

// NewAsteroid creates a new asteroid with the specified position and spritesheet.
func NewAsteroid(pos geom.Vec2, img *ebiten.Image) *Asteroid {
	return &Asteroid{
		Pos: pos,
		Animation: animation.NewWithRandomSpeed(
			img,
			AsteroidSpriteSize,
			AsteroidSpriteCount,
			AsteroidSpriteCols,
			AsteroidMinCycleTime,
			AsteroidMaxCycleTime,
		),
	}
}

// Update updates the asteroid's animation state.
func (a *Asteroid) Update(elapsedTime time.Duration) {
	a.Animation.Update(elapsedTime)
}

// Draw renders the asteroid with an optional camera offset.
func (a *Asteroid) Draw(screen *ebiten.Image, cameraOffset ...geom.Vec2) {
	drawPos := a.Pos
	if len(cameraOffset) > 0 {
		drawPos = a.Pos.Sub(cameraOffset[0])
	}

	drawX, drawY := drawPos.Round().Unpack()
	a.Animation.Draw(screen, drawX, drawY)
}

// Radius returns the radius of the asteroid for collision detection.
func (a *Asteroid) Radius() float64 {
	return float64(AsteroidSpriteSize) / 2
}

// AsteroidField represents a clustered field of asteroids.
type AsteroidField struct {
	Center    geom.Vec2
	Radius    float64
	Asteroids []*Asteroid
}

// NewAsteroidField creates a new empty asteroid field at the specified position.
func NewAsteroidField(center geom.Vec2, radius float64) *AsteroidField {
	return &AsteroidField{
		Center:    center,
		Radius:    radius,
		Asteroids: make([]*Asteroid, 0),
	}
}
