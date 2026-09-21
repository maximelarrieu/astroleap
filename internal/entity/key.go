package entity

import (
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

// SecurityKey represents the golden digital access keycard needed in Sector 4
// to unlock the bulkhead airlock leading to Sector 5's Reactor Room.
type SecurityKey struct {
	X         float64
	Y         float64
	BaseY     float64
	Collected bool
	Timer     float64
}

func NewSecurityKey(x, y float64) *SecurityKey {
	return &SecurityKey{
		X:     x,
		Y:     y,
		BaseY: y,
	}
}

func (k *SecurityKey) Update(dt float64) {
	if k.Collected {
		return
	}
	k.Timer += dt
	k.Y = k.BaseY + math.Sin(k.Timer*4.0)*2.5
}

func (k *SecurityKey) GetRect() physics.Rect {
	if k.Collected {
		return physics.Rect{}
	}
	return physics.Rect{
		X: k.X + 2,
		Y: k.Y + 2,
		W: 12,
		H: 12,
	}
}

func (k *SecurityKey) Draw(screen *ebiten.Image, camX float64) {
	if k.Collected {
		return
	}
	sx := k.X - camX
	if sx < -20 || sx > 330 {
		return
	}

	atlas := assets.Get()
	if atlas.KeyCard != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, k.Y)
		screen.DrawImage(atlas.KeyCard, op)
	}
}

// RepairCore represents a critical warp reactor fuel rod / repair node in Sector 5.
// Collecting all 3 restores the spaceship's systems and powers up the warp drive.
type RepairCore struct {
	Index     int
	Name      string
	X         float64
	Y         float64
	BaseY     float64
	Collected bool
	Timer     float64
}

func NewRepairCore(index int, name string, x, y float64) *RepairCore {
	return &RepairCore{
		Index: index,
		Name:  name,
		X:     x,
		Y:     y,
		BaseY: y,
	}
}

func (rc *RepairCore) Update(dt float64) {
	if rc.Collected {
		return
	}
	rc.Timer += dt
	rc.Y = rc.BaseY + math.Sin(rc.Timer*4.2+float64(rc.Index))*2.5
}

func (rc *RepairCore) GetRect() physics.Rect {
	if rc.Collected {
		return physics.Rect{}
	}
	return physics.Rect{
		X: rc.X + 2,
		Y: rc.Y + 2,
		W: 12,
		H: 12,
	}
}

func (rc *RepairCore) Draw(screen *ebiten.Image, camX float64) {
	if rc.Collected {
		return
	}
	sx := rc.X - camX
	if sx < -20 || sx > 330 {
		return
	}

	atlas := assets.Get()
	if atlas.RepairCore != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, rc.Y)
		screen.DrawImage(atlas.RepairCore, op)
	}
}
