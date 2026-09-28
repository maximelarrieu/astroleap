package entity

import (
	"image/color"
	"math"

	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DataTerminal represents an interactive holographic lore transmission terminal.
type DataTerminal struct {
	X         float64
	Y         float64
	Sector    int
	LogID     string
	Collected bool
	Timer     float64
}

// NewDataTerminal creates a new terminal at the specified world coordinates.
func NewDataTerminal(x, y float64, sector int, logID string) *DataTerminal {
	return &DataTerminal{
		X:      x,
		Y:      y,
		Sector: sector,
		LogID:  logID,
	}
}

// Update updates the holographic pulsation and rotation timers.
func (dt *DataTerminal) Update(delta float64) {
	dt.Timer += delta
}

// GetRect returns the 16x16 bounding box of the terminal.
func (dt *DataTerminal) GetRect() physics.Rect {
	return physics.Rect{
		X: dt.X,
		Y: dt.Y,
		W: 16,
		H: 16,
	}
}

// Draw renders the holographic pedestal and floating data beacon.
func (dt *DataTerminal) Draw(screen *ebiten.Image, camX float64) {
	sx := float32(dt.X - camX)
	if sx < -20 || sx > 330 {
		return
	}
	sy := float32(dt.Y)

	// 1. Sleek metallic pedestal base
	colBase := color.RGBA{22, 28, 42, 255}
	colTrim := color.RGBA{45, 60, 90, 255}
	if dt.Collected {
		colBase = color.RGBA{18, 20, 28, 200}
		colTrim = color.RGBA{35, 45, 60, 200}
	}
	vector.DrawFilledRect(screen, sx+2, sy+11, 12, 5, colBase, false)
	vector.DrawFilledRect(screen, sx+4, sy+8, 8, 3, colTrim, false)
	vector.StrokeRect(screen, sx+2, sy+8, 12, 8, 1, colTrim, false)

	// Pedestal status diode
	diodeCol := color.RGBA{60, 220, 255, 255}
	if dt.Collected {
		diodeCol = color.RGBA{90, 105, 130, 180}
	}
	vector.DrawFilledRect(screen, sx+7, sy+12, 2, 2, diodeCol, false)

	// 2. Holographic floating data prism (if uncollected)
	if !dt.Collected {
		floatOffset := float32(math.Sin(dt.Timer*4.0) * 2.0)
		holoY := sy + 2 + floatOffset

		// Pulsing outer holo aura
		pulse := 0.5 + 0.5*math.Sin(dt.Timer*6.0)
		alphaHolo := uint8(140 + pulse*80)

		// Diamond hologram projector
		vector.DrawFilledRect(screen, sx+6, holoY, 4, 6, color.RGBA{80, 230, 255, alphaHolo}, false)
		vector.DrawFilledRect(screen, sx+4, holoY+1, 8, 4, color.RGBA{50, 180, 240, uint8(alphaHolo * 3 / 4)}, false)
		vector.DrawFilledRect(screen, sx+7, holoY+2, 2, 2, color.RGBA{255, 255, 255, 255}, false)

		// Scanning emitter light beam from pedestal to diamond
		vector.DrawFilledRect(screen, sx+7, holoY+6, 2, float32(sy+8-(holoY+6)), color.RGBA{60, 210, 255, 60}, false)
	}
}
