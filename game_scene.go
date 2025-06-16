package main

import (
	"bytes"
	"image"
	"invaders/assets"
	"log" // Added for logging
	"math"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const gameSceneHeight = 240 // Height of the game area

type GameScene struct {
	sceneManager     *SceneManager
	aliens           []*Alien
	timer            *stopwatch.Stopwatch
	currentDirection Direction
	audioContext     *audio.Context
	player *Player
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
				g.sceneManager.TransitionTo(SceneEndScreen) // Assumes SceneEnd is defined in scene_manager.go
				return nil                            // Transitioning, no more updates for this scene
			}
		}
	}

	if err := g.player.Update(g.audioContext); err != nil {
		return err
	}

	// Check for missile-alien collisions
	g.CheckPlayerMissileCollision()

	return nil
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
}

func (g *GameScene) Layout(outerWidth, outerHeight int) (int, int) {
	return outerWidth, outerHeight
}

func (g *GameScene) Reset() {}

func NewGameScene(sm *SceneManager) *GameScene {
	// Initialize audio context
	audioContext := audio.NewContext(44100)

	g := &GameScene{
		sceneManager:     sm,
		aliens:           SpawnAlienWave(),
		timer:            stopwatch.NewStopwatch(1 * time.Second),
		currentDirection: LEFT,
		audioContext:     audioContext,
		player: NewPlayer(),
	}
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
		log.Printf("Error decoding move sound (audio/move.ogg): %v", err)
		return // Don't proceed if decoding failed
	}

	moveAudioPlayer, err := g.audioContext.NewPlayer(moveStream)
	if err != nil {
		log.Printf("Error creating audio player for move sound: %v", err)
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
