package entity

import (
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

// Crystal represents a floating, collectible energy diamond.
type Crystal struct {
	X         float64
	Y         float64
	BaseY     float64
	Collected bool
	AnimTimer float64
	Frame     int
}

func NewCrystal(x, y float64) *Crystal {
	return &Crystal{
		X:     x,
		Y:     y,
		BaseY: y,
	}
}

func (c *Crystal) Update(dt float64) {
	if c.Collected {
		return
	}
	c.AnimTimer += dt
	// 4 FPS sparkling rotation
	c.Frame = int(c.AnimTimer*7.0) % 4
	// Gentle float bobbing
	c.Y = c.BaseY + math.Sin(c.AnimTimer*4.0)*2.5
}

func (c *Crystal) GetRect() physics.Rect {
	return physics.Rect{
		X: c.X + 3,
		Y: c.Y + 3,
		W: 10,
		H: 10,
	}
}

func (c *Crystal) Draw(screen *ebiten.Image, camX float64) {
	if c.Collected {
		return
	}
	screenX := c.X - camX
	if screenX < -16 || screenX > 320 {
		return
	}

	atlas := assets.Get()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenX, c.Y)
	screen.DrawImage(atlas.Crystals[c.Frame], op)
}
