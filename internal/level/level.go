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
	MaxLevels   = 3
)

const (
	TileEmpty    = 0
	TileSurface  = 1
	TileRock     = 2
	TilePlatform = 3
	TileSpike    = 4
)

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

// Level represents a playable stage with tile grid, hazards, and background layers.
type Level struct {
	SectorIndex     int
	Name            string
	Theme           string // "moon", "crimson", "ice"
	Celestial       string // "earth", "mars", "jupiter"
	Width           int
	Height          int
	Tiles           [][]int
	Crystals        []CrystalData
	EnemySpawns     []EnemySpawn
	WeaponSpawns    []WeaponSpawn
	MovingPlatforms []MovingPlatformData
	Geysers         []GeyserData
	Medals          []MedalData
	LanderX         float64
	LanderY         float64
	HasBoss         bool
	BossX           float64
	BossY           float64
	StarfieldImg    *ebiten.Image
}

func makeGrid(w, h int) [][]int {
	g := make([][]int, h)
	for y := 0; y < h; y++ {
		g[y] = make([]int, w)
	}
	return g
}

// LoadLevel loads the level for the given sector (1, 2, or 3).
func LoadLevel(sector int) *Level {
	switch sector {
	case 2:
		return NewLevel2()
	case 3:
		return NewLevel3()
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
	lvl.Tiles[7][81] = TilePlatform
	lvl.Tiles[7][82] = TilePlatform
	lvl.Tiles[6][86] = TilePlatform
	lvl.Tiles[6][87] = TilePlatform
	lvl.Tiles[5][90] = TilePlatform
	lvl.Tiles[5][91] = TilePlatform

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
		{TileX: 71, TileY: 2}, {TileX: 76, TileY: 7}, {TileX: 82, TileY: 5},
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
		{TileX: 84, TileY: 5, Kind: "hover"},
		{TileX: 89, TileY: 4, Kind: "hover"},
		{TileX: 102, TileY: 6, Kind: "blob"},
		{TileX: 112, TileY: 6, Kind: "blob"},
	}

	// 4. Weapon Capsules
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 55, TileY: 2, Kind: "nova"}, // Nova Cannon on high secret ledge
	}

	// 5. Moving Platforms
	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(80 * TileSize),
			StartY: float64(7 * TileSize),
			EndX:   float64(89 * TileSize),
			EndY:   float64(7 * TileSize),
			W:      32,
			H:      8,
			Speed:  35.0,
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
	lvl.Tiles[7][91] = TilePlatform
	lvl.Tiles[7][92] = TilePlatform
	lvl.Tiles[5][96] = TilePlatform
	lvl.Tiles[5][97] = TilePlatform
	lvl.Tiles[4][101] = TilePlatform
	lvl.Tiles[4][102] = TilePlatform
	lvl.Tiles[6][106] = TilePlatform

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
		{TileX: 92, TileY: 5}, {TileX: 97, TileY: 3}, {TileX: 102, TileY: 2},
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
		{TileX: 94, TileY: 5, Kind: "hover"},
		{TileX: 99, TileY: 3, Kind: "hover"},
		{TileX: 116, TileY: 6, Kind: "blob"},
		{TileX: 124, TileY: 6, Kind: "hover"},
	}

	// 4. Weapon Capsules (Laser and Nova caches for players wanting quick rearm)
	lvl.WeaponSpawns = []WeaponSpawn{
		{TileX: 24, TileY: 5, Kind: "laser"},
		{TileX: 82, TileY: 3, Kind: "nova"},
	}

	// 5. Moving Platforms across the frozen abyss
	lvl.MovingPlatforms = []MovingPlatformData{
		{
			StartX: float64(93 * TileSize),
			StartY: float64(6 * TileSize),
			EndX:   float64(103 * TileSize),
			EndY:   float64(6 * TileSize),
			W:      32,
			H:      8,
			Speed:  40.0,
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
	}

	for i := -1; i < 4; i++ {
		bx := float32(float64(i*160) - ridgeOffset)
		vector.DrawFilledRect(screen, bx, 130, 160, 50, ridgeColor, false)
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
				} else {
					tileImg = atlas.TileSurface
				}
			case TileRock:
				if lvl.Theme == "crimson" {
					tileImg = atlas.TileRockCrimson
				} else if lvl.Theme == "ice" {
					tileImg = atlas.TileRockIce
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

	// Draw Lunar Lander Goal
	landerScreenX := lvl.LanderX - camX
	if landerScreenX > -40 && landerScreenX < 330 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(landerScreenX, lvl.LanderY)
		screen.DrawImage(atlas.Lander, op)
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
