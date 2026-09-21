package level

import (
	"image/color"
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TileSize    = 16
	LevelHeight = 12
	MaxLevels   = 5
)

const (
	TileEmpty    = 0
	TileSurface  = 1
	TileRock     = 2
	TilePlatform = 3
	TileSpike    = 4
)

// SecurityKeyData defines placement of the Sector 4 security keycard.
type SecurityKeyData struct {
	TileX int
	TileY int
}

// RepairCoreData defines placement of Sector 5 warp reactor repair cores.
type RepairCoreData struct {
	Index int
	Name  string
	TileX int
	TileY int
}

// CrystalData holds placement of crystals.
type CrystalData struct {
	TileX int
	TileY int
}

// EnemySpawn holds enemy spawn coordinates and type.
type EnemySpawn struct {
	TileX int
	TileY int
	Kind  string // "blob" or "hover"
}

// WeaponSpawn holds weapon capsule locations.
type WeaponSpawn struct {
	TileX int
	TileY int
	Kind  string // "laser" or "nova"
}

// MovingPlatformData defines moving platform waypoints.
type MovingPlatformData struct {
	StartX float64
	StartY float64
	EndX   float64
	EndY   float64
	W      float64
	H      float64
	Speed  float64
}

// GeyserData defines upward steam vent locations.
type GeyserData struct {
	TileX  int
	TileY  int
	Height float64
}

// MedalData defines hidden collectible gold medal placement.
type MedalData struct {
	TileX int
	TileY int
}

// CrumbleData defines fragile crumbling platform locations.
type CrumbleData struct {
	TileX int
	TileY int
	W     float64
	H     float64
}

// HealthPickupData defines oxygen canister locations.
type HealthPickupData struct {
	TileX int
	TileY int
}

// BoostPickupData defines temporary ion thruster boost core locations.
type BoostPickupData struct {
	TileX int
	TileY int
}

// Level represents a playable stage with tile grid, hazards, and background layers.
type Level struct {
	SectorIndex        int
	Name               string
	Theme              string // "moon", "crimson", "ice", "vessel", "reactor"
	Celestial          string // "earth", "mars", "jupiter", "mothership", "core"
	GoalKind           string // "lander", "airlock", "warp_console"
	Width              int
	Height             int
	Tiles              [][]int
	Crystals           []CrystalData
	EnemySpawns        []EnemySpawn
	WeaponSpawns       []WeaponSpawn
	MovingPlatforms    []MovingPlatformData
	Geysers            []GeyserData
	Medals             []MedalData
	CrumblingPlatforms []CrumbleData
	HealthPickups      []HealthPickupData
	BoostPickups       []BoostPickupData
	SecurityKey        *SecurityKeyData
	RepairCores        []RepairCoreData
	LanderX            float64
	LanderY            float64
	HasBoss            bool
	BossX              float64
	BossY              float64
	StarfieldImg       *ebiten.Image
}

func makeGrid(w, h int) [][]int {
	g := make([][]int, h)
	for y := 0; y < h; y++ {
		g[y] = make([]int, w)
	}
	return g
}

// LoadLevel loads the level for the given sector (1 to 5).
func LoadLevel(sector int) *Level {
	switch sector {
	case 2:
		return NewLevel2()
	case 3:
		return NewLevel3()
	case 4:
		return NewLevel4()
	case 5:
		return NewLevel5()
	default:
		return NewLevel1()
	}
}

