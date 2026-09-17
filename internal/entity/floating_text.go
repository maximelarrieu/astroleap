package entity

import (
	"image/color"

	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

// FloatingText represents an animated drifting score/combo popup.
type FloatingText struct {
	X       float64
	Y       float64
	Text    string
	Color   color.RGBA
	Life    float64
	MaxLife float64
	VY      float64
	Dead    bool
}

func NewFloatingText(text string, x, y float64, col color.RGBA) *FloatingText {
	return &FloatingText{
		X:       x,
		Y:       y,
		Text:    text,
		Color:   col,
		Life:    0.9,
		MaxLife: 0.9,
		VY:      -28.0, // Drifts upward
		Dead:    false,
	}
}

func (ft *FloatingText) Update(dt float64) {
	ft.Life -= dt
	if ft.Life <= 0 {
		ft.Dead = true
		return
	}
	ft.Y += ft.VY * dt
	ft.VY *= 0.96 // Gentle deceleration
}

func (ft *FloatingText) Draw(screen *ebiten.Image, camX float64) {
	if ft.Dead {
		return
	}
	sx := int(ft.X - camX)
	if sx < -60 || sx > 340 {
		return
	}

	alpha := ft.Life / ft.MaxLife
	if alpha > 1.0 {
		alpha = 1.0
	} else if alpha < 0 {
		alpha = 0
	}

	col := color.RGBA{
		R: uint8(float64(ft.Color.R) * alpha),
		G: uint8(float64(ft.Color.G) * alpha),
		B: uint8(float64(ft.Color.B) * alpha),
		A: uint8(255 * alpha),
	}

	ui.DrawText(screen, ft.Text, sx, int(ft.Y), col)
}
