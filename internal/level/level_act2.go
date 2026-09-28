package level

import (
	"astroleap/internal/assets"
)

// NewLevel6 builds Sector 6: Station Orbitale Olympus (Act II - Orbiting above Earth's Quarantine Grid).
func NewLevel6() *Level {
	width := 135
	lvl := &Level{
		SectorIndex:  6,
		Name:         "STATION OLYMPUS",
		Theme:        "olympus",
		Celestial:    "earth",
		GoalKind:     "lander",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
	}

	// 1. Terrain & Structural Station Modules
	// Start Hangar Platform (0..16)
	for x := 0; x <= 16; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}

	// Ceiling girder (0..18)
	for x := 0; x <= 18; x++ {
		lvl.Tiles[0][x] = TileRock
	}

	// Module 1: Stepped observation bridge (20..32)
	for x := 20; x <= 32; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	// Laser conduit hazard on lower platform
	lvl.Tiles[8][25] = TileSpike
	lvl.Tiles[8][26] = TileSpike

	// High floating platform (24..29, y=4)
	for x := 24; x <= 29; x++ {
		lvl.Tiles[4][x] = TilePlatform
	}

	// Mid Section 2: Shattered Solar Wing Hub (36..52)
	for x := 36; x <= 52; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	lvl.Tiles[9][42] = TileSpike
	lvl.Tiles[9][43] = TileSpike

	// High Communications Relay Ledge (50..58, y=3)
	for x := 50; x <= 58; x++ {
		lvl.Tiles[3][x] = TilePlatform
	}

	// Solar Array Chasm Bridge (66..82)
	for x := 66; x <= 82; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	lvl.Tiles[8][72] = TileSpike
	lvl.Tiles[8][73] = TileSpike
	lvl.Tiles[8][74] = TileSpike

	// High observation deck (74..80, y=4)
	for x := 74; x <= 80; x++ {
		lvl.Tiles[4][x] = TilePlatform
	}

	// Section 3: Elevator Transfer Hangar (92..134)
	for x := 92; x < width; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	// Elevator approach stairs (115..122)
	for x := 115; x <= 122; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
	}
	for x := 120; x <= 122; x++ {
		lvl.Tiles[7][x] = TileSurface
		lvl.Tiles[8][x] = TileRock
	}

	// Ceiling girder over airlock
	for x := 115; x < width; x++ {
		lvl.Tiles[0][x] = TileRock
		lvl.Tiles[1][x] = TileRock
	}

	// 2. Energy Crystals
	lvl.Crystals = []CrystalData{
		{TileX: 10, TileY: 7},
		{TileX: 13, TileY: 7},
		{TileX: 26, TileY: 2}, // High on platform
		{TileX: 27, TileY: 2},
		{TileX: 38, TileY: 7},
		{TileX: 46, TileY: 7},
		{TileX: 54, TileY: 1}, // High relay
		{TileX: 56, TileY: 1},
		{TileX: 68, TileY: 6},
		{TileX: 77, TileY: 2},
		{TileX: 95, TileY: 7},
		{TileX: 102, TileY: 7},
		{TileX: 108, TileY: 7},
	}

	// 3. Enemy Spawns (Automated Defense Drones & Mutated Blobs)
	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 22, TileY: 7, Kind: "blob"},
		{TileX: 28, TileY: 7, Kind: "blob"},
		{TileX: 33, TileY: 5, Kind: "hover"},
		{TileX: 40, TileY: 8, Kind: "blob"},
		{TileX: 48, TileY: 8, Kind: "blob"},
		{TileX: 62, TileY: 4, Kind: "hover"}, // Patrolling the chasm
		{TileX: 70, TileY: 7, Kind: "blob"},
		{TileX: 78, TileY: 7, Kind: "blob"},
		{TileX: 86, TileY: 5, Kind: "hover"},
		{TileX: 98, TileY: 8, Kind: "blob"},
		{TileX: 106, TileY: 8, Kind: "blob"},
		{TileX: 112, TileY: 5, Kind: "hover"},
	}

	// 4. Moving Mechanical Gantry Platforms across solar chasm
	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(53 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(65 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
		{
			StartX: float64(83 * TileSize),
			StartY: float64(7 * TileSize),
			EndX:   float64(91 * TileSize),
			EndY:   float64(7 * TileSize),
			W:      32,
			H:      8,
			Speed:  32.0,
		},
	}

	// 5. Fragile Crumbling Catwalks
	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 17, TileY: 8, W: 32, H: 8},
		{TileX: 33, TileY: 8, W: 32, H: 8},
	}

	// 6. Secret Gold Astronaut Medals
	lvl.Medals = []MedalData{
		{TileX: 27, TileY: 1}, // Very high above observation bridge
		{TileX: 58, TileY: 1}, // High above communications relay
		{TileX: 104, TileY: 2}, // High above final corridor
	}

	// 7. Power-ups & Sustenance
	lvl.HealthPickups = []HealthPickupData{
		{TileX: 47, TileY: 5},
		{TileX: 100, TileY: 6},
	}
	lvl.BoostPickups = []BoostPickupData{
		{TileX: 76, TileY: 2},
	}

	// 8. Holographic Data-Log Terminal (LOG-06)
	lvl.Terminals = []TerminalData{
		{TileX: 54, TileY: 2, LogID: "LOG-06"},
	}

	// 9. Weapon Spawns
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 24, TileY: 3, Kind: "laser"},
		{TileX: 75, TileY: 3, Kind: "nova"},
	}

	return lvl
}

