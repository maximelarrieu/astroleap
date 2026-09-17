package entity

import (
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

// WeaponPickup represents a collectible holographic pod that unlocks a weapon.
type WeaponPickup struct {
	X         float64
	Y         float64
	BaseY     float64
	Weapon    WeaponType
	Collected bool
	Timer     float64
}

func NewWeaponPickup(x, y float64, weapon WeaponType) *WeaponPickup {
	return &WeaponPickup{
		X:      x,
		Y:      y,
		BaseY:  y,
		Weapon: weapon,
	}
}

func (w *WeaponPickup) Update(dt float64) {
	if w.Collected {
		return
	}
	w.Timer += dt
	// Gentle hover bobbing
	w.Y = w.BaseY + math.Sin(w.Timer*3.5)*3.0
}

func (w *WeaponPickup) GetRect() physics.Rect {
	if w.Collected {
		return physics.Rect{}
	}
	return physics.Rect{
		X: w.X + 2,
		Y: w.Y + 2,
		W: 12,
		H: 12,
	}
}

func (w *WeaponPickup) Draw(screen *ebiten.Image, camX float64) {
	if w.Collected {
		return
	}
	screenX := w.X - camX
	if screenX < -20 || screenX > 330 {
		return
	}

	atlas := assets.Get()
	var img *ebiten.Image
	if w.Weapon == WeaponLaser {
		img = atlas.CapsuleLaser
	} else if w.Weapon == WeaponNova {
		img = atlas.CapsuleNova
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenX, w.Y)
		screen.DrawImage(img, op)
	}
}
