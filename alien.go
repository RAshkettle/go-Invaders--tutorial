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
	Sprite      []*ebiten.Image
	X           int
	Y           int
	PointsValue int
	AlienType   AlienType
}

func NewAlien(a AlienType) *Alien {
	return &Alien{
		Sprite: GetAlienSpriteByType(a),

		PointsValue: 40,
		AlienType:   a,
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
