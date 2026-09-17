package entity

import (
	"image/color"

	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MovingPlatform represents a floating mechanical platform moving between two waypoints.
type MovingPlatform struct {
	X         float64
	Y         float64
	StartX    float64
	StartY    float64
	EndX      float64
	EndY      float64
	W         float64
	H         float64
	Speed     float64
	Progress  float64
	Direction float64
	DX        float64
	DY        float64
}

func NewMovingPlatform(startX, startY, endX, endY float64, w, h float64, speed float64) *MovingPlatform {
	return &MovingPlatform{
		X:         startX,
		Y:         startY,
		StartX:    startX,
		StartY:    startY,
		EndX:      endX,
		EndY:      endY,
		W:         w,
		H:         h,
		Speed:     speed,
		Progress:  0,
		Direction: 1.0,
	}
}

func (mp *MovingPlatform) GetRect() physics.Rect {
	return physics.Rect{
		X: mp.X,
		Y: mp.Y,
		W: mp.W,
		H: mp.H,
	}
}

func (mp *MovingPlatform) Update(dt float64) {
	// Total distance between endpoints
	dx := mp.EndX - mp.StartX
	dy := mp.EndY - mp.StartY
	totalDist := (dx*dx + dy*dy)
	if totalDist <= 0.001 {
		return
	}

	distPerSec := mp.Speed
	normSpeed := distPerSec / 100.0 // rate of progress

	mp.Progress += mp.Direction * normSpeed * dt
	if mp.Progress >= 1.0 {
		mp.Progress = 1.0
		mp.Direction = -1.0
	} else if mp.Progress <= 0.0 {
		mp.Progress = 0.0
		mp.Direction = 1.0
	}

	nextX := mp.StartX + dx*mp.Progress
	nextY := mp.StartY + dy*mp.Progress

	mp.DX = nextX - mp.X
	mp.DY = nextY - mp.Y

	mp.X = nextX
	mp.Y = nextY
}

func (mp *MovingPlatform) Draw(screen *ebiten.Image, camX float64) {
	sx := float32(mp.X - camX)
	if sx < -float32(mp.W) || sx > 330 {
		return
	}
	sy := float32(mp.Y)
	w := float32(mp.W)
	h := float32(mp.H)

	// Platform body
	vector.DrawFilledRect(screen, sx, sy, w, h, color.RGBA{45, 55, 75, 255}, false)
	// Platform top crest (metal highlight)
	vector.DrawFilledRect(screen, sx, sy, w, 2, color.RGBA{140, 180, 220, 255}, false)
	// Border
	vector.StrokeRect(screen, sx, sy, w, h, 1, color.RGBA{20, 25, 38, 255}, false)

	// Glowing energy thruster beneath
	glowW := w * 0.6
	glowX := sx + (w-glowW)*0.5
	vector.DrawFilledRect(screen, glowX, sy+h-1, glowW, 2, color.RGBA{80, 220, 255, 220}, false)
}
