package main

import (
	"log"
	"time"

	"github.com/ambersignal/blacksunrising/internal/logic"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 640
	height = 360
	mul    = 2
)

// Logic interface defines the contract for game logic operations
type Logic interface {
	Update(shift time.Duration) error
}

// Game represents the main game structure
type Game struct {
	logic Logic
	prev  time.Time
}

// Update updates the game logic
func (g *Game) Update() error {
	// Delegate to the logic implementation
	next := time.Now()
	if err := g.logic.Update(next.Sub(g.prev)); err != nil {
		return err
	}
	g.prev = next

	return nil
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: Implement rendering logic here
}

// Layout returns the game's screen size
func (g *Game) Layout(width, height int) (int, int) {
	// TODO: Define proper screen dimensions
	return width / mul, height / mul
}

func main() {
	gameLogic := &logic.Game{}
	game := &Game{
		logic: gameLogic,
	}

	// TODO: Configure Ebiten settings
	ebiten.SetWindowSize(width*2, height*mul)
	ebiten.SetWindowTitle("Black Sun Rising")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
