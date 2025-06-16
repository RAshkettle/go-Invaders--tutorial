package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	sceneManager *SceneManager
	aliens       []*Alien
}

const (
	STEP = 16
)

func (g *GameScene) Update() error { return nil }

func (g *GameScene) Draw(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	scaledWidth := width / 320.0
	scaledHeight := height / 240.0
	scale := math.Min(float64(scaledWidth), float64(scaledHeight))

	// Calculate centering offsets
	gameWidth := 320.0 * scale
	gameHeight := 240.0 * scale
	offsetX := (float64(width) - gameWidth) / 2.0
	offsetY := (float64(height) - gameHeight) / 2.0

	for _, alien := range g.aliens {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Scale(float64(scale), float64(scale))
		op.GeoM.Translate(float64(alien.X)*scale+offsetX, float64(alien.Y)*scale+offsetY)
		screen.DrawImage(alien.Sprite[0], op)
	}
}

func (g *GameScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}

func (g *GameScene) Reset() {}

func NewGameScene(sm *SceneManager) *GameScene {
	return &GameScene{
		sceneManager: sm,
		aliens:       SpawnAlienWave(),
	}
}

func (g *GameScene) SpawnAliens() []*Alien {
	return SpawnAlienWave()
}
