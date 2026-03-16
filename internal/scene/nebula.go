package scene

import (
	"github.com/ambersignal/blacksunrising/data"
	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cameraCoeff = 0.001
)

type NebulaBackground struct {
	shader *ebiten.Shader
}

func NewNebulaBackground() (*NebulaBackground, error) {
	shader, err := data.LoadShader("nebula")
	if err != nil {
		return nil, err
	}

	return &NebulaBackground{
		shader: shader,
	}, nil
}

func (n *NebulaBackground) Draw(screen *ebiten.Image, camera geom.Vec2, time float64) {
	bounds := screen.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	camera = camera.Mul(cameraCoeff)

	op := &ebiten.DrawRectShaderOptions{}
	op.GeoM.Translate(0, 0)
	op.Uniforms = map[string]any{
		"Time":           float32(time),
		"ScreenSize":     []float32{float32(width), float32(height)},
		"CameraPosition": []float32{float32(camera[0]), float32(camera[1])},
	}

	screen.DrawRectShader(width, height, n.shader, op)
}

func (n *NebulaBackground) Dispose() {
	if n.shader != nil {
		n.shader.Deallocate()
		n.shader = nil
	}
}
