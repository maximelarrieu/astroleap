package entity

import (
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/audio"
	"astroleap/internal/level"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PlayerWidth     = 12.0
	PlayerHeight    = 15.0
	MoonGravity     = 0.16
	JumpImpulse     = -3.8
	JetpackThrust   = -0.28
	MaxRunSpeed     = 1.75
	MaxFallSpeed    = 3.5
	MaxThrustFall   = 1.0
	MaxFuel         = 100.0
	FuelBurnRate    = 32.0 // units per second
	FuelRecharge    = 48.0 // units per second
	MaxHealth       = 3
	InvulnDuration  = 1.2
)

type PlayerState string

const (
	StateIdle   PlayerState = "idle"
	StateRun    PlayerState = "run"
	StateJump   PlayerState = "jump"
	StateThrust PlayerState = "thrust"
)

type Player struct {
	X           float64
	Y           float64
	VX          float64
	VY          float64
	FacingRight bool
	OnGround    bool
	WasOnGround bool

	Health       int
	Fuel         float64
	IsThrusting  bool
	InvulnTimer  float64
	Dead         bool

	AnimTimer      float64
	ThrustSFXTimer float64
	State          PlayerState
	StompCombo     int
	BoostTimer     float64 // Active duration of infinite ion thruster boost

	// Weapon Inventory
	HasLaser      bool
	HasNova       bool
	ActiveWeapon  WeaponType
	ShootCooldown float64
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:            x,
		Y:            y,
		FacingRight:  true,
		Health:       MaxHealth,
		Fuel:         MaxFuel,
		State:        StateIdle,
		ActiveWeapon: WeaponNone,
	}
}

func (p *Player) GetRect() physics.Rect {
	return physics.Rect{
		X: p.X + 2,
		Y: p.Y + 1,
		W: PlayerWidth,
		H: PlayerHeight,
	}
}

