package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"os"

	"github.com/ambersignal/blacksunrising/internal/scene"
	"github.com/ambersignal/blacksunrising/pkg/geom"
	"github.com/ambersignal/blacksunrising/pkg/loader"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameState int

const (
	GameStateMainMenu GameState = iota
	GameStateLaunched
	GameStateTerminated
)

const (
	width  = 640
	height = 360
	mul    = 2
)

var (
	BackgroundColor = color.RGBA{15, 13, 14, 255}
)

// Scene represents a game scene
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// Game represents the main game structure that implements ebiten.Game interface
type Game struct {
	ui    *ebitenui.UI
	scene Scene
	state GameState
}

// NewGame creates a new game instance
func NewGame(scene Scene) *Game {
	return &Game{
		scene: scene,
	}
}

// Update updates the game logic
func (g *Game) Update() error {
	switch g.state {
	case GameStateMainMenu:
		g.ui.Update()
	case GameStateLaunched:
		return g.scene.Update()
	case GameStateTerminated:
		return ebiten.Termination
	}

	return nil
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)

	switch g.state {
	case GameStateMainMenu:
		g.ui.Draw(screen)
	case GameStateLaunched:
		g.scene.Draw(screen)
	}
}

// Layout returns the game's screen size
func (g *Game) Layout(w, h int) (int, int) {
	// TODO: Define proper screen dimensions
	return width, width * h / w
}

func run() error {
	scene, err := scene.NewScene(geom.Vec2{2000, 2000}, geom.Vec2{640, 360})
	if err != nil {
		return fmt.Errorf("scene initialization: %w", err)
	}

	game := NewGame(scene)

	loader := loader.NewLoader("./data")
	rndrCtx := RenderContext{
		Loader: loader,
	}

	root := Root{
		Menu{
			MenuButton{
				Text: "New Game",
				OnPress: func(args *widget.ButtonClickedEventArgs) {
					game.state = GameStateLaunched
				},
			},
			menuButton("Load Game"),
			MenuButton{
				Text: "Exit",
				OnPress: func(args *widget.ButtonClickedEventArgs) {
					game.state = GameStateTerminated
				},
			},
		},
	}

	ui, err := root.Build(rndrCtx)
	if err != nil {
		return fmt.Errorf("prepare UI builder: %w", err)
	}
	game.ui = ui

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
