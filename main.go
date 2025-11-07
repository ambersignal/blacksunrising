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
	Ships     []*Ship
	Target    geom.Vec2 // Target position for ships to move to
	HasTarget bool      // Whether a target has been set

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
	numShips := rand.Intn(5) + 5 // Create 3-5 ships
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
	// Calculate elapsed time since last frame
	elapsedTime := time.Since(g.startTime)
	g.startTime = time.Now() // Update for next frame

	// Check for mouse click to set target position
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.Target = geom.Vec2{float64(x), float64(y)}
		g.HasTarget = true
	}

	// Apply steering behaviors to all ships
	for _, ship := range g.Ships {
		// Calculate steering forces
		alignForce := g.Alignment(ship)
		separateForce := g.Separation(ship)
		cohesionForce := g.Cohesion(ship)

		// If we have a target, add a seek force
		var seekForce geom.Vec2
		if g.HasTarget {
			seekForce = g.Seek(ship, g.Target)
		}

		// Apply forces to ship's acceleration
		// These weights can be adjusted to change behavior
		ship.Accel = ship.Accel.Add(alignForce.Mul(1.0))
		ship.Accel = ship.Accel.Add(separateForce.Mul(2.5))
		ship.Accel = ship.Accel.Add(cohesionForce.Mul(1.0))
		ship.Accel = ship.Accel.Add(seekForce.Mul(2.0)) // Stronger seek force
	}

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
	// Generate random velocity between 10 and 30 pixels per second for better flocking
	velMagnitude := 20*rand.Float64() + 10

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
