package entity

import (
	"astroleap/internal/level"
	"testing"
)

func TestPlayerInitialization(t *testing.T) {
	p := NewPlayer(32, 100)
	if p.Health != MaxHealth {
		t.Errorf("expected player health to be %d, got %d", MaxHealth, p.Health)
	}
	if p.Fuel != MaxFuel {
		t.Errorf("expected player fuel to be %f, got %f", MaxFuel, p.Fuel)
	}
	if p.Dead {
		t.Errorf("expected player to start alive")
	}
}

func TestPlayerTakeDamage(t *testing.T) {
	p := NewPlayer(32, 100)
	p.TakeDamage(1)
	if p.Health != 2 {
		t.Errorf("expected player health to be 2, got %d", p.Health)
	}
	if p.InvulnTimer <= 0 {
		t.Errorf("expected invulnerability timer to be active")
	}

	// Invulnerable player should not take damage immediately
	p.TakeDamage(1)
	if p.Health != 2 {
		t.Errorf("expected health to stay 2 during invulnerability, got %d", p.Health)
	}
}

func TestPlayerFuelRecharge(t *testing.T) {
	p := NewPlayer(32, 100)
	p.Fuel = 10.0
	p.OnGround = true
	lvl := level.NewLevel1()
	ps := NewParticleSystem()

	p.Update(0.1, false, false, false, false, lvl, ps)
	if p.Fuel <= 10.0 {
		t.Errorf("expected fuel to recharge while on ground, got %f", p.Fuel)
	}
}

func TestPlayerGroundedStability(t *testing.T) {
	// Tile Y=9 is at 144. Grounded player should be at Y = 144 - 15 - 1 = 128.
	p := NewPlayer(32, 128)
	lvl := level.NewLevel1()
	ps := NewParticleSystem()

	dt := 1.0 / 60.0
	// Run 60 frames without input
	for frame := 0; frame < 60; frame++ {
		p.Update(dt, false, false, false, false, lvl, ps)
		if !p.OnGround {
			t.Fatalf("frame %d: expected player to remain grounded on floor, but OnGround is false (Y=%f, VY=%f)", frame, p.Y, p.VY)
		}
		if p.Y != 128.0 {
			t.Fatalf("frame %d: expected player Y to be stable at 128.0, got %f", frame, p.Y)
		}
		if p.State != StateIdle {
			t.Fatalf("frame %d: expected player state to be StateIdle, got %s", frame, p.State)
		}
	}
}

func TestPlayerWeaponUnlockAndSwitch(t *testing.T) {
	p := NewPlayer(32, 128)
	ps := NewParticleSystem()

	// Initial: No weapon
	if p.ActiveWeapon != WeaponNone {
		t.Errorf("expected initially no weapon, got %v", p.ActiveWeapon)
	}
	if proj := p.Shoot(ps); proj != nil {
		t.Errorf("expected cannot shoot without weapon")
	}

	// Unlock Laser
	p.UnlockWeapon(WeaponLaser)
	if !p.HasLaser || p.ActiveWeapon != WeaponLaser {
		t.Errorf("expected Laser to be active after unlock")
	}

	proj := p.Shoot(ps)
	if proj == nil || proj.Type != ProjLaser {
		t.Errorf("expected Laser projectile to be fired")
	}

	// Cannot shoot during cooldown
	if p.Shoot(ps) != nil {
		t.Errorf("expected shooting to be blocked by cooldown")
	}

	// Unlock Nova
	p.UnlockWeapon(WeaponNova)
	if !p.HasNova || p.ActiveWeapon != WeaponNova {
		t.Errorf("expected Nova to be equipped on unlock")
	}

	// Switch back to Laser
	p.SwitchWeapon()
	if p.ActiveWeapon != WeaponLaser {
		t.Errorf("expected active weapon to swap to Laser, got %v", p.ActiveWeapon)
	}

	// Switch back to Nova
	p.SwitchWeapon()
	if p.ActiveWeapon != WeaponNova {
		t.Errorf("expected active weapon to swap back to Nova, got %v", p.ActiveWeapon)
	}
}
