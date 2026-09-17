package state

import "github.com/hajimehoshi/ebiten/v2"

// Machine manages scene transitions cleanly.
type Machine struct {
	current State
}

func NewMachine(initial State) *Machine {
	m := &Machine{}
	m.Change(initial)
	return m
}

func (m *Machine) Change(next State) {
	if m.current != nil {
		m.current.Exit()
	}
	m.current = next
	if m.current != nil {
		m.current.Enter()
	}
}

func (m *Machine) Update(dt float64) {
	if m.current != nil {
		m.current.Update(dt)
	}
}

func (m *Machine) Draw(screen *ebiten.Image) {
	if m.current != nil {
		m.current.Draw(screen)
	}
}
