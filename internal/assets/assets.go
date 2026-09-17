package assets

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// TextureAtlas holds all cached game sprites and tiles.
type TextureAtlas struct {
	// Player
	AstronautIdle   *ebiten.Image
	AstronautRun1   *ebiten.Image
	AstronautRun2   *ebiten.Image
	AstronautJump   *ebiten.Image
	AstronautThrust *ebiten.Image
	AstronautHurt   *ebiten.Image

	// Enemies
	Blob1        *ebiten.Image
	Blob2        *ebiten.Image
	BlobSquished *ebiten.Image
	Hover1       *ebiten.Image
	Hover2       *ebiten.Image
	Boss1        *ebiten.Image
	Boss2        *ebiten.Image
	BossHurt     *ebiten.Image

	// Environment & Items
	TileSurface        *ebiten.Image
	TileRock           *ebiten.Image
	TileSurfaceCrimson *ebiten.Image
	TileRockCrimson    *ebiten.Image
	TileSurfaceIce     *ebiten.Image
	TileRockIce        *ebiten.Image
	TilePlatform       *ebiten.Image
	TileSpikes         *ebiten.Image
	Crystals           [4]*ebiten.Image
	Lander             *ebiten.Image
	Earth              *ebiten.Image
	Mars               *ebiten.Image
	Jupiter            *ebiten.Image

	// UI
	HeartFull  *ebiten.Image
	HeartEmpty *ebiten.Image
	ThrusterIcon *ebiten.Image

	// Weapons & Projectiles
	LaserBolt    *ebiten.Image
	NovaBall     *ebiten.Image
	CapsuleLaser *ebiten.Image
	CapsuleNova  *ebiten.Image
	IconLaser    *ebiten.Image
	IconNova     *ebiten.Image
}

var (
	atlasInstance *TextureAtlas
	atlasOnce     sync.Once
)

// Get returns the global TextureAtlas.
func Get() *TextureAtlas {
	atlasOnce.Do(func() {
		atlasInstance = buildAtlas()
	})
	return atlasInstance
}

func newImageFromPixels(w, h int, pixels []color.RGBA) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if idx < len(pixels) {
				img.SetRGBA(x, y, pixels[idx])
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func parsePixelArt(w, h int, rows []string, colormap map[rune]color.RGBA) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	trans := color.RGBA{0, 0, 0, 0}
	for y := 0; y < h; y++ {
		line := ""
		if y < len(rows) {
			line = rows[y]
		}
		for x := 0; x < w; x++ {
			c := trans
			if x < len(line) {
				ch := rune(line[x])
				if val, ok := colormap[ch]; ok {
					c = val
				}
			}
			img.SetRGBA(x, y, c)
		}
	}
	return ebiten.NewImageFromImage(img)
}

