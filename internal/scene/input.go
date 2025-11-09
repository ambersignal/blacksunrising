package scene

import (
	"math"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandler handles all input-related logic for the game
type InputHandler struct {
	state      *State
	isDragging bool
	dragStart  geom.Vec2
	dragEnd    geom.Vec2
}

// NewInputHandler creates a new input handler
func NewInputHandler(state *State) *InputHandler {
	return &InputHandler{
		state: state,
	}
}

// Update processes input for the current frame
func (ih *InputHandler) Update() {
	// Handle right mouse button for selection
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		cursorPos := geom.Vec2{float64(x), float64(y)}

		if !ih.isDragging {
			// Start dragging
			ih.isDragging = true
			ih.dragStart = cursorPos
			ih.dragEnd = cursorPos
		} else {
			// Continue dragging
			ih.dragEnd = cursorPos
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
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ih.state.CurrentGroupIndex >= 0 {
		x, y := ebiten.CursorPosition()
		target := geom.Vec2{float64(x), float64(y)}

		// Set target for current group only
		if ih.state.CurrentGroupIndex < len(ih.state.Groups) {
			ih.state.Groups[ih.state.CurrentGroupIndex].Target = target
			ih.state.Groups[ih.state.CurrentGroupIndex].HasTarget = true
		}
	}
}

// processSelection handles the creation of new groups based on selection
func (ih *InputHandler) processSelection() {
	// Clear current selection
	for ship := range ih.state.Selected {
		delete(ih.state.Selected, ship)
	}

	// Determine bounding box of drag area
	minX := math.Min(ih.dragStart[0], ih.dragEnd[0])
	maxX := math.Max(ih.dragStart[0], ih.dragEnd[0])
	minY := math.Min(ih.dragStart[1], ih.dragEnd[1])
	maxY := math.Max(ih.dragStart[1], ih.dragEnd[1])

	// Remove selected ships from any existing groups first

	// Select ships within the drag area
	for _, ship := range ih.state.Ships {
		if ship.Pos[0] >= minX && ship.Pos[0] <= maxX &&
			ship.Pos[1] >= minY && ship.Pos[1] <= maxY {
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

// DragStart returns the starting position of the drag
func (ih *InputHandler) DragStart() geom.Vec2 {
	return ih.dragStart
}

// DragEnd returns the ending position of the drag
func (ih *InputHandler) DragEnd() geom.Vec2 {
	return ih.dragEnd
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