// NewLevel1 builds Sector 1: Lunar Outpost (Earth celestial, gray regolith).
func NewLevel1() *Level {
	width := 120
	lvl := &Level{
		SectorIndex:  1,
		Name:         "LUNAR OUTPOST",
		Theme:        "moon",
		Celestial:    "earth",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(6 * TileSize),
	}

	// 1. Terrain
	for x := 0; x <= 22; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// First Chasm at 23..26 (Jump required! Floaty jump tutorial)
	for x := 27; x <= 34; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[6][29] = TilePlatform
	lvl.Tiles[6][30] = TilePlatform
	lvl.Tiles[6][31] = TilePlatform

	// Spike pit at 35..38
	for x := 35; x <= 38; x++ {
		lvl.Tiles[11][x] = TileSpike
	}

	// Island 2 at 39..48
	for x := 39; x <= 48; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Floating jetpack staircase (49..59)
	lvl.Tiles[6][50] = TilePlatform
	lvl.Tiles[6][51] = TilePlatform
	lvl.Tiles[4][54] = TilePlatform
	lvl.Tiles[4][55] = TilePlatform
	lvl.Tiles[6][58] = TilePlatform
	lvl.Tiles[6][59] = TilePlatform
	for x := 49; x <= 59; x++ {
		lvl.Tiles[11][x] = TileSpike
	}

	// Plateau 3 (60..74)
	for x := 60; x <= 74; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[5][64] = TilePlatform
	lvl.Tiles[5][65] = TilePlatform
	lvl.Tiles[5][69] = TilePlatform
	lvl.Tiles[5][70] = TilePlatform

	// Wide Crater Basin (75..84)
	for x := 75; x <= 84; x++ {
		lvl.Tiles[10][x] = TileSurface
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[7][78] = TilePlatform
	lvl.Tiles[7][81] = TilePlatform

	// Deep Chasm (85..89)
	for x := 85; x <= 89; x++ {
		lvl.Tiles[11][x] = TileSpike
	}

	// Final Landing Zone (90..width-1)
	for x := 90; x < width; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// 2. Crystals
	lvl.Crystals = []CrystalData{
		{TileX: 8, TileY: 7}, {TileX: 12, TileY: 7}, {TileX: 18, TileY: 6},
		{TileX: 25, TileY: 5}, {TileX: 30, TileY: 4}, {TileX: 42, TileY: 6},
		{TileX: 45, TileY: 6}, {TileX: 51, TileY: 4}, {TileX: 55, TileY: 2},
		{TileX: 59, TileY: 4}, {TileX: 64, TileY: 3}, {TileX: 70, TileY: 3},
		{TileX: 78, TileY: 5}, {TileX: 81, TileY: 5}, {TileX: 87, TileY: 4},
		{TileX: 96, TileY: 6}, {TileX: 100, TileY: 6}, {TileX: 105, TileY: 6},
	}

	// 3. Enemies
	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 15, TileY: 8, Kind: "blob"},
		{TileX: 31, TileY: 8, Kind: "blob"},
		{TileX: 44, TileY: 7, Kind: "blob"},
		{TileX: 53, TileY: 3, Kind: "hover"},
		{TileX: 66, TileY: 8, Kind: "blob"},
		{TileX: 72, TileY: 8, Kind: "blob"},
		{TileX: 79, TileY: 9, Kind: "blob"},
		{TileX: 83, TileY: 5, Kind: "hover"},
		{TileX: 98, TileY: 7, Kind: "blob"},
		{TileX: 104, TileY: 7, Kind: "hover"},
	}

	// 4. Weapon Capsules
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 20, TileY: 8, Kind: "laser"}, // Laser Blaster unlock
	}

	// 5. Secret Gold Astronaut Medals
	lvl.Medals = []MedalData{
		{TileX: 25, TileY: 3},
		{TileX: 55, TileY: 1},
		{TileX: 87, TileY: 2},
	}

	// 6. Fragile Crumbling Platforms
	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 36, TileY: 8, W: 32, H: 8},
		{TileX: 74, TileY: 8, W: 32, H: 8},
	}

	// 7. Power-ups & Sustenance
	lvl.HealthPickups = []HealthPickupData{
		{TileX: 45, TileY: 5},
	}
	lvl.BoostPickups = []BoostPickupData{
		{TileX: 64, TileY: 2},
	}

	return lvl
}

