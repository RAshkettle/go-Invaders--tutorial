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
	PlayerShot = loadImage("player/PlayerShot.png")
	AlienShot = loadImage("invaders/AlienShot.png")

	baseSpriteSheet = loadImage("player/base.png")
	BaseSprites = splitBaseImage(baseSpriteSheet)

	MoveSound = loadAudio("audio/move.ogg")
	PlayerShootSound = loadAudio("audio/laserShoot.ogg")
	AlienExplosionSound = loadAudio("audio/alienexplosion.ogg")
	PlayerDeathSound = loadAudio("audio/playerDeath.ogg")
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

func splitBaseImage(spriteSheet *ebiten.Image) []*ebiten.Image{
	const tileSize = 16
	first := spriteSheet.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image)
	second := spriteSheet.SubImage(image.Rect(tileSize, 0, tileSize*2, tileSize)).(*ebiten.Image)
	third := spriteSheet.SubImage(image.Rect(tileSize, 0, tileSize*3, tileSize)).(*ebiten.Image)
	return []*ebiten.Image{first, second, third}
}
