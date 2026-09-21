package game

import (
	"astroleap/internal/shaders"
	"astroleap/internal/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	VirtualWidth  = 320
	VirtualHeight = 180
	TargetFPS     = 60
)

// Game implements the ebiten.Game interface.
type Game struct {
	machine    *state.Machine
	crtShader  *ebiten.Shader
	crtEnabled bool
	offscreen  *ebiten.Image
}

func NewGame() *Game {
	m := &state.Machine{}
	m.Change(state.NewTitleState(m))
	shader, _ := shaders.NewCRTShader()
	return &Game{
		machine:    m,
		crtShader:  shader,
		crtEnabled: true,
		offscreen:  ebiten.NewImage(VirtualWidth, VirtualHeight),
	}
}

func (g *Game) Update() error {
	// Toggle CRT scanline post-processing with F2
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.crtEnabled = !g.crtEnabled
	}

	// Toggle Fullscreen with F11 or Alt+Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	dt := 1.0 / float64(TargetFPS)
	g.machine.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.crtEnabled && g.crtShader != nil {
		g.offscreen.Clear()
		g.machine.Draw(g.offscreen)

		opts := &ebiten.DrawRectShaderOptions{}
		opts.Images[0] = g.offscreen
		opts.Uniforms = map[string]any{
			"VignetteIntensity": float32(0.6),
			"ScanlineIntensity": float32(0.12),
		}
		screen.DrawRectShader(VirtualWidth, VirtualHeight, g.crtShader, opts)
	} else {
		g.machine.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return VirtualWidth, VirtualHeight
}
