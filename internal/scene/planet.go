package scene

import (
	"github.com/ambersignal/blacksunrising/data"
	"github.com/ambersignal/blacksunrising/internal/scene/state"
	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

type Planet struct {
	shader *ebiten.Shader
}

func NewPlanet() (*Planet, error) {
	shader, err := data.LoadShader("planet")
	if err != nil {
		return nil, err
	}

	return &Planet{
		shader: shader,
	}, nil
}

func (p *Planet) Draw(screen *ebiten.Image, planet *state.Planet, camera geom.Vec2, time float64) {
	width := 512
	height := 512

	op := &ebiten.DrawRectShaderOptions{}
	op.GeoM.Translate(planet.Pos[0], planet.Pos[1])
	op.GeoM.Translate(-camera[0], -camera[1])
	op.Uniforms = map[string]any{
		"Time": float32(time),
		"Size": []float32{float32(width), float32(height)},
	}

	screen.DrawRectShader(width, height, p.shader, op)
}

func (p *Planet) Dispose() {
	if p.shader != nil {
		p.shader.Deallocate()
		p.shader = nil
	}
}
