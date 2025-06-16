package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"invaders/assets"
	"log" // Added for logging
	"math"
	"math/rand"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const gameSceneHeight = 240 // Height of the game area
var audioContext = audio.NewContext(44100)


type GameScene struct {
	sceneManager     *SceneManager
	aliens           []*Alien
	timer            *stopwatch.Stopwatch
	currentDirection Direction
	audioContext     *audio.Context
	player           *Player
	scoreFont        *text.GoTextFace
	waveTimer        *stopwatch.Stopwatch
	alienMissiles    []*AlienMissile
	deathTimer       *stopwatch.Stopwatch
	playerDead       bool
	bases            []*Base
}

const (
	STEP = 16
)

type Direction int 
const(
	LEFT Direction = iota
	RIGHT
)



func (g *GameScene) Update() error {
	currentSpeed := len(g.aliens) * 20

	// Check death timer first
	if g.playerDead {
		g.deathTimer.Update()
		if g.deathTimer.IsDone() {
			g.sceneManager.TransitionTo(SceneEndScreen)
			return nil
		}
		// Don't process other game logic while player is dead
		return nil
	}

	if !g.timer.IsRunning() {
		g.timer.Start()
	}
	g.timer.Update()
	if g.timer.IsDone() {
		// This is when we animate and Move
		g.moveAliens()
		g.timer = stopwatch.NewStopwatch(time.Duration(currentSpeed) * time.Millisecond)
		g.timer.Start()
	}

	// Check for lose condition (aliens reaching bottom)
	if len(g.aliens) > 0 {
		// Get alien height from the sprite. Assumes all alien sprites for CurrentFrame are same height.
		alienHeight := g.aliens[0].Sprite[g.aliens[0].CurrentFrame].Bounds().Dy()
		for _, alien := range g.aliens {
			if alien.Y+alienHeight >= gameSceneHeight {
				g.sceneManager.TransitionTo(SceneEndScreen) // Immediate transition for aliens reaching bottom
				return nil                            // Transitioning, no more updates for this scene
			}
		}
	}

	if err := g.player.Update(g.audioContext); err != nil {
		return err
	}

	// Check for missile-alien collisions
	g.CheckPlayerMissileCollision()
	
	// Check for alien missile-player collisions
	g.CheckAlienMissilePlayerCollision()
	
	// Check for missile-base collisions
	g.CheckMissileBaseCollisions()
	
	// Check for alien-base collisions
	g.CheckAlienBaseCollisions()
	
	g.CheckWaveStatus()
	g.waveTimer.Update()
	if g.waveTimer.IsDone() {
		g.waveTimer.Stop()
		g.aliens = SpawnAlienWave()
	}

	// Update alien missiles
	activeAlienMissiles := make([]*AlienMissile, 0, len(g.alienMissiles))
	for _, missile := range g.alienMissiles {
		missile.Y += 1 // Move missile down at speed 1
		if missile.Y < gameSceneHeight { // Keep missile if still on screen
			activeAlienMissiles = append(activeAlienMissiles, missile)
		}
	}
	g.alienMissiles = activeAlienMissiles

 

	return nil
}

func (g *GameScene)CheckWaveStatus(){
	if len(g.aliens) == 0 && !g.waveTimer.IsRunning(){
		g.waveTimer.Reset()
		g.waveTimer.Start()
	}
}


