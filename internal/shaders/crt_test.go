package shaders

import (
	"testing"
)

func TestNewCRTShader(t *testing.T) {
	s, err := NewCRTShader()
	if err != nil {
		t.Fatalf("Failed to compile CRT shader: %v", err)
	}
	if s == nil {
		t.Fatal("Expected non-nil shader instance")
	}
}
