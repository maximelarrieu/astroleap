package entity

import (
	"image/color"
	"testing"
)

func TestMovingPlatform(t *testing.T) {
	mp := NewMovingPlatform(100, 50, 200, 50, 32, 8, 50.0)

	rect := mp.GetRect()
	if rect.W != 32 || rect.H != 8 {
		t.Errorf("expected platform 32x8, got %+v", rect)
	}

	// Update platform
	mp.Update(0.1)

	if mp.X <= 100 {
		t.Errorf("expected platform to move right from 100, got X=%f", mp.X)
	}
	if mp.DX <= 0 {
		t.Errorf("expected positive DX, got %f", mp.DX)
	}
}

func TestCryoGeyser(t *testing.T) {
	g := NewCryoGeyser(80, 120, 60)
	rect := g.GetRect()

	if rect.X <= 0 || rect.Y != 60 || rect.H != 64 {
		t.Errorf("unexpected geyser rect: %+v", rect)
	}
}

func TestGoldMedal(t *testing.T) {
	m := NewGoldMedal(50, 60)
	if m.Collected {
		t.Errorf("medal should not be collected initially")
	}

	rect := m.GetRect()
	if rect.W != 12 || rect.H != 12 {
		t.Errorf("expected 12x12 medal rect, got %+v", rect)
	}
}

func TestFloatingText(t *testing.T) {
	ft := NewFloatingText("+100", 40, 50, color.RGBA{100, 245, 255, 255})
	if ft.Dead {
		t.Errorf("expected new floating text to be alive")
	}

	startY := ft.Y
	ft.Update(0.1)

	if ft.Y >= startY {
		t.Errorf("expected floating text to drift upward, got Y=%f", ft.Y)
	}

	// Fast-forward lifetime
	ft.Update(2.0)
	if !ft.Dead {
		t.Errorf("expected floating text to be dead after expiring")
	}
}

func TestPlayerStompComboAndFuel(t *testing.T) {
	p := NewPlayer(0, 0)
	p.Fuel = 10.0

	if p.StompCombo != 0 {
		t.Errorf("expected initial combo 0, got %d", p.StompCombo)
	}

	// First bounce
	p.Bounce()
	if p.StompCombo != 1 {
		t.Errorf("expected combo 1, got %d", p.StompCombo)
	}
	if p.Fuel <= 10.0 {
		t.Errorf("expected fuel to increase on bounce, got %f", p.Fuel)
	}

	// Second bounce in the air
	p.Bounce()
	if p.StompCombo != 2 {
		t.Errorf("expected combo 2, got %d", p.StompCombo)
	}
}
