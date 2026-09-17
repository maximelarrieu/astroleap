package ui

import (
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	text := "YEAR 2084: RESEARCH EXPEDITION LEAP-1 CRIPPLED OVER THE MOON"
	lines := WrapText(text, 120) // 120 / 6 = 20 chars max per line

	if len(lines) == 0 {
		t.Fatalf("expected non-empty wrapped lines")
	}

	for i, l := range lines {
		if len(l) > 20 {
			t.Errorf("line %d exceeds 20 characters: '%s' (len %d)", i, l, len(l))
		}
	}

	reconstructed := strings.Join(lines, " ")
	if reconstructed != text {
		t.Errorf("reconstructed text does not match original: got '%s'", reconstructed)
	}
}
