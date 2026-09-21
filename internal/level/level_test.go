package level

import "testing"

func TestLoadAllSectors(t *testing.T) {
	for sector := 1; sector <= MaxLevels; sector++ {
		lvl := LoadLevel(sector)
		if lvl == nil {
			t.Fatalf("LoadLevel(%d) returned nil", sector)
		}

		if lvl.SectorIndex != sector {
			t.Errorf("expected SectorIndex %d, got %d", sector, lvl.SectorIndex)
		}

		if lvl.Width < 100 || lvl.Height != LevelHeight {
			t.Errorf("sector %d invalid dimensions: %dx%d", sector, lvl.Width, lvl.Height)
		}

		if len(lvl.Crystals) == 0 {
			t.Errorf("sector %d has no crystals", sector)
		}

		if len(lvl.EnemySpawns) == 0 {
			t.Errorf("sector %d has no enemy spawns", sector)
		}

		if lvl.LanderX <= 0 || lvl.LanderX >= float64(lvl.Width*TileSize) {
			t.Errorf("sector %d LanderX (%f) out of bounds", sector, lvl.LanderX)
		}

		landerRect := lvl.GetLanderRect()
		if landerRect.W <= 0 || landerRect.H <= 0 {
			t.Errorf("sector %d invalid lander rect: %+v", sector, landerRect)
		}

		// Landing spawn point (x=2, y=9) must be solid ground
		if !lvl.IsSolid(2, 9) {
			t.Errorf("sector %d landing spawn ground (2, 9) not solid", sector)
		}

		if len(lvl.CrumblingPlatforms) == 0 {
			t.Errorf("sector %d has no crumbling platforms", sector)
		}

		if len(lvl.HealthPickups) == 0 {
			t.Errorf("sector %d has no health pickups", sector)
		}

		if len(lvl.BoostPickups) == 0 {
			t.Errorf("sector %d has no boost pickups", sector)
		}
	}
}

func TestSectorThemesAndWeapons(t *testing.T) {
	s1 := LoadLevel(1)
	if s1.Theme != "moon" || s1.Celestial != "earth" {
		t.Errorf("sector 1 unexpected theme/celestial: %s/%s", s1.Theme, s1.Celestial)
	}
	hasLaserSpawn := false
	for _, w := range s1.WeaponSpawns {
		if w.Kind == "laser" {
			hasLaserSpawn = true
			break
		}
	}
	if !hasLaserSpawn {
		t.Errorf("sector 1 missing laser weapon spawn")
	}

	s2 := LoadLevel(2)
	if s2.Theme != "crimson" || s2.Celestial != "mars" {
		t.Errorf("sector 2 unexpected theme/celestial: %s/%s", s2.Theme, s2.Celestial)
	}
	hasNovaSpawn := false
	for _, w := range s2.WeaponSpawns {
		if w.Kind == "nova" {
			hasNovaSpawn = true
			break
		}
	}
	if !hasNovaSpawn {
		t.Errorf("sector 2 missing nova weapon spawn")
	}

	s3 := LoadLevel(3)
	if s3.Theme != "ice" || s3.Celestial != "jupiter" {
		t.Errorf("sector 3 unexpected theme/celestial: %s/%s", s3.Theme, s3.Celestial)
	}

	s4 := LoadLevel(4)
	if s4.Theme != "vessel" || s4.GoalKind != "airlock" {
		t.Errorf("sector 4 unexpected theme/goal: %s/%s", s4.Theme, s4.GoalKind)
	}
	if s4.SecurityKey == nil {
		t.Errorf("sector 4 missing SecurityKey")
	}

	s5 := LoadLevel(5)
	if s5.Theme != "reactor" || s5.GoalKind != "warp_console" {
		t.Errorf("sector 5 unexpected theme/goal: %s/%s", s5.Theme, s5.GoalKind)
	}
	if len(s5.RepairCores) != 3 {
		t.Errorf("sector 5 expected 3 repair cores, got %d", len(s5.RepairCores))
	}
}

func TestMovingPlatformsDoNotOverlapSolidTiles(t *testing.T) {
	for sector := 1; sector <= MaxLevels; sector++ {
		lvl := LoadLevel(sector)
		for idx, mp := range lvl.MovingPlatforms {
			// Check bounding box along travel from Start to End
			minX := int(mp.StartX / float64(TileSize))
			maxX := int((mp.EndX + mp.W - 0.001) / float64(TileSize))
			if minX > maxX {
				minX, maxX = maxX, minX
			}
			minY := int(mp.StartY / float64(TileSize))
			maxY := int((mp.EndY + mp.H - 0.001) / float64(TileSize))
			if minY > maxY {
				minY, maxY = maxY, minY
			}

			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {
					if lvl.IsSolid(x, y) {
						t.Errorf("sector %d moving platform %d at lane (%d..%d, %d..%d) overlaps solid tile at (%d, %d)",
							sector, idx, minX, maxX, minY, maxY, x, y)
					}
				}
			}

			// Check that crumbling platforms do not overlap moving platforms
			for cIdx, cp := range lvl.CrumblingPlatforms {
				cpX := float64(cp.TileX * TileSize)
				cpY := float64(cp.TileY * TileSize)
				// Check horizontal and vertical range
				if (cpX+cp.W > mp.StartX && cpX < mp.EndX+mp.W) && (cpY+cp.H > mp.StartY && cpY < mp.EndY+mp.H) {
					t.Errorf("sector %d moving platform %d overlaps crumbling platform %d at (%f, %f)", sector, idx, cIdx, cpX, cpY)
				}
			}
		}
	}
}

func TestFallbackLevel(t *testing.T) {
	fallback := LoadLevel(999)
	if fallback.SectorIndex != 1 {
		t.Errorf("expected fallback sector 1, got %d", fallback.SectorIndex)
	}
}
