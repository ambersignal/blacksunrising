package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// LoadShipImage loads the ship image from file
func LoadShipImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("data/fighter.png")
	return img, err
}

// LoadAsteroidImage loads the asteroid spritesheet from file
func LoadAsteroidImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("data/asteroids1.png")
	return img, err
}
