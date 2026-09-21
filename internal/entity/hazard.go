package entity

import (
	"image/color"
	"math"
	"math/rand"

	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CrumbleState int

const (
	CrumbleIdle CrumbleState = iota
	CrumbleShaking
	CrumbleFalling
	CrumbleHidden
)

// CrumblingPlatform represents a fragile lunar rock slab that fractures and drops
// shortly after the astronaut lands on it, respawning after a delay.
type CrumblingPlatform struct {
	StartX       float64
	StartY       float64
	X            float64
	Y            float64
	W            float64
	H            float64
	State        CrumbleState
	StateTimer   float64
	ShakeOffsetX float64
}

func NewCrumblingPlatform(x, y, w, h float64) *CrumblingPlatform {
	return &CrumblingPlatform{
		StartX: x,
		StartY: y,
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
		State:  CrumbleIdle,
	}
}

func (cp *CrumblingPlatform) GetRect() physics.Rect {
	if cp.State == CrumbleHidden {
		return physics.Rect{}
	}
	return physics.Rect{
		X: cp.X,
		Y: cp.Y,
		W: cp.W,
		H: cp.H,
	}
}

func (cp *CrumblingPlatform) Touch() {
	if cp.State == CrumbleIdle {
		cp.State = CrumbleShaking
		cp.StateTimer = 0.75 // 0.75s before collapse
	}
}

func (cp *CrumblingPlatform) Update(dt float64) {
	switch cp.State {
	case CrumbleShaking:
		cp.StateTimer -= dt
		cp.ShakeOffsetX = (rand.Float64()*2.0 - 1.0) * 1.5
		if cp.StateTimer <= 0 {
			cp.State = CrumbleFalling
			cp.StateTimer = 0.45 // Fall downward for 0.45s before disappearing
			cp.ShakeOffsetX = 0
		}
	case CrumbleFalling:
		cp.Y += 130.0 * dt
		cp.StateTimer -= dt
		if cp.StateTimer <= 0 {
			cp.State = CrumbleHidden
			cp.StateTimer = 3.0 // Respawn after 3.0s
		}
	case CrumbleHidden:
		cp.StateTimer -= dt
		if cp.StateTimer <= 0 {
			cp.X = cp.StartX
			cp.Y = cp.StartY
			cp.State = CrumbleIdle
			cp.ShakeOffsetX = 0
		}
	}
}

func (cp *CrumblingPlatform) Draw(screen *ebiten.Image, camX float64) {
	if cp.State == CrumbleHidden {
		return
	}
	sx := float32(cp.X + cp.ShakeOffsetX - camX)
	if sx < -float32(cp.W) || sx > 330 {
		return
	}
	sy := float32(cp.Y)
	w := float32(cp.W)
	h := float32(cp.H)

	fillColor := color.RGBA{65, 70, 85, 255}
	accentColor := color.RGBA{130, 135, 150, 255}
	if cp.State == CrumbleShaking {
		fillColor = color.RGBA{95, 60, 65, 255}
		accentColor = color.RGBA{235, 95, 85, 255}
	} else if cp.State == CrumbleFalling {
		fillColor = color.RGBA{45, 45, 55, 180}
	}

	vector.DrawFilledRect(screen, sx, sy, w, h, fillColor, false)
	vector.DrawFilledRect(screen, sx, sy, w, 2, accentColor, false)
	vector.StrokeRect(screen, sx, sy, w, h, 1, color.RGBA{25, 25, 35, 255}, false)
}

// HealthPickup represents an emergency Oxygen cylinder restoring astronaut life.
type HealthPickup struct {
	X         float64
	Y         float64
	BaseY     float64
	Timer     float64
	Collected bool
}

func NewHealthPickup(x, y float64) *HealthPickup {
	return &HealthPickup{
		X:     x,
		Y:     y,
		BaseY: y,
	}
}

func (h *HealthPickup) Update(dt float64) {
	if h.Collected {
		return
	}
	h.Timer += dt
	h.Y = h.BaseY + math.Sin(h.Timer*4.0)*2.5
}

func (h *HealthPickup) GetRect() physics.Rect {
	if h.Collected {
		return physics.Rect{}
	}
	return physics.Rect{
		X: h.X + 2,
		Y: h.Y + 2,
		W: 12,
		H: 12,
	}
}

func (h *HealthPickup) Draw(screen *ebiten.Image, camX float64) {
	if h.Collected {
		return
	}
	sx := float32(h.X - camX)
	if sx < -16 || sx > 330 {
		return
	}
	sy := float32(h.Y)

	// Canister body (cyan cylinder with white/red cross)
	vector.DrawFilledRect(screen, sx+3, sy+1, 8, 12, color.RGBA{40, 180, 200, 255}, false)
	vector.DrawFilledRect(screen, sx+5, sy-1, 4, 3, color.RGBA{180, 210, 230, 255}, false)
	vector.DrawFilledRect(screen, sx+5, sy+5, 4, 4, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, sx+6, sy+4, 2, 6, color.RGBA{230, 40, 60, 255}, false)
	vector.DrawFilledRect(screen, sx+4, sy+6, 6, 2, color.RGBA{230, 40, 60, 255}, false)
}

// ThrusterBoostPickup represents an Ion Core cell giving supercharged infinite thrust.
type ThrusterBoostPickup struct {
	X         float64
	Y         float64
	BaseY     float64
	Timer     float64
	Collected bool
}

func NewThrusterBoostPickup(x, y float64) *ThrusterBoostPickup {
	return &ThrusterBoostPickup{
		X:     x,
		Y:     y,
		BaseY: y,
	}
}

func (b *ThrusterBoostPickup) Update(dt float64) {
	if b.Collected {
		return
	}
	b.Timer += dt
	b.Y = b.BaseY + math.Sin(b.Timer*4.5)*2.5
}

func (b *ThrusterBoostPickup) GetRect() physics.Rect {
	if b.Collected {
		return physics.Rect{}
	}
	return physics.Rect{
		X: b.X + 2,
		Y: b.Y + 2,
		W: 12,
		H: 12,
	}
}

func (b *ThrusterBoostPickup) Draw(screen *ebiten.Image, camX float64) {
	if b.Collected {
		return
	}
	sx := float32(b.X - camX)
	if sx < -16 || sx > 330 {
		return
	}
	sy := float32(b.Y)

	// Ion core glowing diamond / cell
	vector.DrawFilledRect(screen, sx+2, sy+2, 10, 10, color.RGBA{255, 170, 30, 255}, false)
	vector.DrawFilledRect(screen, sx+4, sy+4, 6, 6, color.RGBA{255, 240, 100, 255}, false)
	vector.StrokeRect(screen, sx+2, sy+2, 10, 10, 1, color.RGBA{255, 90, 20, 255}, false)
}