// NewLevel7 builds Sector 7: L'Ascenseur Spatial Céleste (Act II - High-speed descent along the planetary tether).
func NewLevel7() *Level {
	width := 135
	lvl := &Level{
		SectorIndex:  7,
		Name:         "CELESTIAL ELEVATOR",
		Theme:        "vessel",
		Celestial:    "earth",
		GoalKind:     "lander",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
	}

	// High altitude platforms along the elevator ribbon
	for x := 0; x <= 18; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 24; x <= 50; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 56; x <= 85; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 92; x < width; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}

	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(19 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(22 * TileSize),
			EndY:   float64(4 * TileSize),
			W:      32,
			H:      8,
			Speed:  28.0,
		},
		{
			StartX: float64(51 * TileSize),
			StartY: float64(7 * TileSize),
			EndX:   float64(54 * TileSize),
			EndY:   float64(5 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
	}

	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 86, TileY: 8, W: 32, H: 8},
		{TileX: 89, TileY: 8, W: 32, H: 8},
	}

	lvl.Crystals = []CrystalData{
		{TileX: 12, TileY: 7},
		{TileX: 30, TileY: 6},
		{TileX: 42, TileY: 6},
		{TileX: 65, TileY: 7},
		{TileX: 78, TileY: 7},
		{TileX: 105, TileY: 6},
	}

	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 28, TileY: 7, Kind: "blob"},
		{TileX: 36, TileY: 5, Kind: "hover"},
		{TileX: 62, TileY: 8, Kind: "blob"},
		{TileX: 72, TileY: 5, Kind: "hover"},
		{TileX: 100, TileY: 7, Kind: "blob"},
	}

	lvl.Medals = []MedalData{
		{TileX: 38, TileY: 2},
		{TileX: 70, TileY: 2},
		{TileX: 115, TileY: 2},
	}

	lvl.HealthPickups = []HealthPickupData{
		{TileX: 48, TileY: 5},
	}

	lvl.BoostPickups = []BoostPickupData{
		{TileX: 82, TileY: 4},
	}

	lvl.Terminals = []TerminalData{
		{TileX: 68, TileY: 4, LogID: "LOG-07"},
	}

	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 32, TileY: 6, Kind: "laser"},
		{TileX: 76, TileY: 7, Kind: "nova"},
	}

	return lvl
}

// NewLevel8 builds Sector 8: Bouclier de Défense Aegis (Act II - Heavy security laser barriers).
func NewLevel8() *Level {
	width := 135
	lvl := &Level{
		SectorIndex:  8,
		Name:         "AEGIS DEFENSE GRID",
		Theme:        "vessel",
		Celestial:    "earth",
		GoalKind:     "airlock",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
	}

	for x := 0; x <= 16; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 22; x <= 45; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	lvl.Tiles[8][32] = TileSpike
	lvl.Tiles[8][33] = TileSpike

	for x := 52; x <= 80; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	lvl.Tiles[9][64] = TileSpike
	lvl.Tiles[9][65] = TileSpike

	for x := 88; x < width; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}

	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(17 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(20 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
		{
			StartX: float64(46 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(50 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
	}

	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 82, TileY: 8, W: 32, H: 8},
		{TileX: 85, TileY: 8, W: 32, H: 8},
	}

	lvl.SecurityKey = &SecurityKeyData{
		TileX: 68,
		TileY: 2,
	}

	lvl.Crystals = []CrystalData{
		{TileX: 10, TileY: 7},
		{TileX: 28, TileY: 6},
		{TileX: 60, TileY: 7},
		{TileX: 72, TileY: 7},
		{TileX: 95, TileY: 6},
	}

	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 26, TileY: 7, Kind: "blob"},
		{TileX: 38, TileY: 5, Kind: "hover"},
		{TileX: 58, TileY: 8, Kind: "blob"},
		{TileX: 70, TileY: 5, Kind: "hover"},
		{TileX: 102, TileY: 7, Kind: "blob"},
	}

	lvl.Medals = []MedalData{
		{TileX: 36, TileY: 2},
		{TileX: 75, TileY: 2},
		{TileX: 110, TileY: 2},
	}

	lvl.HealthPickups = []HealthPickupData{
		{TileX: 42, TileY: 5},
	}

	lvl.BoostPickups = []BoostPickupData{
		{TileX: 68, TileY: 3},
	}

	lvl.Terminals = []TerminalData{
		{TileX: 55, TileY: 4, LogID: "LOG-08"},
	}

	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 25, TileY: 6, Kind: "laser"},
		{TileX: 62, TileY: 7, Kind: "nova"},
	}

	return lvl
}

