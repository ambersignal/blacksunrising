package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	selectionColor  = color.RGBA{30, 188, 115, 255}
	backgroundColor = color.RGBA{15, 13, 14, 255}
)

const (
	width  = 640
	height = 360
	mul    = 2
)

// Game represents the main game structure
type Game struct {
	ships             []*Ship
	groups            []*Group       // Groups of ships
	selected          map[*Ship]bool // Map of selected ships
	currentGroupIndex int            // Index of the current group being formed
	inputHandler      *InputHandler  // Handles input logic

	startTime time.Time
}

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		startTime:         time.Now(),
		selected:          make(map[*Ship]bool),
		groups:            make([]*Group, 0),
		currentGroupIndex: -1, // No group selected initially
	}

	// Initialize input handler
	game.inputHandler = NewInputHandler(game)

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
		game.ships = append(game.ships, ship)
	}

	return game
}

// Update updates the game logic
func (g *Game) Update() error {
	// Calculate elapsed time since last frame
	elapsedTime := time.Since(g.startTime)
	g.startTime = time.Now() // Update for next frame

	// Handle input
	g.inputHandler.Update()

	// Sync ship selection state
	for _, ship := range g.ships {
		ship.IsSelected = g.selected[ship]
	}

	// Clean up empty groups periodically
	g.cleanupEmptyGroups()

	// Apply steering behaviors
	for _, ship := range g.ships {
		// Calculate steering forces
		alignForce := g.Alignment(ship)
		separateForce := g.Separation(ship)
		cohesionForce := g.Cohesion(ship)

		// If we have a target and ship is in a group, add a seek force
		var seekForce geom.Vec2
		group := g.getGroupForShip(ship)
		if group != nil && group.HasTarget {
			seekForce = g.Seek(ship, group.Target)
		}

		// Apply forces to ship's acceleration
		// These weights can be adjusted to change behavior
		ship.Accel = ship.Accel.Add(alignForce.Mul(1.0))
		ship.Accel = ship.Accel.Add(separateForce.Mul(2.5))
		ship.Accel = ship.Accel.Add(cohesionForce.Mul(1.0))
		ship.Accel = ship.Accel.Add(seekForce.Mul(2.0)) // Stronger seek force
	}

	// Update all ships
	for _, ship := range g.ships {
		if err := ship.Update(elapsedTime); err != nil {
			return err
		}
	}

	return nil
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	// Draw all ships
	for _, ship := range g.ships {
		ship.Draw(screen)
	}

	// Draw selection rectangle if dragging
	if g.inputHandler.IsDragging() {
		// Draw a rectangle from DragStart to DragEnd
		dragStart := g.inputHandler.DragStart()
		dragEnd := g.inputHandler.DragEnd()
		minX := math.Min(dragStart[0], dragEnd[0])
		maxX := math.Max(dragStart[0], dragEnd[0])
		minY := math.Min(dragStart[1], dragEnd[1])
		maxY := math.Max(dragStart[1], dragEnd[1])

		// Create a simple rectangle visualization using vector.StrokeLine
		vector.StrokeRect(screen, float32(minX), float32(minY),
			float32(maxX-minX), float32(maxY-minY), 1, selectionColor, false)
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

// getGroupForShip returns the group that contains the ship, or nil if not in any group
func (g *Game) getGroupForShip(ship *Ship) *Group {
	for _, group := range g.groups {
		if group.Contains(ship) {
			return group
		}
	}
	return nil
}

// cleanupEmptyGroups removes groups that have no ships
func (g *Game) cleanupEmptyGroups() {
	// Iterate backwards to safely remove elements
	for i := len(g.groups) - 1; i >= 0; i-- {
		if g.groups[i].IsEmpty() {
			// Remove the group
			g.groups = append(g.groups[:i], g.groups[i+1:]...)

			// Adjust current group index if needed
			if g.currentGroupIndex >= i && g.currentGroupIndex > 0 {
				g.currentGroupIndex--
			} else if g.currentGroupIndex >= len(g.groups) {
				g.currentGroupIndex = len(g.groups) - 1
			}
		}
	}

	// If no groups left, reset current group index
	if len(g.groups) == 0 {
		g.currentGroupIndex = -1
	}
}

func main() {
	game := NewGame()

	// ebiten.SetWindowSize(width*mul, height*mul)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Black Sun Rising")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