// NewLevel2 builds Sector 2: Phobos Ridge (Mars celestial, crimson regolith, Nova Cannon unlock).
func NewLevel2() *Level {
	width := 130
	lvl := &Level{
		SectorIndex:  2,
		Name:         "PHOBOS RIDGE",
		Theme:        "crimson",
		Celestial:    "mars",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
	}

	// 1. Terrain: Steeper red cliffs and longer chasm jumps
	// Starting ridge (0..18)
	for x := 0; x <= 18; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Floating platforms across lava/spike trench (19..28)
	for x := 19; x <= 28; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[7][21] = TilePlatform
	lvl.Tiles[7][22] = TilePlatform
	lvl.Tiles[5][25] = TilePlatform
	lvl.Tiles[5][26] = TilePlatform

	// High red plateau (29..45)
	for x := 29; x <= 45; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Secret high sky sanctuary holding Nova Cannon (46..58)
	for x := 46; x <= 58; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[7][48] = TilePlatform
	lvl.Tiles[7][49] = TilePlatform
	// High secret tech platform (jetpack ascent required)
	lvl.Tiles[3][54] = TilePlatform
	lvl.Tiles[3][55] = TilePlatform
	lvl.Tiles[3][56] = TilePlatform
	lvl.Tiles[6][57] = TilePlatform
	lvl.Tiles[6][58] = TilePlatform

	// Canyon basin (59..78)
	for x := 59; x <= 78; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[6][64] = TilePlatform
	lvl.Tiles[6][65] = TilePlatform
	lvl.Tiles[4][70] = TilePlatform
	lvl.Tiles[4][71] = TilePlatform

	// Double spike trench (79..92)
	for x := 79; x <= 92; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	// Mid-trench resting ledge
	lvl.Tiles[6][87] = TilePlatform
	lvl.Tiles[6][88] = TilePlatform

	// High cliff summit (93..width-1)
	for x := 93; x < width; x++ {
		lvl.Tiles[7][x] = TileSurface
		lvl.Tiles[8][x] = TileRock
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// 2. Crystals
	lvl.Crystals = []CrystalData{
		{TileX: 6, TileY: 7}, {TileX: 12, TileY: 7}, {TileX: 22, TileY: 5},
		{TileX: 26, TileY: 3}, {TileX: 33, TileY: 6}, {TileX: 38, TileY: 6},
		{TileX: 43, TileY: 6}, {TileX: 49, TileY: 5}, {TileX: 53, TileY: 1}, // High secrets
		{TileX: 55, TileY: 1}, {TileX: 57, TileY: 1}, {TileX: 65, TileY: 4},
		{TileX: 71, TileY: 2}, {TileX: 76, TileY: 7}, {TileX: 83, TileY: 5},
		{TileX: 87, TileY: 4}, {TileX: 91, TileY: 3}, {TileX: 98, TileY: 5},
		{TileX: 106, TileY: 5}, {TileX: 114, TileY: 5}, {TileX: 120, TileY: 5},
	}

	// 3. Enemies
	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 14, TileY: 8, Kind: "blob"},
		{TileX: 24, TileY: 4, Kind: "hover"}, // Hover over spike trench
		{TileX: 34, TileY: 7, Kind: "blob"},
		{TileX: 41, TileY: 7, Kind: "blob"},
		{TileX: 50, TileY: 5, Kind: "hover"},
		{TileX: 63, TileY: 8, Kind: "blob"},
		{TileX: 69, TileY: 8, Kind: "blob"},
		{TileX: 73, TileY: 3, Kind: "hover"},
		{TileX: 84, TileY: 3, Kind: "hover"},
		{TileX: 89, TileY: 4, Kind: "hover"},
		{TileX: 102, TileY: 6, Kind: "blob"},
		{TileX: 112, TileY: 6, Kind: "blob"},
	}

	// 4. Weapon Capsules
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 55, TileY: 2, Kind: "nova"}, // Nova Cannon on high secret ledge
	}

	// 5. Moving Platforms - Clean dedicated airspace across trench
	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(80 * TileSize),
			StartY: float64(7 * TileSize),
			EndX:   float64(85 * TileSize),
			EndY:   float64(7 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
	}

	// 6. Thermal Cryo Geysers
	lvl.Geysers = []GeyserData{
		{TileX: 47, TileY: 10, Height: 75.0}, // Launch to secret Nova Cannon ledge
	}

	// 7. Secret Gold Astronaut Medals
	lvl.Medals = []MedalData{
		{TileX: 23, TileY: 3},
		{TileX: 56, TileY: 1},
		{TileX: 88, TileY: 2},
	}

	// 8. Fragile Crumbling Platforms
	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 62, TileY: 7, W: 32, H: 8},
		{TileX: 84, TileY: 6, W: 32, H: 8},
	}

	// 9. Power-ups & Sustenance
	lvl.HealthPickups = []HealthPickupData{
		{TileX: 52, TileY: 2},
	}
	lvl.BoostPickups = []BoostPickupData{
		{TileX: 78, TileY: 5},
	}

	return lvl
}

