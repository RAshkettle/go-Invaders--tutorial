package main

import (
	"invaders/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type AlienType int

const (
	SquidAlien AlienType = iota
	ArmAlien
	FootAlien
)

type Alien struct {
	Sprite       []*ebiten.Image
	X            int
	Y            int
	PointsValue  int
	AlienType    AlienType
	CurrentFrame int
}

func NewAlien(a AlienType) *Alien {
	return &Alien{
		Sprite:       GetAlienSpriteByType(a),
		PointsValue:  40,
		AlienType:    a,
		CurrentFrame: 0,
	}
}

func GetAlienSpriteByType(a AlienType) []*ebiten.Image {
	switch a {
	case SquidAlien:
		return assets.TopInvaderAnimation
	case ArmAlien:
		return assets.MiddleInvaderAnimation
	default:
		return assets.BottomInvaderAnimation

	}
}

func (a *Alien) Update() error {
	return nil
}

// ToggleFrame switches between animation frames (0 and 1)
func (a *Alien) ToggleFrame() {
	a.CurrentFrame = (a.CurrentFrame + 1) % 2
}
