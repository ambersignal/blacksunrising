package animation

import (
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Image           *ebiten.Image
	SpriteSize      int
	SpriteCount     int
	Cols            int
	FrameDuration   time.Duration
	currentFrame    int
	animationTimer  float64
}

func New(img *ebiten.Image, spriteSize, spriteCount, cols int, frameDuration time.Duration) *Animation {
	return &Animation{
		Image:          img,
		SpriteSize:     spriteSize,
		SpriteCount:    spriteCount,
		Cols:           cols,
		FrameDuration:  frameDuration,
		currentFrame:   rand.Intn(spriteCount),
		animationTimer: 0,
	}
}

func NewWithRandomSpeed(img *ebiten.Image, spriteSize, spriteCount, cols int, minCycle, maxCycle time.Duration) *Animation {
	minFrameDur := minCycle.Seconds() / float64(spriteCount)
	maxFrameDur := maxCycle.Seconds() / float64(spriteCount)
	frameDuration := time.Duration((minFrameDur + rand.Float64()*(maxFrameDur-minFrameDur)) * float64(time.Second))

	return &Animation{
		Image:          img,
		SpriteSize:     spriteSize,
		SpriteCount:    spriteCount,
		Cols:           cols,
		FrameDuration:  frameDuration,
		currentFrame:   rand.Intn(spriteCount),
		animationTimer: 0,
	}
}

func (a *Animation) Update(elapsedTime time.Duration) {
	deltaSeconds := elapsedTime.Seconds()

	a.animationTimer += deltaSeconds
	if a.animationTimer >= a.FrameDuration.Seconds() {
		a.animationTimer -= a.FrameDuration.Seconds()
		a.currentFrame++
		if a.currentFrame >= a.SpriteCount {
			a.currentFrame = 0
		}
	}
}

func (a *Animation) Draw(screen *ebiten.Image, x, y float64) {
	if a.Image == nil {
		return
	}

	cols := a.Cols
	srcX := (a.currentFrame % cols) * a.SpriteSize
	srcY := (a.currentFrame / cols) * a.SpriteSize

	srcRect := image.Rect(
		srcX, srcY,
		srcX+a.SpriteSize, srcY+a.SpriteSize,
	)

	frameImg := a.Image.SubImage(srcRect).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)

	screen.DrawImage(frameImg, opts)
}

func (a *Animation) CurrentFrame() int {
	return a.currentFrame
}