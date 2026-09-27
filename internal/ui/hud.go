package ui

import (
	"image/color"
	"strconv"
	"strings"

	"astroleap/internal/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func cleanChar(ch rune) rune {
	switch ch {
	case 'É', 'È', 'Ê', 'Ë':
		return 'E'
	case 'À', 'Â', 'Ä':
		return 'A'
	case 'Î', 'Ï':
		return 'I'
	case 'Ô', 'Ö':
		return 'O'
	case 'Ù', 'Û', 'Ü':
		return 'U'
	case 'Ç':
		return 'C'
	default:
		return ch
	}
}

// DrawText draws text using a built-in 5x7 retro bitmap font.
func DrawText(screen *ebiten.Image, str string, x, y int, col color.RGBA) {
	str = strings.ToUpper(str)
	cx := x
	cy := y

	// Drop shadow for legibility
	drawTextRaw(screen, str, cx+1, cy+1, color.RGBA{10, 10, 20, 200})
	drawTextRaw(screen, str, cx, cy, col)
}

func drawTextRaw(screen *ebiten.Image, str string, startX, startY int, col color.RGBA) {
	cx := startX
	cy := startY

	for _, rawCh := range str {
		if rawCh == '\n' {
			cx = startX
			cy += 9
			continue
		}
		if rawCh == ' ' {
			cx += 5
			continue
		}

		ch := cleanChar(rawCh)

		glyph, ok := fontData[ch]
		if !ok {
			// Only draw '?' if the actual string character was '?'
			if rawCh == '?' {
				glyph = fontData['?']
			} else {
				cx += 4
				continue
			}
		}

		for row := 0; row < 7; row++ {
			bits := glyph[row]
			for colIdx := 0; colIdx < 5; colIdx++ {
				if (bits & (1 << (4 - colIdx))) != 0 {
					screen.Set(cx+colIdx, cy+row, col)
				}
			}
		}
		cx += 6
	}
}

