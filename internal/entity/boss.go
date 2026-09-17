package entity

import (
	"image/color"
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/audio"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	BossWidth     = 24.0
	BossHeight    = 18.0
	BossMaxHealth = 8
)

// BossProjectile represents energy bolts fired by the Boss Overlord Mech.
type BossProjectile struct {
	X    float64
	Y    float64
	VX   float64
	VY   float64
	Dead bool
}

func (bp *BossProjectile) Update(dt float64) {
	bp.X += bp.VX * dt
	bp.Y += bp.VY * dt

	// Dead if out of screen / world bounds
	if bp.Y < -50 || bp.Y > 250 {
		bp.Dead = true
	}
}

func (bp *BossProjectile) GetRect() physics.Rect {
	return physics.Rect{
		X: bp.X + 1,
		Y: bp.Y + 1,
		W: 6,
		H: 6,
	}
}

func (bp *BossProjectile) Draw(screen *ebiten.Image, camX float64) {
	sx := bp.X - camX
	if sx < -10 || sx > 330 {
		return
	}
	// Glowing dual-tone plasma orb
	vector.DrawFilledCircle(screen, float32(sx+4), float32(bp.Y+4), 4, color.RGBA{255, 60, 80, 255}, false)
	vector.DrawFilledCircle(screen, float32(sx+4), float32(bp.Y+4), 2, color.RGBA{255, 240, 180, 255}, false)
}

// Boss represents the Europa Overlord Mech guarding the final sector extraction.
type Boss struct {
	X           float64
	Y           float64
	BaseY       float64
	VX          float64
	Health      int
	Dead        bool
	InvulnTimer float64
	ShootTimer  float64
	AnimTimer   float64
	FacingLeft  bool
	MinX        float64
	MaxX        float64
}

func NewBoss(x, y float64) *Boss {
	return &Boss{
		X:          x,
		Y:          y,
		BaseY:      y,
		VX:         -45.0,
		Health:     BossMaxHealth,
		ShootTimer: 1.5,
		MinX:       x - 120.0,
		MaxX:       x + 40.0,
		FacingLeft: true,
	}
}

func (b *Boss) GetRect() physics.Rect {
	return physics.Rect{
		X: b.X + 2,
		Y: b.Y + 2,
		W: BossWidth - 4,
		H: BossHeight - 4,
	}
}

func (b *Boss) TakeDamage(dmg int, particles *ParticleSystem) bool {
	if b.InvulnTimer > 0 || b.Dead {
		return false
	}
	b.Health -= dmg
	b.InvulnTimer = 0.45
	audio.Get().PlayStomp()
	particles.SpawnStompBurst(b.X+BossWidth/2, b.Y+BossHeight/2)

	if b.Health <= 0 {
		b.Health = 0
		b.Dead = true
		audio.Get().PlayPowerup()
		// Massive explosion fireworks
		for dx := -12; dx <= 12; dx += 8 {
			particles.SpawnStompBurst(b.X+BossWidth/2+float64(dx), b.Y+BossHeight/2)
			particles.SpawnCrystalSparkles(b.X+BossWidth/2+float64(dx), b.Y+BossHeight/2)
		}
	}
	return true
}

func (b *Boss) Update(dt float64, playerX, playerY float64, particles *ParticleSystem) *BossProjectile {
	if b.Dead {
		return nil
	}

	b.AnimTimer += dt
	if b.InvulnTimer > 0 {
		b.InvulnTimer -= dt
	}

	// Phase 2 triggers below 50% health: faster movement and attack rate
	speedMul := 1.0
	shootInterval := 2.0
	if b.Health <= 4 {
		speedMul = 1.6
		shootInterval = 1.2
	}

	// Horizontal patrol
	b.X += b.VX * speedMul * dt
	if b.X < b.MinX {
		b.X = b.MinX
		b.VX = math.Abs(b.VX)
	} else if b.X > b.MaxX {
		b.X = b.MaxX
		b.VX = -math.Abs(b.VX)
	}

	// Smooth hovering sine wave motion
	hoverAmp := 12.0
	hoverFreq := 2.5
	if b.Health <= 4 {
		hoverAmp = 18.0
		hoverFreq = 3.8
	}
	b.Y = b.BaseY + math.Sin(b.AnimTimer*hoverFreq)*hoverAmp

	// Face the player
	b.FacingLeft = (playerX < b.X)

	// Thruster exhaust particles
	if int(b.AnimTimer*40)%3 == 0 {
		particles.SpawnJetpackSparks(b.X+6, b.Y+BossHeight, false)
		particles.SpawnJetpackSparks(b.X+BossWidth-6, b.Y+BossHeight, false)
	}

	// Shooting behavior
	b.ShootTimer -= dt
	if b.ShootTimer <= 0 {
		b.ShootTimer = shootInterval

		// Calculate aim vector toward player
		dx := playerX - (b.X + BossWidth/2)
		dy := playerY - (b.Y + BossHeight/2)
		dist := math.Hypot(dx, dy)
		if dist < 1.0 {
			dist = 1.0
		}
		projSpeed := 95.0
		if b.Health <= 4 {
			projSpeed = 125.0
		}

		audio.Get().PlayLaser()
		return &BossProjectile{
			X:    b.X + BossWidth/2 - 4,
			Y:    b.Y + BossHeight/2,
			VX:   (dx / dist) * projSpeed,
			VY:   (dy / dist) * projSpeed,
			Dead: false,
		}
	}

	return nil
}

func (b *Boss) Draw(screen *ebiten.Image, camX float64) {
	if b.Dead {
		return
	}

	screenX := b.X - camX
	if screenX < -40 || screenX > 340 {
		return
	}

	atlas := assets.Get()
	var img *ebiten.Image

	if b.InvulnTimer > 0 && int(b.InvulnTimer*30)%2 == 0 {
		img = atlas.BossHurt
	} else if int(b.AnimTimer*8)%2 == 0 {
		img = atlas.Boss1
	} else {
		img = atlas.Boss2
	}

	op := &ebiten.DrawImageOptions{}
	if !b.FacingLeft {
		// Flip horizontally
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(BossWidth, 0)
	}
	op.GeoM.Translate(screenX, b.Y)
	screen.DrawImage(img, op)
}
