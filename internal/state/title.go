package state

import (
	"image/color"
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/audio"
	"astroleap/internal/entity"
	"astroleap/internal/input"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type TitleState struct {
	machine      *Machine
	timer        float64
	starfield    *ebiten.Image
	particles    *entity.ParticleSystem
	astroX       float64
	astroY       float64
}

func NewTitleState(m *Machine) *TitleState {
	return &TitleState{
		machine:   m,
		starfield: assets.GenerateStarfield(320, 180),
		particles: entity.NewParticleSystem(),
		astroX:    152.0,
		astroY:    84.0,
	}
}

func (s *TitleState) Enter() {
	audio.Get().StartBGM()
}

func (s *TitleState) Update(dt float64) {
	s.timer += dt
	inp := input.PollInput()

	if inp.ToggleMute {
		audio.Get().ToggleMute()
	}

	// Floating astronaut motion
	s.astroY = 78.0 + math.Sin(s.timer*2.5)*5.0

	// Thruster sparks beneath mascot
	if int(s.timer*30)%4 == 0 {
		s.particles.SpawnJetpackSparks(s.astroX+8, s.astroY+14, true)
	}
	s.particles.Update(dt)

	if inp.JumpPressed {
		audio.Get().PlayJump()
		s.machine.Change(NewIntroCinematicState(s.machine))
		return
	} else if inp.Restart {
		audio.Get().PlayJump()
		s.machine.Change(NewPlayState(s.machine))
		return
	}
}

func (s *TitleState) Draw(screen *ebiten.Image) {
	// 1. Starfield background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.starfield, op)

	// 2. Distant Earth
	earthOp := &ebiten.DrawImageOptions{}
	earthOp.GeoM.Translate(230, 20)
	screen.DrawImage(assets.Get().Earth, earthOp)

	// 3. Particles
	s.particles.Draw(screen, 0)

	// 4. Floating Mascot Astronaut
	atlas := assets.Get()
	astroOp := &ebiten.DrawImageOptions{}
	astroOp.GeoM.Translate(s.astroX, s.astroY)
	screen.DrawImage(atlas.AstronautThrust, astroOp)

	// 5. Title & Typography
	ui.DrawText(screen, "ASTROLEAP", 124, 26, color.RGBA{100, 240, 255, 255})
	ui.DrawText(screen, "LUNAR ODYSSEY", 112, 38, color.RGBA{255, 200, 60, 255})

	// 6. Blinking Start prompts
	if math.Sin(s.timer*5.0) > -0.2 {
		ui.DrawText(screen, "PRESS SPACE FOR STORY PROLOGUE", 68, 108, color.RGBA{255, 255, 255, 255})
		ui.DrawText(screen, "PRESS ENTER FOR QUICK START", 78, 120, color.RGBA{255, 220, 100, 255})
	}

	// 7. Instructions / Controls
	ui.DrawText(screen, "A / D : MOVE", 48, 142, color.RGBA{170, 180, 210, 255})
	ui.DrawText(screen, "SPACE : JUMP / THRUSTER", 134, 142, color.RGBA{170, 180, 210, 255})
	ui.DrawText(screen, "J : SHOOT WEAPON", 54, 154, color.RGBA{60, 220, 255, 255})
	ui.DrawText(screen, "Q : SWAP WEAPON", 168, 154, color.RGBA{255, 210, 60, 255})
	ui.DrawText(screen, "M : MUTE AUDIO", 114, 166, color.RGBA{120, 135, 170, 255})
}

func (s *TitleState) Exit() {}
