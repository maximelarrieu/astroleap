package ui

import (
	"image/color"
	"strings"

	"astroleap/internal/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

	for _, ch := range str {
		if ch == '\n' {
			cx = startX
			cy += 9
			continue
		}
		if ch == ' ' {
			cx += 5
			continue
		}

		glyph, ok := fontData[ch]
		if !ok {
			glyph = fontData['?']
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

// DrawHUD draws the top status bar (hearts, score, sector, crystals, medals, timer, thruster fuel, active weapon).
func DrawHUD(screen *ebiten.Image, health int, maxHealth int, crystals int, score int, fuel float64, maxFuel float64, activeWeapon int, hasMultipleWeapons bool, sector int, elapsedTime float64, medals int) {
	atlas := assets.Get()

	// 1. Health (Top-Left)
	for i := 0; i < maxHealth; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(6+i*9), 6)
		if i < health {
			screen.DrawImage(atlas.HeartFull, op)
		} else {
			screen.DrawImage(atlas.HeartEmpty, op)
		}
	}

	// 2. Score
	scoreStr := formatScore(score)
	DrawText(screen, scoreStr, 36, 7, color.RGBA{240, 240, 255, 255})

	// 3. Sector Indicator (e.g. "SEC 1/3")
	secDigit := '1'
	if sector >= 1 && sector <= 3 {
		secDigit = rune('0' + sector)
	}
	secStr := "S" + string(secDigit) + "/3"
	DrawText(screen, secStr, 76, 7, color.RGBA{255, 215, 60, 255})

	// 4. Energy Crystals (Top-Center)
	crystalOp := &ebiten.DrawImageOptions{}
	crystalOp.GeoM.Translate(108, 3)
	screen.DrawImage(atlas.Crystals[0], crystalOp)
	crystStr := "X" + formatTwoDigits(crystals)
	DrawText(screen, crystStr, 124, 7, color.RGBA{100, 245, 255, 255})

	// 5. Secret Medals (M: X/3)
	medDigit := rune('0' + medals)
	if medals < 0 {
		medDigit = '0'
	} else if medals > 3 {
		medDigit = '3'
	}
	medStr := "M:" + string(medDigit) + "/3"
	DrawText(screen, medStr, 150, 7, color.RGBA{255, 205, 50, 255})

	// 6. Live Speedrun Timer (MM:SS)
	timeStr := formatTimer(elapsedTime)
	DrawText(screen, timeStr, 192, 7, color.RGBA{200, 230, 255, 255})

	// 7. Jetpack Thruster Fuel Gauge (Top-Right)
	// Fuel icon
	iconOp := &ebiten.DrawImageOptions{}
	iconOp.GeoM.Translate(238, 6)
	screen.DrawImage(atlas.ThrusterIcon, iconOp)

	// Fuel bar background & border
	barX := float32(250)
	barY := float32(7)
	barW := float32(62)
	barH := float32(6)

	// Background
	vector.DrawFilledRect(screen, barX, barY, barW, barH, color.RGBA{18, 22, 35, 220}, false)
	// Outline
	vector.StrokeRect(screen, barX, barY, barW, barH, 1, color.RGBA{80, 100, 140, 255}, false)

	// Fill proportion
	ratio := float32(fuel / maxFuel)
	if ratio < 0 {
		ratio = 0
	} else if ratio > 1 {
		ratio = 1
	}

	fillW := (barW - 2) * ratio
	if fillW > 0 {
		fillColor := color.RGBA{60, 220, 255, 255}
		if ratio < 0.25 {
			fillColor = color.RGBA{255, 70, 60, 255} // Red warning
		} else if ratio < 0.5 {
			fillColor = color.RGBA{255, 190, 40, 255} // Orange
		}
		vector.DrawFilledRect(screen, barX+1, barY+1, fillW, barH-2, fillColor, false)
	}

	// 5. Active Weapon Indicator (Below Health, X: 8, Y: 18)
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

		if icon != nil {
			iconOp := &ebiten.DrawImageOptions{}
			iconOp.GeoM.Translate(8, 17)
			screen.DrawImage(icon, iconOp)
		}
		DrawText(screen, wName, 18, 18, wCol)

		if hasMultipleWeapons {
			DrawText(screen, "Q:SWAP", 72, 18, color.RGBA{170, 180, 210, 220})
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