// NewLevel9 builds Sector 9: Stratosphère & Chute Libre (Act II - Free fall atmospheric entry).
func NewLevel9() *Level {
	width := 140
	lvl := &Level{
		SectorIndex:  9,
		Name:         "STRATOSPHERE DESCENT",
		Theme:        "crimson",
		Celestial:    "earth",
		GoalKind:     "lander",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
	}

	for x := 0; x <= 20; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 26; x <= 55; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 62; x <= 90; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 98; x < width; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}

	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(21 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(23 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
	}

	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 92, TileY: 8, W: 32, H: 8},
		{TileX: 95, TileY: 8, W: 32, H: 8},
	}

	lvl.Geysers = []GeyserData{
		{TileX: 45, TileY: 8, Height: 75.0},
		{TileX: 80, TileY: 9, Height: 80.0},
	}

	lvl.Crystals = []CrystalData{
		{TileX: 12, TileY: 7},
		{TileX: 35, TileY: 6},
		{TileX: 50, TileY: 3},
		{TileX: 72, TileY: 7},
		{TileX: 110, TileY: 6},
	}

	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 30, TileY: 7, Kind: "blob"},
		{TileX: 48, TileY: 5, Kind: "hover"},
		{TileX: 68, TileY: 8, Kind: "blob"},
		{TileX: 85, TileY: 5, Kind: "hover"},
		{TileX: 108, TileY: 7, Kind: "blob"},
	}

	lvl.Medals = []MedalData{
		{TileX: 50, TileY: 1},
		{TileX: 85, TileY: 1},
		{TileX: 120, TileY: 2},
	}

	lvl.HealthPickups = []HealthPickupData{
		{TileX: 52, TileY: 5},
	}

	lvl.BoostPickups = []BoostPickupData{
		{TileX: 76, TileY: 3},
	}

	lvl.Terminals = []TerminalData{
		{TileX: 70, TileY: 4, LogID: "LOG-09"},
	}

	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 28, TileY: 6, Kind: "nova"},
		{TileX: 75, TileY: 7, Kind: "laser"},
	}

	return lvl
}

// NewLevel10 builds Sector 10: Le Cœur de l'Arche Terrestre (Act II Final Boss: A.E.G.I.S. Core).
func NewLevel10() *Level {
	width := 145
	lvl := &Level{
		SectorIndex:  10,
		Name:         "TERRAN ARK CORE",
		Theme:        "reactor",
		Celestial:    "core",
		GoalKind:     "warp_console",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
		HasBoss:      true,
		BossX:        float64(125 * TileSize),
		BossY:        float64(4 * TileSize),
	}

	for x := 0; x <= 22; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	for x := 28; x <= 60; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	lvl.Tiles[8][42] = TileSpike
	lvl.Tiles[8][43] = TileSpike

	for x := 68; x <= 100; x++ {
		lvl.Tiles[9][x] = TileSurface
		for y := 10; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}
	lvl.Tiles[9][82] = TileSpike
	lvl.Tiles[9][83] = TileSpike

	for x := 108; x < width; x++ {
		lvl.Tiles[8][x] = TileSurface
		for y := 9; y < LevelHeight; y++ {
			lvl.Tiles[y][x] = TileRock
		}
	}

	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(23 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(25 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
		{
			StartX: float64(61 * TileSize),
			StartY: float64(7 * TileSize),
			EndX:   float64(65 * TileSize),
			EndY:   float64(7 * TileSize),
			W:      32,
			H:      8,
			Speed:  32.0,
		},
	}

	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 102, TileY: 8, W: 32, H: 8},
		{TileX: 105, TileY: 8, W: 32, H: 8},
	}

	lvl.Crystals = []CrystalData{
		{TileX: 14, TileY: 7},
		{TileX: 38, TileY: 6},
		{TileX: 52, TileY: 5},
		{TileX: 78, TileY: 7},
		{TileX: 95, TileY: 7},
		{TileX: 114, TileY: 6},
	}

	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 32, TileY: 7, Kind: "blob"},
		{TileX: 48, TileY: 5, Kind: "hover"},
		{TileX: 74, TileY: 8, Kind: "blob"},
		{TileX: 90, TileY: 5, Kind: "hover"},
		{TileX: 112, TileY: 7, Kind: "blob"},
	}

	lvl.Medals = []MedalData{
		{TileX: 50, TileY: 1},
		{TileX: 92, TileY: 1},
		{TileX: 135, TileY: 2},
	}

	lvl.HealthPickups = []HealthPickupData{
		{TileX: 54, TileY: 5},
		{TileX: 104, TileY: 5},
	}

	lvl.BoostPickups = []BoostPickupData{
		{TileX: 84, TileY: 3},
	}

	lvl.Terminals = []TerminalData{
		{TileX: 75, TileY: 4, LogID: "LOG-10"},
	}

	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 18, TileY: 7, Kind: "laser"},
		{TileX: 72, TileY: 7, Kind: "nova"},
	}

	return lvl
}
