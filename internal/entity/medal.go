package entity

import (
	"image/color"
	"math"

	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// GoldMedal represents a hidden collectible astronaut medal.
type GoldMedal struct {
	X         float64
	Y         float64
	Collected bool
	Timer     float64
}

func NewGoldMedal(x, y float64) *GoldMedal {
	return &GoldMedal{
		X: x,
		Y: y,
	}
}

func (m *GoldMedal) GetRect() physics.Rect {
	return physics.Rect{
		X: m.X + 2,
		Y: m.Y + 2,
		W: 12,
		H: 12,
	}
}

func (m *GoldMedal) Update(dt float64) {
	m.Timer += dt
}

func (m *GoldMedal) Draw(screen *ebiten.Image, camX float64) {
	if m.Collected {
		return
	}

	sx := float32(m.X - camX)
	if sx < -20 || sx > 330 {
		return
	}

	// Floating gentle bob
	bobY := float32(m.Y + math.Sin(m.Timer*3.5)*3.0)

	// Outer golden ring
	vector.DrawFilledCircle(screen, sx+8, bobY+8, 7, color.RGBA{255, 200, 30, 240}, false)
	vector.StrokeCircle(screen, sx+8, bobY+8, 7, 1, color.RGBA{255, 245, 120, 255}, false)

	// Inner star insignia
	vector.DrawFilledCircle(screen, sx+8, bobY+8, 4, color.RGBA{255, 120, 20, 255}, false)
	vector.DrawFilledCircle(screen, sx+8, bobY+8, 2, color.RGBA{255, 255, 240, 255}, false)
}
