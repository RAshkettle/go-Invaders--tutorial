package main

import (
	"invaders/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type BaseBlock struct {
	Sprite      *ebiten.Image
	X           int
	Y           int
	DamageLevel int // 0 = no damage, 1-2 = damaged, 3 = destroyed (removed)
	Exists      bool
}

type Base struct {
	Blocks []*BaseBlock
	X      int
	Y      int
}

func NewBaseBlock(x, y int) *BaseBlock {
	return &BaseBlock{
		Sprite:      assets.BaseSprites[0], // Start with no damage frame
		X:           x,
		Y:           y,
		DamageLevel: 0,
		Exists:      true,
	}
}

func NewBase(baseX, baseY int) *Base {
	base := &Base{
		Blocks: make([]*BaseBlock, 0),
		X:      baseX,
		Y:      baseY,
	}

	// Create 4x4 grid of blocks (8x8 pixels each, scaled down 50%)
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			blockX := baseX + (col * 8)
			blockY := baseY + (row * 8)

			// Skip bottom center blocks (archway effect)
			if row == 3 && (col == 1 || col == 2) {
				continue
			}

			block := NewBaseBlock(blockX, blockY)
			base.Blocks = append(base.Blocks, block)
		}
	}

	return base
}

func (b *BaseBlock) TakeDamage() {
	if !b.Exists {
		return
	}

	b.DamageLevel++

	if b.DamageLevel >= 3 {
		b.Exists = false
	} else {
		// Update sprite to show damage (0=first sprite, 1=second, 2=third)
		b.Sprite = assets.BaseSprites[b.DamageLevel]
	}
}

func CreateBases(playerY int) []*Base {
	bases := make([]*Base, 4)

	// Calculate base positioning
	baseWidth := 4 * 8 // 4 blocks * 8 pixels each (scaled down 50%)
	screenWidth := 320
	spacing := (screenWidth - (4 * baseWidth)) / 5 // Equal spacing between and around bases

	baseY := playerY - 8 - (4 * 8) // 8 pixels above player, minus base height

	for i := 0; i < 4; i++ {
		baseX := spacing + (i * (baseWidth + spacing))
		bases[i] = NewBase(baseX, baseY)
	}

	return bases
}
