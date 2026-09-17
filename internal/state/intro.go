package state

import (
	"math"

	"astroleap/internal/assets"
	"astroleap/internal/audio"
	"astroleap/internal/entity"
	"astroleap/internal/input"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type introActData struct {
	title string
	text  string
}

var introActs = []introActData{
	{
		title: "ACT 1: DEEP SPACE DISASTER",
		text:  "YEAR 2084: RESEARCH EXPEDITION LEAP-1. A SUDDEN GRAVITATIONAL ANOMALY CRIPPLES YOUR VESSEL OVER THE MOON. YOUR LANDER CRASH-LANDS ON THE DESOLATE REGOLITH...",
	},
	{
		title: "ACT 2: HOSTILE LIFEFORMS",
		text:  "ALIEN HORDE HAS OVERRUN THE SOLAR OUTPOSTS, HOARDING THE HYPER-CRYSTALS REQUIRED TO REFUEL YOUR ESCAPE SHUTTLES ACROSS THE SECTORS...",
	},
	{
		title: "ACT 3: PROTOTYPE ARSENAL",
		text:  "RECOVER LOST HIGH-TECH WEAPON CAPSULES: THE HIGH-VELOCITY LASER BLASTER (PRESS J) AND THE BOUNCING PLASMA NOVA CANNON (PRESS Q TO SWAP)...",
	},
	{
		title: "ACT 4: THE ODYSSEY BEGINS",
		text:  "TRAVERSE 3 HOSTILE SECTORS: MOON, PHOBOS, AND EUROPA. RECLAIM THE CRYSTALS, PURGE THE INVASION, AND RETURN SAFELY HOME TO EARTH!",
	},
}

// IntroCinematicState presents a 4-act retro cinematic prologue using game sprites.
type IntroCinematicState struct {
	machine        *Machine
	act            int
	actTimer       float64
	typewriterIdx  int
	typewriterTime float64
	autoAdvanceWait float64
	starfield      *ebiten.Image
	particles      *entity.ParticleSystem

	// Dynamic sprite staging
	landerX        float64
	landerY        float64
	astroX         float64
	astroY         float64
	blobX          float64
	blobY          float64
	blobVY         float64
	hoverX         float64
	hoverY         float64
	crystalRot     float64
	demoLaserX     float64
	demoNovaX      float64
	demoNovaY      float64
	demoNovaVX     float64
	demoNovaVY     float64
}

func NewIntroCinematicState(m *Machine) *IntroCinematicState {
	s := &IntroCinematicState{
		machine:   m,
		starfield: assets.GenerateStarfield(320, 180),
		particles: entity.NewParticleSystem(),
	}
	s.resetAct(0)
	return s
}

func (s *IntroCinematicState) resetAct(act int) {
	s.act = act
	s.actTimer = 0
	s.typewriterIdx = 0
	s.typewriterTime = 0
	s.autoAdvanceWait = 0

	switch act {
	case 0:
		s.landerX = 30
		s.landerY = 20
	case 1:
		s.blobX = 70
		s.blobY = 88
		s.blobVY = 0
		s.hoverX = 170
		s.hoverY = 46
	case 2:
		s.demoLaserX = 85
		s.demoNovaX = 180
		s.demoNovaY = 70
		s.demoNovaVX = 75
		s.demoNovaVY = 40
	case 3:
		s.astroX = 40
		s.astroY = 86
	}
}

func (s *IntroCinematicState) Enter() {
	audio.Get().StartBGM()
}

func (s *IntroCinematicState) Update(dt float64) {
	s.actTimer += dt
	s.crystalRot += dt * 7.0
	s.particles.Update(dt)

	inp := input.PollInput()

	// Instant skip to gameplay
	if inp.Pause {
		s.machine.Change(NewPlayState(s.machine))
		return
	}

	curData := introActs[s.act]

	// Typewriter text progression
	if s.typewriterIdx < len(curData.text) {
		s.typewriterTime += dt
		if s.typewriterTime >= 0.022 {
			s.typewriterTime = 0
			s.typewriterIdx++
			if s.typewriterIdx%3 == 0 {
				audio.Get().PlayTypewriter()
			}
		}
	} else {
		// Text complete, count auto-advance timer
		s.autoAdvanceWait += dt
		if s.autoAdvanceWait >= 4.0 {
			s.advanceAct()
			return
		}
	}

	// Player interaction
	if inp.JumpPressed || inp.Shoot || inp.Restart {
		if s.typewriterIdx < len(curData.text) {
			// Fast-forward text
			s.typewriterIdx = len(curData.text)
		} else {
			s.advanceAct()
			return
		}
	}

	// Update per-act diorama physics
	switch s.act {
	case 0:
		// Act 1: Spiraling crash lander
		if s.landerX < 145 {
			s.landerX += 28 * dt
			s.landerY += 15 * dt
			s.particles.SpawnJetpackSparks(s.landerX+16, s.landerY+12, false)
		}
	case 1:
		// Act 2: Bouncing alien blob & hovering drone
		s.blobVY += 380 * dt
		s.blobY += s.blobVY * dt
		if s.blobY >= 88 {
			s.blobY = 88
			s.blobVY = -120
		}
		s.hoverX = 180 + math.Sin(s.actTimer*2.2)*35.0
		s.hoverY = 46 + math.Cos(s.actTimer*3.0)*6.0
	case 2:
		// Act 3: Laser shooting demo and Nova cannon bounce
		s.demoLaserX += 160 * dt
		if s.demoLaserX > 165 {
			s.demoLaserX = 75
			s.particles.SpawnCrystalSparkles(160, 52)
		}

		s.demoNovaX += s.demoNovaVX * dt
		s.demoNovaY += s.demoNovaVY * dt
		if s.demoNovaX > 250 {
			s.demoNovaVX = -math.Abs(s.demoNovaVX)
		} else if s.demoNovaX < 185 {
			s.demoNovaVX = math.Abs(s.demoNovaVX)
		}
		if s.demoNovaY > 88 {
			s.demoNovaY = 88
			s.demoNovaVY = -110
			s.particles.SpawnStompBurst(s.demoNovaX+6, s.demoNovaY+6)
		} else {
			s.demoNovaVY += 220 * dt
		}
	case 3:
		// Act 4: Astronaut hero advance and thruster pulse
		if s.astroX < 140 {
			s.astroX += 30 * dt
		}
		s.astroY = 84 + math.Sin(s.actTimer*3.0)*3.0
		if int(s.actTimer*30)%4 == 0 {
			s.particles.SpawnJetpackSparks(s.astroX+6, s.astroY+14, true)
		}
	}
}

func (s *IntroCinematicState) advanceAct() {
	if s.act < len(introActs)-1 {
		audio.Get().PlayCrystal()
		s.resetAct(s.act + 1)
	} else {
		audio.Get().PlayWin()
		s.machine.Change(NewPlayState(s.machine))
	}
}

func (s *IntroCinematicState) Draw(screen *ebiten.Image) {
	atlas := assets.Get()

	// 1. Starfield background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.starfield, op)

	// 2. Act Diorama Visuals
	switch s.act {
	case 0:
		// Distant Earth
		earthOp := &ebiten.DrawImageOptions{}
		earthOp.GeoM.Translate(225, 18)
		screen.DrawImage(atlas.Earth, earthOp)

		// Lunar Surface silhouette
		vector.DrawFilledRect(screen, 0, 102, 320, 20, color.RGBA{45, 48, 62, 255}, false)
		vector.DrawFilledRect(screen, 0, 104, 320, 18, color.RGBA{30, 32, 42, 255}, false)

		// Particles (smoke / sparks from damaged lander)
		s.particles.Draw(screen, 0)

		// Descending damaged lander
		landerOp := &ebiten.DrawImageOptions{}
		landerOp.GeoM.Translate(s.landerX, s.landerY)
		screen.DrawImage(atlas.Lander, landerOp)

	case 1:
		// Mars celestial
		marsOp := &ebiten.DrawImageOptions{}
		marsOp.GeoM.Translate(225, 16)
		screen.DrawImage(atlas.Mars, marsOp)

		// Crimson Phobos ridge
		vector.DrawFilledRect(screen, 0, 102, 320, 20, color.RGBA{85, 32, 38, 255}, false)
		vector.DrawFilledRect(screen, 0, 104, 320, 18, color.RGBA{52, 20, 24, 255}, false)

		// Bouncing Alien Blob
		blobOp := &ebiten.DrawImageOptions{}
		blobOp.GeoM.Translate(s.blobX, s.blobY)
		if int(s.actTimer*8)%2 == 0 {
			screen.DrawImage(atlas.Blob1, blobOp)
		} else {
			screen.DrawImage(atlas.Blob2, blobOp)
		}

		// Hover drone
		hoverOp := &ebiten.DrawImageOptions{}
		hoverOp.GeoM.Translate(s.hoverX, s.hoverY)
		if int(s.actTimer*8)%2 == 0 {
			screen.DrawImage(atlas.Hover1, hoverOp)
		} else {
			screen.DrawImage(atlas.Hover2, hoverOp)
		}

		// Floating Hyper-Crystal
		crystFrame := int(s.crystalRot) % 4
		if crystFrame < 0 {
			crystFrame = 0
		}
		crystOp := &ebiten.DrawImageOptions{}
		crystOp.GeoM.Translate(135, 68+math.Sin(s.actTimer*2.5)*4.0)
		screen.DrawImage(atlas.Crystals[crystFrame], crystOp)

	case 2:
		// Weapon Technology Laboratory Stage
		vector.DrawFilledRect(screen, 0, 102, 320, 20, color.RGBA{22, 32, 48, 255}, false)
		vector.DrawFilledRect(screen, 0, 104, 320, 18, color.RGBA{14, 22, 34, 255}, false)

		// Laser Capsule
		podLaserOp := &ebiten.DrawImageOptions{}
		podLaserOp.GeoM.Translate(45, 52+math.Sin(s.actTimer*2.0)*3.0)
		screen.DrawImage(atlas.CapsuleLaser, podLaserOp)
		ui.DrawText(screen, "LASER BLASTER", 25, 78, color.RGBA{80, 230, 255, 255})

		// Laser Bolt demonstration
		boltOp := &ebiten.DrawImageOptions{}
		boltOp.GeoM.Translate(s.demoLaserX, 56)
		screen.DrawImage(atlas.LaserBolt, boltOp)

		// Nova Capsule
		podNovaOp := &ebiten.DrawImageOptions{}
		podNovaOp.GeoM.Translate(220, 52+math.Cos(s.actTimer*2.0)*3.0)
		screen.DrawImage(atlas.CapsuleNova, podNovaOp)
		ui.DrawText(screen, "NOVA CANNON", 204, 78, color.RGBA{255, 215, 60, 255})

		// Bouncing Nova Ball demonstration
		novaOp := &ebiten.DrawImageOptions{}
		novaOp.GeoM.Translate(s.demoNovaX, s.demoNovaY)
		screen.DrawImage(atlas.NovaBall, novaOp)

		s.particles.Draw(screen, 0)

	case 3:
		// Ringed Jupiter celestial
		jupOp := &ebiten.DrawImageOptions{}
		jupOp.GeoM.Translate(210, 14)
		screen.DrawImage(atlas.Jupiter, jupOp)

		// Glacial Europa Ice shelf
		vector.DrawFilledRect(screen, 0, 102, 320, 20, color.RGBA{28, 62, 92, 255}, false)
		vector.DrawFilledRect(screen, 0, 104, 320, 18, color.RGBA{16, 38, 58, 255}, false)

		// Rescue lander waiting at extraction point
		landerOp := &ebiten.DrawImageOptions{}
		landerOp.GeoM.Translate(230, 72)
		screen.DrawImage(atlas.Lander, landerOp)

		s.particles.Draw(screen, 0)

		// Heroic astronaut stepping forward
		astroOp := &ebiten.DrawImageOptions{}
		astroOp.GeoM.Translate(s.astroX, s.astroY)
		screen.DrawImage(atlas.AstronautThrust, astroOp)
	}

	// 3. Narrative Dialogue Box (Bottom modal)
	curData := introActs[s.act]
	displayedText := curData.text[:s.typewriterIdx]

	prompt := "[SPACE] CONTINUE  |  [ESC] SKIP"
	if s.act == len(introActs)-1 && s.typewriterIdx >= len(curData.text) {
		prompt = "[PRESS SPACE TO LAUNCH MISSION!]"
	}

	ui.DrawDialogueBox(screen, 16, 114, 288, 60, curData.title, displayedText, prompt)
}

func (s *IntroCinematicState) Exit() {}
