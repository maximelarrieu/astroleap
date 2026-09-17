package game

import (
	"astroleap/internal/state"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	VirtualWidth  = 320
	VirtualHeight = 180
	TargetFPS     = 60
)

// Game implements the ebiten.Game interface.
type Game struct {
	machine *state.Machine
}

func NewGame() *Game {
	m := &state.Machine{}
	m.Change(state.NewTitleState(m))
	return &Game{
		machine: m,
	}
}

func (g *Game) Update() error {
	dt := 1.0 / float64(TargetFPS)
	g.machine.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.machine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return VirtualWidth, VirtualHeight
}
