package state

import (
	"testing"
)

func TestIntroCinematicState(t *testing.T) {
	m := NewMachine(nil)
	intro := NewIntroCinematicState(m)

	if intro == nil {
		t.Fatalf("expected non-nil IntroCinematicState")
	}

	if intro.act != 0 {
		t.Errorf("expected starting act 0, got %d", intro.act)
	}

	// Verify all 4 acts have valid titles and texts
	if len(introActs) != 4 {
		t.Errorf("expected 4 acts, got %d", len(introActs))
	}

	for i, act := range introActs {
		if act.title == "" || act.text == "" {
			t.Errorf("act %d has empty title or text: %+v", i, act)
		}
	}

	// Advance through all acts
	for i := 0; i < len(introActs)-1; i++ {
		curAct := intro.act
		intro.advanceAct()
		if intro.act != curAct+1 {
			t.Errorf("expected act %d after advance, got %d", curAct+1, intro.act)
		}
	}

	// Advance past last act transitions machine to PlayState
	intro.advanceAct()
	if _, ok := m.current.(*PlayState); !ok {
		t.Errorf("expected machine to change to PlayState after completing intro, got %T", m.current)
	}
}
