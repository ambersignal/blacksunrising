package main

import (
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 640
	height = 360
	mul    = 2
)

// Game represents the main game structure
type Game struct {
	Ships []*Ship

	startTime time.Time
}

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		startTime: time.Now(),
	}

	// Load the ship image
	shipImg, err := LoadShipImage()
	if err != nil {
		log.Fatal(err)
	}

	// Create multiple ships in random positions with random angles
	numShips := rand.Intn(2) + 3 // Create 3-5 ships
	for i := 0; i < numShips; i++ {
		// Generate random position within screen bounds
		x := rand.Intn(width-shipImg.Bounds().Dx()) + shipImg.Bounds().Dx()/2
		y := rand.Intn(height-shipImg.Bounds().Dy()) + shipImg.Bounds().Dy()/2
		pos := geom.FromPoint(image.Pt(x, y))

		// Generate random angle (in radians)
		angle := geom.Angle(rand.Float64() * 2 * math.Pi)

		// Create ship with random position and angle
		ship := NewShip(pos, angle, shipImg)
		game.Ships = append(game.Ships, ship)
	}

	return game
}

// Update updates the game logic
func (g *Game) Update() error {
	elapsedTime := time.Since(g.startTime)
	// Update all ships
	for _, ship := range g.Ships {
		if err := ship.Update(elapsedTime); err != nil {
			return err
		}
	}

	return nil
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw all ships
	for _, ship := range g.Ships {
		ship.Draw(screen)
	}
}

// Layout returns the game's screen size
func (g *Game) Layout(width, height int) (int, int) {
	// TODO: Define proper screen dimensions
	return width / mul, height / mul
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(width*mul, height*mul)
	ebiten.SetWindowTitle("Black Sun Rising")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
