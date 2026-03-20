package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ambersignal/blacksunrising/data"
	"github.com/ambersignal/blacksunrising/internal/scene/state"
	"github.com/ambersignal/blacksunrising/pkg/draw"
)

var (
	MinimapViewportColor = color.RGBA{255, 255, 0, 255}   // Yellow for camera viewport
	MinimapShipColor     = color.RGBA{255, 255, 255, 255} // White for unselected ships
	MinimapSelectedColor = color.RGBA{30, 188, 115, 255}  // Green for selected ships
	MinimapAsteroidColor = color.RGBA{220, 140, 90, 255}  // Brown for asteroids
)

// MiniMap represents the minimap component.
type MiniMap struct {
	state  *state.State
	shader *ebiten.Shader
}

// NewMiniMap creates a new minimap with shader.
func NewMiniMap(state *state.State) (*MiniMap, error) {
	shader, err := data.LoadShader("map")
	if err != nil {
		return nil, err
	}

	return &MiniMap{
		state:  state,
		shader: shader,
	}, nil
}

// Draw renders the minimap on the screen.
func (m *MiniMap) Draw(screen *ebiten.Image, time float32) {
	minimapPos := m.state.MiniMap.Min
	minimapSize := m.state.MiniMap.Size()

	// Draw the minimap shader
	op := &ebiten.DrawRectShaderOptions{}
	op.GeoM.Translate(minimapPos[0], minimapPos[1])

	op.Uniforms = map[string]any{
		"Time": time,
		"Size": minimapSize.AsFloat32Slice(),
	}

	screen.DrawRectShader(int(minimapSize[0]), int(minimapSize[1]), m.shader, op)

	// Draw camera viewport indicator
	// Convert camera rectangle to normalized coordinates
	cameraNormalized := m.state.Camera.HadamardDivide(m.state.WorldSize)

	cameraMinimap := cameraNormalized.HadamardProduct(m.state.MiniMap.Size()).
		Add(m.state.MiniMap.Min)
	draw.StrokeRect(screen, cameraMinimap, 1, MinimapViewportColor)

	// Draw all ships on minimap
	for _, ship := range m.state.Ships {
		// Convert world coordinates to normalized coordinates
		normalizedPos := ship.Pos.HadamardDevide(m.state.WorldSize).
			HadamardProduct(m.state.MiniMap.Size()).Add(m.state.MiniMap.Min)

		// Choose color based on ship selection
		shipColor := MinimapShipColor // White for unselected
		if ship.IsSelected {
			shipColor = MinimapSelectedColor // Green for selected
		}

		// Draw small dot for ship
		draw.Pixel(screen, normalizedPos, shipColor)
	}

	// Draw all asteroids on minimap
	for _, asteroid := range m.state.Asteroids {
		// Convert world coordinates to normalized coordinates
		normalizedPos := asteroid.Pos.HadamardDevide(m.state.WorldSize).
			HadamardProduct(m.state.MiniMap.Size()).Add(m.state.MiniMap.Min)

		// Draw small dot for asteroid using brown color
		draw.Pixel(screen, normalizedPos, MinimapAsteroidColor)
	}
}

// Dispose releases the shader resources.
func (m *MiniMap) Dispose() {
	if m.shader != nil {
		m.shader.Deallocate()
		m.shader = nil
	}
}
