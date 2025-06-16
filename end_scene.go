package main

import "github.com/hajimehoshi/ebiten/v2"

type endScene struct {
	sceneManager *SceneManager
}

func (e *endScene) Update() error {
	return nil
}

func (e *endScene) Draw(screen *ebiten.Image) {}

func (e *endScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}

func NewEndScene(sm *SceneManager) *endScene {
	return &endScene{
		sceneManager: sm,
	}
}