// DrawHUD draws the top status bar with clear grouping, breathing room, and sharp typography.
func DrawHUD(screen *ebiten.Image, health int, maxHealth int, crystals int, score int, fuel float64, maxFuel float64, activeWeapon int, hasMultipleWeapons bool, sector int, elapsedTime float64, medals int, hasKey bool, repairCores int) {
	atlas := assets.Get()

	// 0. Translucent HUD backdrop strip (Y: 0..16) for clean readability against bright starfields
	vector.DrawFilledRect(screen, 0, 0, 320, 16, color.RGBA{8, 12, 22, 190}, false)
	vector.StrokeLine(screen, 0, 16, 320, 16, 1, color.RGBA{45, 60, 95, 140}, false)

	// 1. Health Hearts (X: 8, 18, 28)
	for i := 0; i < maxHealth; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(8+i*10), 4)
		if i < health {
			screen.DrawImage(atlas.HeartFull, op)
		} else {
			screen.DrawImage(atlas.HeartEmpty, op)
		}
	}

	// 2. Score (X: 44)
	scoreStr := formatScore(score)
	DrawText(screen, scoreStr, 44, 5, color.RGBA{240, 240, 255, 255})

	// 3. Sector Indicator (X: 84)
	secStr := "S" + strconv.Itoa(sector) + "/5"
	DrawText(screen, secStr, 84, 5, color.RGBA{255, 215, 60, 255})

	// 4. Energy Crystals (X: 118 icon, X: 132 count)
	if atlas.Crystals[0] != nil {
		crystalOp := &ebiten.DrawImageOptions{}
		crystalOp.GeoM.Scale(0.75, 0.75)
		crystalOp.GeoM.Translate(118, 2)
		screen.DrawImage(atlas.Crystals[0], crystalOp)
	}
	crystStr := "X" + formatTwoDigits(crystals)
	DrawText(screen, crystStr, 132, 5, color.RGBA{80, 240, 255, 255})

	// 5. Secret Medals (X: 158)
	medStr := "M:" + strconv.Itoa(medals)
	DrawText(screen, medStr, 158, 5, color.RGBA{255, 205, 50, 255})

	// 6. Sector Objectives or Speedrun Timer (X: 188..236)
	if sector == 4 {
		// Sector 4: Keycard indicator
		if hasKey {
			if atlas.IconKey != nil {
				opKey := &ebiten.DrawImageOptions{}
				opKey.GeoM.Translate(186, 4)
				screen.DrawImage(atlas.IconKey, opKey)
			}
			DrawText(screen, "KEY", 196, 5, color.RGBA{255, 225, 60, 255})
		} else {
			DrawText(screen, "NO KEY", 186, 5, color.RGBA{130, 140, 160, 255})
		}
		timeStr := formatTimer(elapsedTime)
		DrawText(screen, timeStr, 226, 5, color.RGBA{180, 210, 240, 255})
	} else if sector == 5 {
		// Sector 5: Repair Cores indicator
		repStr := "REP:" + strconv.Itoa(repairCores) + "/3"
		repCol := color.RGBA{255, 180, 50, 255}
		if repairCores >= 3 {
			repCol = color.RGBA{80, 255, 120, 255}
		}
		DrawText(screen, repStr, 186, 5, repCol)
		timeStr := formatTimer(elapsedTime)
		DrawText(screen, timeStr, 226, 5, color.RGBA{180, 210, 240, 255})
	} else {
		// Standard live speedrun timer
		timeStr := formatTimer(elapsedTime)
		DrawText(screen, timeStr, 196, 5, color.RGBA{190, 225, 255, 255})
	}

	// 7. Jetpack Thruster Fuel Gauge (X: 250..312)
	iconOp := &ebiten.DrawImageOptions{}
	iconOp.GeoM.Translate(250, 4)
	screen.DrawImage(atlas.ThrusterIcon, iconOp)

	// Fuel bar
	barX := float32(260)
	barY := float32(5)
	barW := float32(52)
	barH := float32(6)

	vector.DrawFilledRect(screen, barX, barY, barW, barH, color.RGBA{14, 18, 30, 220}, false)
	vector.StrokeRect(screen, barX, barY, barW, barH, 1, color.RGBA{70, 90, 130, 255}, false)

	ratio := float32(fuel / maxFuel)
	if ratio < 0 {
		ratio = 0
	} else if ratio > 1 {
		ratio = 1
	}

	fillW := (barW - 2) * ratio
	if fillW > 0 {
		fillColor := color.RGBA{50, 220, 255, 255}
		if ratio < 0.25 {
			fillColor = color.RGBA{255, 60, 50, 255} // Red warning
		} else if ratio < 0.5 {
			fillColor = color.RGBA{255, 180, 30, 255} // Orange
		}
		vector.DrawFilledRect(screen, barX+1, barY+1, fillW, barH-2, fillColor, false)
	}

	// 8. Active Weapon Indicator (Clean floating badge, Y: 18)
	if activeWeapon > 0 {
		var icon *ebiten.Image
		wName := ""
		wCol := color.RGBA{60, 220, 255, 255}
		if activeWeapon == 1 {
			icon = atlas.IconLaser
			wName = "LASER [J]"
		} else if activeWeapon == 2 {
			icon = atlas.IconNova
			wName = "NOVA [J]"
			wCol = color.RGBA{255, 200, 40, 255}
		}

		badgeW := float32(78)
		if hasMultipleWeapons {
			badgeW = float32(118)
		}
		vector.DrawFilledRect(screen, 6, 18, badgeW, 11, color.RGBA{10, 14, 24, 200}, false)
		vector.StrokeRect(screen, 6, 18, badgeW, 11, 1, color.RGBA{45, 60, 95, 160}, false)

		if icon != nil {
			iconOp := &ebiten.DrawImageOptions{}
			iconOp.GeoM.Translate(9, 20)
			screen.DrawImage(icon, iconOp)
		}
		DrawText(screen, wName, 20, 20, wCol)

		if hasMultipleWeapons {
			DrawText(screen, "Q:SWAP", 76, 20, color.RGBA{170, 185, 220, 230})
		}
	}
}

func formatScore(s int) string {
	str := ""
	if s == 0 {
		return "00000"
	}
	for s > 0 {
		digit := s % 10
		str = string(rune('0'+digit)) + str
		s /= 10
	}
	for len(str) < 5 {
		str = "0" + str
	}
	return str
}

func formatTwoDigits(n int) string {
	d1 := (n / 10) % 10
	d2 := n % 10
	return string([]rune{rune('0' + d1), rune('0' + d2)})
}

func formatTimer(t float64) string {
	if t < 0 {
		t = 0
	}
	mins := int(t) / 60
	secs := int(t) % 60
	return string([]rune{
		rune('0' + (mins/10)%10),
		rune('0' + mins%10),
		':',
		rune('0' + (secs/10)%10),
		rune('0' + secs%10),
	})
}

