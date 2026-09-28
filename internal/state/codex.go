package state

import (
	"fmt"
	"image/color"
	"strconv"

	"astroleap/internal/assets"
	"astroleap/internal/audio"
	"astroleap/internal/input"
	"astroleap/internal/narrative"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// CodexState presents an archive viewer for discovered narrative data logs.
type CodexState struct {
	machine      *Machine
	starfield    *ebiten.Image
	logs         []narrative.LogEntry
	selectedIdx  int
	timer        float64
}

// NewCodexState creates a new archive browser.
func NewCodexState(m *Machine) *CodexState {
	return &CodexState{
		machine:   m,
		starfield: assets.GenerateStarfield(320, 180),
		logs:      narrative.GetAllLogs(),
	}
}

func (s *CodexState) Enter() {}

func (s *CodexState) Update(dt float64) {
	s.timer += dt
	inp := input.PollInput()

	// Navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if s.selectedIdx > 0 {
			s.selectedIdx--
			audio.Get().PlayJump()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if s.selectedIdx < len(s.logs)-1 {
			s.selectedIdx++
			audio.Get().PlayJump()
		}
	}

	// Exit
	if inp.JumpPressed || inp.Pause || inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		audio.Get().PlayJump()
		s.machine.Change(NewTitleState(s.machine))
	}
}

func (s *CodexState) Draw(screen *ebiten.Image) {
	// 1. Starfield background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.starfield, op)

	// CRT scanline ambient tint
	vector.DrawFilledRect(screen, 0, 0, 320, 180, color.RGBA{10, 16, 28, 220}, false)

	// 2. Top Header
	vector.DrawFilledRect(screen, 12, 8, 296, 20, color.RGBA{18, 24, 40, 240}, false)
	vector.StrokeRect(screen, 12, 8, 296, 20, 1, color.RGBA{60, 210, 255, 255}, false)
	ui.DrawText(screen, "EXPEDITION ARCHIVES - DATA LOGS", 22, 14, color.RGBA{255, 215, 60, 255})

	unlockedCount := narrative.GetUnlockedCount()
	countStr := fmt.Sprintf("%d/%d RECOVERED", unlockedCount, len(s.logs))
	ui.DrawText(screen, countStr, 224, 14, color.RGBA{100, 240, 255, 255})

	// 3. Left Panel: Log Selector (X: 12, W: 86)
	vector.DrawFilledRect(screen, 12, 32, 86, 116, color.RGBA{14, 20, 34, 220}, false)
	vector.StrokeRect(screen, 12, 32, 86, 116, 1, color.RGBA{45, 65, 100, 200}, false)

	for i, l := range s.logs {
		itemY := 36 + i*22
		isSel := i == s.selectedIdx
		isUnlocked := narrative.IsLogUnlocked(l.ID)

		if isSel {
			vector.DrawFilledRect(screen, 14, float32(itemY-2), 82, 18, color.RGBA{35, 60, 100, 240}, false)
			vector.StrokeRect(screen, 14, float32(itemY-2), 82, 18, 1, color.RGBA{80, 230, 255, 255}, false)
		}

		icon := "SEC." + strconv.Itoa(l.Sector)
		colText := color.RGBA{140, 160, 190, 255}
		if isUnlocked {
			colText = color.RGBA{220, 235, 255, 255}
			if isSel {
				colText = color.RGBA{255, 230, 80, 255}
			}
		} else {
			colText = color.RGBA{80, 95, 120, 200}
		}

		ui.DrawText(screen, l.ID, 18, itemY, colText)
		ui.DrawText(screen, icon, 54, itemY+8, color.RGBA{100, 180, 220, 200})
	}

	// 4. Right Panel: Content Reader (X: 102, W: 206)
	vector.DrawFilledRect(screen, 102, 32, 206, 116, color.RGBA{12, 18, 30, 230}, false)
	vector.StrokeRect(screen, 102, 32, 206, 116, 1, color.RGBA{60, 210, 255, 255}, false)

	curLog := s.logs[s.selectedIdx]
	isCurUnlocked := narrative.IsLogUnlocked(curLog.ID)

	if isCurUnlocked {
		// Title
		ui.DrawText(screen, curLog.Title, 110, 38, color.RGBA{255, 215, 60, 255})
		// Meta
		metaStr := curLog.Author + " | " + curLog.Date
		ui.DrawText(screen, metaStr, 110, 48, color.RGBA{100, 230, 255, 240})

		// Divider
		vector.DrawFilledRect(screen, 110, 56, 190, 1, color.RGBA{45, 75, 120, 200}, false)

		// Body lines
		for lineIdx, line := range curLog.Content {
			ui.DrawText(screen, line, 110, 62+lineIdx*11, color.RGBA{225, 235, 255, 255})
		}
	} else {
		// Encrypted Terminal Notice
		ui.DrawText(screen, "ENCRYPTED TRANSMISSION", 136, 62, color.RGBA{255, 80, 80, 255})
		ui.DrawText(screen, curLog.ID+" REMAINS LOCKED", 146, 76, color.RGBA{180, 195, 220, 220})
		hint := fmt.Sprintf("RECOVER DATA-TERMINAL IN SECTOR %d", curLog.Sector)
		ui.DrawText(screen, hint, 114, 94, color.RGBA{100, 220, 255, 240})
	}

	// 5. Footer Navigation Controls
	vector.DrawFilledRect(screen, 12, 152, 296, 20, color.RGBA{14, 20, 34, 230}, false)
	vector.StrokeRect(screen, 12, 152, 296, 20, 1, color.RGBA{45, 65, 100, 200}, false)
	ui.DrawText(screen, "UP/DOWN: SELECT LOG", 24, 158, color.RGBA{160, 185, 220, 240})
	ui.DrawText(screen, "ESC/SPACE: RETURN TO MENU", 154, 158, color.RGBA{255, 215, 60, 255})
}

func (s *CodexState) Exit() {}