// NewLevel3 builds Sector 3: Europa Glacial Core (Ringed Jupiter, crystalline ice caverns).
func NewLevel3() *Level {
	width := 140
	lvl := &Level{
		SectorIndex:  3,
		Name:         "EUROPA ICE CORE",
		Theme:        "ice",
		Celestial:    "jupiter",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
		HasBoss:      true,
		BossX:        float64(120 * TileSize),
		BossY:        float64(4 * TileSize),
	}

	// 1. Terrain: Crystalline ice bridges, narrow vertical columns, floating platforms
	// Starting glacier (0..16)
	for x := 0; x <= 16; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Ice crevasses with narrow stepping stones (17..32)
	for x := 17; x <= 32; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[8][19] = TilePlatform
	lvl.Tiles[6][23] = TilePlatform
	lvl.Tiles[6][24] = TilePlatform
	lvl.Tiles[7][28] = TilePlatform
	lvl.Tiles[7][29] = TilePlatform

	// Glacial Pillar 1 (33..44)
	for x := 33; x <= 44; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[5][37] = TilePlatform
	lvl.Tiles[5][38] = TilePlatform

	// Vertical Gauntlet with alternating floating platforms (45..68)
	for x := 45; x <= 68; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[7][47] = TilePlatform
	lvl.Tiles[5][51] = TilePlatform
	lvl.Tiles[5][52] = TilePlatform
	lvl.Tiles[3][56] = TilePlatform
	lvl.Tiles[3][57] = TilePlatform
	lvl.Tiles[5][61] = TilePlatform
	lvl.Tiles[5][62] = TilePlatform
	lvl.Tiles[7][66] = TilePlatform

	// Glacial Island 2 (69..88)
	for x := 69; x <= 88; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[6][74] = TilePlatform
	lvl.Tiles[6][75] = TilePlatform
	lvl.Tiles[4][81] = TilePlatform
	lvl.Tiles[4][82] = TilePlatform

	// Final Frozen Abyss (89..108)
	for x := 89; x <= 108; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	// Stepping platform before moving platform
	lvl.Tiles[7][91] = TilePlatform
	lvl.Tiles[7][92] = TilePlatform
	// Receiving platform after crumbling ledge
	lvl.Tiles[6][106] = TilePlatform
	lvl.Tiles[6][107] = TilePlatform

	// Final Extraction Citadel (109..width-1)
	for x := 109; x < width; x++ {
		lvl.Tiles[7][x] = TileSurface
		lvl.Tiles[8][x] = TileRock
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// 2. Crystals
	lvl.Crystals = []CrystalData{
		{TileX: 5, TileY: 7}, {TileX: 11, TileY: 7}, {TileX: 20, TileY: 6},
		{TileX: 24, TileY: 4}, {TileX: 29, TileY: 5}, {TileX: 35, TileY: 6},
		{TileX: 38, TileY: 3}, {TileX: 42, TileY: 6}, {TileX: 48, TileY: 5},
		{TileX: 52, TileY: 3}, {TileX: 56, TileY: 1}, {TileX: 57, TileY: 1},
		{TileX: 62, TileY: 3}, {TileX: 66, TileY: 5}, {TileX: 72, TileY: 7},
		{TileX: 75, TileY: 4}, {TileX: 82, TileY: 2}, {TileX: 86, TileY: 7},
		{TileX: 92, TileY: 5}, {TileX: 96, TileY: 4}, {TileX: 101, TileY: 4},
		{TileX: 114, TileY: 5}, {TileX: 120, TileY: 5}, {TileX: 126, TileY: 5},
	}

	// 3. Enemies
	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 13, TileY: 8, Kind: "blob"},
		{TileX: 21, TileY: 5, Kind: "hover"},
		{TileX: 27, TileY: 5, Kind: "hover"},
		{TileX: 36, TileY: 7, Kind: "blob"},
		{TileX: 41, TileY: 7, Kind: "blob"},
		{TileX: 49, TileY: 5, Kind: "hover"},
		{TileX: 54, TileY: 3, Kind: "hover"},
		{TileX: 60, TileY: 3, Kind: "hover"},
		{TileX: 73, TileY: 8, Kind: "blob"},
		{TileX: 83, TileY: 8, Kind: "blob"},
		{TileX: 94, TileY: 3, Kind: "hover"},
		{TileX: 99, TileY: 3, Kind: "hover"},
		{TileX: 116, TileY: 6, Kind: "blob"},
		{TileX: 124, TileY: 6, Kind: "hover"},
	}

	// 4. Weapon Capsules (Laser and Nova caches for players wanting quick rearm)
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 24, TileY: 5, Kind: "laser"},
		{TileX: 82, TileY: 3, Kind: "nova"},
	}

	// 5. Moving Platforms across the frozen abyss - clean open lane
	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(93 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(100 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  35.0,
		},
	}

	// 6. Thermal Cryo Geysers
	lvl.Geysers = []GeyserData{
		{TileX: 46, TileY: 10, Height: 80.0},
		{TileX: 89, TileY: 10, Height: 75.0},
	}

	// 7. Secret Gold Astronaut Medals
	lvl.Medals = []MedalData{
		{TileX: 20, TileY: 4},
		{TileX: 57, TileY: 1},
		{TileX: 102, TileY: 1},
	}

	// 8. Fragile Crumbling Platforms - placed after moving platform path
	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 20, TileY: 7, W: 32, H: 8},
		{TileX: 102, TileY: 6, W: 32, H: 8},
	}

	// 9. Power-ups & Sustenance
	lvl.HealthPickups = []HealthPickupData{
		{TileX: 74, TileY: 5},
		{TileX: 110, TileY: 6},
	}
	lvl.BoostPickups = []BoostPickupData{
		{TileX: 48, TileY: 4},
	}

	return lvl
}

