package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ambersignal/blacksunrising/internal/scene"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 640
	height = 360
	mul    = 2
)

// Scene represents a game scene
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// Game represents the main game structure that implements ebiten.Game interface
type Game struct {
	scene Scene
}

// NewGame creates a new game instance
func NewGame(scene Scene) *Game {
	return &Game{
		scene: scene,
	}
}

// Update updates the game logic
func (g *Game) Update() error {
	return g.scene.Update()
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

// Layout returns the game's screen size
func (g *Game) Layout(w, h int) (int, int) {
	// TODO: Define proper screen dimensions
	return width, width * h / w
}

func run() error {
	scene, err := scene.NewScene()
	if err != nil {
		return fmt.Errorf("scene initialization: %w", err)
	}

	game := NewGame(scene)

	// ebiten.SetWindowSize(width*mul, height*mul)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Black Sun Rising")

	if err := ebiten.RunGame(game); err != nil {
		return fmt.Errorf("game run: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("Unable to run the game", "err", err)
		os.Exit(-1)
	}
}