func buildAtlas() *TextureAtlas {
	a := &TextureAtlas{}

	// Common Sci-Fi Palette
	// Astronaut colors
	cOut := color.RGBA{18, 20, 32, 255}    // Dark navy outline
	cSuit := color.RGBA{235, 238, 245, 255} // White suit base
	cSuitSh := color.RGBA{160, 170, 195, 255} // Suit shadow
	cGold := color.RGBA{255, 190, 40, 255}  // Visor gold
	cGoldHi := color.RGBA{255, 240, 140, 255} // Visor highlight
	cPack := color.RGBA{110, 120, 145, 255} // Jetpack metal
	cRed := color.RGBA{225, 45, 60, 255}    // Red badge/accents
	cFlame1 := color.RGBA{255, 220, 50, 255} // Flame core
	cFlame2 := color.RGBA{255, 90, 20, 255}  // Flame outer

	astroMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': cOut,
		'S': cSuit,
		's': cSuitSh,
		'G': cGold,
		'H': cGoldHi,
		'P': cPack,
		'R': cRed,
		'F': cFlame1,
		'f': cFlame2,
	}

	// 16x16 Astronaut Idle
	idleRows := []string{
		"....######......",
		"...#SSSSH#......",
		"..#SSGGGGH#.....",
		"..#SGGGGGH#.....",
		"..#SSGGGG##.....",
		"...#SSSS##......",
		"..##SPPRSS#.....",
		".#SSSPPRSSS#....",
		".#SSSSSSSSS#....",
		".#SSssssssS#....",
		"..##ssssss##....",
		"...#SS##SS#.....",
		"...#ss##ss#.....",
		"...#ss##ss#.....",
		"...####.####....",
		"................",
	}
	a.AstronautIdle = parsePixelArt(16, 16, idleRows, astroMap)

	// Astronaut Run Frame 1
	run1Rows := []string{
		"....######......",
		"...#SSSSH#......",
		"..#SSGGGGH#.....",
		"..#SGGGGGH#.....",
		"..#SSGGGG##.....",
		"...#SSSS##......",
		"..##SPPRSS#.....",
		".#SSSPPRSSS#....",
		".#SSSSSSSSS#....",
		"..#SssssssS#....",
		"...#sss###......",
		"...#SS#..#SS#...",
		"...#ss#...#ss#..",
		"...####....####.",
		"................",
		"................",
	}
	a.AstronautRun1 = parsePixelArt(16, 16, run1Rows, astroMap)

	// Astronaut Run Frame 2
	run2Rows := []string{
		"....######......",
		"...#SSSSH#......",
		"..#SSGGGGH#.....",
		"..#SGGGGGH#.....",
		"..#SSGGGG##.....",
		"...#SSSS##......",
		"..##SPPRSS#.....",
		".#SSSPPRSSS#....",
		".#SSSSSSSSS#....",
		"..#SssssssS#....",
		"...###sss#......",
		"..#SS#..#SS#....",
		"..#ss#...#ss#...",
		".####....####...",
		"................",
		"................",
	}
	a.AstronautRun2 = parsePixelArt(16, 16, run2Rows, astroMap)

	// Astronaut Jump (floating tuck)
	jumpRows := []string{
		"....######......",
		"...#SSSSH#......",
		"..#SSGGGGH#.....",
		"..#SGGGGGH#.....",
		"..#SSGGGG##.....",
		"...#SSSS##......",
		"..##SPPRSS#.....",
		".#SSSPPRSSS#....",
		".#SSSSSSSSS#....",
		"..#SsssssS#.....",
		"...###sss###....",
		"..#SS#...#SS#...",
		"..####...####...",
		"................",
		"................",
		"................",
	}
	a.AstronautJump = parsePixelArt(16, 16, jumpRows, astroMap)

	// Astronaut Thrust (Jetpack firing)
	thrustRows := []string{
		"....######......",
		"...#SSSSH#......",
		"..#SSGGGGH#.....",
		"..#SGGGGGH#.....",
		"..#SSGGGG##.....",
		"...#SSSS##......",
		"..##SPPRSS#.....",
		".#PPSPPRSSS#....",
		".#PPSSSSSSS#....",
		".#ffSsssssS#....",
		".#Ff##sss###....",
		".#F#.#SS##SS#...",
		"..#..#ss##ss#...",
		".....####.####..",
		"................",
		"................",
	}
	a.AstronautThrust = parsePixelArt(16, 16, thrustRows, astroMap)

	// Astronaut Hurt (Flashing / damaged)
	hurtRows := []string{
		"....######......",
		"...#RRRRR##.....",
		"..#RRGGGGH#.....",
		"..#RGGGGGH#.....",
		"..#RRGGGG##.....",
		"...#RRRR##......",
		"..##RPPRRR#.....",
		".#RRRPPRRRR#....",
		".#RRRRRRRRR#....",
		".#RRssssssR#....",
		"..##ssssss##....",
		"...#RR##RR#.....",
		"...#ss##ss#.....",
		"...####.####....",
		"................",
		"................",
	}
	a.AstronautHurt = parsePixelArt(16, 16, hurtRows, astroMap)

	// -------------------------------------------------------------
	// ENEMIES
	// -------------------------------------------------------------
	// Alien Blob (Purple/Cyan Cute Slime)
	cBlobOut := color.RGBA{22, 10, 36, 255}
	cBlobBase := color.RGBA{160, 45, 215, 255}
	cBlobLight := color.RGBA{210, 95, 255, 255}
	cBlobShad := color.RGBA{95, 20, 135, 255}
	cEye := color.RGBA{240, 255, 240, 255}
	cPupil := color.RGBA{20, 10, 30, 255}

	blobMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': cBlobOut,
		'B': cBlobBase,
		'L': cBlobLight,
		'S': cBlobShad,
		'E': cEye,
		'P': cPupil,
	}

	blob1Rows := []string{
		"......####......",
		"....##LLLL##....",
		"...#LLBBBBBL#...",
		"..#LBBEEBEEBL#..",
		"..#LBPEEBPEBL#..",
		".#LBBBEEBEEBBL#.",
		".#BBBBBBBBBBBB#.",
		"#BBBBBBBBBBBBBB#",
		"#BBBBBBBBBBBBBB#",
		"#BBBBBBBBBBBBBB#",
		"#SBBBBBBBBBBBBS#",
		"#SSBBBBBBBBBBSS#",
		".#SSSSSSSSSSSS#.",
		"..#SSSSSSSSSS#..",
		"...##########...",
		"................",
	}
	a.Blob1 = parsePixelArt(16, 16, blob1Rows, blobMap)

	blob2Rows := []string{
		"................",
		".....######.....",
		"...##LLLLLL##...",
		"..#LLBBBBBBBL#..",
		".#LBBEEBBEEBBL#.",
		".#LBPEEBBPEEBBL#",
		"#LBBBEEBBEEBBBL#",
		"#BBBBBBBBBBBBBB#",
		"#BBBBBBBBBBBBBB#",
		"#BBBBBBBBBBBBBB#",
		"#SBBBBBBBBBBBBS#",
		"#SSBBBBBBBBBBSS#",
		"#SSSSSSSSSSSSSS#",
		".##############.",
		"................",
		"................",
	}
	a.Blob2 = parsePixelArt(16, 16, blob2Rows, blobMap)

	blobSquishRows := []string{
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"................",
		"....##########..",
		"..##LLBBBBBBLL##",
		".#LBBEEBBEEBBLS#",
		"#BBBBPEBBPEBBBBB#",
		"#SSSSSSSSSSSSSS#",
		".##############.",
		"................",
		"................",
		"................",
	}
	a.BlobSquished = parsePixelArt(16, 16, blobSquishRows, blobMap)

	// Alien Hover (Floating Cyclops / Jelly)
	cHoverOut := color.RGBA{15, 30, 45, 255}
	cHoverBase := color.RGBA{30, 195, 180, 255}
	cHoverLight := color.RGBA{110, 245, 230, 255}
	cHoverShad := color.RGBA{15, 115, 110, 255}
	cRedIris := color.RGBA{240, 40, 70, 255}

	hoverMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': cHoverOut,
		'H': cHoverBase,
		'L': cHoverLight,
		'S': cHoverShad,
		'E': cEye,
		'R': cRedIris,
		'P': cPupil,
	}

	hover1Rows := []string{
		".....######.....",
		"...##LLLLLL##...",
		"..#LLHHHHHHHL#..",
		".#LLHH####HHHL#.",
		".#LHH#EEEE#HHL#.",
		"#LHH#EERRPE#HHL#",
		"#HHH#EERRPE#HHH#",
		"#HHH#EERRPE#HHH#",
		"#HHH#EEEEEE#HHH#",
		"#SHH#E####E#HHS#",
		".#SHH######HHS#.",
		"..#SSHHHHHHSS#..",
		"...#SS####SS#...",
		"..#S#.#SS#.#S#..",
		"..#.#..##..#.#..",
		"................",
	}
	a.Hover1 = parsePixelArt(16, 16, hover1Rows, hoverMap)

	hover2Rows := []string{
		".....######.....",
		"...##LLLLLL##...",
		"..#LLHHHHHHHL#..",
		".#LLHH####HHHL#.",
		".#LHH#EEEE#HHL#.",
		"#LHH#EERRPE#HHL#",
		"#HHH#EERRPE#HHH#",
		"#HHH#EERRPE#HHH#",
		"#HHH#EEEEEE#HHH#",
		"#SHH#E####E#HHS#",
		".#SHH######HHS#.",
		"..#SSHHHHHHSS#..",
		"...#SS####SS#...",
		"...#S#.#S#.#S#..",
		"...#.#.###.#.#..",
		"................",
	}
	a.Hover2 = parsePixelArt(16, 16, hover2Rows, hoverMap)

	// Boss Overlord Mech (24x18 Heavy Armored Drone)
	cBossOut := color.RGBA{18, 12, 30, 255}
	cBossBase := color.RGBA{170, 35, 75, 255}
	cBossLight := color.RGBA{235, 80, 120, 255}
	cBossShad := color.RGBA{95, 18, 42, 255}
	cBossVisor := color.RGBA{255, 230, 60, 255}
	cBossCore := color.RGBA{255, 60, 40, 255}
	cBossThruster := color.RGBA{70, 230, 255, 255}
	cBossThrusterHi := color.RGBA{240, 255, 255, 255}

	bossMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': cBossOut,
		'B': cBossBase,
		'L': cBossLight,
		'S': cBossShad,
		'V': cBossVisor,
		'C': cBossCore,
		'T': cBossThruster,
		'W': cBossThrusterHi,
	}

	boss1Rows := []string{
		"........########........",
		"......##LLLLLLLL##......",
		"....##LLBBBBBBBBLL##....",
		"..##LLBB########BBLL##..",
		".#LLBB##VVVVVVVV##BBLL#.",
		"#LLBB#VVVVCCCCVVVV#BBLL#",
		"#LBBB#VVVCCCCCCVVV#BBBL#",
		"#BBBB#VVVVCCCCVVVV#BBBB#",
		"#BBBB##VVVVVVVVVV##BBBB#",
		"#SBBB###VVVVVVVV###BBBS#",
		"#SSBBBB##########BBBBSS#",
		".#SSBBBBBBBBBBBBBBBBSS#.",
		"..#SSSSBBBBBBBBSSSS#....",
		"...#SSSSSSSSSSSSSS#.....",
		"....#SS########SS#......",
		"....#TT#......#TT#......",
		"....#WW#......#WW#......",
		".....##........##.......",
	}
	a.Boss1 = parsePixelArt(24, 18, boss1Rows, bossMap)

	boss2Rows := []string{
		"........########........",
		"......##LLLLLLLL##......",
		"....##LLBBBBBBBBLL##....",
		"..##LLBB########BBLL##..",
		".#LLBB##VVVVVVVV##BBLL#.",
		"#LLBB#VVVVCCCCVVVV#BBLL#",
		"#LBBB#VVVCCCCCCVVV#BBBL#",
		"#BBBB#VVVVCCCCVVVV#BBBB#",
		"#BBBB##VVVVVVVVVV##BBBB#",
		"#SBBB###VVVVVVVV###BBBS#",
		"#SSBBBB##########BBBBSS#",
		".#SSBBBBBBBBBBBBBBBBSS#.",
		"..#SSSSBBBBBBBBSSSS#....",
		"...#SSSSSSSSSSSSSS#.....",
		"....#SS########SS#......",
		"....#WW#......#WW#......",
		"....#TT#......#TT#......",
		"....#TT#......#TT#......",
	}
	a.Boss2 = parsePixelArt(24, 18, boss2Rows, bossMap)

	// Flashing White/Red Hurt State
	bossHurtMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{255, 255, 255, 255},
		'B': color.RGBA{255, 220, 220, 255},
		'L': color.RGBA{255, 255, 255, 255},
		'S': color.RGBA{255, 180, 180, 255},
		'V': color.RGBA{255, 255, 255, 255},
		'C': color.RGBA{255, 100, 100, 255},
		'T': color.RGBA{255, 255, 255, 255},
		'W': color.RGBA{255, 255, 255, 255},
	}
	a.BossHurt = parsePixelArt(24, 18, boss1Rows, bossHurtMap)

	// -------------------------------------------------------------
	// TILES (16x16)
	// -------------------------------------------------------------
	// Moon Surface Tile: Gray regolith with highlighted top crest
	cTopRim := color.RGBA{215, 225, 240, 255}
	cSurfBase := color.RGBA{140, 150, 170, 255}
	cSurfCrater := color.RGBA{95, 105, 125, 255}
	cSurfDark := color.RGBA{65, 70, 90, 255}

	tileSurfMap := map[rune]color.RGBA{
		'T': cTopRim,
		'B': cSurfBase,
		'C': cSurfCrater,
		'D': cSurfDark,
	}
	tileSurfRows := []string{
		"TTTTTTTTTTTTTTTT",
		"BBBBBBBBBBBBBBBB",
		"BBBCDBBBBBBCDBBB",
		"BBCDDDBBBBCDDDBB",
		"BBCDDDBBBCDDDDDB",
		"BBBCDBBBBBCDDDDB",
		"BBBBBBBBBBBBCDBB",
		"BBBBBBCDBBBBBBBB",
		"BBBBCDDDBBBBBCDB",
		"BBBCDDDDDBBBCDDB",
		"BBBCDDDDDBBCDDDB",
		"BBBBCDDDDBBCDDDB",
		"DDDDDDDDDDDDDDDD",
		"DDDDDDDDDDDDDDDD",
		"DDDDDDDDDDDDDDDD",
		"DDDDDDDDDDDDDDDD",
	}
	a.TileSurface = parsePixelArt(16, 16, tileSurfRows, tileSurfMap)

	// Moon Underground / Basalt Rock
	tileRockRows := []string{
		"DDDDDDDDDDDDDDDD",
		"DDCDBBDDDDDCDBDD",
		"DCDBBBDDDDCDBBDD",
		"DDBBBDDDDDBBDDDD",
		"DDDDDDDDCDBDDDDD",
		"DDDDCDBDCDBBDDDD",
		"DDDCDBBBDDBBDDDD",
		"DDCDBBBDDDDDDDDD",
		"DDDDBBDDDDDCDBDD",
		"DDDDDDDDDDCDBBDD",
		"DDCDBDDDDDCDBBDD",
		"DCDBBDDDDDDBBDDD",
		"DDBBBDDDDDDDDDDD",
		"DDDDDDDDDDDDCDBD",
		"DDDDCDBDDDDDCDBD",
		"DDDDDDDDDDDDDDDD",
	}
	a.TileRock = parsePixelArt(16, 16, tileRockRows, tileSurfMap)

	// Sector 2: Crimson Martian / Phobos Tiles
	tileCrimsonMap := map[rune]color.RGBA{
		'T': color.RGBA{255, 175, 90, 255}, // Amber/gold highlighted rim
		'B': color.RGBA{175, 65, 45, 255},  // Rust red base
		'C': color.RGBA{115, 35, 25, 255},  // Deep crimson crater
		'D': color.RGBA{70, 20, 20, 255},   // Dark volcanic basalt
	}
	a.TileSurfaceCrimson = parsePixelArt(16, 16, tileSurfRows, tileCrimsonMap)
	a.TileRockCrimson = parsePixelArt(16, 16, tileRockRows, tileCrimsonMap)

	// Sector 3: Glacial Europa Ice Tiles
	tileIceMap := map[rune]color.RGBA{
		'T': color.RGBA{245, 255, 255, 255}, // Pure white frost rim
		'B': color.RGBA{130, 215, 245, 255}, // Pale icy blue
		'C': color.RGBA{60, 140, 190, 255},  // Glacial cyan crater
		'D': color.RGBA{25, 55, 90, 255},    // Deep abyss navy
	}
	a.TileSurfaceIce = parsePixelArt(16, 16, tileSurfRows, tileIceMap)
	a.TileRockIce = parsePixelArt(16, 16, tileRockRows, tileIceMap)

	// Floating Sci-Fi Metal Platform
	cPlatTop := color.RGBA{75, 215, 255, 255} // Glowing cyan rim
	cPlatMet1 := color.RGBA{80, 95, 120, 255}
	cPlatMet2 := color.RGBA{45, 55, 75, 255}
	cPlatBolt := color.RGBA{220, 230, 250, 255}
	platMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'C': cPlatTop,
		'M': cPlatMet1,
		'D': cPlatMet2,
		'B': cPlatBolt,
		'#': color.RGBA{20, 25, 35, 255},
	}
	platRows := []string{
		"CCCCCCCCCCCCCCCC",
		"################",
		"#BDMMMMMDDMMMMDB#",
		"#MDMMMMMDDMMMMDM#",
		"#MDMMMDDDDMMMMDM#",
		"#MDMMDDDDDMMMMDM#",
		"#MDDMMMMMMMMDDDM#",
		"#BDMMMMMDDMMMMDB#",
		"################",
		"....#DD#....#DD#",
		"....#DD#....#DD#",
		"....####....####",
		"................",
		"................",
		"................",
		"................",
	}
	a.TilePlatform = parsePixelArt(16, 16, platRows, platMap)

	// Cosmic Spikes / Hazards
	cSpike1 := color.RGBA{255, 45, 110, 255}
	cSpike2 := color.RGBA{180, 20, 70, 255}
	cSpikeHi := color.RGBA{255, 180, 220, 255}
	spikeMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{25, 10, 20, 255},
		'S': cSpike1,
		'D': cSpike2,
		'H': cSpikeHi,
	}
	spikeRows := []string{
		"................",
		"....#......#....",
		"...#H#....#H#...",
		"...#SH#...#SH#..",
		"..#SSDH#.#SSDH#.",
		"..#SSDSH##SSDSH#",
		".#SSSSDD##SSSSDD",
		".#SSSDDD##SSSDDD",
		"#SSSSSDDDSSSSSDD",
		"#SSSSDDDDSSSSD",
		"#SSSDDDDDSSSDDDD",
		"################",
		"################",
		"................",
		"................",
		"................",
	}
	a.TileSpikes = parsePixelArt(16, 16, spikeRows, spikeMap)

	// -------------------------------------------------------------
	// ENERGY CRYSTALS (4-frame rotation sparkle)
	// -------------------------------------------------------------
	cGemCore := color.RGBA{80, 245, 255, 255}
	cGemHi := color.RGBA{230, 255, 255, 255}
	cGemShad := color.RGBA{15, 140, 200, 255}
	cGemOut := color.RGBA{5, 45, 75, 255}
	gemMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': cGemOut,
		'C': cGemCore,
		'H': cGemHi,
		'S': cGemShad,
	}
	gemF1 := []string{
		".......##.......",
		"......#HH#......",
		".....#HHCC#.....",
		"....#HHCCCC#....",
		"...#HHCCCCCC#...",
		"..#HHCCCCCCC#...",
		"..#HCCCCCCCC#...",
		"..#CCCCCCCCC#...",
		"..#SSSSSSSS#....",
		"...#SSSSSS#.....",
		"....#SSSS#......",
		".....#SS#.......",
		"......##........",
		"................",
		"................",
		"................",
	}
	gemF2 := []string{
		".......##.......",
		"......#HC#......",
		".....#HCC#......",
		"....#HCCC#......",
		"...#HCCCC#......",
		"...#HCCCC#......",
		"...#CCCCC#......",
		"...#SCCCC#......",
		"...#SSCCC#......",
		"....#SSC#.......",
		".....#SC#.......",
		"......##........",
		"................",
		"................",
		"................",
		"................",
	}
	gemF3 := []string{
		".......##.......",
		"......#HH#......",
		"......#HC#......",
		"......#CC#......",
		"......#CC#......",
		"......#CC#......",
		"......#CC#......",
		"......#SC#......",
		"......#SS#......",
		"......#SS#......",
		"......#S#.......",
		".......#........",
		"................",
		"................",
		"................",
		"................",
	}
	gemF4 := []string{
		".......##.......",
		"......#CH#......",
		"......#CCH#.....",
		"......#CCCH#....",
		"......#CCCCH#...",
		"......#CCCCH#...",
		"......#CCCCC#...",
		"......#CCCCS#...",
		"......#CCCSS#...",
		".......#CSS#....",
		".......#CS#.....",
		"........##......",
		"................",
		"................",
		"................",
		"................",
	}
	a.Crystals[0] = parsePixelArt(16, 16, gemF1, gemMap)
	a.Crystals[1] = parsePixelArt(16, 16, gemF2, gemMap)
	a.Crystals[2] = parsePixelArt(16, 16, gemF3, gemMap)
	a.Crystals[3] = parsePixelArt(16, 16, gemF4, gemMap)

	// -------------------------------------------------------------
	// LUNAR LANDER (32x32)
	// -------------------------------------------------------------
	cLanderBody := color.RGBA{220, 225, 235, 255}
	cLanderShad := color.RGBA{140, 150, 170, 255}
	cLanderGold := color.RGBA{235, 180, 45, 255}
	cLanderGoldSh := color.RGBA{165, 115, 20, 255}
	cLanderGlass := color.RGBA{45, 195, 255, 255}
	cLanderRed := color.RGBA{240, 50, 60, 255}
	cLanderLeg := color.RGBA{95, 105, 120, 255}
	landerMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{15, 20, 30, 255},
		'B': cLanderBody,
		'S': cLanderShad,
		'G': cLanderGold,
		'g': cLanderGoldSh,
		'W': cLanderGlass,
		'R': cLanderRed,
		'L': cLanderLeg,
		'F': color.RGBA{255, 255, 255, 255}, // Flag white
	}
	landerRows := []string{
		"...............#................",
		"..............#R#...............",
		".............#RRR#..............",
		"..............###...............",
		".............#####..............",
		"............#BBBBB#.............",
		"...........#BBBBBBB#............",
		"..........#BBWWWWBBB#...........",
		"..........#BWWWWWWBB#...........",
		".........#BBWWWWWWBBB#..........",
		".........#BBBBBBBBBBB#..........",
		".........#BBBBSSSSBBB#..........",
		".........#############..........",
		"........#GGGGGGGGGGGGG#.........",
		".......#GGGGGGgGGGGGGGG#........",
		"......#GGGGGGgggGGGGGGGG#.......",
		"......#GGgGGGgggGGGgGGGG#.......",
		"......#GGgGGGGGGGGGgGGGG#.......",
		"......#ggggggggggggggggg#.......",
		".......#################........",
		".......#...#SSSSSSS#...#........",
		"......#L#..#SSSSSSS#..#L#.......",
		".....#L.#...#######...#.L#......",
		"....#L..#.............#..L#.....",
		"...#L...#.............#...L#....",
		"..#L.....#...........#.....L#...",
		".#L......#...........#......L#..",
		"#L........#.........#........L#.",
		"###......###.......###......###.",
		"................................",
		"................................",
		"................................",
	}
	a.Lander = parsePixelArt(32, 32, landerRows, landerMap)

	// -------------------------------------------------------------
	// EARTH (24x24) Glowing blue marble in sky
	// -------------------------------------------------------------
	cSpaceTrans := color.RGBA{0, 0, 0, 0}
	cOcean := color.RGBA{25, 80, 195, 255}
	cOceanHi := color.RGBA{65, 140, 245, 255}
	cLand := color.RGBA{45, 160, 85, 255}
	cLandHi := color.RGBA{110, 205, 105, 255}
	cCloud := color.RGBA{240, 245, 255, 220}
	cAtmo := color.RGBA{90, 180, 255, 120}

	earthMap := map[rune]color.RGBA{
		'.': cSpaceTrans,
		'A': cAtmo,
		'O': cOcean,
		'o': cOceanHi,
		'L': cLand,
		'l': cLandHi,
		'C': cCloud,
	}
	earthRows := []string{
		"........AAAAAAA.........",
		".....AAoooooooooAA......",
		"...AAoooollCCoooooAA....",
		"..AoooollLLCCCCoooooA...",
		".AooolLLLLLCCCCOOOOOOA..",
		".AooLLLLLLCOCCOOOOOOOA..",
		"AoooolLLLCCCOCOOOOOOOOA.",
		"AooooooollCCCOOCOOOOOOA.",
		"AooooCoooCCCCCOCCOOOOOA.",
		"AoCCCLlCCoCCCCCOOOOOOOA.",
		"ACCLLLLlCoOOOOOOOOOOOOA.",
		"ACLLLLLLCCOOOOOOOOOOOOA.",
		"AoLLLLLLLCCOOOOOOOOOOOA.",
		"AooLLLLLCCCOOOOOOOOOOOA.",
		"AooCCLLLCCOOOOOOOOOOOOA.",
		"AooooCCCCCOOOOOOOOOOOOA.",
		"AooooooCCOOOOOOOOOOOOOA.",
		".AooooooCOOOOOOOOOOOOA..",
		".AoooooOOOOOOOOOOOOOOA..",
		"..AooooOOOOOOOOOOOOOA...",
		"...AAooOOOOOOOOOOOAA....",
		".....AAOOOOOOOOOAA......",
		"........AAAAAAA.........",
		"........................",
	}
	a.Earth = parsePixelArt(24, 24, earthRows, earthMap)

	// Mars (24x24): Crimson planet with polar ice cap and rust basalt
	cMarsAtmo := color.RGBA{240, 140, 100, 100}
	cMarsRust := color.RGBA{210, 85, 45, 255}
	cMarsLight := color.RGBA{245, 125, 75, 255}
	cMarsDark := color.RGBA{135, 45, 25, 255}
	cMarsIce := color.RGBA{250, 250, 255, 240}
	marsMap := map[rune]color.RGBA{
		'.': color.RGBA{0, 0, 0, 0},
		'A': cMarsAtmo,
		'R': cMarsRust,
		'L': cMarsLight,
		'D': cMarsDark,
		'I': cMarsIce,
	}
	marsRows := []string{
		"........AAAAAAA.........",
		".....AAIIIILIIIIAA......",
		"...AAIIIILLLLLLIIIAA....",
		"..AIILLLLLRRRRRLLLLIA...",
		".ALLRRRRRRRRRRRRRRRLLA..",
		".ALRRRDDDRRRRRRRRRRRLLA.",
		"ALRRRDDDDDRRRRRRRRRRRRA.",
		"ALRRDDDDDDDRRRRRRDDDRRA.",
		"ALRRRRDDDDDRRRRDDDDDDRA.",
		"ALRRRRRRRRRRRRRDDDDDDRA.",
		"ALRRRRLLRRRRRRRDDDDDLLA.",
		"ALRRRLLLLRRRRRRRDDLLLLA.",
		"ALRRRLLLLRRRRRRRRLLLLLA.",
		"ALRRRRLLLRRRRRRRRRLLLLA.",
		"ALRRRRRRRRRRRRRRRRRLLLA.",
		"ALRRRRRRDDDDDRRRRRRLLLA.",
		"ALRRRRRDDDDDDRRRRRRRLLA.",
		".ALRRRDDDDDDDRRRRRRRLLA.",
		".ALRRRRDDDDDRRRRRRRRLLA.",
		"..ALRRRRRRRRRRRRRRRLLA..",
		"...AALRRRRRRRRRRRLLAA...",
		".....AALLLLLLLLLAA......",
		"........AAAAAAA.........",
		"........................",
	}
	a.Mars = parsePixelArt(24, 24, marsRows, marsMap)

	// Jupiter (32x24): Majestic gas giant with swirling bands and rings
	cJupRing := color.RGBA{240, 220, 170, 140}
	cJupCream := color.RGBA{240, 230, 205, 255}
	cJupAmber := color.RGBA{215, 150, 80, 255}
	cJupRust := color.RGBA{170, 85, 45, 255}
	cJupSpot := color.RGBA{225, 60, 40, 255} // Great Red Spot
	jupMap := map[rune]color.RGBA{
		'.': color.RGBA{0, 0, 0, 0},
		'~': cJupRing,
		'C': cJupCream,
		'A': cJupAmber,
		'R': cJupRust,
		'S': cJupSpot,
	}
	jupRows := []string{
		"............CCCCCCCC............",
		"........CCCCCAAAAACCCCC.........",
		"......CCCAAAAAAAAAAAAACCC.......",
		"....CCCAAAAAAAAAAAAAAAAACCC.....",
		"...CCRRRRRRRRRRRRRRRRRRRRCCC....",
		"..CCRRRRRRRRRRRRRRRRRRRRRRCCC...",
		".CCCAAAAAAAAAAAAACCAAAAAAAACCC..",
		".CAAAAAAAAAAAAAAAAAAAAAAAAAAARC.",
		"CAAAAAAAAAAAAAAAAAAAAAAAAAAAAARC",
		"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR",
		"RRRRRRRRRRRRRRRRRRRSSSSRRRRRRRRR",
		"~~~~~~~~~~~~~~~~~~RSSSSR~~~~~~~~",
		"~~~~~~~~~~~~~~~~~~RSSSSR~~~~~~~~",
		"CAAAAAAAAAAAAAAAAAASSSSAAAAAAARC",
		"CAAAAAAAAAAAAAAAAAAAAAAAAAAAAARC",
		".CAAAAAAAAAAAAAAAAAAAAAAAAAAARC.",
		".CCCAAAAAAAAAAAAACCAAAAAAAACCC..",
		"..CCRRRRRRRRRRRRRRRRRRRRRRCCC...",
		"...CCRRRRRRRRRRRRRRRRRRRRCCC....",
		"....CCCAAAAAAAAAAAAAAAAACCC.....",
		"......CCCAAAAAAAAAAAAACCC.......",
		"........CCCCCAAAAACCCCC.........",
		"............CCCCCCCC............",
		"................................",
	}
	a.Jupiter = parsePixelArt(32, 24, jupRows, jupMap)

	// -------------------------------------------------------------
	// UI ICONS
	// -------------------------------------------------------------
	// Heart Full & Empty (9x8)
	heartRowsFull := []string{
		".##.##..",
		"#RR#RR#.",
		"#RRRRR#.",
		"#RRRRR#.",
		".#RRR#..",
		"..#R#...",
		"...#....",
		"........",
	}
	heartRowsEmpty := []string{
		".##.##..",
		"#..#..#.",
		"#.....#.",
		"#.....#.",
		".#...#..",
		"..#.#...",
		"...#....",
		"........",
	}
	heartMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{20, 10, 25, 255},
		'R': color.RGBA{245, 45, 75, 255},
	}
	a.HeartFull = parsePixelArt(8, 8, heartRowsFull, heartMap)
	a.HeartEmpty = parsePixelArt(8, 8, heartRowsEmpty, heartMap)

	// Thruster Gauge Icon (8x8)
	thrustIconRows := []string{
		"..####..",
		".#FFFF#.",
		"#FFffFF#",
		"#FffffF#",
		"#ffffff#",
		".#ffff#.",
		"..#ff#..",
		"...##...",
	}
	thrustIconMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{30, 20, 10, 255},
		'F': color.RGBA{255, 220, 60, 255},
		'f': color.RGBA{255, 80, 20, 255},
	}
	a.ThrusterIcon = parsePixelArt(8, 8, thrustIconRows, thrustIconMap)

	// -------------------------------------------------------------
	// WEAPONS & PROJECTILES
	// -------------------------------------------------------------
	// 1. Laser Bolt (8x4): High-speed cyan photon beam
	laserMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{10, 50, 100, 255},
		'C': color.RGBA{40, 210, 255, 255},
		'H': color.RGBA{180, 250, 255, 255},
		'W': color.RGBA{255, 255, 255, 255},
	}
	laserRows := []string{
		".######.",
		"#CHHHWW#",
		"#CHHHWW#",
		".######.",
	}
	a.LaserBolt = parsePixelArt(8, 4, laserRows, laserMap)

	// 2. Nova Ball (8x8): Bouncing fiery plasma star
	novaMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{90, 20, 10, 255},
		'O': color.RGBA{255, 95, 20, 255},
		'Y': color.RGBA{255, 210, 40, 255},
		'W': color.RGBA{255, 255, 210, 255},
	}
	novaRows := []string{
		"..####..",
		".#OYYO#.",
		"#OYYYYO#",
		"#YYYYWW#",
		"#YYYYWW#",
		"#OYYYYO#",
		".#OYYO#.",
		"..####..",
	}
	a.NovaBall = parsePixelArt(8, 8, novaRows, novaMap)

	// 3. Holographic Weapon Pod: Laser Blaster (16x16)
	capsuleLaserMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{20, 30, 50, 255},
		'D': color.RGBA{60, 140, 220, 180},  // Hologram glass dome
		'G': color.RGBA{120, 220, 255, 255}, // Glass highlight
		'M': color.RGBA{100, 115, 140, 255}, // Metal base
		'B': color.RGBA{40, 210, 255, 255},  // Gun laser barrel
		'W': color.RGBA{255, 255, 255, 255},
	}
	capsuleLaserRows := []string{
		".....######.....",
		"...##DDGGDD##...",
		"..#DDGGGGDDDD#..",
		".#DDGGBBBBBDDD#.",
		".#DDGBBBBBBDDD#.",
		"#DDGG...#BBDDDD#",
		"#DDGG...#BBDDDD#",
		"#DDGG...#BBDDDD#",
		".#DDG...#BDDDD#.",
		".#DDGGGGDDDDDD#.",
		"...##DDDDDD##...",
		"....#MMMMMM#....",
		"...#MMMMMMMM#...",
		"..#MMMMMMMMMM#..",
		"..############..",
		"................",
	}
	a.CapsuleLaser = parsePixelArt(16, 16, capsuleLaserRows, capsuleLaserMap)

	// 4. Holographic Weapon Pod: Nova Cannon (16x16)
	capsuleNovaMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{30, 20, 40, 255},
		'D': color.RGBA{180, 100, 220, 180}, // Violet hologram glass
		'G': color.RGBA{230, 160, 255, 255},
		'M': color.RGBA{110, 100, 130, 255}, // Metal base
		'N': color.RGBA{255, 190, 30, 255},  // Golden nova orb inside
		'O': color.RGBA{255, 80, 20, 255},
		'W': color.RGBA{255, 255, 200, 255},
	}
	capsuleNovaRows := []string{
		".....######.....",
		"...##DDGGDD##...",
		"..#DDGGGGDDDD#..",
		".#DDGG..N.DDDD#.",
		".#DDG.NWWN.DDD#.",
		"#DDGG.NOWO.DDDD#",
		"#DDGG..N..DDDDD#",
		"#DDGG.....DDDDD#",
		".#DDGG...DDDDD#.",
		".#DDGGGGDDDDDD#.",
		"...##DDDDDD##...",
		"....#MMMMMM#....",
		"...#MMMMMMMM#...",
		"..#MMMMMMMMMM#..",
		"..############..",
		"................",
	}
	a.CapsuleNova = parsePixelArt(16, 16, capsuleNovaRows, capsuleNovaMap)

	// 5. HUD Weapon Icons (8x8)
	iconLaserMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{10, 30, 60, 255},
		'C': color.RGBA{60, 215, 255, 255},
		'W': color.RGBA{240, 255, 255, 255},
		'M': color.RGBA{140, 160, 190, 255},
	}
	iconLaserRows := []string{
		"........",
		"...#####",
		"..#CCCW#",
		"..#MMMM#",
		"....#M#.",
		"....#M#.",
		"....###.",
		"........",
	}
	a.IconLaser = parsePixelArt(8, 8, iconLaserRows, iconLaserMap)

	iconNovaMap := map[rune]color.RGBA{
		'.': {0, 0, 0, 0},
		'#': color.RGBA{50, 10, 10, 255},
		'Y': color.RGBA{255, 210, 40, 255},
		'O': color.RGBA{255, 80, 20, 255},
		'W': color.RGBA{255, 255, 220, 255},
	}
	iconNovaRows := []string{
		"...##...",
		"..#YY#..",
		".#OWWYO.",
		"#YYYYYY#",
		"#YYYYYY#",
		".#OYYO#.",
		"..#YY#..",
		"...##...",
	}
	a.IconNova = parsePixelArt(8, 8, iconNovaRows, iconNovaMap)

	return a
}

