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
	Ships             []*Ship
	Groups            [][]*Ship      // Multiple groups of ships
	Targets           []geom.Vec2    // Targets for each group
	HasTarget         []bool         // Whether each group has a target
	Selected          map[*Ship]bool // Map of selected ships
	CurrentGroupIndex int            // Index of the current group being formed
	IsDragging        bool           // Whether we're currently dragging for selection
	DragStart         geom.Vec2      // Starting position of drag selection
	DragEnd           geom.Vec2      // Ending position of drag selection

	startTime time.Time
}

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		startTime:         time.Now(),
		Selected:          make(map[*Ship]bool),
		Groups:            make([][]*Ship, 0),
		Targets:           make([]geom.Vec2, 0),
		HasTarget:         make([]bool, 0),
		CurrentGroupIndex: -1, // No group selected initially
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
			// Add selected ships to a new group
			g.updateSelection()

			// Remove selected ships from any existing groups
			g.removeShipsFromAllGroups()

			// Create a new group with selected ships
			newGroup := make([]*Ship, 0, len(g.Selected))
			for ship := range g.Selected {
				newGroup = append(newGroup, ship)
			}

			// Add the new group to our groups slice
			g.Groups = append(g.Groups, newGroup)
			g.Targets = append(g.Targets, geom.Vec2{0, 0})
			g.HasTarget = append(g.HasTarget, false)

			// Set this as the current group
			g.CurrentGroupIndex = len(g.Groups) - 1

			g.IsDragging = false
		}
	}

	// Sync ship selection state
	for _, ship := range g.Ships {
		ship.IsSelected = g.Selected[ship]
	}

	// Clean up empty groups periodically
	g.cleanupEmptyGroups()

	// Check for left mouse click to set target position for the current group
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.CurrentGroupIndex >= 0 {
		x, y := ebiten.CursorPosition()
		target := geom.Vec2{float64(x), float64(y)}

		// Set target for current group only
		if g.CurrentGroupIndex < len(g.Targets) {
			g.Targets[g.CurrentGroupIndex] = target
			g.HasTarget[g.CurrentGroupIndex] = true
		}
	}

	// Apply steering behaviors
	for _, ship := range g.Ships {
		// Calculate steering forces
		alignForce := g.Alignment(ship)
		separateForce := g.Separation(ship)
		cohesionForce := g.Cohesion(ship)

		// If we have a target and ship is in a group, add a seek force
		var seekForce geom.Vec2
		groupIndex := g.getGroupIndexForShip(ship)
		if groupIndex >= 0 && groupIndex < len(g.HasTarget) && g.HasTarget[groupIndex] {
			seekForce = g.Seek(ship, g.Targets[groupIndex])
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
	screen.Fill(backgroundColor)
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

// isInAnyGroup checks if a ship is part of any group
func (g *Game) isInAnyGroup(ship *Ship) bool {
	for _, group := range g.Groups {
		for _, groupShip := range group {
			if groupShip == ship {
				return true
			}
		}
	}
	return false
}

// getGroupIndexForShip returns the index of the group that contains the ship, or -1 if not in any group
func (g *Game) getGroupIndexForShip(ship *Ship) int {
	for i, group := range g.Groups {
		for _, groupShip := range group {
			if groupShip == ship {
				return i
			}
		}
	}
	return -1
}

// removeShipsFromAllGroups removes selected ships from all existing groups
func (g *Game) removeShipsFromAllGroups() {
	// For each selected ship, remove it from any group it might be in
	for ship := range g.Selected {
		for i := range g.Groups {
			// Create a new slice without the selected ship
			newGroup := make([]*Ship, 0, len(g.Groups[i]))
			for _, groupShip := range g.Groups[i] {
				if groupShip != ship {
					newGroup = append(newGroup, groupShip)
				}
			}
			g.Groups[i] = newGroup
		}
	}
}

// cleanupEmptyGroups removes groups that have no ships
func (g *Game) cleanupEmptyGroups() {
	// Iterate backwards to safely remove elements
	for i := len(g.Groups) - 1; i >= 0; i-- {
		if len(g.Groups[i]) == 0 {
			// Remove the group
			g.Groups = append(g.Groups[:i], g.Groups[i+1:]...)
			g.Targets = append(g.Targets[:i], g.Targets[i+1:]...)
			g.HasTarget = append(g.HasTarget[:i], g.HasTarget[i+1:]...)

			// Adjust current group index if needed
			if g.CurrentGroupIndex >= i && g.CurrentGroupIndex > 0 {
				g.CurrentGroupIndex--
			} else if g.CurrentGroupIndex >= len(g.Groups) {
				g.CurrentGroupIndex = len(g.Groups) - 1
			}
		}
	}

	// If no groups left, reset current group index
	if len(g.Groups) == 0 {
		g.CurrentGroupIndex = -1
	}
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

	// ebiten.SetWindowSize(width*mul, height*mul)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Black Sun Rising")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