func (p *Player) Update(dt float64, moveLeft, moveRight, jumpPressed, jumpHeld bool, lvl *level.Level, ps *ParticleSystem) {
	if p.Dead {
		return
	}

	p.AnimTimer += dt
	p.WasOnGround = p.OnGround

	// 1. Timers update
	if p.InvulnTimer > 0 {
		p.InvulnTimer -= dt
	}
	if p.ShootCooldown > 0 {
		p.ShootCooldown -= dt
	}

	// 2. Horizontal movement input
	accel := 0.22
	if !p.OnGround {
		accel = 0.14 // Air control in low gravity
	}

	if moveLeft {
		p.VX -= accel
		p.FacingRight = false
	} else if moveRight {
		p.VX += accel
		p.FacingRight = true
	} else {
		// Friction
		fric := 0.82
		if !p.OnGround {
			fric = 0.94 // Floaty space drift
		}
		p.VX *= fric
		if p.VX > -0.05 && p.VX < 0.05 {
			p.VX = 0
		}
	}

	// Clamp VX
	if p.VX > MaxRunSpeed {
		p.VX = MaxRunSpeed
	} else if p.VX < -MaxRunSpeed {
		p.VX = -MaxRunSpeed
	}

	// 3. Jump & Jetpack Thruster logic
	p.IsThrusting = false

	// Countdown active Ion Boost timer
	if p.BoostTimer > 0 {
		p.BoostTimer -= dt
		if p.BoostTimer < 0 {
			p.BoostTimer = 0
		}
	}

	if p.OnGround {
		p.Fuel += FuelRecharge * dt
		if p.Fuel > MaxFuel {
			p.Fuel = MaxFuel
		}

		if jumpPressed {
			p.VY = JumpImpulse
			p.OnGround = false
			audio.Get().PlayJump()
			ps.SpawnDustKick(p.X+8, p.Y+16)
		}
	} else {
		// Mid-air: Can use Jetpack if holding jump and has fuel (or ion boost active)
		if jumpHeld && (p.Fuel > 0 || p.BoostTimer > 0) && p.VY > -2.2 {
			p.IsThrusting = true
			p.VY += JetpackThrust
			if p.BoostTimer <= 0 {
				p.Fuel -= FuelBurnRate * dt
				if p.Fuel < 0 {
					p.Fuel = 0
				}
			} else {
				p.Fuel = MaxFuel // Supercharged infinite fuel
			}

			// Particle exhaust sparks
			thrusterX := p.X + 4
			if !p.FacingRight {
				thrusterX = p.X + 12
			}
			ps.SpawnJetpackSparks(thrusterX, p.Y+12, p.FacingRight)

			// Audio hiss
			p.ThrustSFXTimer -= dt
			if p.ThrustSFXTimer <= 0 {
				audio.Get().PlayThrust()
				p.ThrustSFXTimer = 0.12
			}
		}
	}

	// 4. Gravity application
	p.VY += MoonGravity
	maxFall := MaxFallSpeed
	if p.IsThrusting {
		maxFall = MaxThrustFall
	}
	if p.VY > maxFall {
		p.VY = maxFall
	}

	// 5. Physics Collision Resolution
	// 5a. Horizontal Step
	dx := p.VX * dt * 60.0
	if dx != 0 {
		newX := p.X + dx
		boxX := newX + 2.0
		boxY := p.Y + 1.0
		boxW := PlayerWidth
		boxH := PlayerHeight

		minTX := int(boxX / float64(level.TileSize))
		maxTX := int((boxX + boxW - 0.001) / float64(level.TileSize))
		minTY := int(boxY / float64(level.TileSize))
		maxTY := int((boxY + boxH - 0.001) / float64(level.TileSize))

		collidedX := false
		for ty := minTY; ty <= maxTY; ty++ {
			for tx := minTX; tx <= maxTX; tx++ {
				if lvl.IsSolid(tx, ty) {
					collidedX = true
					if p.VX > 0 {
						p.X = float64(tx*level.TileSize) - boxW - 2.0
					} else if p.VX < 0 {
						p.X = float64((tx+1)*level.TileSize) - 2.0
					}
					p.VX = 0
					break
				}
			}
			if collidedX {
				break
			}
		}
		if !collidedX {
			p.X = newX
		}
	}

	// Prevent walking off the left edge of the level
	if p.X < 0 {
		p.X = 0
		p.VX = 0
	}

	// 5b. Vertical Step
	dy := p.VY * dt * 60.0
	newY := p.Y + dy
	boxX := p.X + 2.0
	boxY := newY + 1.0
	boxW := PlayerWidth
	boxH := PlayerHeight

	minTX := int(boxX / float64(level.TileSize))
	maxTX := int((boxX + boxW - 0.001) / float64(level.TileSize))
	minTY := int(boxY / float64(level.TileSize))
	maxTY := int((boxY + boxH - 0.001) / float64(level.TileSize))

	collidedY := false
	p.OnGround = false

	for ty := minTY; ty <= maxTY; ty++ {
		for tx := minTX; tx <= maxTX; tx++ {
			if lvl.IsSolid(tx, ty) {
				collidedY = true
				if p.VY > 0 {
					// Falling down: landed on top of solid tile
					p.Y = float64(ty*level.TileSize) - boxH - 1.0
					p.OnGround = true
					p.StompCombo = 0
					p.VY = 0
					if !p.WasOnGround {
						ps.SpawnDustKick(p.X+8, p.Y+16)
					}
				} else if p.VY < 0 {
					// Moving up: hit ceiling
					p.Y = float64((ty+1)*level.TileSize) - 1.0
					p.VY = 0
				}
				break
			}
		}
		if collidedY {
			break
		}
	}
	if !collidedY {
		p.Y = newY
	}

	// Check spikes / hazards at feet
	feetTileX := int((p.X + 8.0) / float64(level.TileSize))
	feetTileY := int((p.Y + 16.5) / float64(level.TileSize))
	if lvl.IsHazard(feetTileX, feetTileY) {
		p.TakeDamage(1)
		p.VY = JumpImpulse * 0.8
	}

	// Fall out of level bottom check
	if p.Y > float64(level.LevelHeight*level.TileSize+32) {
		p.Health = 0
		p.Dead = true
	}

	// 6. Animation State Selection
	if !p.OnGround {
		if p.IsThrusting {
			p.State = StateThrust
		} else {
			p.State = StateJump
		}
	} else {
		if p.VX != 0 {
			p.State = StateRun
		} else {
			p.State = StateIdle
		}
	}
}

