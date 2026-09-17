package entity

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Particle struct {
	X       float64
	Y       float64
	VX      float64
	VY      float64
	Life    float64
	MaxLife float64
	Color   color.RGBA
	Size    float32
}

type ParticleSystem struct {
	particles []*Particle
	rng       *rand.Rand
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make([]*Particle, 0, 256),
		rng:       rand.New(rand.NewSource(1234)),
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	alive := ps.particles[:0]
	for _, p := range ps.particles {
		p.Life -= dt
		if p.Life > 0 {
			p.X += p.VX * dt * 60.0
			p.Y += p.VY * dt * 60.0
			alive = append(alive, p)
		}
	}
	ps.particles = alive
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, camX float64) {
	for _, p := range ps.particles {
		screenX := float32(p.X - camX)
		screenY := float32(p.Y)
		if screenX < -10 || screenX > 330 {
			continue
		}

		alpha := uint8(255.0 * (p.Life / p.MaxLife))
		col := p.Color
		col.A = alpha

		sz := p.Size * float32(p.Life/p.MaxLife)
		if sz < 1 {
			sz = 1
		}
		vector.DrawFilledRect(screen, screenX-sz*0.5, screenY-sz*0.5, sz, sz, col, false)
	}
}

// SpawnJetpackSparks spawns exhaust particles from the astronaut's backpack.
func (ps *ParticleSystem) SpawnJetpackSparks(x, y float64, facingRight bool) {
	for i := 0; i < 3; i++ {
		angleOffset := (ps.rng.Float64() - 0.5) * 0.4
		speed := 1.2 + ps.rng.Float64()*1.8
		p := &Particle{
			X:       x,
			Y:       y,
			VX:      angleOffset,
			VY:      speed, // Shoot downward
			Life:    0.25 + ps.rng.Float64()*0.15,
			MaxLife: 0.4,
			Color:   color.RGBA{255, uint8(120 + ps.rng.Intn(130)), 30, 255},
			Size:    2.0,
		}
		ps.particles = append(ps.particles, p)
	}
}

// SpawnDustKick spawns soft gray lunar dust puffs when hitting the ground.
func (ps *ParticleSystem) SpawnDustKick(x, y float64) {
	for i := 0; i < 5; i++ {
		p := &Particle{
			X:       x,
			Y:       y,
			VX:      (ps.rng.Float64() - 0.5) * 1.5,
			VY:      -0.4 - ps.rng.Float64()*0.6,
			Life:    0.3 + ps.rng.Float64()*0.2,
			MaxLife: 0.5,
			Color:   color.RGBA{180, 190, 210, 200},
			Size:    2.0,
		}
		ps.particles = append(ps.particles, p)
	}
}

// SpawnCrystalSparkles spawns cyan twinkling stars upon collecting a crystal.
func (ps *ParticleSystem) SpawnCrystalSparkles(x, y float64) {
	for i := 0; i < 8; i++ {
		p := &Particle{
			X:       x,
			Y:       y,
			VX:      (ps.rng.Float64() - 0.5) * 2.5,
			VY:      (ps.rng.Float64() - 0.5) * 2.5,
			Life:    0.4 + ps.rng.Float64()*0.3,
			MaxLife: 0.7,
			Color:   color.RGBA{80, 240, 255, 255},
			Size:    2.5,
		}
		ps.particles = append(ps.particles, p)
	}
}

// SpawnStompBurst spawns alien pop particles.
func (ps *ParticleSystem) SpawnStompBurst(x, y float64) {
	for i := 0; i < 10; i++ {
		p := &Particle{
			X:       x,
			Y:       y,
			VX:      (ps.rng.Float64() - 0.5) * 3.0,
			VY:      (ps.rng.Float64() - 0.5) * 3.0,
			Life:    0.35 + ps.rng.Float64()*0.2,
			MaxLife: 0.55,
			Color:   color.RGBA{210, 80, 255, 255},
			Size:    2.5,
		}
		ps.particles = append(ps.particles, p)
	}
}
