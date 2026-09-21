package state

import (
	"image/color"
	"math"
	"strconv"

	"astroleap/internal/assets"
	"astroleap/internal/entity"
	"astroleap/internal/input"
	"astroleap/internal/level"
	"astroleap/internal/records"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// SectorClearState handles the inter-level transition cutscene and score breakdown.
type SectorClearState struct {
	machine     *Machine
	sector      int
	score       int
	crystals    int
	playerRef   *entity.Player
	elapsedTime float64
	medals      int
	timer       float64
	landerY     float64
	particles   *entity.ParticleSystem
	starfield   *ebiten.Image
	sectorName  string
}

func NewSectorClearState(m *Machine, sector int, score int, crystals int, p *entity.Player) *SectorClearState {
	return NewSectorClearStateWithRecords(m, sector, score, crystals, p, 0, 0)
}

func NewSectorClearStateWithRecords(m *Machine, sector int, score int, crystals int, p *entity.Player, elapsedTime float64, medals int) *SectorClearState {
	secName := "LUNAR OUTPOST"
	if sector == 2 {
		secName = "PHOBOS RIDGE"
	} else if sector == 3 {
		secName = "EUROPA ICE CORE"
	} else if sector == 4 {
		secName = "MOTHERSHIP CORRIDOR"
	} else if sector == 5 {
		secName = "REACTOR BAY"
	}

	return &SectorClearState{
		machine:     m,
		sector:      sector,
		score:       score,
		crystals:    crystals,
		playerRef:   p,
		elapsedTime: elapsedTime,
		medals:      medals,
		landerY:     110.0,
		particles:   entity.NewParticleSystem(),
		starfield:   assets.GenerateStarfield(320, 180),
		sectorName:  secName,
	}
}

func (s *SectorClearState) Enter() {}

func (s *SectorClearState) Update(dt float64) {
	s.timer += dt
	inp := input.PollInput()

	// Animate Lander ascending into the stars
	if s.landerY > -40 {
		s.landerY -= 35.0 * dt
		s.particles.SpawnJetpackSparks(160, s.landerY+26, true)
	}
	s.particles.Update(dt)

	if s.timer > 0.8 && (inp.JumpPressed || inp.Restart) {
		nextSector := s.sector + 1
		if nextSector > level.MaxLevels {
			s.machine.Change(NewWinStateWithRecords(s.machine, s.score, s.crystals, s.elapsedTime, s.medals))
		} else {
			s.machine.Change(NewPlayStateWithRecords(s.machine, nextSector, s.score, s.crystals, s.playerRef, s.elapsedTime, s.medals))
		}
	}
}

func (s *SectorClearState) Draw(screen *ebiten.Image) {
	// Background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.starfield, op)

	// Celestial body for the cleared sector
	celestialOp := &ebiten.DrawImageOptions{}
	celestialOp.GeoM.Translate(230, 20)
	atlas := assets.Get()
	if s.sector == 1 {
		screen.DrawImage(atlas.Earth, celestialOp)
	} else if s.sector == 2 {
		screen.DrawImage(atlas.Mars, celestialOp)
	} else if s.sector == 3 {
		screen.DrawImage(atlas.Jupiter, celestialOp)
	} else {
		screen.DrawImage(atlas.Earth, celestialOp)
	}

	// Rocket exhaust particles
	s.particles.Draw(screen, 0)

	// Ascending Lander or Shuttle
	landerOp := &ebiten.DrawImageOptions{}
	landerOp.GeoM.Translate(144, s.landerY)
	screen.DrawImage(atlas.Lander, landerOp)

	// Intermission Modal Box
	vector.DrawFilledRect(screen, 36, 40, 248, 102, color.RGBA{14, 20, 36, 230}, false)
	vector.StrokeRect(screen, 36, 40, 248, 102, 1, color.RGBA{100, 235, 255, 255}, false)

	header := "SECTOR " + strconv.Itoa(s.sector) + " COMPLETE!"
	ui.DrawText(screen, header, 96, 48, color.RGBA{255, 220, 60, 255})
	ui.DrawText(screen, s.sectorName+" SECURED", 86, 60, color.RGBA{140, 240, 255, 255})

	scoreStr := "TOTAL SCORE: " + strconv.Itoa(s.score)
	ui.DrawText(screen, scoreStr, 88, 74, color.RGBA{255, 255, 255, 255})

	crystStr := "CRYSTALS: " + strconv.Itoa(s.crystals) + "   MEDALS: " + strconv.Itoa(s.medals) + "/15"
	ui.DrawText(screen, crystStr, 76, 86, color.RGBA{100, 245, 255, 255})

	timeStr := "MISSION TIME: " + records.FormatTime(s.elapsedTime)
	ui.DrawText(screen, timeStr, 82, 98, color.RGBA{255, 215, 60, 255})

	if s.timer > 0.8 && math.Sin(s.timer*6.0) > -0.2 {
		nextSector := s.sector + 1
		nextPrompt := "PRESS SPACE FOR SECTOR " + strconv.Itoa(nextSector)
		if nextSector == 5 {
			nextPrompt = "PRESS SPACE FOR REACTOR BAY"
		}
		ui.DrawText(screen, nextPrompt, 68, 122, color.RGBA{255, 200, 80, 255})
	}
}

func (s *SectorClearState) Exit() {}
