package scene

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	SelectionColor = color.RGBA{30, 188, 115, 255}
)

// Scene represents the main game scene
type Scene struct {
	state        *State
	inputHandler *InputHandler
	minimap      *MiniMap

	startTime time.Time
}

// NewScene creates a new game scene
func NewScene(worldSize geom.Vec2, cameraSize geom.Vec2) (*Scene, error) {
	state := NewState()
	state.Camera = geom.Rectangle{
		Max: cameraSize,
	}
	state.WorldSize = worldSize

	scene := &Scene{
		startTime: time.Now(),
		state:     state,
		minimap:   NewMiniMap(),
	}

	// Initialize input handler
	scene.inputHandler = NewInputHandler(state, scene.minimap)

	// Load the ship image
	shipImg, err := LoadShipImage()
	if err != nil {
		return nil, err
	}

	// Load the minimap image
	minimapImg, err := LoadMiniMapImage()
	if err != nil {
		return nil, fmt.Errorf("failed to load minimap image: %w", err)
	}
	scene.minimap.Image = minimapImg

	// Create multiple ships in random positions with random velocities
	numShips := rand.Intn(11) + 10 // Create 10-20 ships
	for i := 0; i < numShips; i++ {
		// Generate random position within world bounds
		pos := generateRandomPosition(worldSize,
			geom.FromPoint(shipImg.Bounds().Min))

		// Generate random velocity
		velocity := GenerateRandomVelocity()

		// Create ship with random position and velocity
		ship := NewShip(pos, velocity, shipImg)
		scene.state.AddShip(ship)
	}

	return scene, nil
}

// Update updates the game logic
func (g *Scene) Update() error {
	// Calculate elapsed time since last frame
	elapsedTime := time.Since(g.startTime)
	g.startTime = time.Now() // Update for next frame

	// Handle input
	if err := g.inputHandler.Update(); err != nil {
		return err
	}

	// Sync ship selection state
	for _, ship := range g.state.Ships {
		_, ship.IsSelected = g.state.Selected[ship]
	}

	// Clean up empty groups periodically
	g.cleanupEmptyGroups()

	// Apply steering behaviors
	for _, ship := range g.state.Ships {
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
	for _, ship := range g.state.Ships {
		if err := ship.Update(elapsedTime); err != nil {
			return err
		}
	}

	return nil
}

// Draw renders the game screen
func (g *Scene) Draw(screen *ebiten.Image) {
	// Draw all ships that are within the camera view
	for _, ship := range g.state.Ships {
		// Check if ship is within camera view
		if g.isShipInView(ship, g.state.Camera) {
			// Draw ship with camera offset
			ship.Draw(screen, g.state.Camera.Min)
		}
	}

	// Draw selection rectangle if dragging (with camera adjustment)
	if g.inputHandler.IsDragging() {
		selection := g.inputHandler.selection.Sub(g.state.Camera.Min)

		// Create a simple rectangle visualization using vector.StrokeLine
		vector.StrokeRect(screen, float32(selection.Min[0]), float32(selection.Min[1]),
			float32(selection.Size()[0]), float32(selection.Size()[1]),
			1, SelectionColor, false)
	}

	// Draw minimap in the top-right corner
	g.minimap.Draw(screen, g.state)
}

// isShipInView checks if a ship is within the camera view
func (g *Scene) isShipInView(ship *Ship, cameraRect geom.Rectangle) bool {
	if ship.Image == nil {
		return false
	}

	// Get ship size
	bounds := ship.Image.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// Calculate the radius of the ship for circular collision detection
	radius := math.Max(width, height) / 2

	// Check if ship's circular boundary intersects with camera rectangle
	return cameraRect.IntersectsCircle(ship.Pos, float64(radius))
}

// generateRandomPosition creates a random position within world bounds
func generateRandomPosition(worldSize geom.Vec2, shipSize geom.Vec2) geom.Vec2 {
	return worldSize.Sub(shipSize).
		HadamardProduct(geom.RandVec2()).
		Add(shipSize.Mul(0.5))
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
func (g *Scene) getGroupForShip(ship *Ship) *Group {
	for _, group := range g.state.Groups {
		if group.Contains(ship) {
			return group
		}
	}
	return nil
}

// cleanupEmptyGroups removes groups that have no ships
func (g *Scene) cleanupEmptyGroups() {
	// Iterate backwards to safely remove elements
	for i := len(g.state.Groups) - 1; i >= 0; i-- {
		if g.state.Groups[i].IsEmpty() {
			// Remove the group
			g.state.Groups = append(g.state.Groups[:i], g.state.Groups[i+1:]...)

			// Adjust current group index if needed
			if g.state.CurrentGroupIndex >= i && g.state.CurrentGroupIndex > 0 {
				g.state.CurrentGroupIndex--
			} else if g.state.CurrentGroupIndex >= len(g.state.Groups) {
				g.state.CurrentGroupIndex = len(g.state.Groups) - 1
			}
		}
	}

	// If no groups left, reset current group index
	if len(g.state.Groups) == 0 {
		g.state.CurrentGroupIndex = -1
	}
}

func SmoothStep(edge0, edge1, x float64) float64 {
	t := Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)

	return t * t * (3 - 2*t)
}

func Clamp(x, low, high float64) float64 {
	if x < low {
		return low
	}

	if x > high {
		return high
	}

	return x
}
