package entity

import (
	"image/color"
	"math"

	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// CryoGeyser represents an upward thermal/cryo steam vent that propels the player into the sky.
type CryoGeyser struct {
	X      float64
	Y      float64 // Base of the vent
	Height float64 // Upward reach of the jet
	Timer  float64
}

func NewCryoGeyser(x, y, height float64) *CryoGeyser {
	return &CryoGeyser{
		X:      x,
		Y:      y,
		Height: height,
	}
}

func (g *CryoGeyser) GetRect() physics.Rect {
	return physics.Rect{
		X: g.X + 2,
		Y: g.Y - g.Height,
		W: 12,
		H: g.Height + 4,
	}
}

func (g *CryoGeyser) Update(dt float64, particles *ParticleSystem) {
	g.Timer += dt

	// Periodically emit upward ice/steam particles
	if int(g.Timer*60)%4 == 0 {
		emitY := g.Y - math.Mod(g.Timer*80.0, g.Height)
		particles.SpawnCrystalSparkles(g.X+8, emitY)
	}
}

func (g *CryoGeyser) Draw(screen *ebiten.Image, camX float64) {
	sx := float32(g.X - camX)
	if sx < -30 || sx > 340 {
		return
	}

	topY := float32(g.Y - g.Height)
	baseY := float32(g.Y)
	h := float32(g.Height)

	// Upward shimmering energy/steam beam
	pulse := float32(math.Sin(g.Timer*8.0))*0.15 + 0.85
	beamAlpha := uint8(60.0 * pulse)

	// Outer beam
	vector.DrawFilledRect(screen, sx+2, topY, 12, h, color.RGBA{100, 240, 255, beamAlpha}, false)
	// Inner core
	vector.DrawFilledRect(screen, sx+5, topY, 6, h, color.RGBA{220, 250, 255, beamAlpha + 40}, false)

	// Ground vent casing
	vector.DrawFilledRect(screen, sx, baseY, 16, 4, color.RGBA{40, 48, 65, 255}, false)
	vector.StrokeRect(screen, sx, baseY, 16, 4, 1, color.RGBA{100, 220, 255, 255}, false)
	vector.DrawFilledRect(screen, sx+4, baseY+1, 8, 2, color.RGBA{80, 240, 255, 255}, false)
}
