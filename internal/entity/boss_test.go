package entity

import (
	"testing"
)

func TestBossInitializationAndDamage(t *testing.T) {
	boss := NewBoss(200, 100)
	particles := NewParticleSystem()

	if boss.Health != BossMaxHealth {
		t.Errorf("expected boss health %d, got %d", BossMaxHealth, boss.Health)
	}

	if boss.Dead {
		t.Errorf("expected boss to be alive initially")
	}

	rect := boss.GetRect()
	if rect.W <= 0 || rect.H <= 0 {
		t.Errorf("invalid boss rect: %+v", rect)
	}

	// Deal 2 damage
	damaged := boss.TakeDamage(2, particles)
	if !damaged || boss.Health != BossMaxHealth-2 {
		t.Errorf("expected health %d after 2 dmg, got %d", BossMaxHealth-2, boss.Health)
	}

	// Invulnerability prevents instant damage on the same frame
	secondHit := boss.TakeDamage(1, particles)
	if secondHit {
		t.Errorf("expected second hit to be rejected during invulnerability")
	}

	// Fast-forward invulnerability
	boss.InvulnTimer = 0

	// Fatal damage
	boss.TakeDamage(BossMaxHealth, particles)
	if !boss.Dead || boss.Health != 0 {
		t.Errorf("expected boss to be dead with 0 health, got dead=%v, hp=%d", boss.Dead, boss.Health)
	}
}

func TestBossAIAndProjectiles(t *testing.T) {
	boss := NewBoss(200, 100)
	particles := NewParticleSystem()

	// Update boss and advance shoot timer
	boss.ShootTimer = 0.01
	proj := boss.Update(0.02, 50, 100, particles)

	if proj == nil {
		t.Fatalf("expected boss to fire a projectile when shoot timer elapses")
	}

	if proj.VX >= 0 {
		t.Errorf("expected projectile to aim left toward player (playerX=50, bossX=200), got VX=%f", proj.VX)
	}

	projRect := proj.GetRect()
	if projRect.W <= 0 || projRect.H <= 0 {
		t.Errorf("invalid boss projectile rect: %+v", projRect)
	}
}