// NewLevel4 builds Sector 4: Mothership Corridor (Spaceship interior, Security Key required to open Airlock).
func NewLevel4() *Level {
	width := 130
	lvl := &Level{
		SectorIndex:  4,
		Name:         "MOTHERSHIP CORRIDOR",
		Theme:        "vessel",
		Celestial:    "mothership",
		GoalKind:     "airlock",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
		SecurityKey:  &SecurityKeyData{TileX: 68, TileY: 2},
	}

	// 1. Terrain: Metallic ship deck, laser conduit trenches, elevator shafts, security gantry
	// Starting hangar deck (0..18)
	for x := 0; x <= 18; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// High-voltage laser conduit trench (19..28)
	for x := 19; x <= 28; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[7][21] = TilePlatform
	lvl.Tiles[7][22] = TilePlatform
	lvl.Tiles[5][25] = TilePlatform
	lvl.Tiles[5][26] = TilePlatform

	// Security Checkpoint Corridor (29..44)
	for x := 29; x <= 44; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Vertical elevator shaft (45..54) across deep hazard drop
	for x := 45; x <= 54; x++ {
		lvl.Tiles[11][x] = TileSpike
	}

	// High Security Bridge & Keycard Vault (55..75)
	// Lower service corridor
	for x := 55; x <= 75; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	// Upper security gantry holding the Access Keycard!
	for x := 62; x <= 73; x++ {
		lvl.Tiles[4][x] = TilePlatform
	}

	// Plasma Surge Chamber (76..92) with horizontal cargo platform
	for x := 76; x <= 92; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[6][89] = TilePlatform
	lvl.Tiles[6][90] = TilePlatform

	// Airlock Decompression Chamber & Exit Gate (93..width-1)
	for x := 93; x < width; x++ {
		lvl.Tiles[7][x] = TileSurface
		lvl.Tiles[8][x] = TileRock
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// 2. Crystals
	lvl.Crystals = []CrystalData{
		{TileX: 6, TileY: 7}, {TileX: 12, TileY: 7}, {TileX: 22, TileY: 5},
		{TileX: 26, TileY: 3}, {TileX: 33, TileY: 6}, {TileX: 39, TileY: 6},
		{TileX: 47, TileY: 2}, {TileX: 63, TileY: 2}, {TileX: 66, TileY: 2},
		{TileX: 70, TileY: 2}, {TileX: 73, TileY: 2}, {TileX: 60, TileY: 7},
		{TileX: 68, TileY: 7}, {TileX: 80, TileY: 5}, {TileX: 85, TileY: 5},
		{TileX: 90, TileY: 4}, {TileX: 96, TileY: 5}, {TileX: 102, TileY: 5},
		{TileX: 108, TileY: 5}, {TileX: 114, TileY: 5}, {TileX: 120, TileY: 5},
	}

	// 3. Enemies
	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 14, TileY: 8, Kind: "blob"},
		{TileX: 23, TileY: 4, Kind: "hover"},
		{TileX: 35, TileY: 7, Kind: "blob"},
		{TileX: 42, TileY: 7, Kind: "blob"},
		{TileX: 50, TileY: 3, Kind: "hover"},
		{TileX: 65, TileY: 3, Kind: "blob"},
		{TileX: 71, TileY: 3, Kind: "hover"},
		{TileX: 84, TileY: 4, Kind: "hover"},
		{TileX: 98, TileY: 6, Kind: "blob"},
		{TileX: 106, TileY: 6, Kind: "hover"},
		{TileX: 116, TileY: 6, Kind: "blob"},
	}

	// 4. Weapon Capsules
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 36, TileY: 7, Kind: "laser"},
	}

	// 5. Moving Platforms - Clean dedicated flight corridors
	lvl.MovingPlatforms = []MovingPlatformData{
		// Vertical elevator lift to security balcony (47)
		{
			StartX: float64(47 * TileSize),
			StartY: float64(8 * TileSize),
			EndX:   float64(47 * TileSize),
			EndY:   float64(4 * TileSize),
			W:      32,
			H:      8,
			Speed:  28.0,
		},
		// Horizontal cargo transport over surge trench (78..86)
		{
			StartX: float64(78 * TileSize),
			StartY: float64(7 * TileSize),
			EndX:   float64(86 * TileSize),
			EndY:   float64(7 * TileSize),
			W:      32,
			H:      8,
			Speed:  32.0,
		},
	}

	// 6. Thermal Cryo / Plasma Geysers
	lvl.Geysers = []GeyserData{
		{TileX: 52, TileY: 10, Height: 75.0},
	}

	// 7. Secret Gold Astronaut Medals
	lvl.Medals = []MedalData{
		{TileX: 25, TileY: 3},
		{TileX: 71, TileY: 1}, // High above key vault
		{TileX: 90, TileY: 2},
	}

	// 8. Fragile Crumbling Platforms
	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 38, TileY: 7, W: 32, H: 8},
		{TileX: 92, TileY: 6, W: 32, H: 8},
	}

	// 9. Power-ups & Sustenance
	lvl.HealthPickups = []HealthPickupData{
		{TileX: 40, TileY: 6},
	}
	lvl.BoostPickups = []BoostPickupData{
		{TileX: 75, TileY: 2},
	}

	return lvl
}

