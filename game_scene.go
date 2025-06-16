package main

import "github.com/hajimehoshi/ebiten/v2"

type gameScene struct {
	sceneManager *SceneManager
}

func (g *gameScene) Update() error { return nil }

func (g *gameScene) Draw(screen *ebiten.Image) {}

func (g *gameScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}

func NewGameScene(sm *SceneManager) *gameScene {
	return &gameScene{
		sceneManager: sm,
	}
}
