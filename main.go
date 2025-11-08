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

const (
	width  = 640
	height = 360
	mul    = 2
)

// Game represents the main game structure
type Game struct {
	Ships      []*Ship
	Target     geom.Vec2      // Target position for ships to move to
	HasTarget  bool           // Whether a target has been set
	Selected   map[*Ship]bool // Map of selected ships
	Group      []*Ship        // Group of ships for collective movement
	IsDragging bool           // Whether we're currently dragging for selection
	DragStart  geom.Vec2      // Starting position of drag selection
	DragEnd    geom.Vec2      // Ending position of drag selection

	startTime time.Time
}

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		startTime: time.Now(),
		Selected:  make(map[*Ship]bool),
		Group:     make([]*Ship, 0),
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

	// Handle right mouse button for selection
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		cursorPos := geom.Vec2{float64(x), float64(y)}

		if !g.IsDragging {
			// Start dragging
			g.IsDragging = true
			g.DragStart = cursorPos
			g.DragEnd = cursorPos
		} else {
			// Continue dragging
			g.DragEnd = cursorPos
		}
	} else {
		// Right mouse button released
		if g.IsDragging {
			// Add selected ships to group
			g.updateSelection()
			g.Group = make([]*Ship, 0, len(g.Selected))
			for ship := range g.Selected {
				g.Group = append(g.Group, ship)
			}
			g.IsDragging = false
		}
	}

	// Sync ship selection state
	for _, ship := range g.Ships {
		ship.IsSelected = g.Selected[ship]
	}

	// Check for left mouse click to set target position for the group
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.Target = geom.Vec2{float64(x), float64(y)}
		g.HasTarget = true
	}

	// Apply steering behaviors
	for _, ship := range g.Ships {
		// Calculate steering forces
		alignForce := g.Alignment(ship)
		separateForce := g.Separation(ship)
		cohesionForce := g.Cohesion(ship)

		// If we have a target and ship is in the group, add a seek force
		var seekForce geom.Vec2
		if g.HasTarget && g.isInGroup(ship) {
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

	// Draw selection rectangle if dragging
	if g.IsDragging {
		// Draw a rectangle from DragStart to DragEnd
		minX := math.Min(g.DragStart[0], g.DragEnd[0])
		maxX := math.Max(g.DragStart[0], g.DragEnd[0])
		minY := math.Min(g.DragStart[1], g.DragEnd[1])
		maxY := math.Max(g.DragStart[1], g.DragEnd[1])

		// Create a simple rectangle visualization using vector.StrokeLine
		vector.StrokeLine(screen, float32(minX), float32(minY), float32(maxX), float32(minY), 2, color.RGBA{0, 255, 0, 255}, false) // Top
		vector.StrokeLine(screen, float32(minX), float32(maxY), float32(maxX), float32(maxY), 2, color.RGBA{0, 255, 0, 255}, false) // Bottom
		vector.StrokeLine(screen, float32(minX), float32(minY), float32(minX), float32(maxY), 2, color.RGBA{0, 255, 0, 255}, false) // Left
		vector.StrokeLine(screen, float32(maxX), float32(minY), float32(maxX), float32(maxY), 2, color.RGBA{0, 255, 0, 255}, false) // Right
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

// isInGroup checks if a ship is part of the current group
func (g *Game) isInGroup(ship *Ship) bool {
	for _, groupShip := range g.Group {
		if groupShip == ship {
			return true
		}
	}
	return false
}

// updateSelection updates the selected ships based on the drag area
func (g *Game) updateSelection() {
	// Clear current selection
	for ship := range g.Selected {
		delete(g.Selected, ship)
	}

	// Determine bounding box of drag area
	minX := math.Min(g.DragStart[0], g.DragEnd[0])
	maxX := math.Max(g.DragStart[0], g.DragEnd[0])
	minY := math.Min(g.DragStart[1], g.DragEnd[1])
	maxY := math.Max(g.DragStart[1], g.DragEnd[1])

	// Select ships within the drag area
	for _, ship := range g.Ships {
		if ship.Pos[0] >= minX && ship.Pos[0] <= maxX &&
			ship.Pos[1] >= minY && ship.Pos[1] <= maxY {
			g.Selected[ship] = true
		}
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
