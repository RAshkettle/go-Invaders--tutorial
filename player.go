package main

import (
	"invaders/assets"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gameWidth  = 320
	gameHeight = 240
	playerSpeed = 2
)

type Player struct {
	Sprite     *ebiten.Image
	X          int
	Y          int
	ShootTimer *stopwatch.Stopwatch
}

func NewPlayer() *Player {
	playerWidth := assets.Player.Bounds().Dx()
	playerHeight := assets.Player.Bounds().Dy()
	return &Player{
		Sprite:     assets.Player,
		X:          (gameWidth - playerWidth) / 2,
		Y:          gameHeight - playerHeight - 8, // 8 pixels from the bottom
		ShootTimer: stopwatch.NewStopwatch(2 * time.Second),
	}
}

func (p *Player) Update() error{
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += playerSpeed
	}

	// Keep player within screen bounds
	playerWidth := p.Sprite.Bounds().Dx()
	if p.X < 0 {
		p.X = 0
	}
	if p.X + playerWidth > gameWidth {
		p.X = gameWidth - playerWidth
	}
	return nil
}