// GenerateStarfield creates a 320x180 deep space background with stars and nebula.
func GenerateStarfield(w, h int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Deep space gradient (Navy/Purple black)
			ny := float64(y) / float64(h)
			r := uint8(8 + 14*ny)
			g := uint8(6 + 10*ny)
			b := uint8(22 + 28*ny)

			// Subtle nebula cloud patch
			distNebula := math.Hypot(float64(x)-float64(w)*0.7, float64(y)-float64(h)*0.3)
			if distNebula < 90 {
				factor := (1.0 - distNebula/90) * 0.35
				r = uint8(float64(r) + 80*factor)
				b = uint8(float64(b) + 110*factor)
			}

			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Scatter twinkling stars deterministically
	starRNG := uint32(1337)
	nextRand := func() uint32 {
		starRNG = starRNG*1664525 + 1013904223
		return starRNG
	}

	numStars := 120
	for i := 0; i < numStars; i++ {
		sx := int(nextRand() % uint32(w))
		sy := int(nextRand() % uint32(h))
		brightness := uint8(140 + (nextRand() % 115))
		c := color.RGBA{brightness, brightness, uint8(math.Min(255, float64(brightness)+20)), 255}
		img.SetRGBA(sx, sy, c)
		if i%8 == 0 && sx+1 < w && sy+1 < h {
			// Cross sparkle on brighter stars
			img.SetRGBA(sx+1, sy, color.RGBA{brightness / 2, brightness / 2, brightness / 2, 255})
			img.SetRGBA(sx-1, sy, color.RGBA{brightness / 2, brightness / 2, brightness / 2, 255})
			img.SetRGBA(sx, sy+1, color.RGBA{brightness / 2, brightness / 2, brightness / 2, 255})
			img.SetRGBA(sx, sy-1, color.RGBA{brightness / 2, brightness / 2, brightness / 2, 255})
		}
	}

	return ebiten.NewImageFromImage(img)
}
