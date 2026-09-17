package state

import "github.com/hajimehoshi/ebiten/v2"

// State represents a game scene lifecycle.
type State interface {
	Enter()
	Update(dt float64)
	Draw(screen *ebiten.Image)
	Exit()
}
