package state

import (
	"image/color"
	"math"
	"strconv"

	"astroleap/internal/assets"
	"astroleap/internal/entity"
	"astroleap/internal/input"
	"astroleap/internal/records"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type WinState struct {
	machine     *Machine
	score       int
	crystals    int
	timeSeconds float64
	medals      int
	rank        string
	recordRes   records.RecordResult
	timer       float64
	landerY     float64
	particles   *entity.ParticleSystem
	starfield   *ebiten.Image
}

func NewWinState(m *Machine, score int, crystals int) *WinState {
	return NewWinStateWithRecords(m, score, crystals, 0, 0)
}

func NewWinStateWithRecords(m *Machine, score int, crystals int, timeSec float64, medals int) *WinState {
	rank := records.EvaluateRank(score, timeSec, medals)
	res := records.SubmitRun(score, timeSec, medals)

	return &WinState{
		machine:     m,
		score:       score,
		crystals:    crystals,
		timeSeconds: timeSec,
		medals:      medals,
		rank:        rank,
		recordRes:   res,
		landerY:     110.0,
		particles:   entity.NewParticleSystem(),
		starfield:   assets.GenerateStarfield(320, 180),
	}
}

func (s *WinState) Enter() {}

func (s *WinState) Update(dt float64) {
	s.timer += dt
	inp := input.PollInput()

	// Animate Lander blastoff ascending into the stars
	if s.landerY > -40 {
		s.landerY -= 35.0 * dt
		s.particles.SpawnJetpackSparks(160, s.landerY+26, true)
	}
	s.particles.Update(dt)

	if s.timer > 1.0 && (inp.JumpPressed || inp.Restart) {
		s.machine.Change(NewTitleState(s.machine))
	}
}

func (s *WinState) Draw(screen *ebiten.Image) {
	// Background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.starfield, op)

	// Distant Earth
	earthOp := &ebiten.DrawImageOptions{}
	earthOp.GeoM.Translate(230, 20)
	screen.DrawImage(assets.Get().Earth, earthOp)

	// Rocket exhaust particles
	s.particles.Draw(screen, 0)

	// Ascending Lander
	atlas := assets.Get()
	landerOp := &ebiten.DrawImageOptions{}
	landerOp.GeoM.Translate(144, s.landerY)
	screen.DrawImage(atlas.Lander, landerOp)

	// Expanded Banner Modal
	vector.DrawFilledRect(screen, 32, 38, 256, 108, color.RGBA{15, 20, 35, 230}, false)
	vector.StrokeRect(screen, 32, 38, 256, 108, 1, color.RGBA{100, 230, 255, 255}, false)

	ui.DrawText(screen, "MISSION ACCOMPLISHED!", 94, 46, color.RGBA{255, 220, 60, 255})
	ui.DrawText(screen, "ESCAPE ROCKET LAUNCHED!", 88, 58, color.RGBA{140, 240, 255, 255})

	scoreStr := "FINAL SCORE: " + strconv.Itoa(s.score)
	ui.DrawText(screen, scoreStr, 56, 72, color.RGBA{255, 255, 255, 255})

	rankCol := color.RGBA{255, 215, 50, 255}
	if s.rank == "S" {
		rankCol = color.RGBA{255, 80, 140, 255}
	}
	rankStr := "RANK: [" + s.rank + "]"
	ui.DrawText(screen, rankStr, 196, 72, rankCol)

	statsStr := "CRYSTALS: " + strconv.Itoa(s.crystals) + "   MEDALS: " + strconv.Itoa(s.medals) + "/9"
	ui.DrawText(screen, statsStr, 56, 84, color.RGBA{100, 245, 255, 255})

	timeStr := "TOTAL TIME: " + records.FormatTime(s.timeSeconds)
	ui.DrawText(screen, timeStr, 56, 96, color.RGBA{255, 220, 100, 255})

	if s.recordRes.NewHighScore || s.recordRes.NewBestTime {
		if int(s.timer*6)%2 == 0 {
			ui.DrawText(screen, "★ NEW BEST RECORD! ★", 92, 110, color.RGBA{255, 240, 80, 255})
		}
	} else if s.timer > 1.0 && math.Sin(s.timer*6.0) > -0.2 {
		ui.DrawText(screen, "PRESS SPACE TO RETURN", 94, 130, color.RGBA{200, 220, 255, 255})
	}
}

func (s *WinState) Exit() {}
