package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Invaders")
	ebiten.SetWindowSize(640, 480)

	sceneManager := NewSceneManager()

	err := ebiten.RunGame(sceneManager)
	if err != nil {
		panic(err)
	}
}
