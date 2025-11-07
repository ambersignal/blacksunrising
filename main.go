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

	// Create multiple ships in random positions with random velocities
	numShips := rand.Intn(2) + 3 // Create 3-5 ships
	for i := 0; i < numShips; i++ {
		// Generate random position within screen bounds
		pos := GenerateRandomPosition(width, height, shipImg.Bounds().Dx(), shipImg.Bounds().Dy())

		// Generate random velocity
		velocity := GenerateRandomVelocity()

		// Create ship with random position and velocity
		ship := NewShip(pos, velocity, shipImg)
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

// GenerateRandomPosition creates a random position within screen bounds
func GenerateRandomPosition(screenWidth, screenHeight, imgWidth, imgHeight int) geom.Vec2 {
	x := rand.Intn(screenWidth-imgWidth) + imgWidth/2
	y := rand.Intn(screenHeight-imgHeight) + imgHeight/2
	return geom.FromPoint(image.Pt(x, y))
}

// GenerateRandomVelocity creates a random velocity vector with magnitude between 0 and 50 pixels per second
func GenerateRandomVelocity() geom.Vec2 {
	// Generate random velocity between 0 and 50 pixels per second
	velMagnitude := 0.5*rand.Float64() + 0.01

	// Generate random direction for velocity vector
	velAngle := rand.Float64() * 2 * math.Pi
	return geom.Vec2{
		velMagnitude * math.Cos(velAngle),
		velMagnitude * math.Sin(velAngle),
	}
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(width*mul, height*mul)
	ebiten.SetWindowTitle("Black Sun Rising")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
