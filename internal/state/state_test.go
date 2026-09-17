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
}
