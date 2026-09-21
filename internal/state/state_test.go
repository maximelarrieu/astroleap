package state

import (
	"testing"

	"astroleap/internal/entity"
)

func TestPlayStateWithProgress(t *testing.T) {
	m := NewMachine(nil)

	carriedPlayer := entity.NewPlayer(0, 0)
	carriedPlayer.UnlockWeapon(entity.WeaponLaser)
	carriedPlayer.UnlockWeapon(entity.WeaponNova)
	carriedPlayer.ActiveWeapon = entity.WeaponNova

	ps := NewPlayStateWithProgress(m, 2, 1250, 14, carriedPlayer)

	if ps.currentSector != 2 {
		t.Errorf("expected currentSector 2, got %d", ps.currentSector)
	}

	if ps.score != 1250 {
		t.Errorf("expected score 1250, got %d", ps.score)
	}

	if ps.crystalsCollected != 14 {
		t.Errorf("expected crystals 14, got %d", ps.crystalsCollected)
	}

	if !ps.player.HasLaser || !ps.player.HasNova {
		t.Errorf("expected carried weapons to be preserved")
	}

	if ps.player.ActiveWeapon != entity.WeaponNova {
		t.Errorf("expected active weapon to be WeaponNova, got %v", ps.player.ActiveWeapon)
	}

	if ps.lvl == nil || ps.lvl.SectorIndex != 2 {
		t.Errorf("expected level sector 2 to be loaded")
	}
}

func TestGameOverStatePreservesSector(t *testing.T) {
	m := NewMachine(nil)
	p := entity.NewPlayer(0, 0)
	p.UnlockWeapon(entity.WeaponLaser)

	gos := NewGameOverState(m, 500, 5, 2, p)
	if gos.sector != 2 {
		t.Errorf("expected game over sector 2, got %d", gos.sector)
	}
	if gos.score != 500 {
		t.Errorf("expected score 500, got %d", gos.score)
	}
	if gos.playerRef == nil || !gos.playerRef.HasLaser {
		t.Errorf("expected playerRef to retain laser unlock")
	}
}

func TestSectorClearState(t *testing.T) {
	m := NewMachine(nil)
	p := entity.NewPlayer(0, 0)

	scs := NewSectorClearState(m, 1, 800, 8, p)
	if scs.sector != 1 {
		t.Errorf("expected sector 1, got %d", scs.sector)
	}
	if scs.sectorName != "LUNAR OUTPOST" {
		t.Errorf("expected sector name LUNAR OUTPOST, got %s", scs.sectorName)
	}

	scs2 := NewSectorClearState(m, 2, 1500, 16, p)
	if scs2.sectorName != "PHOBOS RIDGE" {
		t.Errorf("expected sector name PHOBOS RIDGE, got %s", scs2.sectorName)
	}

	scs3 := NewSectorClearState(m, 3, 2500, 24, p)
	if scs3.sectorName != "EUROPA ICE CORE" {
		t.Errorf("expected sector name EUROPA ICE CORE, got %s", scs3.sectorName)
	}

	scs4 := NewSectorClearState(m, 4, 3500, 32, p)
	if scs4.sectorName != "MOTHERSHIP CORRIDOR" {
		t.Errorf("expected sector name MOTHERSHIP CORRIDOR, got %s", scs4.sectorName)
	}
}

func TestPlayStateHazardsAndPickups(t *testing.T) {
	m := NewMachine(nil)
	ps := NewPlayStateWithProgress(m, 1, 0, 0, nil)

	if len(ps.crumblingPlatforms) == 0 {
		t.Errorf("expected sector 1 to have crumbling platforms initialized")
	}
	if len(ps.healthPickups) == 0 {
		t.Errorf("expected sector 1 to have health pickups initialized")
	}
	if len(ps.boostPickups) == 0 {
		t.Errorf("expected sector 1 to have boost pickups initialized")
	}
}

func TestSector4KeyAndAirlockMechanic(t *testing.T) {
	m := NewMachine(nil)
	ps := NewPlayStateWithProgress(m, 4, 1000, 10, nil)

	if ps.securityKey == nil {
		t.Fatalf("expected sector 4 to have security key")
	}
	if ps.hasSecurityKey {
		t.Errorf("hasSecurityKey should start false")
	}

	// Move player to key
	ps.player.X = ps.securityKey.X
	ps.player.Y = ps.securityKey.Y
	ps.player.VY = 0

	ps.Update(0.016)

	if !ps.hasSecurityKey {
		t.Errorf("expected hasSecurityKey to become true after overlap")
	}
	if !ps.securityKey.Collected {
		t.Errorf("expected securityKey.Collected to be true")
	}
}

func TestSector5RepairCoreMechanic(t *testing.T) {
	m := NewMachine(nil)
	ps := NewPlayStateWithProgress(m, 5, 2000, 20, nil)

	if len(ps.repairCores) != 3 {
		t.Fatalf("expected 3 repair cores in sector 5, got %d", len(ps.repairCores))
	}
	if ps.repairCoresCount != 0 {
		t.Errorf("expected 0 initial repair cores collected")
	}

	// Collect first core
	ps.player.X = ps.repairCores[0].X
	ps.player.Y = ps.repairCores[0].Y
	ps.player.VY = 0
	ps.Update(0.016)

	if ps.repairCoresCount != 1 {
		t.Errorf("expected 1 repair core collected, got %d", ps.repairCoresCount)
	}

	// Collect remaining cores
	ps.repairCores[1].Collected = true
	ps.repairCores[2].Collected = true
	ps.repairCoresCount = 3

	// Test reaching goal with all 3 cores
	ps.player.X = ps.lvl.LanderX
	ps.player.Y = ps.lvl.LanderY
	ps.player.VY = 0
	ps.Update(0.016)

	if _, ok := m.Current().(*WinState); !ok {
		t.Errorf("expected transition to WinState after repairing all 3 cores, got %T", m.Current())
	}
}

func TestMovingPlatformCarriage(t *testing.T) {
	m := NewMachine(nil)
	ps := NewPlayStateWithProgress(m, 2, 0, 0, nil)

	if len(ps.movingPlatforms) == 0 {
		t.Fatalf("expected sector 2 to have moving platform")
	}

	mp := ps.movingPlatforms[0]
	// Position player on platform top
	ps.player.X = mp.X + 4
	ps.player.Y = mp.Y - entity.PlayerHeight
	ps.player.VY = 0
	ps.player.OnGround = true

	// Perform snap
	ps.Update(0.016)

	if ps.ridingPlatform != mp {
		t.Errorf("expected player to ride platform, got %+v", ps.ridingPlatform)
	}

	startX := ps.player.X
	startY := ps.player.Y

	// Next tick: platform moves, player should be carried
	ps.Update(0.016)

	if ps.player.X == startX && mp.DX != 0 {
		t.Errorf("expected player X to move with platform DX")
	}
	if ps.player.Y > startY+10 {
		t.Errorf("player fell through platform: Y=%f startY=%f", ps.player.Y, startY)
	}
}
