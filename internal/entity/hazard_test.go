package entity

import (
	"testing"
)

func TestCrumblingPlatformCycle(t *testing.T) {
	cp := NewCrumblingPlatform(100, 100, 32, 8)
	if cp.State != CrumbleIdle {
		t.Fatalf("expected CrumbleIdle, got %v", cp.State)
	}

	rect := cp.GetRect()
	if rect.W != 32 || rect.H != 8 {
		t.Errorf("unexpected rect: %+v", rect)
	}

	// Trigger touch
	cp.Touch()
	if cp.State != CrumbleShaking {
		t.Fatalf("expected CrumbleShaking, got %v", cp.State)
	}

	// Progress through shaking (0.75s)
	cp.Update(0.8)
	if cp.State != CrumbleFalling {
		t.Fatalf("expected CrumbleFalling, got %v", cp.State)
	}

	// Progress through falling (0.45s)
	cp.Update(0.5)
	if cp.State != CrumbleHidden {
		t.Fatalf("expected CrumbleHidden, got %v", cp.State)
	}
	if !cp.GetRect().IsEmpty() {
		t.Errorf("expected empty rect while hidden, got %+v", cp.GetRect())
	}

	// Progress through hidden respawn (3.0s)
	cp.Update(3.1)
	if cp.State != CrumbleIdle {
		t.Fatalf("expected CrumbleIdle after respawn, got %v", cp.State)
	}
	if cp.X != cp.StartX || cp.Y != cp.StartY {
		t.Errorf("expected position reset to (%f, %f), got (%f, %f)", cp.StartX, cp.StartY, cp.X, cp.Y)
	}
}

func TestHealthPickup(t *testing.T) {
	hp := NewHealthPickup(50, 60)
	if hp.Collected {
		t.Errorf("expected Collected=false")
	}

	r := hp.GetRect()
	if r.W != 12 || r.H != 12 {
		t.Errorf("expected 12x12 rect, got %+v", r)
	}

	hp.Collected = true
	if !hp.GetRect().IsEmpty() {
		t.Errorf("expected empty rect when collected")
	}
}

func TestThrusterBoostPickup(t *testing.T) {
	bp := NewThrusterBoostPickup(50, 60)
	if bp.Collected {
		t.Errorf("expected Collected=false")
	}

	r := bp.GetRect()
	if r.W != 12 || r.H != 12 {
		t.Errorf("expected 12x12 rect, got %+v", r)
	}

	bp.Collected = true
	if !bp.GetRect().IsEmpty() {
		t.Errorf("expected empty rect when collected")
	}
}

func TestPlayerHealAndBoost(t *testing.T) {
	p := NewPlayer(0, 0)
	p.Health = 2

	if !p.Heal(1) {
		t.Errorf("expected Heal(1) to succeed when health=2")
	}
	if p.Health != 3 {
		t.Errorf("expected health=3, got %d", p.Health)
	}

	// Heal at max health should return false
	if p.Heal(1) {
		t.Errorf("expected Heal to return false when at MaxHealth")
	}

	// Apply boost
	p.ApplyBoost(4.0)
	if p.BoostTimer != 4.0 {
		t.Errorf("expected BoostTimer 4.0, got %f", p.BoostTimer)
	}
	if p.Fuel != MaxFuel {
		t.Errorf("expected Fuel to be MaxFuel (%f), got %f", MaxFuel, p.Fuel)
	}
}