// 5x7 Minimal Retro Font Definition
var fontData = map[rune][7]byte{
	'A': {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'B': {0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E},
	'C': {0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E},
	'D': {0x1C, 0x12, 0x11, 0x11, 0x11, 0x12, 0x1C},
	'E': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F},
	'F': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10},
	'G': {0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0E},
	'H': {0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'I': {0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E},
	'J': {0x07, 0x02, 0x02, 0x02, 0x02, 0x12, 0x0C},
	'K': {0x11, 0x12, 0x14, 0x18, 0x14, 0x12, 0x11},
	'L': {0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x1F},
	'M': {0x11, 0x1B, 0x15, 0x11, 0x11, 0x11, 0x11},
	'N': {0x11, 0x19, 0x15, 0x13, 0x11, 0x11, 0x11},
	'O': {0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'P': {0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10},
	'Q': {0x0E, 0x11, 0x11, 0x11, 0x15, 0x12, 0x0D},
	'R': {0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11},
	'S': {0x0E, 0x11, 0x10, 0x0E, 0x01, 0x11, 0x0E},
	'T': {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04},
	'U': {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'V': {0x11, 0x11, 0x11, 0x11, 0x11, 0x0A, 0x04},
	'W': {0x11, 0x11, 0x11, 0x15, 0x15, 0x15, 0x0A},
	'X': {0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11},
	'Y': {0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04},
	'Z': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F},
	'0': {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E},
	'1': {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E},
	'2': {0x0E, 0x11, 0x01, 0x06, 0x08, 0x10, 0x1F},
	'3': {0x1F, 0x01, 0x02, 0x06, 0x01, 0x11, 0x0E},
	'4': {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02},
	'5': {0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E},
	'6': {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E},
	'7': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08},
	'8': {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E},
	'9': {0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C},
	':': {0x00, 0x0C, 0x0C, 0x00, 0x0C, 0x0C, 0x00},
	'!': {0x04, 0x04, 0x04, 0x04, 0x04, 0x00, 0x04},
	'?': {0x0E, 0x11, 0x01, 0x02, 0x04, 0x00, 0x04},
	'-': {0x00, 0x00, 0x00, 0x1F, 0x00, 0x00, 0x00},
	'+': {0x00, 0x04, 0x04, 0x1F, 0x04, 0x04, 0x00},
	'/': {0x01, 0x02, 0x04, 0x08, 0x10, 0x00, 0x00},
	'.': {0x00, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x0C},
	',': {0x00, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x08},
	'\'': {0x04, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00},
	'"': {0x0A, 0x0A, 0x14, 0x00, 0x00, 0x00, 0x00},
	'(': {0x02, 0x04, 0x08, 0x08, 0x08, 0x04, 0x02},
	')': {0x08, 0x04, 0x02, 0x02, 0x02, 0x04, 0x08},
	'[': {0x0E, 0x08, 0x08, 0x08, 0x08, 0x08, 0x0E},
	']': {0x0E, 0x02, 0x02, 0x02, 0x02, 0x02, 0x0E},
	'<': {0x02, 0x04, 0x08, 0x10, 0x08, 0x04, 0x02},
	'>': {0x08, 0x04, 0x02, 0x01, 0x02, 0x04, 0x08},
	'*': {0x00, 0x15, 0x0E, 0x1F, 0x0E, 0x15, 0x00},
	'=': {0x00, 0x00, 0x1F, 0x00, 0x1F, 0x00, 0x00},
	'_': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1F},
	'%': {0x19, 0x1A, 0x04, 0x08, 0x10, 0x0B, 0x13},
	'#': {0x0A, 0x0A, 0x1F, 0x0A, 0x1F, 0x0A, 0x0A},
	';': {0x00, 0x0C, 0x0C, 0x00, 0x0C, 0x08, 0x00},
	'★': {0x04, 0x04, 0x1F, 0x0E, 0x0E, 0x15, 0x11},
}

// WrapText wraps text into lines that fit within maxWidthPx (each char is 6px).
func WrapText(str string, maxWidthPx int) []string {
	maxChars := maxWidthPx / 6
	if maxChars <= 0 {
		maxChars = 1
	}
	words := strings.Fields(str)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	currentLine := ""

	for _, w := range words {
		if currentLine == "" {
			currentLine = w
		} else if len(currentLine)+1+len(w) <= maxChars {
			currentLine += " " + w
		} else {
			lines = append(lines, currentLine)
			currentLine = w
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

// DrawDialogueBox draws a stylish semi-transparent bordered modal dialog with wrapped text.
func DrawDialogueBox(screen *ebiten.Image, x, y, w, h float32, title string, body string, prompt string) {
	// Dark sci-fi slate box
	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{12, 16, 28, 235}, false)
	vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{80, 220, 255, 255}, false)

	lineY := int(y) + 8
	if title != "" {
		DrawText(screen, title, int(x)+10, lineY, color.RGBA{255, 220, 60, 255})
		lineY += 12
	}

	lines := WrapText(body, int(w)-20)
	for _, l := range lines {
		DrawText(screen, l, int(x)+10, lineY, color.RGBA{215, 230, 250, 255})
		lineY += 9
	}

	if prompt != "" {
		DrawText(screen, prompt, int(x)+10, int(y+h)-11, color.RGBA{255, 205, 80, 255})
	}
}
