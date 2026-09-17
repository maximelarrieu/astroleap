package entity

import (
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/level"
	"astroleap/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyType string

const (
	TypeBlob  EnemyType = "blob"
	TypeHover EnemyType = "hover"
)

type Enemy struct {
	X            float64
	Y            float64
	VX           float64
	VY           float64
	OriginX      float64
	Type         EnemyType
	Dead         bool
	Squished     bool
	SquishTimer  float64
	AnimTimer    float64
	FacingRight  bool
	HoverRange   float64
	HoverBaseY   float64
}

func NewEnemy(tileX, tileY int, kind string) *Enemy {
	x := float64(tileX * level.TileSize)
	y := float64(tileY * level.TileSize)
	e := &Enemy{
		X:          x,
		Y:          y,
		OriginX:    x,
		HoverBaseY: y,
		Type:       TypeBlob,
		VX:         -0.6,
		HoverRange: 40.0,
	}
	if kind == "hover" {
		e.Type = TypeHover
		e.VX = 0.5
	}
	return e
}

func (e *Enemy) Update(dt float64, lvl *level.Level) {
	if e.Dead {
		return
	}

	if e.Squished {
		e.SquishTimer -= dt
		if e.SquishTimer <= 0 {
			e.Dead = true
		}
		return
	}

	e.AnimTimer += dt

	if e.Type == TypeBlob {
		// 1. Move horizontally
		newX := e.X + e.VX*dt*60.0

		// Check wall collision
		checkTileX := int(newX / level.TileSize)
		if e.VX > 0 {
			checkTileX = int((newX + 15) / level.TileSize)
		}
		tileY := int((e.Y + 8) / level.TileSize)

		// Check ledge below ahead
		belowTileY := int((e.Y + 17) / level.TileSize)

		if lvl.IsSolid(checkTileX, tileY) || !lvl.IsSolid(checkTileX, belowTileY) {
			// Reverse direction
			e.VX = -e.VX
		} else {
			e.X = newX
		}

		e.FacingRight = (e.VX > 0)

	} else if e.Type == TypeHover {
		// Air hover patrol
		e.X += e.VX * dt * 60.0
		if math.Abs(e.X-e.OriginX) > e.HoverRange {
			e.VX = -e.VX
		}
		// Sine wave vertical bobbing
		e.Y = e.HoverBaseY + math.Sin(e.AnimTimer*3.0)*6.0
		e.FacingRight = (e.VX > 0)
	}
}

func (e *Enemy) GetRect() physics.Rect {
	if e.Squished || e.Dead {
		return physics.Rect{}
	}
	return physics.Rect{
		X: e.X + 2,
		Y: e.Y + 3,
		W: 12,
		H: 12,
	}
}

// Stomp triggers the squished state when stomped by player.
func (e *Enemy) Stomp() {
	e.Squished = true
	e.SquishTimer = 0.35
}

func (e *Enemy) Draw(screen *ebiten.Image, camX float64) {
	if e.Dead {
		return
	}
	screenX := e.X - camX
	if screenX < -20 || screenX > 330 {
		return
	}

	atlas := assets.Get()
	var img *ebiten.Image

	if e.Type == TypeBlob {
		if e.Squished {
			img = atlas.BlobSquished
		} else {
			// 4 FPS squash-stretch
			frame := int(e.AnimTimer*4.0) % 2
			if frame == 0 {
				img = atlas.Blob1
			} else {
				img = atlas.Blob2
			}
		}
	} else if e.Type == TypeHover {
		if e.Squished {
			img = atlas.BlobSquished
		} else {
			frame := int(e.AnimTimer*5.0) % 2
			if frame == 0 {
				img = atlas.Hover1
			} else {
				img = atlas.Hover2
			}
		}
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		if e.FacingRight {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(16, 0)
		}
		op.GeoM.Translate(screenX, e.Y)
		screen.DrawImage(img, op)
	}
}
