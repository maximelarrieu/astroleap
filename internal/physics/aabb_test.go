package physics

import "testing"

func TestRectOverlaps(t *testing.T) {
	r1 := Rect{X: 10, Y: 10, W: 20, H: 20}
	r2 := Rect{X: 15, Y: 15, W: 20, H: 20}
	r3 := Rect{X: 40, Y: 40, W: 10, H: 10}

	if !r1.Overlaps(r2) {
		t.Errorf("expected r1 and r2 to overlap")
	}
	if r1.Overlaps(r3) {
		t.Errorf("expected r1 and r3 not to overlap")
	}
}

func TestRectContainsPoint(t *testing.T) {
	r := Rect{X: 10, Y: 10, W: 20, H: 20}

	if !r.ContainsPoint(15, 15) {
		t.Errorf("expected (15, 15) to be inside rect")
	}
	if r.ContainsPoint(5, 5) {
		t.Errorf("expected (5, 5) to be outside rect")
	}
}

func TestRectIsEmpty(t *testing.T) {
	if !(Rect{}).IsEmpty() {
		t.Errorf("expected zero rect to be empty")
	}
	if (Rect{W: 10, H: 10}).IsEmpty() {
		t.Errorf("expected 10x10 rect not to be empty")
	}
}