func (g *GameScene) Draw(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	scaledWidth := width / 320.0
	scaledHeight := height / 240.0
	scale := math.Min(float64(scaledWidth), float64(scaledHeight))

	// Calculate centering offsets
	gameWidth := 320.0 * scale
	gameHeight := 240.0 * scale
	offsetX := (float64(width) - gameWidth) / 2.0
	offsetY := (float64(height) - gameHeight) / 2.0

	for _, alien := range g.aliens {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Scale(float64(scale), float64(scale))
		op.GeoM.Translate(float64(alien.X)*scale+offsetX, float64(alien.Y)*scale+offsetY)
		screen.DrawImage(alien.Sprite[alien.CurrentFrame], op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.GeoM.Translate(float64(g.player.X)*scale+offsetX, float64(g.player.Y)*scale+offsetY)

	screen.DrawImage(g.player.Sprite,op)

	// Draw player missiles
	for _, missile := range g.player.Missiles {
		missileOp := &ebiten.DrawImageOptions{}
		missileOp.GeoM.Scale(float64(scale), float64(scale))
		missileOp.GeoM.Translate(float64(missile.X)*scale+offsetX, float64(missile.Y)*scale+offsetY)
		screen.DrawImage(missile.Sprite, missileOp)
	}

	// Draw alien missiles
	for _, missile := range g.alienMissiles {
		missileOp := &ebiten.DrawImageOptions{}
		missileOp.GeoM.Scale(float64(scale), float64(scale))
		missileOp.GeoM.Translate(float64(missile.X)*scale+offsetX, float64(missile.Y)*scale+offsetY)
		screen.DrawImage(missile.Sprite, missileOp)
	}

	// Draw bases
	for _, base := range g.bases {
		for _, block := range base.Blocks {
			if block.Exists {
				blockOp := &ebiten.DrawImageOptions{}
				blockOp.GeoM.Scale(float64(scale)*0.5, float64(scale)*0.5) // Scale blocks down by 50%
				blockOp.GeoM.Translate(float64(block.X)*scale+offsetX, float64(block.Y)*scale+offsetY)
				screen.DrawImage(block.Sprite, blockOp)
			}
		}
	}

	// Draw score
	scoreText := fmt.Sprintf("SCORE: %d", g.player.Points)
	textOp := &text.DrawOptions{}
	textOp.GeoM.Scale(float64(scale), float64(scale))
	textOp.GeoM.Translate(offsetX+15*scale, offsetY+15*scale) // Increased padding for better positioning
	textOp.ColorScale.ScaleWithColor(color.RGBA{220, 220, 255, 255}) // Light blue-white color for better contrast
	text.Draw(screen, scoreText, g.scoreFont, textOp)
}

func (g *GameScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}



func NewGameScene(sm *SceneManager) *GameScene {
	// Initialize audio context
	

	// Create score font
	scoreFontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	scoreFont := &text.GoTextFace{
		Source: scoreFontSource,
		Size:   8,
	}

	g := &GameScene{
		sceneManager:     sm,
		aliens:           SpawnAlienWave(),
		timer:            stopwatch.NewStopwatch(1 * time.Second),
		currentDirection: LEFT,
		audioContext:     audioContext,
		player:           NewPlayer(),
		scoreFont:        scoreFont,
		waveTimer:        stopwatch.NewStopwatch(3 * time.Second),
		alienMissiles:    make([]*AlienMissile, 0),
		deathTimer:       stopwatch.NewStopwatch(1500 * time.Millisecond), // 1.5 seconds
		playerDead:       false,
	}
	
	// Create bases positioned above the player
	g.bases = CreateBases(g.player.Y)
	
	return g
}

func (g *GameScene) SpawnAliens() []*Alien {
	return SpawnAlienWave()
}

func toggleDirection(current Direction)Direction{
	if current == LEFT{
		return RIGHT
	}
	return LEFT
}

func (g *GameScene) moveAliens() {
	// Play Move Sound
	moveStream, err := vorbis.DecodeWithSampleRate(g.audioContext.SampleRate(), bytes.NewReader(assets.MoveSound))
	if err != nil {

		return // Don't proceed if decoding failed
	}

	moveAudioPlayer, err := g.audioContext.NewPlayer(moveStream)
	if err != nil {

		return // Don't proceed if player creation failed
	}
	moveAudioPlayer.Play()

	// Check if any alien will hit the screen boundaries
	shouldReverse := false
	for _, alien := range g.aliens {
		if g.currentDirection == LEFT && alien.X-8 <= 0 {
			shouldReverse = true
			break
		} else if g.currentDirection == RIGHT && alien.X+8 >= 320-ALIEN_SIZE {
			shouldReverse = true
			break
		}
	}

	// If we need to reverse direction, do it and move down
	if shouldReverse {
		g.currentDirection = toggleDirection(g.currentDirection)
		for _, alien := range g.aliens {
			alien.Y += 8 // Move down when reversing direction
			alien.ToggleFrame() // Toggle animation frame
		}
	} else {
		// Move aliens horizontally
		for _, alien := range g.aliens {
			if g.currentDirection == LEFT {
				alien.X -= 8
			} else {
				alien.X += 8
			}
			alien.ToggleFrame() // Toggle animation frame
		}
	}

	// Check for SquidAlien shooting (10% chance per movement)
	for _, alien := range g.aliens {
		// Only allow shooting if we have less than 3 missiles active
		if alien.AlienType == SquidAlien && rand.Float64() < 0.1 && len(g.alienMissiles) < 3 {
			// Create new alien missile
			missileX := alien.X + alien.Sprite[alien.CurrentFrame].Bounds().Dx()/2 - assets.AlienShot.Bounds().Dx()/2
			missileY := alien.Y + alien.Sprite[alien.CurrentFrame].Bounds().Dy()
			
			newAlienMissile := &AlienMissile{
				Sprite: assets.AlienShot,
				X:      missileX,
				Y:      missileY,
			}
			g.alienMissiles = append(g.alienMissiles, newAlienMissile)
		}
	}
}


func(g *GameScene) CheckPlayerMissileCollision() {
	activeMissiles := make([]*PlayerMissile, 0, len(g.player.Missiles))
	activeAliens := make([]*Alien, 0, len(g.aliens))

	// Track which aliens were hit
	aliensHit := make(map[*Alien]bool)

	for _, missile := range g.player.Missiles {
		hit := false
		
		// Get missile center point (only center 2 pixels)
		missileX := missile.X + missile.Sprite.Bounds().Dx()/2 - 1
		missileY := missile.Y + missile.Sprite.Bounds().Dy()/2 - 1
		missileRect := image.Rect(missileX, missileY, missileX+2, missileY+2)

		for _, alien := range g.aliens {
			// Skip if this alien was already hit
			if aliensHit[alien] {
				continue
			}

			// Get alien sprite bounds
			alienRect := image.Rect(alien.X, alien.Y, 
				alien.X+alien.Sprite[alien.CurrentFrame].Bounds().Dx(), 
				alien.Y+alien.Sprite[alien.CurrentFrame].Bounds().Dy())

			// Check if missile center intersects with alien
			if missileRect.Overlaps(alienRect) {
				// Add alien points to player
				g.player.Points += alien.PointsValue
				hit = true
				aliensHit[alien] = true

				// Play alien explosion sound
				explosionStream, err := vorbis.DecodeWithSampleRate(g.audioContext.SampleRate(), bytes.NewReader(assets.AlienExplosionSound))
				if err != nil {
					log.Printf("Error decoding alien explosion sound: %v", err)
				} else {
					explosionAudioPlayer, err := g.audioContext.NewPlayer(explosionStream)
					if err != nil {
						log.Printf("Error creating audio player for explosion sound: %v", err)
					} else {
						explosionAudioPlayer.Play()
					}
				}

				break // This missile hit an alien, don't check other aliens
			}
		}

		// Only keep missile if it didn't hit anything
		if !hit {
			activeMissiles = append(activeMissiles, missile)
		}
	}

	// Build active aliens list (only aliens that weren't hit)
	for _, alien := range g.aliens {
		if !aliensHit[alien] {
			activeAliens = append(activeAliens, alien)
		}
	}

	// Update the slices with only active (non-collided) objects
	g.player.Missiles = activeMissiles
	g.aliens = activeAliens
}

func (g *GameScene) CheckAlienMissilePlayerCollision() {
	// Don't check collisions if player is already dead
	if g.playerDead {
		return
	}

	activeAlienMissiles := make([]*AlienMissile, 0, len(g.alienMissiles))

	// Get player bounds
	playerRect := image.Rect(g.player.X, g.player.Y, 
		g.player.X+g.player.Sprite.Bounds().Dx(), 
		g.player.Y+g.player.Sprite.Bounds().Dy())

	for _, missile := range g.alienMissiles {
		// Get missile bounds
		missileRect := image.Rect(missile.X, missile.Y,
			missile.X+missile.Sprite.Bounds().Dx(),
			missile.Y+missile.Sprite.Bounds().Dy())

		// Check if missile intersects with player
		if missileRect.Overlaps(playerRect) {
			// Player is hit - start death timer and play death sound
			g.playerDead = true
			g.deathTimer.Reset()
			g.deathTimer.Start()

			// Play player death sound
			deathStream, err := vorbis.DecodeWithSampleRate(g.audioContext.SampleRate(), bytes.NewReader(assets.PlayerDeathSound))
			if err != nil {
				log.Printf("Error decoding player death sound: %v", err)
			} else {
				deathAudioPlayer, err := g.audioContext.NewPlayer(deathStream)
				if err != nil {
					log.Printf("Error creating audio player for death sound: %v", err)
				} else {
					deathAudioPlayer.Play()
				}
			}

			return // No need to process remaining missiles
		} else {
			// Keep missile if no collision
			activeAlienMissiles = append(activeAlienMissiles, missile)
		}
	}

	// Update alien missiles slice
	g.alienMissiles = activeAlienMissiles
}

func (g *GameScene) CheckMissileBaseCollisions() {
	// Check player missiles vs bases
	activeMissiles := make([]*PlayerMissile, 0, len(g.player.Missiles))
	for _, missile := range g.player.Missiles {
		hit := false
		
		// Get missile center 4 pixels on X-axis for more precise collision
		missileWidth := missile.Sprite.Bounds().Dx()
		missileCenterX := missile.X + missileWidth/2 - 2 // Center minus 2 pixels
		missileRect := image.Rect(missileCenterX, missile.Y,
			missileCenterX+4, // Only 4 pixels wide
			missile.Y+missile.Sprite.Bounds().Dy())
		
		for _, base := range g.bases {
			for _, block := range base.Blocks {
				if !block.Exists {
					continue
				}
				
				// Get block bounds (accounting for 50% scale)
				blockRect := image.Rect(block.X, block.Y, block.X+8, block.Y+8)
				
				if missileRect.Overlaps(blockRect) {
					block.TakeDamage()
					hit = true
					
					// Play alien explosion sound for base hit
					explosionStream, err := vorbis.DecodeWithSampleRate(g.audioContext.SampleRate(), bytes.NewReader(assets.AlienExplosionSound))
					if err != nil {
						log.Printf("Error decoding base hit sound: %v", err)
					} else {
						explosionAudioPlayer, err := g.audioContext.NewPlayer(explosionStream)
						if err != nil {
							log.Printf("Error creating audio player for base hit sound: %v", err)
						} else {
							explosionAudioPlayer.Play()
						}
					}
					break
				}
			}
			if hit {
				break
			}
		}
		
		if !hit {
			activeMissiles = append(activeMissiles, missile)
		}
	}
	g.player.Missiles = activeMissiles
	
	// Check alien missiles vs bases
	activeAlienMissiles := make([]*AlienMissile, 0, len(g.alienMissiles))
	for _, missile := range g.alienMissiles {
		hit := false
		
		// Get missile bounds
		missileRect := image.Rect(missile.X, missile.Y,
			missile.X+missile.Sprite.Bounds().Dx(),
			missile.Y+missile.Sprite.Bounds().Dy())
		
		for _, base := range g.bases {
			for _, block := range base.Blocks {
				if !block.Exists {
					continue
				}
				
				// Get block bounds (accounting for 50% scale)
				blockRect := image.Rect(block.X, block.Y, block.X+8, block.Y+8)
				
				if missileRect.Overlaps(blockRect) {
					block.TakeDamage()
					hit = true
					
					// Play alien explosion sound for base hit
					explosionStream, err := vorbis.DecodeWithSampleRate(g.audioContext.SampleRate(), bytes.NewReader(assets.AlienExplosionSound))
					if err != nil {
						log.Printf("Error decoding base hit sound: %v", err)
					} else {
						explosionAudioPlayer, err := g.audioContext.NewPlayer(explosionStream)
						if err != nil {
							log.Printf("Error creating audio player for base hit sound: %v", err)
						} else {
							explosionAudioPlayer.Play()
						}
					}
					break
				}
			}
			if hit {
				break
			}
		}
		
		if !hit {
			activeAlienMissiles = append(activeAlienMissiles, missile)
		}
	}
	g.alienMissiles = activeAlienMissiles
}

func (g *GameScene) CheckAlienBaseCollisions() {
	for _, alien := range g.aliens {
		// Get alien bounds
		alienRect := image.Rect(alien.X, alien.Y,
			alien.X+alien.Sprite[alien.CurrentFrame].Bounds().Dx(),
			alien.Y+alien.Sprite[alien.CurrentFrame].Bounds().Dy())
		
		for _, base := range g.bases {
			for _, block := range base.Blocks {
				if !block.Exists {
					continue
				}
				
				// Get block bounds (accounting for 50% scale)
				blockRect := image.Rect(block.X, block.Y, block.X+8, block.Y+8)
				
				if alienRect.Overlaps(blockRect) {
					// Alien collides with base block - destroy the block immediately
					block.Exists = false
					// Alien is not harmed and continues moving
				}
			}
		}
	}
}