// NewLevel5 builds Sector 5: Reactor Bay (Ship interior repair room, 3 Repair Cores required to fix ship).
func NewLevel5() *Level {
	width := 135
	lvl := &Level{
		SectorIndex:  5,
		Name:         "REACTOR BAY",
		Theme:        "reactor",
		Celestial:    "core",
		GoalKind:     "warp_console",
		Width:        width,
		Height:       LevelHeight,
		Tiles:        makeGrid(width, LevelHeight),
		StarfieldImg: assets.GenerateStarfield(320, 180),
		LanderX:      float64((width - 9) * TileSize),
		LanderY:      float64(5 * TileSize),
		RepairCores: []RepairCoreData{
			{Index: 1, Name: "COOLANT REGULATOR", TileX: 28, TileY: 3},
			{Index: 2, Name: "PLASMA STABILIZER", TileX: 68, TileY: 2},
			{Index: 3, Name: "HYPER-FLUX ROD", TileX: 104, TileY: 4},
		},
	}

	// 1. Terrain: High-temp reactor deck, coolant trenches, magnetic conduits
	// Starting maintenance bay (0..16)
	for x := 0; x <= 16; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Coolant flooded trench (17..32) holding Core 1 on high vent platform
	for x := 17; x <= 32; x++ {
		lvl.Tiles[11][x] = TileSpike
	}
	lvl.Tiles[7][21] = TilePlatform
	lvl.Tiles[5][27] = TilePlatform
	lvl.Tiles[5][28] = TilePlatform
	lvl.Tiles[5][29] = TilePlatform

	// Primary Containment Wall & Generator Bay (33..52)
	for x := 33; x <= 52; x++ {
		lvl.Tiles[8][x] = TileSurface
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// Vertical Central Piston Chasm (53..64)
	for x := 53; x <= 64; x++ {
		lvl.Tiles[11][x] = TileSpike
	}

	// High Reactor Containment Arch (65..80) holding Core 2
	for x := 65; x <= 80; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	for x := 66; x <= 74; x++ {
		lvl.Tiles[4][x] = TilePlatform
	}

	// Secondary Coolant Abyss (81..95)
	for x := 81; x <= 95; x++ {
		lvl.Tiles[11][x] = TileSpike
	}

	// Magnetic Flux Conduit Corridor (96..110) holding Core 3
	for x := 96; x <= 110; x++ {
		lvl.Tiles[9][x] = TileSurface
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}
	lvl.Tiles[6][103] = TilePlatform
	lvl.Tiles[6][104] = TilePlatform
	lvl.Tiles[6][105] = TilePlatform

	// Final Warp Drive Engine Room & Bridge Console (111..width-1)
	for x := 111; x < width; x++ {
		lvl.Tiles[7][x] = TileSurface
		lvl.Tiles[8][x] = TileRock
		lvl.Tiles[9][x] = TileRock
		lvl.Tiles[10][x] = TileRock
		lvl.Tiles[11][x] = TileRock
	}

	// 2. Crystals
	lvl.Crystals = []CrystalData{
		{TileX: 6, TileY: 7}, {TileX: 12, TileY: 7}, {TileX: 21, TileY: 5},
		{TileX: 27, TileY: 4}, {TileX: 35, TileY: 6}, {TileX: 40, TileY: 6},
		{TileX: 46, TileY: 6}, {TileX: 58, TileY: 2}, {TileX: 67, TileY: 2},
		{TileX: 70, TileY: 2}, {TileX: 73, TileY: 2}, {TileX: 78, TileY: 7},
		{TileX: 84, TileY: 4}, {TileX: 90, TileY: 4}, {TileX: 98, TileY: 7},
		{TileX: 103, TileY: 4}, {TileX: 106, TileY: 4}, {TileX: 115, TileY: 5},
		{TileX: 120, TileY: 5}, {TileX: 124, TileY: 5}, {TileX: 128, TileY: 5},
	}

	// 3. Enemies
	lvl.EnemySpawns = []EnemySpawn{
		{TileX: 14, TileY: 8, Kind: "blob"},
		{TileX: 23, TileY: 4, Kind: "hover"},
		{TileX: 36, TileY: 7, Kind: "blob"},
		{TileX: 43, TileY: 7, Kind: "blob"},
		{TileX: 56, TileY: 3, Kind: "hover"},
		{TileX: 69, TileY: 3, Kind: "blob"},
		{TileX: 75, TileY: 7, Kind: "blob"},
		{TileX: 87, TileY: 3, Kind: "hover"},
		{TileX: 100, TileY: 7, Kind: "blob"},
		{TileX: 107, TileY: 4, Kind: "hover"},
		{TileX: 118, TileY: 6, Kind: "blob"},
		{TileX: 126, TileY: 6, Kind: "hover"},
	}

	// 4. Weapon Capsules
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 42, TileY: 6, Kind: "nova"},
	}

	// 5. Moving Platforms - Clean dedicated airspace
	lvl.MovingPlatforms = []MovingPlatformData{
		// Vertical piston lift to Core 2 on containment arch (58)
		{
			StartX: float64(58 * TileSize),
			StartY: float64(8 * TileSize),
			EndX:   float64(58 * TileSize),
			EndY:   float64(4 * TileSize),
			W:      32,
			H:      8,
			Speed:  30.0,
		},
		// Horizontal containment gantry over secondary coolant abyss (83..91)
		{
			StartX: float64(83 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(91 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  34.0,
		},
	}

	// 6. Thermal Cryo / Plasma Geysers
	lvl.Geysers = []GeyserData{
		{TileX: 24, TileY: 10, Height: 85.0}, // Propel to Core 1
		{TileX: 81, TileY: 10, Height: 80.0},
	}

	// 7. Secret Gold Astronaut Medals
	lvl.Medals = []MedalData{
		{TileX: 28, TileY: 1}, // High above Core 1
		{TileX: 74, TileY: 1}, // High above Core 2
		{TileX: 108, TileY: 2},
	}

	// 8. Fragile Crumbling Platforms
	lvl.CrumblingPlatforms = []CrumbleData{
		{TileX: 48, TileY: 7, W: 32, H: 8},
		{TileX: 94, TileY: 6, W: 32, H: 8},
	}

	// 9. Power-ups & Sustenance
	lvl.HealthPickups = []HealthPickupData{
		{TileX: 45, TileY: 6},
		{TileX: 98, TileY: 7},
	}
	lvl.BoostPickups = []BoostPickupData{
		{TileX: 72, TileY: 2},
	}

	return lvl
}

// IsSolid returns true if the tile at (tileX, tileY) blocks movement.
func (lvl *Level) IsSolid(tileX, tileY int) bool {
	if tileX < 0 || tileX >= lvl.Width || tileY < 0 || tileY >= lvl.Height {
		return false
	}
	t := lvl.Tiles[tileY][tileX]
	return t == TileSurface || t == TileRock || t == TilePlatform
}

// IsHazard returns true if the tile causes damage (e.g. spikes).
func (lvl *Level) IsHazard(tileX, tileY int) bool {
	if tileX < 0 || tileX >= lvl.Width || tileY < 0 || tileY >= lvl.Height {
		return false
	}
	return lvl.Tiles[tileY][tileX] == TileSpike
}

// DrawBackground renders the parallax starfield, distant celestial body, and thematic ridges.
func (lvl *Level) DrawBackground(screen *ebiten.Image, camX float64) {
	atlas := assets.Get()

	// 1. Starfield base (wraps horizontally with slow 0.08x parallax)
	sfOffset := math.Mod(camX*0.08, 320.0)
	if sfOffset < 0 {
		sfOffset += 320.0
	}
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(-sfOffset, 0)
	screen.DrawImage(lvl.StarfieldImg, op1)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(320.0-sfOffset, 0)
	screen.DrawImage(lvl.StarfieldImg, op2)

	// 2. Distant Celestial body in the sky (slow 0.03x parallax)
	celestialOp := &ebiten.DrawImageOptions{}
	celestialX := 220.0 - camX*0.03
	celestialOp.GeoM.Translate(celestialX, 16.0)

	var celestialImg *ebiten.Image
	switch lvl.Celestial {
	case "mars":
		celestialImg = atlas.Mars
	case "jupiter":
		celestialImg = atlas.Jupiter
	case "mothership", "core":
		celestialImg = atlas.Earth
	default: // earth
		celestialImg = atlas.Earth
	}
	if celestialImg != nil {
		screen.DrawImage(celestialImg, celestialOp)
	}

	// 3. Thematic mountain ridge silhouettes (0.25x parallax)
	ridgeOffset := math.Mod(camX*0.25, 160.0)
	ridgeColor := color.RGBA{32, 36, 52, 255} // Moon gray
	if lvl.Theme == "crimson" {
		ridgeColor = color.RGBA{58, 24, 28, 255} // Martian crimson
	} else if lvl.Theme == "ice" {
		ridgeColor = color.RGBA{18, 42, 65, 255} // Europa glacial navy
	} else if lvl.Theme == "vessel" {
		ridgeColor = color.RGBA{18, 24, 38, 255} // Mothership deep hull
	} else if lvl.Theme == "reactor" {
		ridgeColor = color.RGBA{34, 20, 16, 255} // Reactor thermal bronze
	}

	for i := -1; i < 4; i++ {
		bx := float32(float64(i*160) - ridgeOffset)
		vector.DrawFilledRect(screen, bx, 130, 160, 50, ridgeColor, false)
	}

	// For spaceship interior themes: draw structural pillars and neon conduits
	if lvl.Theme == "vessel" || lvl.Theme == "reactor" {
		pillarOffset := math.Mod(camX*0.15, 80.0)
		for i := -1; i < 6; i++ {
			px := float32(float64(i*80) - pillarOffset)
			colPillar := color.RGBA{18, 22, 32, 210}
			colNeon := color.RGBA{60, 200, 255, 140}
			if lvl.Theme == "reactor" {
				colPillar = color.RGBA{28, 18, 14, 210}
				colNeon = color.RGBA{255, 140, 30, 140}
			}
			vector.DrawFilledRect(screen, px, 0, 14, 180, colPillar, false)
			vector.DrawFilledRect(screen, px+6, 0, 2, 180, colNeon, false)
		}
	}
}

// DrawTiles renders visible tiles within camera frustum using the sector's theme.
func (lvl *Level) DrawTiles(screen *ebiten.Image, camX float64) {
	atlas := assets.Get()

	startCol := int(camX / TileSize)
	if startCol < 0 {
		startCol = 0
	}
	endCol := startCol + (320 / TileSize) + 2
	if endCol > lvl.Width {
		endCol = lvl.Width
	}

	for y := 0; y < lvl.Height; y++ {
		for x := startCol; x < endCol; x++ {
			t := lvl.Tiles[y][x]
			if t == TileEmpty {
				continue
			}

			screenX := float64(x*TileSize) - camX
			screenY := float64(y * TileSize)

			var tileImg *ebiten.Image
			switch t {
			case TileSurface:
				if lvl.Theme == "crimson" {
					tileImg = atlas.TileSurfaceCrimson
				} else if lvl.Theme == "ice" {
					tileImg = atlas.TileSurfaceIce
				} else if lvl.Theme == "vessel" {
					tileImg = atlas.TileSurfaceVessel
				} else if lvl.Theme == "reactor" {
					tileImg = atlas.TileSurfaceReactor
				} else {
					tileImg = atlas.TileSurface
				}
			case TileRock:
				if lvl.Theme == "crimson" {
					tileImg = atlas.TileRockCrimson
				} else if lvl.Theme == "ice" {
					tileImg = atlas.TileRockIce
				} else if lvl.Theme == "vessel" {
					tileImg = atlas.TileRockVessel
				} else if lvl.Theme == "reactor" {
					tileImg = atlas.TileRockReactor
				} else {
					tileImg = atlas.TileRock
				}
			case TilePlatform:
				tileImg = atlas.TilePlatform
			case TileSpike:
				tileImg = atlas.TileSpikes
			}

			if tileImg != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(screenX, screenY)
				screen.DrawImage(tileImg, op)
			}
		}
	}

	// Draw Goal (Lunar Lander, Airlock Door, or Warp Console)
	landerScreenX := lvl.LanderX - camX
	if landerScreenX > -40 && landerScreenX < 330 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(landerScreenX, lvl.LanderY)
		if lvl.GoalKind == "airlock" && atlas.AirlockDoor != nil {
			screen.DrawImage(atlas.AirlockDoor, op)
		} else if lvl.GoalKind == "warp_console" && atlas.WarpConsole != nil {
			screen.DrawImage(atlas.WarpConsole, op)
		} else {
			screen.DrawImage(atlas.Lander, op)
		}
	}
}

// GetLanderRect returns the collision rect of the Lander escape vehicle.
func (lvl *Level) GetLanderRect() physics.Rect {
	return physics.Rect{
		X: lvl.LanderX + 6,
		Y: lvl.LanderY + 8,
		W: 20,
		H: 24,
	}
}
