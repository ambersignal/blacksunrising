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

// LoadMiniMapImage loads the minimap image from file
func LoadMiniMapImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("data/map.png")
	return img, err
}
