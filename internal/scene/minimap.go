package scene

import (
	"image/color"
	"math"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	MinimapViewportColor = color.RGBA{255, 255, 0, 255}   // Yellow for camera viewport
	MinimapShipColor     = color.RGBA{255, 255, 255, 255} // White for unselected ships
	MinimapSelectedColor = color.RGBA{30, 188, 115, 255}  // Green for selected ships
)

// MiniMap represents the minimap component
type MiniMap struct {
	Size       float64
	BorderSize float64
	Image      *ebiten.Image
	InnerSize  float64
}

// NewMiniMap creates a new minimap with default size
func NewMiniMap() *MiniMap {
	return &MiniMap{
		Size:       50.0,
		BorderSize: 3.0, // 56x56 image with 50x50 inner area = 3px border on each side
		InnerSize:  50.0,
	}
}

// ScreenToWorld converts minimap screen coordinates to world coordinates
func (m *MiniMap) ScreenToWorld(screenPos geom.Vec2, screenSize geom.Vec2, worldSize geom.Vec2) geom.Vec2 {
	// Calculate minimap position (top-right corner)
	minimapPos := geom.Vec2{
		screenSize[0] - m.Size - 10,
		10.0,
	}

	// Calculate the position within the minimap's inner area
	innerPos := screenPos.Sub(minimapPos).Sub(geom.Vec2{m.BorderSize, m.BorderSize})

	// Normalize the position within the inner area (0-1 range)
	normalizedPos := innerPos.HadamardDevide(geom.Vec2{m.InnerSize, m.InnerSize})

	// Clamp to valid range [0, 1]
	normalizedPos[0] = math.Max(0, math.Min(1, normalizedPos[0]))
	normalizedPos[1] = math.Max(0, math.Min(1, normalizedPos[1]))

	// Convert to world coordinates
	return normalizedPos.HadamardProduct(worldSize)
}

// LoadMiniMapImage loads the minimap image from file
func LoadMiniMapImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("data/map.png")
	return img, err
}

// Draw renders the minimap on the screen
func (m *MiniMap) Draw(screen *ebiten.Image, state *State) {
	screenWidth := float64(screen.Bounds().Dx())

	// Position minimap in top-right corner with some padding
	minimapPos := geom.Vec2{
		screenWidth - m.Size - 10,
		10.0,
	}

	// Draw the minimap image
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(minimapPos[0], minimapPos[1])
	screen.DrawImage(m.Image, opts)

	// Draw camera viewport indicator
	// Convert camera rectangle to normalized coordinates
	cameraNormalized := state.Camera.HadamardDivide(state.WorldSize)

	borderSize := geom.Vec2{
		m.BorderSize,
		m.BorderSize,
	}

	cameraMinimap := cameraNormalized.Mul(m.InnerSize).
		Add(borderSize).Add(minimapPos)

	// Draw camera viewport as a transparent rectangle with border
	vector.StrokeRect(screen, float32(cameraMinimap.Min[0]), float32(cameraMinimap.Min[1]),
		float32(cameraMinimap.Size()[0]), float32(cameraMinimap.Size()[1]),
		1, MinimapViewportColor, false) // Yellow border for camera viewport

	// Draw all ships on minimap within the inner area
	for _, ship := range state.Ships {
		// Convert world coordinates to normalized coordinates
		normalizedPos := ship.Pos.HadamardDevide(state.WorldSize).
			Mul(m.InnerSize).Add(borderSize).Add(minimapPos)

		// Choose color based on ship selection
		shipColor := MinimapShipColor // White for unselected
		if ship.IsSelected {
			shipColor = MinimapSelectedColor // Green for selected
		}

		// Draw small dot for ship
		screen.Set(
			int(math.Round(normalizedPos[0])),
			int(math.Round(normalizedPos[1])),
			shipColor)
	}
}
