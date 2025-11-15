package scene

import (
	"image"
	"math"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandler handles all input-related logic for the game
type InputHandler struct {
	state      *State
	isDragging bool
	selection  geom.Rectangle
}

// NewInputHandler creates a new input handler
func NewInputHandler(state *State) *InputHandler {
	return &InputHandler{
		state: state,
	}
}

// Update processes input for the current frame
func (ih *InputHandler) Update() error {
	// Check for ESC key to exit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	ih.updateCamera()

	// Handle right mouse button for selection
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorPos := ih.CursorPosition().Add(ih.state.Camera.Min)

		if !ih.isDragging {
			// Start dragging
			ih.isDragging = true
			ih.selection.Min = cursorPos
			ih.selection.Max = cursorPos
		} else {
			// Continue dragging
			ih.selection.Max = cursorPos
			ih.selection = ih.selection.Normalize()
		}
	} else {
		// Right mouse button released
		if ih.isDragging {
			// Process the selection
			ih.processSelection()
			ih.isDragging = false
		}
	}

	// Check for left mouse click to set target position for the current group
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ih.state.CurrentGroupIndex >= 0 {
		cursorPos := ih.CursorPosition()

		// Adjust target for camera offset
		cameraOffset := ih.state.Camera.Min
		worldTarget := cursorPos.Add(cameraOffset)

		// Set target for current group only
		if ih.state.CurrentGroupIndex < len(ih.state.Groups) {
			ih.state.Groups[ih.state.CurrentGroupIndex].Target = worldTarget
			ih.state.Groups[ih.state.CurrentGroupIndex].HasTarget = true
		}
	}

	return nil
}

// updateCamera moves the camera when mouse is near screen bounds
func (ih *InputHandler) updateCamera() {
	cursorPosition := ih.CursorPosition()

	// Define border margin (in pixels) where camera starts moving
	margin := 50.0

	// Camera movement speed (pixels per frame)
	moveSpeed := 2.5

	// For both x and y
	var cameraShift geom.Vec2
	for i := range 2 {
		if cursorPosition[i] < margin {
			cameraShift[i] = -moveSpeed
		} else if cursorPosition[i] > ih.state.Camera.Size()[i]-margin {
			cameraShift[i] = moveSpeed
		}
	}
	ih.state.Camera = ih.state.Camera.Add(cameraShift)

	for i := range 2 {
		if ih.state.Camera.Min[i] < 0 {
			cameraShift[i] = -ih.state.Camera.Min[i]
		} else if ih.state.Camera.Max[i] > ih.state.WorldSize[i] {
			cameraShift[i] = ih.state.Camera.Max[i] - ih.state.WorldSize[i]
		}
	}
	ih.state.Camera = ih.state.Camera.Add(cameraShift)
}

func (ih *InputHandler) CursorPosition() geom.Vec2 {
	return geom.FromPoint(image.Pt(ebiten.CursorPosition()))
}

// processSelection handles the creation of new groups based on selection
func (ih *InputHandler) processSelection() {
	// Clear current selection
	for ship := range ih.state.Selected {
		delete(ih.state.Selected, ship)
	}

	selection := ih.selection.Normalize()

	// Select ships within the drag area
	for _, ship := range ih.state.Ships {
		shipSize := ship.Image.Bounds().Size()
		radius := math.Max(float64(shipSize.X), float64(shipSize.Y))

		if selection.IntersectsCircle(ship.Pos, radius) {
			ih.state.Selected[ship] = struct{}{}
		}
	}

	// Remove selected ships from any existing groups first
	ih.removeFromAllGroups()

	// Create a new group with selected ships
	newGroup := NewGroup()
	for ship := range ih.state.Selected {
		newGroup.AddShip(ship)
	}

	// Add the new group to our groups slice
	ih.state.Groups = append(ih.state.Groups, newGroup)

	// Set this as the current group
	ih.state.CurrentGroupIndex = len(ih.state.Groups) - 1
}

// IsDragging returns whether we're currently dragging for selection
func (ih *InputHandler) IsDragging() bool {
	return ih.isDragging
}

// removeFromAllGroups removes selected ships from all existing groups
func (ih *InputHandler) removeFromAllGroups() {
	// For each selected ship, remove it from any group it might be in
	for ship := range ih.state.Selected {
		for _, group := range ih.state.Groups {
			group.RemoveShip(ship)
		}
	}
}
