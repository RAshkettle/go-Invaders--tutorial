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
const (
	NUMBER_OF_ALIENS_IN_ROW = 12
	ALIEN_SIZE              = 16
	PADDING                 = 64
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
		PointsValue:  getAlienPointsByType(a),
		AlienType:    AlienType(a),
		CurrentFrame: 0,
	}

}

func getAlienPointsByType(a AlienType) int {
	switch a {
	case SquidAlien:
		return 40
	case ArmAlien:
		return 20
	case FootAlien:
		return 10
	default:
		return 10 // Default to FootAlien points
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


// ToggleFrame switches between animation frames (0 and 1)
func (a *Alien) ToggleFrame() {
	a.CurrentFrame = (a.CurrentFrame + 1) % 2
}

func SpawnAlienWave() []*Alien {
	aliens := make([]*Alien, 0)

	for i := range NUMBER_OF_ALIENS_IN_ROW {
		//Make the top row dude
		alien := NewAlien(SquidAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE
		aliens = append(aliens, alien)

		//Make the Middle Row Dudes
		alien = NewAlien(ArmAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 2
		aliens = append(aliens, alien)

		alien = NewAlien(ArmAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 3
		aliens = append(aliens, alien)
		//Make the bottom row dudes
		alien = NewAlien(FootAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 4
		aliens = append(aliens, alien)

		alien = NewAlien(FootAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE * 5
		aliens = append(aliens, alien)
	}
	return aliens
}
