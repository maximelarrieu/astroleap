package state

import (
	"image/color"
	"math"
	"strconv"

	"astroleap/internal/assets"
	"astroleap/internal/entity"
	"astroleap/internal/input"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameOverState struct {
	machine   *Machine
	score     int
	crystals  int
	sector    int
	playerRef *entity.Player
	timer     float64
	starfield *ebiten.Image
}

func NewGameOverState(m *Machine, score int, crystals int, sector int, p *entity.Player) *GameOverState {
	if sector < 1 {
		sector = 1
	}
	return &GameOverState{
		machine:   m,
		score:     score,
		crystals:  crystals,
		sector:    sector,
		playerRef: p,
		starfield: assets.GenerateStarfield(320, 180),
	}
}

func (s *GameOverState) Enter() {}

func (s *GameOverState) Update(dt float64) {
	s.timer += dt
	inp := input.PollInput()

	if s.timer > 0.8 && (inp.JumpPressed || inp.Restart) {
		s.machine.Change(NewPlayStateWithProgress(s.machine, s.sector, s.score, s.crystals, s.playerRef))
	}
}

func (s *GameOverState) Draw(screen *ebiten.Image) {
	// Background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.starfield, op)

	// Modal Box
	vector.DrawFilledRect(screen, 40, 48, 240, 84, color.RGBA{22, 14, 25, 225}, false)
	vector.StrokeRect(screen, 40, 48, 240, 84, 1, color.RGBA{240, 50, 70, 255}, false)

	ui.DrawText(screen, "MISSION FAILED", 112, 58, color.RGBA{255, 60, 80, 255})
	ui.DrawText(screen, "ASTRONAUT LOST IN DEEP SPACE", 68, 72, color.RGBA{190, 195, 215, 255})

	scoreStr := "SCORE: " + strconv.Itoa(s.score)
	ui.DrawText(screen, scoreStr, 126, 88, color.RGBA{255, 255, 255, 255})

	crystStr := "CRYSTALS: " + strconv.Itoa(s.crystals)
	ui.DrawText(screen, crystStr, 118, 100, color.RGBA{100, 240, 255, 255})

	if s.timer > 0.8 && math.Sin(s.timer*6.0) > -0.2 {
		retryPrompt := "PRESS SPACE TO RETRY SECTOR " + strconv.Itoa(s.sector)
		ui.DrawText(screen, retryPrompt, 73, 116, color.RGBA{255, 220, 100, 255})
	}
}

func (s *GameOverState) Exit() {}
