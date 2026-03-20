package data

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.kage
var Shaders embed.FS

func LoadShader(name string) (*ebiten.Shader, error) {
	data, err := Shaders.ReadFile(name + ".kage")
	if err != nil {
		return nil, fmt.Errorf("read shader %q: %w", name, err)
	}

	shader, err := ebiten.NewShader(data)
	if err != nil {
		return nil, fmt.Errorf("compile shader %q: %w", name, err)
	}

	return shader, nil
}