// TakeDamage reduces health and triggers invulnerability frames.
func (p *Player) TakeDamage(dmg int) {
	if p.InvulnTimer > 0 || p.Dead {
		return
	}
	p.Health -= dmg
	audio.Get().PlayHurt()
	p.InvulnTimer = InvulnDuration
	if p.Health <= 0 {
		p.Health = 0
		p.Dead = true
	}
}

// Bounce is called after stomping an enemy to pop the astronaut into the air.
func (p *Player) Bounce() {
	p.VY = JumpImpulse * 1.15
	p.OnGround = false
	p.StompCombo++
	p.Fuel = math.Min(MaxFuel, p.Fuel+30.0)
}

// UnlockWeapon unlocks and equips a new weapon tier.
func (p *Player) UnlockWeapon(w WeaponType) {
	if w == WeaponLaser {
		p.HasLaser = true
		p.ActiveWeapon = WeaponLaser
	} else if w == WeaponNova {
		p.HasNova = true
		p.ActiveWeapon = WeaponNova
	}
}

// SwitchWeapon toggles between unlocked weapons.
func (p *Player) SwitchWeapon() {
	if !p.HasLaser && !p.HasNova {
		return
	}
	if p.ActiveWeapon == WeaponLaser {
		if p.HasNova {
			p.ActiveWeapon = WeaponNova
		}
	} else if p.ActiveWeapon == WeaponNova {
		if p.HasLaser {
			p.ActiveWeapon = WeaponLaser
		}
	} else {
		if p.HasLaser {
			p.ActiveWeapon = WeaponLaser
		} else if p.HasNova {
			p.ActiveWeapon = WeaponNova
		}
	}
}

// Shoot fires a projectile from the active weapon.
func (p *Player) Shoot(ps *ParticleSystem) *Projectile {
	if p.Dead || p.ShootCooldown > 0 || p.ActiveWeapon == WeaponNone {
		return nil
	}

	spawnX := p.X + 13.0
	if !p.FacingRight {
		spawnX = p.X - 5.0
	}
	spawnY := p.Y + 7.0

	switch p.ActiveWeapon {
	case WeaponLaser:
		p.ShootCooldown = 0.20 // 5 shots/sec
		audio.Get().PlayLaser()
		ps.SpawnJetpackSparks(spawnX, spawnY, p.FacingRight)
		return NewLaserProjectile(spawnX, spawnY, p.FacingRight)

	case WeaponNova:
		p.ShootCooldown = 0.38 // Heavier, bouncing explosive
		audio.Get().PlayNova()
		ps.SpawnJetpackSparks(spawnX, spawnY, p.FacingRight)
		return NewNovaProjectile(spawnX, spawnY, p.FacingRight)
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image, camX float64) {
	if p.Dead {
		return
	}

	// Hurt blinking: blink every 0.1s during invulnerability
	if p.InvulnTimer > 0 {
		blinkPhase := int(p.InvulnTimer*15.0) % 2
		if blinkPhase == 0 {
			return
		}
	}

	screenX := p.X - camX
	screenY := p.Y

	atlas := assets.Get()
	var img *ebiten.Image

	switch p.State {
	case StateThrust:
		img = atlas.AstronautThrust
	case StateJump:
		img = atlas.AstronautJump
	case StateRun:
		frame := int(p.AnimTimer*8.0) % 2
		if frame == 0 {
			img = atlas.AstronautRun1
		} else {
			img = atlas.AstronautRun2
		}
	default: // Idle
		img = atlas.AstronautIdle
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		if !p.FacingRight {
			// Mirror horizontally
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(16, 0)
		}
		op.GeoM.Translate(screenX, screenY)
		screen.DrawImage(img, op)
	}
}

// Heal restores health up to MaxHealth. Returns true if health was restored.
func (p *Player) Heal(amount int) bool {
	if p.Health >= MaxHealth {
		return false
	}
	p.Health += amount
	if p.Health > MaxHealth {
		p.Health = MaxHealth
	}
	return true
}

// ApplyBoost grants temporary infinite ion thruster fuel.
func (p *Player) ApplyBoost(duration float64) {
	p.BoostTimer = duration
	p.Fuel = MaxFuel
}
