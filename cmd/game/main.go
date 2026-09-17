package main

import (
	"log"

	"astroleap/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("AstroLeap: Lunar Odyssey")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game exited with error: %v", err)
	}
}
