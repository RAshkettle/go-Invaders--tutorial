package main

import (
	"bytes"
	"invaders/assets"
	"math"
	"time"

	stopwatch "github.com/RAshkettle/Stopwatch"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type GameScene struct {
	sceneManager     *SceneManager
	aliens           []*Alien
	timer            *stopwatch.Stopwatch
	currentDirection Direction
	audioContext     *audio.Context
}

const (
	STEP = 16
)

type Direction int

const (
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
	}
	return g
}

func (g *GameScene) SpawnAliens() []*Alien {
	return SpawnAlienWave()
}

func toggleDirection(current Direction) Direction {
	if current == LEFT {
		return RIGHT
	}
	return LEFT
}

func (g *GameScene) moveAliens() {
	// Play Move Sound - create fresh player for clean audio
	moveStream, err := vorbis.Decode(g.audioContext, bytes.NewReader(assets.MoveSound))
	if err == nil {
		moveAudioPlayer, err := g.audioContext.NewPlayer(moveStream)
		if err == nil {
			moveAudioPlayer.Play()
		}
	}

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
			alien.Y += 8        // Move down when reversing direction
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
