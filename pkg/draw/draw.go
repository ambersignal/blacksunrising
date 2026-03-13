package draw

import (
	"image/color"

	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func StrokeRect(screen *ebiten.Image, rect geom.Rectangle, border float64, clr color.Color) {
	vector.StrokeRect(screen, float32(rect.Min[0]), float32(rect.Min[1]),
		float32(rect.Size()[0]), float32(rect.Size()[1]),
		float32(border), clr, false)
}

func StrokeCircle(screen *ebiten.Image, center geom.Vec2, radius float64, clr color.Color) {
	vector.StrokeCircle(screen, float32(center[0]),
		float32(center[1]),
		float32(radius), 1, clr, false)
}

func Pixel(screen *ebiten.Image, pos geom.Vec2, clr color.Color) {
	imgPos := pos.ToImagePoint()

	screen.Set(imgPos.X, imgPos.Y, clr)
}
