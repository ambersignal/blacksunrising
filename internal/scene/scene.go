package scene

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/ambersignal/blacksunrising/internal/scene/state"
	"github.com/ambersignal/blacksunrising/pkg/draw"
	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	SelectionColor = color.RGBA{30, 188, 115, 255}
)

const (
	AsteroidFieldMinRadius = 100
	AsteroidFieldMaxRadius = 300
	AsteroidFieldMinCount  = 10
	AsteroidFieldMaxCount  = 20
)

// Scene represents the main game scene
type Scene struct {
	state        *state.State
	inputHandler *InputHandler
	minimap      *MiniMap
	nebula       *NebulaBackground
	planet       *Planet

	startTime   time.Time
	updatedTime time.Time
}

// NewScene creates a new game scene
func NewScene(worldSize geom.Vec2, cameraSize geom.Vec2) (*Scene, error) {
	st := state.NewState()
	st.Camera = geom.Rectangle{
		Max: cameraSize,
	}
	st.WorldSize = worldSize

	// FIXME(evgenii.omelchenko): hardcoded miniMapSizes
	miniMapPos := geom.Vec2{
		cameraSize[0] - 110,
		10,
	}

	st.MiniMap = geom.Rectangle{
		Min: miniMapPos,
		Max: miniMapPos.Add(geom.Vec2{100, 100}),
	}

	// Create the minimap
	minimap, err := NewMiniMap(st)
	if err != nil {
		return nil, fmt.Errorf("failed to load minimap: %w", err)
	}

	scene := &Scene{
		startTime:   time.Now(),
		updatedTime: time.Now(),

		state:   st,
		minimap: minimap,
	}

	// Initialize input handler
	scene.inputHandler = NewInputHandler(st)

	// Load the ship image
	shipImg, err := LoadShipImage()
	if err != nil {
		return nil, err
	}

	// Load the asteroid image
	asteroidImg, err := LoadAsteroidImage()
	if err != nil {
		return nil, err
	}

	// Load the nebula background
	nebulaBg, err := NewNebulaBackground()
	if err != nil {
		return nil, fmt.Errorf("failed to load nebula background: %w", err)
	}
	scene.nebula = nebulaBg

	// Create planet in the bottom-right corner of the map
	//	planetPos := geom.Vec2{worldSize[0] - 600, worldSize[1] - 600}
	planetPos := geom.Vec2{100, 100}
	st.Planet = state.NewPlanet(planetPos)

	planetRenderer, err := NewPlanet()
	if err != nil {
		return nil, fmt.Errorf("failed to load planet: %w", err)
	}
	scene.planet = planetRenderer

	// Create multiple ships in random positions with random velocities
	numShips := rand.Intn(11) + 10 // Create 10-20 ships
	for i := 0; i < numShips; i++ {
		// Generate random position within world bounds
		pos := generateRandomPosition(worldSize,
			geom.FromPoint(shipImg.Bounds().Min))

		// Generate random velocity
		velocity := GenerateRandomVelocity()

		// Create ship with random position and velocity
		ship := state.NewShip(pos, velocity, shipImg)
		scene.state.AddShip(ship)
	}

	// Create asteroid fields using two-step generation
	// Step 1: Generate asteroid fields at random positions (empty)
	numAsteroidFields := rand.Intn(3) + 2 // Create 2-4 asteroid fields
	for range numAsteroidFields {
		pos := generateRandomPosition(worldSize, geom.FromPoint(asteroidImg.Bounds().Size()))
		radius := AsteroidFieldMinRadius + rand.Float64()*(AsteroidFieldMaxRadius-AsteroidFieldMinRadius)
		field := state.NewAsteroidField(pos, radius)

		st.AddAsteroidField(field)
	}

	// Step 2: Populate each field with asteroids around its center
	for _, field := range st.AsteroidFields {
		// Generate random number of asteroids (10-20)
		numAsteroids := rand.Intn(AsteroidFieldMaxCount-AsteroidFieldMinCount+1) + AsteroidFieldMinCount

		// Generate random field radius
		for range numAsteroids {
			pos := generateAsteroidPosition(worldSize, field)
			asteroid := state.NewAsteroid(pos, asteroidImg)

			st.AddAsteroid(asteroid)
			field.Asteroids = append(field.Asteroids, asteroid)
		}
	}

	return scene, nil
}

// Update updates the game logic
func (g *Scene) Update() error {
	// Calculate elapsed time since last frame
	elapsedTime := time.Since(g.updatedTime)
	g.updatedTime = time.Now() // Update for next frame

	// Handle input
	if err := g.inputHandler.Update(); err != nil {
		return err
	}

	// Sync ship selection state
	for _, ship := range g.state.Ships {
		_, ship.IsSelected = g.state.Selected[ship]
	}

	// Clean up empty groups periodically
	g.state.CleanupEmptyGroups()

	// Apply steering behaviors
	for _, ship := range g.state.Ships {
		// Calculate steering forces
		alignForce := g.Alignment(ship)
		separateForce := g.Separation(ship)
		cohesionForce := g.Cohesion(ship)

		// If we have a target and ship is in a group, add a seek force
		var seekForce geom.Vec2
		group := g.state.GetGroupForShip(ship)
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

	// Update all asteroids
	for _, asteroid := range g.state.Asteroids {
		asteroid.Update(elapsedTime)
	}

	return nil
}

// Draw renders the game screen
func (g *Scene) Draw(screen *ebiten.Image) {
	elapsed := time.Since(g.startTime)
	// Draw nebula background first
	if g.nebula != nil {
		g.nebula.Draw(screen, g.state.Camera.Min, float32(elapsed.Seconds()))
	}

	// Draw planet
	if g.planet != nil && g.state.Planet != nil {
		elapsed := time.Since(g.startTime)
		g.planet.Draw(screen, g.state.Planet, g.state.Camera.Min, float64(elapsed.Seconds()))
	}

	// Draw all asteroids that are within the camera view
	for _, asteroid := range g.state.Asteroids {
		if g.isAsteroidInView(asteroid, g.state.Camera) {
			asteroid.Draw(screen, g.state.Camera.Min)
		}
	}

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

		// Create a simple rectangle visualization using draw.StrokeLine
		draw.StrokeRect(screen, selection, 1, SelectionColor)
	}

	// Draw minimap in the top-right corner
	g.minimap.Draw(screen, float32(elapsed.Seconds()))
}

// isShipInView checks if a ship is within the camera view
func (g *Scene) isShipInView(ship *state.Ship, cameraRect geom.Rectangle) bool {
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

// isAsteroidInView checks if an asteroid is within the camera view
func (g *Scene) isAsteroidInView(asteroid *state.Asteroid, cameraRect geom.Rectangle) bool {
	if asteroid.Animation == nil || asteroid.Animation.Image == nil {
		return false
	}

	// Calculate the radius of the asteroid for circular collision detection
	radius := asteroid.Radius()

	// Check if asteroid's circular boundary intersects with camera rectangle
	return cameraRect.IntersectsCircle(asteroid.Pos, radius)
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

// generateAsteroidPosition around the field center
func generateAsteroidPosition(worldSize geom.Vec2, field *state.AsteroidField) geom.Vec2 {
	// Generate random direction and distance from field center
	angle := rand.Float64() * 2 * math.Pi
	distance := rand.Float64() * field.Radius

	// Calculate offset from center
	offset := geom.Vec2{
		distance * math.Cos(angle),
		distance * math.Sin(angle),
	}

	// Position relative to field center
	pos := field.Center.Add(offset)

	return pos
}
