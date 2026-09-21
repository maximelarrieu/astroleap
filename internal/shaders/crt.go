package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// CRTShaderSrc is an Ebitengine Kage pixel shader that simulates subtle CRT scanlines
// and lens vignette curvature for a rich 16-bit arcade CRT aesthetic.
const CRTShaderSrc = `//kage:unit pixels
package main

var VignetteIntensity float
var ScanlineIntensity float

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	clr := imageSrc0At(srcPos)
	if clr.a == 0.0 {
		return clr
	}
	
	// Scanline modulation across alternating horizontal scanlines
	scan := 1.0 - ScanlineIntensity * (0.5 + 0.5 * sin(srcPos.y * 3.14159265))
	
	// Lens vignette darkening towards corners
	size := imageSrcTextureSize()
	uv := srcPos / size
	d := distance(uv, vec2(0.5, 0.5))
	vig := clamp(1.0 - (d * d * VignetteIntensity * 1.5), 0.0, 1.0)
	
	return vec4(clr.rgb * scan * vig, clr.a)
}
`

// NewCRTShader compiles the CRT and vignette post-processing shader.
func NewCRTShader() (*ebiten.Shader, error) {
	return ebiten.NewShader([]byte(CRTShaderSrc))
}
