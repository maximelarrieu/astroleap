package entity

import (
	"astroleap/internal/assets"
	"astroleap/internal/level"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type WeaponType int

const (
	WeaponNone WeaponType = iota
	WeaponLaser
	WeaponNova
)

func (w WeaponType) String() string {
	switch w {
	case WeaponLaser:
		return "LASER BLASTER"
	case WeaponNova:
		return "NOVA CANNON"
	default:
		return "NONE"
	}
}

type ProjectileType int

const (
	ProjLaser ProjectileType = iota
	ProjNova
)

type Projectile struct {
	X           float64
	Y           float64
	VX          float64
	VY          float64
	Type        ProjectileType
	Life        float64
	Dead        bool
	Bounces     int
	FacingRight bool
}

func NewLaserProjectile(x, y float64, facingRight bool) *Projectile {
	vx := 4.8
	if !facingRight {
		vx = -4.8
	}
	return &Projectile{
		X:           x,
		Y:           y,
		VX:          vx,
		VY:          0,
		Type:        ProjLaser,
		Life:        1.4,
		FacingRight: facingRight,
	}
}

func NewNovaProjectile(x, y float64, facingRight bool) *Projectile {
	vx := 2.6
	if !facingRight {
		vx = -2.6
	}
	return &Projectile{
		X:           x,
		Y:           y,
		VX:          vx,
		VY:          -1.4,
		Type:        ProjNova,
		Life:        3.2,
		FacingRight: facingRight,
	}
}

func (p *Projectile) GetRect() physics.Rect {
	if p.Dead {
		return physics.Rect{}
	}
	if p.Type == ProjLaser {
		return physics.Rect{
			X: p.X,
			Y: p.Y,
			W: 8,
			H: 4,
		}
	}
	return physics.Rect{
		X: p.X + 1,
		Y: p.Y + 1,
		W: 6,
		H: 6,
	}
}

func (p *Projectile) Update(dt float64, lvl *level.Level, ps *ParticleSystem) {
	if p.Dead {
		return
	}

	p.Life -= dt
	if p.Life <= 0 {
		p.Dead = true
		return
	}

	if p.Type == ProjLaser {
		// Linear rapid flight
		newX := p.X + p.VX*dt*60.0
		// Check wall collision
		tileX := int((newX + 4.0) / float64(level.TileSize))
		tileY := int((p.Y + 2.0) / float64(level.TileSize))

		if lvl.IsSolid(tileX, tileY) {
			p.Dead = true
			ps.SpawnCrystalSparkles(newX, p.Y)
			return
		}
		p.X = newX

	} else if p.Type == ProjNova {
		// Bouncing projectile with low gravity
		p.VY += MoonGravity * 0.9
		if p.VY > 3.0 {
			p.VY = 3.0
		}

		// Horizontal movement
		newX := p.X + p.VX*dt*60.0
		tileX := int((newX + 4.0) / float64(level.TileSize))
		tileY := int((p.Y + 4.0) / float64(level.TileSize))

		if lvl.IsSolid(tileX, tileY) {
			// Bounce off wall
			p.VX = -p.VX * 0.85
			p.FacingRight = (p.VX > 0)
			ps.SpawnJetpackSparks(p.X+4, p.Y+4, p.FacingRight)
		} else {
			p.X = newX
		}

		// Vertical movement
		newY := p.Y + p.VY*dt*60.0
		checkTileX := int((p.X + 4.0) / float64(level.TileSize))
		checkTileY := int((newY + 7.0) / float64(level.TileSize))

		if lvl.IsSolid(checkTileX, checkTileY) {
			// Bounce off floor
			p.VY = -2.2
			p.Bounces++
			ps.SpawnJetpackSparks(p.X+4, newY+6, true)
			if p.Bounces > 4 {
				p.Dead = true
			}
		} else {
			p.Y = newY
		}
	}
}

func (p *Projectile) Draw(screen *ebiten.Image, camX float64) {
	if p.Dead {
		return
	}
	screenX := p.X - camX
	if screenX < -16 || screenX > 330 {
		return
	}

	atlas := assets.Get()
	op := &ebiten.DrawImageOptions{}

	if p.Type == ProjLaser {
		if !p.FacingRight {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(8, 0)
		}
		op.GeoM.Translate(screenX, p.Y)
		screen.DrawImage(atlas.LaserBolt, op)
	} else if p.Type == ProjNova {
		op.GeoM.Translate(screenX, p.Y)
		screen.DrawImage(atlas.NovaBall, op)
	}
}
