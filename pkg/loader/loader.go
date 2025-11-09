package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Loader struct {
	imgCache  map[string]*ebiten.Image
	fontCache map[string]*text.GoTextFaceSource
	dataPath  string
}

func NewLoader(dataPath string) *Loader {
	return &Loader{
		dataPath:  dataPath,
		imgCache:  make(map[string]*ebiten.Image),
		fontCache: make(map[string]*text.GoTextFaceSource),
	}
}

func (l *Loader) LoadImage(path string) (*ebiten.Image, error) {
	img, ok := l.imgCache[path]
	if !ok {
		var err error
		img, _, err = ebitenutil.NewImageFromFile(l.path(path))
		if err != nil {
			return nil, fmt.Errorf("load %q image: %w", path, err)
		}

		l.imgCache[path] = img
	}

	return img, nil
}

func (l *Loader) path(path string) string {
	return filepath.Join(l.dataPath, path)
}

func (l *Loader) LoadFont(path string, size float64) (text.Face, error) {
	source, ok := l.fontCache[path]
	if !ok {
		file, err := os.Open(l.path(path))
		if err != nil {
			return nil, fmt.Errorf("open %q file: %w", path, err)
		}

		source, err = text.NewGoTextFaceSource(file)
		if err != nil {
			return nil, fmt.Errorf("load %q font: %w", path, err)
		}

		l.fontCache[path] = source
	}

	return &text.GoTextFace{
		Source: source,
		Size:   size,
	}, nil
}
