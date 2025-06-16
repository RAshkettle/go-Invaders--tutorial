package main

import (
	"bytes" // Required for audio decoding
	"invaders/assets"
	"log" // Required for audio error logging
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"        // Required for audio context
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis" // Required for ogg decoding
	"github.com/hajimehoshi/ebiten/v2/inpututil"    // Required for IsKeyJustPressed
)

const (
	gameWidth          = 320
	gameHeight         = 240
	playerSpeed        = 2
	playerMissileSpeed = 3
	playerShootCooldown = 500 * time.Millisecond // Cooldown for shooting
)

type PlayerMissile struct {
	Sprite *ebiten.Image
	X      int
	Y      int
}

type Player struct {
	Sprite     *ebiten.Image
	X          int
	Y          int
	ShootTimer *stopwatch.Stopwatch
	Missiles   []*PlayerMissile // Slice to hold active missiles
	Points int
}

func NewPlayer() *Player {
	playerWidth := assets.Player.Bounds().Dx()
	playerHeight := assets.Player.Bounds().Dy()
	return &Player{
		Sprite:     assets.Player,
		X:          (gameWidth - playerWidth) / 2,
		Y:          gameHeight - playerHeight - 8, // 8 pixels from the bottom
		ShootTimer: stopwatch.NewStopwatch(playerShootCooldown),
		Missiles:   make([]*PlayerMissile, 0), // Initialize missile slice
		Points: 0,
	}
}

func NewPlayerMissile(p *Player) *PlayerMissile {
	// Center missile on player
	missileWidth := assets.PlayerShot.Bounds().Dx()
	playerWidth := p.Sprite.Bounds().Dx()
	return &PlayerMissile{
		Sprite: assets.PlayerShot,
		X:      p.X + (playerWidth/2) - (missileWidth/2),
		Y:      p.Y,
	}
}

// Update now accepts audio.Context to play sounds
func (p *Player) Update(audioContext *audio.Context) error {
	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += playerSpeed
	}

	// Keep player within screen bounds
	playerSpriteWidth := p.Sprite.Bounds().Dx()
	if p.X < 0 {
		p.X = 0
	}
	if p.X+playerSpriteWidth > gameWidth {
		p.X = gameWidth - playerSpriteWidth
	}

	// Shooting logic
	p.ShootTimer.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if !p.ShootTimer.IsRunning() || p.ShootTimer.IsDone() {
			newMissile := NewPlayerMissile(p)
			p.Missiles = append(p.Missiles, newMissile)
			p.ShootTimer.Reset() // Reset and start the timer
			p.ShootTimer.Start()

			// Play shoot sound
			if audioContext != nil {
				shootSoundBytes := assets.PlayerShootSound
				shootStream, err := vorbis.DecodeWithSampleRate(audioContext.SampleRate(), bytes.NewReader(shootSoundBytes))
				if err != nil {
					log.Printf("Error decoding player shoot sound: %v", err)
				} else {
					shootAudioPlayer, err := audioContext.NewPlayer(shootStream)
					if err != nil {
						log.Printf("Error creating audio player for shoot sound: %v", err)
					} else {
						shootAudioPlayer.Play()
					}
				}
			}
		}
	}

	// Update missiles
	activeMissiles := make([]*PlayerMissile, 0, len(p.Missiles))
	for _, missile := range p.Missiles {
		missile.Y -= playerMissileSpeed
		if missile.Y+missile.Sprite.Bounds().Dy() > 0 { // Check if missile is still on screen (top edge)
			activeMissiles = append(activeMissiles, missile)
		}
	}
	p.Missiles = activeMissiles

	return nil
}




