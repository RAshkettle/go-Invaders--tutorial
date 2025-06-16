package assets

import (
	"bytes"
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *
var assets embed.FS

var (
	topInvaderSpriteSheet    = loadImage("invaders/topInvader.png")
	middleInvaderSpriteSheet = loadImage("invaders/middleInvader.png")
	bottomInvaderSpriteSheet = loadImage("invaders/bottomInvader.png")

	TopInvaderAnimation    = splitImage(topInvaderSpriteSheet)
	MiddleInvaderAnimation = splitImage(middleInvaderSpriteSheet)
	BottomInvaderAnimation = splitImage(bottomInvaderSpriteSheet)

	Player = loadImage("player/Player.png")

	MoveSound = loadAudio("audio/move.ogg")
)

func loadImage(filePath string) *ebiten.Image {
	data, err := assets.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	ebitenImg := ebiten.NewImageFromImage(img)
	return ebitenImg
}

func loadAudio(filePath string) []byte {
	data, err := assets.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return data
}

func splitImage(spriteSheet *ebiten.Image) []*ebiten.Image {
	const invaderSize = 16

	firstFrame := spriteSheet.SubImage(image.Rect(0, 0, invaderSize, invaderSize)).(*ebiten.Image)
	secondFrame := spriteSheet.SubImage(image.Rect(invaderSize, 0, invaderSize*2, invaderSize)).(*ebiten.Image)

	return []*ebiten.Image{firstFrame, secondFrame}
}
