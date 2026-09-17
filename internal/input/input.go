package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputState holds current abstract action states for the frame.
type InputState struct {
	MoveLeft     bool
	MoveRight    bool
	JumpPressed  bool
	JumpHeld     bool
	Shoot        bool
	SwitchWeapon bool
	Pause        bool
	Restart      bool
	ToggleMute   bool
}

// PollInput samples keyboard, gamepad, and touch inputs.
func PollInput() InputState {
	state := InputState{}

	// 1. Keyboard & Mouse
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		state.MoveLeft = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		state.MoveRight = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		state.JumpPressed = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		state.JumpHeld = true
	}

	// Shooting: J, Z, F, or Left Mouse Click
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyZ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyF) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		state.Shoot = true
	}

	// Switching weapons: Q, E, Tab
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyE) ||
		inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		state.SwitchWeapon = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		state.Pause = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		state.Restart = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		state.ToggleMute = true
	}

	// 2. Gamepad (Standard IDs)
	gamepadIDs := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepadIDs {
		if !ebiten.IsStandardGamepadLayoutAvailable(id) {
			continue
		}
		// D-pad or sticks
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftLeft) ||
			ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) < -0.3 {
			state.MoveLeft = true
		}
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftRight) ||
			ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) > 0.3 {
			state.MoveRight = true
		}

		// South button (A / Cross): Jump
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) {
			state.JumpPressed = true
		}
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonRightBottom) {
			state.JumpHeld = true
		}

		// West button (X / Square): Shoot
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightLeft) {
			state.Shoot = true
		}

		// North button (Y / Triangle) or R1 / Right Shoulder: Switch Weapon
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightTop) ||
			inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonFrontTopRight) {
			state.SwitchWeapon = true
		}

		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterRight) { // Start
			state.Pause = true
		}
	}

	// 3. Touch Screen (WASM Mobile Support)
	touchIDs := ebiten.AppendTouchIDs(nil)
	for _, tid := range touchIDs {
		tx, _ := ebiten.TouchPosition(tid)
		if tx < 100 {
			state.MoveLeft = true
		} else if tx < 200 {
			state.MoveRight = true
		} else {
			state.JumpHeld = true
		}
	}
	justPressedTouches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, tid := range justPressedTouches {
		tx, _ := ebiten.TouchPosition(tid)
		if tx >= 200 {
			state.JumpPressed = true
		}
	}

	// Fullscreen toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return state
}
