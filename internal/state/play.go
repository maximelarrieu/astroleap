package state

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"astroleap/internal/audio"
	"astroleap/internal/entity"
	"astroleap/internal/input"
	"astroleap/internal/level"
	"astroleap/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayState struct {
	machine   *Machine
	lvl       *level.Level
	player    *entity.Player
	enemies     []*entity.Enemy
	crystals    []*entity.Crystal
	pickups     []*entity.WeaponPickup
	projectiles []*entity.Projectile
	particles   *entity.ParticleSystem

	floatingTexts   []*entity.FloatingText
	movingPlatforms []*entity.MovingPlatform
	ridingPlatform  *entity.MovingPlatform
	geysers         []*entity.CryoGeyser
	medals          []*entity.GoldMedal
	medalsCollected int
	elapsedTime     float64

	securityKey      *entity.SecurityKey
	hasSecurityKey   bool
	repairCores      []*entity.RepairCore
	repairCoresCount int

	crumblingPlatforms []*entity.CrumblingPlatform
	healthPickups      []*entity.HealthPickup
	boostPickups       []*entity.ThrusterBoostPickup
	hitStopTimer       float64

	boss            *entity.Boss
	bossProjectiles []*entity.BossProjectile
	shakeTimer      float64
	shakeIntensity  float64

	currentSector     int
	score             int
	crystalsCollected int
	camX              float64
	paused            bool
	startGrace        float64
	unlockBannerText  string
	unlockBannerTimer float64
}

func NewPlayState(m *Machine) *PlayState {
	return NewPlayStateWithProgress(m, 1, 0, 0, nil)
}

func NewPlayStateWithProgress(m *Machine, sector int, score int, crystals int, carriedPlayer *entity.Player) *PlayState {
	return NewPlayStateWithRecords(m, sector, score, crystals, carriedPlayer, 0, 0)
}

func NewPlayStateWithRecords(m *Machine, sector int, score int, crystals int, carriedPlayer *entity.Player, elapsedTime float64, medals int) *PlayState {
	if sector < 1 {
		sector = 1
	}
	lvl := level.LoadLevel(sector)
	player := entity.NewPlayer(32.0, 128.0)

	// Carry over unlocked weapons
	if carriedPlayer != nil {
		player.HasLaser = carriedPlayer.HasLaser
		player.HasNova = carriedPlayer.HasNova
		player.ActiveWeapon = carriedPlayer.ActiveWeapon
	}

	ps := &PlayState{
		machine:           m,
		lvl:               lvl,
		player:            player,
		particles:         entity.NewParticleSystem(),
		currentSector:     sector,
		score:             score,
		crystalsCollected: crystals,
		medalsCollected:   medals,
		elapsedTime:       elapsedTime,
		startGrace:        0.2, // Ignore residual keypress from title screen
	}

	// Spawn Crystals from level template
	for _, c := range ps.lvl.Crystals {
		ps.crystals = append(ps.crystals, entity.NewCrystal(
			float64(c.TileX*level.TileSize),
			float64(c.TileY*level.TileSize),
		))
	}

	// Spawn Enemies from level template
	for _, e := range ps.lvl.EnemySpawns {
		ps.enemies = append(ps.enemies, entity.NewEnemy(e.TileX, e.TileY, e.Kind))
	}

	// Spawn Weapon Capsules from level template
	for _, ws := range ps.lvl.WeaponSpawns {
		var w entity.WeaponType
		if ws.Kind == "laser" {
			w = entity.WeaponLaser
		} else if ws.Kind == "nova" {
			w = entity.WeaponNova
		}
		ps.pickups = append(ps.pickups, entity.NewWeaponPickup(
			float64(ws.TileX*level.TileSize),
			float64(ws.TileY*level.TileSize),
			w,
		))
	}

	// Spawn Moving Platforms from level template
	for _, mp := range ps.lvl.MovingPlatforms {
		ps.movingPlatforms = append(ps.movingPlatforms, entity.NewMovingPlatform(
			mp.StartX, mp.StartY, mp.EndX, mp.EndY, mp.W, mp.H, mp.Speed,
		))
	}

	// Spawn Cryo Geysers from level template
	for _, g := range ps.lvl.Geysers {
		ps.geysers = append(ps.geysers, entity.NewCryoGeyser(
			float64(g.TileX*level.TileSize),
			float64(g.TileY*level.TileSize),
			g.Height,
		))
	}

	// Spawn Secret Gold Medals from level template
	for _, md := range ps.lvl.Medals {
		ps.medals = append(ps.medals, entity.NewGoldMedal(
			float64(md.TileX*level.TileSize),
			float64(md.TileY*level.TileSize),
		))
	}

	// Spawn Fragile Crumbling Platforms from level template
	for _, cp := range ps.lvl.CrumblingPlatforms {
		ps.crumblingPlatforms = append(ps.crumblingPlatforms, entity.NewCrumblingPlatform(
			float64(cp.TileX*level.TileSize),
			float64(cp.TileY*level.TileSize),
			cp.W,
			cp.H,
		))
	}

	// Spawn Oxygen Health Canisters from level template
	for _, hp := range ps.lvl.HealthPickups {
		ps.healthPickups = append(ps.healthPickups, entity.NewHealthPickup(
			float64(hp.TileX*level.TileSize),
			float64(hp.TileY*level.TileSize),
		))
	}

	// Spawn Ion Thruster Boost Cores from level template
	for _, bp := range ps.lvl.BoostPickups {
		ps.boostPickups = append(ps.boostPickups, entity.NewThrusterBoostPickup(
			float64(bp.TileX*level.TileSize),
			float64(bp.TileY*level.TileSize),
		))
	}

	// Spawn Security Keycard if level features one (Sector 4)
	if ps.lvl.SecurityKey != nil {
		ps.securityKey = entity.NewSecurityKey(
			float64(ps.lvl.SecurityKey.TileX*level.TileSize),
			float64(ps.lvl.SecurityKey.TileY*level.TileSize),
		)
	}

	// Spawn Reactor Repair Cores if level features them (Sector 5)
	for _, rc := range ps.lvl.RepairCores {
		ps.repairCores = append(ps.repairCores, entity.NewRepairCore(
			rc.Index,
			rc.Name,
			float64(rc.TileX*level.TileSize),
			float64(rc.TileY*level.TileSize),
		))
	}

	// Spawn Boss Overlord Mech if level features a boss encounter
	if ps.lvl.HasBoss {
		ps.boss = entity.NewBoss(ps.lvl.BossX, ps.lvl.BossY)
	}

	return ps
}

func (s *PlayState) isPlayerBlocked(x, y float64) bool {
	boxX := x + 2.0
	boxY := y + 1.0
	boxW := entity.PlayerWidth
	boxH := entity.PlayerHeight

	minTX := int(boxX / float64(level.TileSize))
	maxTX := int((boxX + boxW - 0.001) / float64(level.TileSize))
	minTY := int(boxY / float64(level.TileSize))
	maxTY := int((boxY + boxH - 0.001) / float64(level.TileSize))

	for ty := minTY; ty <= maxTY; ty++ {
		for tx := minTX; tx <= maxTX; tx++ {
			if s.lvl.IsSolid(tx, ty) {
				return true
			}
		}
	}
	return false
}

func (s *PlayState) SpawnFloatingText(text string, x, y float64, col color.RGBA) {
	s.floatingTexts = append(s.floatingTexts, entity.NewFloatingText(text, x, y, col))
}

func (s *PlayState) TriggerShake(intensity, duration float64) {
	if intensity > s.shakeIntensity || s.shakeTimer <= 0 {
		s.shakeIntensity = intensity
		s.shakeTimer = duration
	}
}

func (s *PlayState) Enter() {
	audio.Get().StartBGM()
}

func (s *PlayState) Update(dt float64) {
	inp := input.PollInput()

	if s.startGrace > 0 {
		s.startGrace -= dt
		inp.JumpPressed = false
		inp.JumpHeld = false
	}

	if s.shakeTimer > 0 {
		s.shakeTimer -= dt
	}

	if inp.ToggleMute {
		audio.Get().ToggleMute()
	}

	if inp.Pause {
		s.paused = !s.paused
	}

	if s.paused {
		return
	}

	// Hit-Stop micro-pause (freezes gameplay frames on impact for mechanical punch)
	if s.hitStopTimer > 0 {
		s.hitStopTimer -= dt
		if s.shakeTimer > 0 {
			s.shakeTimer -= dt
		}
		return
	}

	// 0. Update Moving Platforms FIRST and carry riding player
	for _, mp := range s.movingPlatforms {
		mp.Update(dt)
	}

	if s.ridingPlatform != nil {
		newPX := s.player.X + s.ridingPlatform.DX
		if !s.isPlayerBlocked(newPX, s.player.Y) {
			s.player.X = newPX
		}
		newPY := s.player.Y + s.ridingPlatform.DY
		if !s.isPlayerBlocked(s.player.X, newPY) {
			s.player.Y = newPY
		}
	}

	// 1. Update Player & Timer
	s.elapsedTime += dt
	s.player.Update(dt, inp.MoveLeft, inp.MoveRight, inp.JumpPressed, inp.JumpHeld, s.lvl, s.particles)

	// Shooting input
	if inp.Shoot {
		if proj := s.player.Shoot(s.particles); proj != nil {
			s.projectiles = append(s.projectiles, proj)
			if s.player.ActiveWeapon == entity.WeaponNova {
				s.TriggerShake(2.5, 0.12)
			}
		}
	}

	// Switch weapon input
	if inp.SwitchWeapon {
		s.player.SwitchWeapon()
	}

	// Update weapon unlock announcement banner
	if s.unlockBannerTimer > 0 {
		s.unlockBannerTimer -= dt
	}

	// Check Player Death
	if s.player.Dead {
		audio.Get().PlayGameOver()
		s.machine.Change(NewGameOverState(s.machine, s.score, s.crystalsCollected, s.currentSector, s.player))
		return
	}

	// 2. Update Crystals & Check Pickups
	playerRect := s.player.GetRect()
	for _, c := range s.crystals {
		c.Update(dt)
		if !c.Collected && playerRect.Overlaps(c.GetRect()) {
			c.Collected = true
			s.score += 100
			s.crystalsCollected++
			audio.Get().PlayCrystal()
			s.particles.SpawnCrystalSparkles(c.X+8, c.Y+8)
			s.SpawnFloatingText("+100", c.X, c.Y-4, color.RGBA{100, 245, 255, 255})
		}
	}

	// 2b. Update Secret Gold Medals
	for _, m := range s.medals {
		m.Update(dt)
		if !m.Collected && playerRect.Overlaps(m.GetRect()) {
			m.Collected = true
			s.medalsCollected++
			s.score += 1000
			audio.Get().PlayPowerup()
			s.particles.SpawnCrystalSparkles(m.X+8, m.Y+8)
			s.SpawnFloatingText("MEDAL! +1000", m.X-8, m.Y-6, color.RGBA{255, 215, 50, 255})
			s.TriggerShake(2.5, 0.12)
		}
	}

	// 2c. Update Moving Platforms Top Landing & Snap
	s.ridingPlatform = nil
	playerRect = s.player.GetRect()
	playerFeet := s.player.Y + entity.PlayerHeight
	for _, mp := range s.movingPlatforms {
		platRect := mp.GetRect()
		if playerRect.X+playerRect.W > platRect.X && playerRect.X < platRect.X+platRect.W {
			if s.player.VY >= 0 && playerFeet >= platRect.Y-2.0 && playerFeet <= platRect.Y+8.0 {
				s.player.Y = platRect.Y - entity.PlayerHeight
				s.player.VY = 0
				s.player.OnGround = true
				s.player.StompCombo = 0
				s.ridingPlatform = mp
				break
			}
		}
	}

	// 2c-2. Update Fragile Crumbling Platforms
	for _, cp := range s.crumblingPlatforms {
		cp.Update(dt)
		cRect := cp.GetRect()
		if cRect.IsEmpty() {
			continue
		}
		playerFeet := s.player.Y + entity.PlayerHeight
		if playerRect.X+playerRect.W > cRect.X && playerRect.X < cRect.X+cRect.W {
			if playerFeet >= cRect.Y-1.0 && playerFeet <= cRect.Y+7.0 && s.player.VY >= 0 {
				cp.Touch()
				s.player.Y = cRect.Y - entity.PlayerHeight
				s.player.VY = 0
				s.player.OnGround = true
				s.player.StompCombo = 0
				s.ridingPlatform = nil
			}
		}
	}

	// 2d. Update Cryo Geysers & Vertical Propulsion
	for _, g := range s.geysers {
		g.Update(dt, s.particles)
		if playerRect.Overlaps(g.GetRect()) {
			s.player.VY = -3.4
			s.player.OnGround = false
			s.player.Fuel = math.Min(entity.MaxFuel, s.player.Fuel+dt*25.0)
			s.particles.SpawnCrystalSparkles(s.player.X+8, s.player.Y+12)
		}
	}

	// 3. Update Weapon Pickups
	for _, pk := range s.pickups {
		pk.Update(dt)
		if !pk.Collected && playerRect.Overlaps(pk.GetRect()) {
			pk.Collected = true
			s.player.UnlockWeapon(pk.Weapon)
			audio.Get().PlayPowerup()
			s.particles.SpawnCrystalSparkles(pk.X+8, pk.Y+8)
			s.unlockBannerText = pk.Weapon.String() + " UNLOCKED!"
			s.unlockBannerTimer = 2.8
			s.score += 500
			s.SpawnFloatingText("+500", pk.X, pk.Y-4, color.RGBA{255, 215, 60, 255})
		}
	}

	// 3b. Update Oxygen Health Canisters
	for _, hp := range s.healthPickups {
		hp.Update(dt)
		if !hp.Collected && playerRect.Overlaps(hp.GetRect()) {
			if s.player.Heal(1) {
				hp.Collected = true
				audio.Get().PlayPowerup()
				s.particles.SpawnCrystalSparkles(hp.X+8, hp.Y+8)
				s.score += 200
				s.SpawnFloatingText("+1 O2 HEALTH", hp.X-8, hp.Y-6, color.RGBA{60, 230, 255, 255})
			}
		}
	}

	// 3c. Update Ion Thruster Boost Cores
	for _, bp := range s.boostPickups {
		bp.Update(dt)
		if !bp.Collected && playerRect.Overlaps(bp.GetRect()) {
			bp.Collected = true
			s.player.ApplyBoost(5.0) // 5 seconds of infinite supercharged fuel
			audio.Get().PlayPowerup()
			s.particles.SpawnCrystalSparkles(bp.X+8, bp.Y+8)
			s.score += 300
			s.TriggerShake(2.0, 0.12)
			s.SpawnFloatingText("ION BOOST 5s!", bp.X-10, bp.Y-6, color.RGBA{255, 190, 40, 255})
		}
	}

	// 3d. Update Security Keycard (Sector 4)
	if s.securityKey != nil && !s.securityKey.Collected {
		s.securityKey.Update(dt)
		if playerRect.Overlaps(s.securityKey.GetRect()) {
			s.securityKey.Collected = true
			s.hasSecurityKey = true
			s.score += 500
			audio.Get().PlayPowerup()
			s.particles.SpawnCrystalSparkles(s.securityKey.X+8, s.securityKey.Y+8)
			s.TriggerShake(2.5, 0.14)
			s.unlockBannerText = "SECURITY KEY ACQUIRED! AIRLOCK UNLOCKED!"
			s.unlockBannerTimer = 3.0
			s.SpawnFloatingText("KEY ACQUIRED! +500", s.securityKey.X-16, s.securityKey.Y-6, color.RGBA{255, 220, 50, 255})
		}
	}

	// 3e. Update Reactor Repair Cores (Sector 5)
	for _, rc := range s.repairCores {
		if !rc.Collected {
			rc.Update(dt)
			if playerRect.Overlaps(rc.GetRect()) {
				rc.Collected = true
				s.repairCoresCount++
				s.score += 500
				audio.Get().PlayPowerup()
				s.particles.SpawnCrystalSparkles(rc.X+8, rc.Y+8)
				s.TriggerShake(3.0, 0.16)
				if s.repairCoresCount == 1 {
					s.unlockBannerText = "REPAIR 1/3: COOLANT SYSTEM ONLINE!"
					s.SpawnFloatingText("COOLANT RESTORED! +500", rc.X-20, rc.Y-6, color.RGBA{100, 240, 255, 255})
				} else if s.repairCoresCount == 2 {
					s.unlockBannerText = "REPAIR 2/3: PLASMA STABILIZER ONLINE!"
					s.SpawnFloatingText("STABILIZER ONLINE! +500", rc.X-20, rc.Y-6, color.RGBA{255, 180, 50, 255})
				} else {
					s.score += 500 // 1000 total for final core
					s.unlockBannerText = "SHIP REPAIRED! WARP DRIVE ENGAGED!"
					s.SpawnFloatingText("WARP DRIVE READY! +1000", rc.X-24, rc.Y-6, color.RGBA{80, 255, 120, 255})
				}
				s.unlockBannerTimer = 3.2
			}
		}
	}

	// 4. Update Projectiles & Check Hits against Enemies / Boss
	aliveProjs := s.projectiles[:0]
	for _, proj := range s.projectiles {
		proj.Update(dt, s.lvl, s.particles)
		if proj.Dead {
			continue
		}

		projRect := proj.GetRect()
		hit := false

		for _, e := range s.enemies {
			if e.Dead || e.Squished {
				continue
			}
			if projRect.Overlaps(e.GetRect()) {
				hit = true
				e.Stomp()
				s.score += 250
				s.particles.SpawnStompBurst(e.X+8, e.Y+8)
				s.TriggerShake(2.0, 0.10)
				audio.Get().PlayStomp()
				s.SpawnFloatingText("+250", e.X, e.Y-4, color.RGBA{100, 245, 255, 255})
				break
			}
		}

		if !hit && s.boss != nil && !s.boss.Dead {
			if projRect.Overlaps(s.boss.GetRect()) {
				hit = true
				dmg := 1
				if proj.Type == entity.ProjNova {
					dmg = 2
				}
				if s.boss.TakeDamage(dmg, s.particles) {
					s.TriggerShake(3.5, 0.14)
					s.hitStopTimer = 0.04
					if s.boss.Dead {
						s.score += 5000
						s.TriggerShake(7.0, 0.4)
						s.hitStopTimer = 0.10
						s.SpawnFloatingText("OVERLORD DOWN! +5000", s.boss.X-16, s.boss.Y-10, color.RGBA{255, 60, 80, 255})
					} else {
						s.SpawnFloatingText(fmt.Sprintf("HIT! -%d", dmg), s.boss.X+6, s.boss.Y-6, color.RGBA{255, 200, 80, 255})
					}
				}
			}
		}

		if !hit {
			aliveProjs = append(aliveProjs, proj)
		} else {
			proj.Dead = true
		}
	}
	s.projectiles = aliveProjs

	// 5. Update Enemies & Check Stomp / Damage Collisions
	for _, e := range s.enemies {
		e.Update(dt, s.lvl)
		if e.Dead || e.Squished {
			continue
		}

		enemyRect := e.GetRect()
		if playerRect.Overlaps(enemyRect) {
			// Check stomp: Player moving down and feet near top of enemy
			playerFeet := s.player.Y + entity.PlayerHeight
			enemyTop := e.Y + 4
			if s.player.VY > 0 && playerFeet <= enemyTop+7 {
				// Satisfying Stomp with combo multiplier!
				e.Stomp()
				s.player.Bounce()
				combo := s.player.StompCombo
				if combo < 1 {
					combo = 1
				}
				stompScore := 200 * combo
				s.score += stompScore
				s.TriggerShake(2.0, 0.10)
				s.hitStopTimer = 0.05
				audio.Get().PlayStompCombo(combo)
				s.particles.SpawnStompBurst(e.X+8, e.Y+8)
				if combo > 1 {
					s.SpawnFloatingText(fmt.Sprintf("COMBO x%d! +%d", combo, stompScore), e.X-8, e.Y-6, color.RGBA{255, 225, 40, 255})
				} else {
					s.SpawnFloatingText("+200", e.X, e.Y-4, color.RGBA{255, 255, 255, 255})
				}
			} else {
				// Player takes damage
				s.player.TakeDamage(1)
				s.TriggerShake(4.5, 0.22)
				if s.player.Dead {
					audio.Get().PlayGameOver()
					s.machine.Change(NewGameOverState(s.machine, s.score, s.crystalsCollected, s.currentSector, s.player))
					return
				}
			}
		}
	}

	// 5b. Update Boss AI, Collisions & Boss Projectiles
	if s.boss != nil && !s.boss.Dead {
		bossRect := s.boss.GetRect()
		if playerRect.Overlaps(bossRect) {
			playerFeet := s.player.Y + entity.PlayerHeight
			bossTop := s.boss.Y + 4
			if s.player.VY > 0 && playerFeet <= bossTop+8 {
				s.boss.TakeDamage(2, s.particles)
				s.player.Bounce()
				s.TriggerShake(4.0, 0.15)
				s.hitStopTimer = 0.06
				audio.Get().PlayStompCombo(3)
				if s.boss.Dead {
					s.score += 5000
					s.TriggerShake(7.0, 0.4)
					s.hitStopTimer = 0.12
					s.SpawnFloatingText("OVERLORD DOWN! +5000", s.boss.X-16, s.boss.Y-10, color.RGBA{255, 60, 80, 255})
				} else {
					s.SpawnFloatingText("CRITICAL STOMP! -2", s.boss.X-8, s.boss.Y-6, color.RGBA{255, 220, 60, 255})
				}
			} else {
				s.player.TakeDamage(1)
				s.TriggerShake(5.0, 0.25)
				if s.player.Dead {
					audio.Get().PlayGameOver()
					s.machine.Change(NewGameOverState(s.machine, s.score, s.crystalsCollected, s.currentSector, s.player))
					return
				}
			}
		}

		if bp := s.boss.Update(dt, s.player.X, s.player.Y, s.particles); bp != nil {
			s.bossProjectiles = append(s.bossProjectiles, bp)
		}
	}

	aliveBossProjs := s.bossProjectiles[:0]
	for _, bp := range s.bossProjectiles {
		bp.Update(dt)
		if bp.Dead {
			continue
		}
		if playerRect.Overlaps(bp.GetRect()) {
			s.player.TakeDamage(1)
			s.TriggerShake(4.5, 0.22)
			s.particles.SpawnStompBurst(bp.X+4, bp.Y+4)
			if s.player.Dead {
				audio.Get().PlayGameOver()
				s.machine.Change(NewGameOverState(s.machine, s.score, s.crystalsCollected, s.currentSector, s.player))
				return
			}
			continue
		}
		aliveBossProjs = append(aliveBossProjs, bp)
	}
	s.bossProjectiles = aliveBossProjs

	// 5c. Update Floating Text Popups
	aliveTexts := s.floatingTexts[:0]
	for _, ft := range s.floatingTexts {
		ft.Update(dt)
		if !ft.Dead {
			aliveTexts = append(aliveTexts, ft)
		}
	}
	s.floatingTexts = aliveTexts

	// 6. Check Goal: Escape Lander, Airlock, or Warp Console Reached
	landerLocked := s.lvl.HasBoss && s.boss != nil && !s.boss.Dead
	if playerRect.Overlaps(s.lvl.GetLanderRect()) {
		if s.lvl.GoalKind == "airlock" && !s.hasSecurityKey {
			// Sector 4: Key required
			s.player.X = s.lvl.LanderX - 22
			s.SpawnFloatingText("AIRLOCK LOCKED - KEY REQUIRED!", s.player.X-24, s.player.Y-8, color.RGBA{255, 70, 70, 255})
			s.TriggerShake(1.5, 0.10)
		} else if s.lvl.GoalKind == "warp_console" && s.repairCoresCount < 3 {
			// Sector 5: 3 Cores required
			s.player.X = s.lvl.LanderX - 22
			s.SpawnFloatingText(fmt.Sprintf("REPAIR CORES NEEDED (%d/3)!", s.repairCoresCount), s.player.X-28, s.player.Y-8, color.RGBA{255, 160, 50, 255})
			s.TriggerShake(1.5, 0.10)
		} else if !landerLocked {
			audio.Get().PlayWin()
			s.score += 1000 + s.player.Health*500
			clearMsg := "SECTOR SECURED! +1000"
			if s.lvl.GoalKind == "airlock" {
				clearMsg = "AIRLOCK CLEARED! +1000"
			} else if s.lvl.GoalKind == "warp_console" {
				clearMsg = "WARP DRIVE ENGAGED! +2000"
				s.score += 1000
			}
			s.SpawnFloatingText(clearMsg, s.player.X-16, s.player.Y-8, color.RGBA{100, 255, 120, 255})
			if s.currentSector < level.MaxLevels {
				s.machine.Change(NewSectorClearStateWithRecords(s.machine, s.currentSector, s.score, s.crystalsCollected, s.player, s.elapsedTime, s.medalsCollected))
			} else {
				s.machine.Change(NewWinStateWithRecords(s.machine, s.score, s.crystalsCollected, s.elapsedTime, s.medalsCollected))
			}
			return
		}
	}

	// 7. Update Particles
	s.particles.Update(dt)

	// 8. Smooth Camera Tracking with Dynamic Velocity Lookahead
	lookahead := s.player.VX * 18.0
	targetCamX := s.player.X - 160.0 + 8.0 + lookahead
	maxCamX := float64(s.lvl.Width*level.TileSize - 320)
	if maxCamX < 0 {
		maxCamX = 0
	}
	if targetCamX < 0 {
		targetCamX = 0
	} else if targetCamX > maxCamX {
		targetCamX = maxCamX
	}
	s.camX += (targetCamX - s.camX) * 0.12
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	// Calculate Screen Shake Rumble
	effCamX := s.camX
	if s.shakeTimer > 0 {
		effCamX += (rand.Float64()*2.0 - 1.0) * s.shakeIntensity
	}

	// 1. Background (Parallax starfield, earth/mars/jupiter, hills)
	s.lvl.DrawBackground(screen, effCamX)

	// 2. Level Tiles & Lander
	s.lvl.DrawTiles(screen, effCamX)

	// Lander Forcefield Shield (Active while Boss is alive in Sector 3)
	if s.lvl.HasBoss && s.boss != nil && !s.boss.Dead {
		landerRect := s.lvl.GetLanderRect()
		lsx := landerRect.X - effCamX
		if lsx > -40 && lsx < 330 {
			vector.DrawFilledRect(screen, float32(lsx-3), float32(landerRect.Y-3), float32(landerRect.W+6), float32(landerRect.H+6), color.RGBA{255, 40, 70, 110}, false)
			vector.StrokeRect(screen, float32(lsx-3), float32(landerRect.Y-3), float32(landerRect.W+6), float32(landerRect.H+6), 1, color.RGBA{255, 120, 140, 240}, false)
		}
	}

	// Sector 4: Airlock Security Barrier (Active until Keycard is found)
	if s.lvl.GoalKind == "airlock" && !s.hasSecurityKey {
		landerRect := s.lvl.GetLanderRect()
		lsx := landerRect.X - effCamX
		if lsx > -40 && lsx < 330 {
			vector.DrawFilledRect(screen, float32(lsx-3), float32(landerRect.Y-3), float32(landerRect.W+6), float32(landerRect.H+6), color.RGBA{255, 30, 60, 120}, false)
			vector.StrokeRect(screen, float32(lsx-3), float32(landerRect.Y-3), float32(landerRect.W+6), float32(landerRect.H+6), 1, color.RGBA{255, 90, 120, 240}, false)
			ui.DrawText(screen, "LOCKED", int(lsx-4), int(landerRect.Y-9), color.RGBA{255, 80, 80, 255})
		}
	}

	// Sector 5: Warp Engine Containment Shield (Active until all 3 cores repaired)
	if s.lvl.GoalKind == "warp_console" && s.repairCoresCount < 3 {
		landerRect := s.lvl.GetLanderRect()
		lsx := landerRect.X - effCamX
		if lsx > -40 && lsx < 330 {
			vector.DrawFilledRect(screen, float32(lsx-3), float32(landerRect.Y-3), float32(landerRect.W+6), float32(landerRect.H+6), color.RGBA{255, 140, 30, 120}, false)
			vector.StrokeRect(screen, float32(lsx-3), float32(landerRect.Y-3), float32(landerRect.W+6), float32(landerRect.H+6), 1, color.RGBA{255, 200, 80, 240}, false)
			ui.DrawText(screen, fmt.Sprintf("%d/3 CORES", s.repairCoresCount), int(lsx-10), int(landerRect.Y-9), color.RGBA{255, 190, 50, 255})
		}
	}

	// 2b. Moving Platforms
	for _, mp := range s.movingPlatforms {
		mp.Draw(screen, effCamX)
	}

	// 2b-2. Fragile Crumbling Platforms
	for _, cp := range s.crumblingPlatforms {
		cp.Draw(screen, effCamX)
	}

	// 2c. Cryo Geysers
	for _, g := range s.geysers {
		g.Draw(screen, effCamX)
	}

	// 2d. Secret Medals
	for _, m := range s.medals {
		m.Draw(screen, effCamX)
	}

	// 2e. Pickups (Oxygen & Ion Cores)
	for _, hp := range s.healthPickups {
		hp.Draw(screen, effCamX)
	}
	for _, bp := range s.boostPickups {
		bp.Draw(screen, effCamX)
	}

	// 2f. Security Access Keycard (Sector 4)
	if s.securityKey != nil {
		s.securityKey.Draw(screen, effCamX)
	}

	// 2g. Reactor Repair Cores (Sector 5)
	for _, rc := range s.repairCores {
		rc.Draw(screen, effCamX)
	}

	// 3. Crystals
	for _, c := range s.crystals {
		c.Draw(screen, effCamX)
	}

	// 4. Weapon Pickups (Holographic Pods)
	for _, pk := range s.pickups {
		pk.Draw(screen, effCamX)
	}

	// 5. Enemies
	for _, e := range s.enemies {
		e.Draw(screen, effCamX)
	}

	// 5b. Boss Overlord Mech & Projectiles
	if s.boss != nil && !s.boss.Dead {
		s.boss.Draw(screen, effCamX)
	}
	for _, bp := range s.bossProjectiles {
		bp.Draw(screen, effCamX)
	}

	// 6. Projectiles
	for _, proj := range s.projectiles {
		proj.Draw(screen, effCamX)
	}

	// 7. Particles
	s.particles.Draw(screen, effCamX)

	// 8. Player
	s.player.Draw(screen, effCamX)

	// 8b. Floating Score/Combo Popups
	for _, ft := range s.floatingTexts {
		ft.Draw(screen, effCamX)
	}

	// 9. Top HUD (Health, Score, Crystals, Medals, Timer, Fuel, Active Weapon, Key/Cores)
	hasMultiple := s.player.HasLaser && s.player.HasNova
	ui.DrawHUD(screen, s.player.Health, entity.MaxHealth, s.crystalsCollected, s.score, s.player.Fuel, entity.MaxFuel, int(s.player.ActiveWeapon), hasMultiple, s.currentSector, s.elapsedTime, s.medalsCollected, s.hasSecurityKey, s.repairCoresCount)

	// Active Ion Thruster Boost Badge
	if s.player.BoostTimer > 0 {
		vector.DrawFilledRect(screen, 240, 20, 75, 11, color.RGBA{45, 25, 10, 220}, false)
		vector.StrokeRect(screen, 240, 20, 75, 11, 1, color.RGBA{255, 180, 30, 255}, false)
		ui.DrawText(screen, fmt.Sprintf("BOOST %.1fs", s.player.BoostTimer), 244, 23, color.RGBA{255, 215, 60, 255})
	}

	// Boss Overlord Health Bar UI
	if s.boss != nil && !s.boss.Dead {
		vector.DrawFilledRect(screen, 95, 20, 130, 12, color.RGBA{14, 18, 28, 230}, false)
		vector.StrokeRect(screen, 95, 20, 130, 12, 1, color.RGBA{255, 60, 80, 255}, false)
		ui.DrawText(screen, "OVERLORD", 99, 23, color.RGBA{255, 80, 100, 255})
		barW := float32(64) * (float32(s.boss.Health) / float32(entity.BossMaxHealth))
		if barW < 0 {
			barW = 0
		}
		vector.DrawFilledRect(screen, 155, 23, barW, 6, color.RGBA{255, 40, 60, 255}, false)
		vector.StrokeRect(screen, 155, 23, 64, 6, 1, color.RGBA{255, 180, 190, 255}, false)
	}

	// 10. Weapon Unlock Announcement Banner
	if s.unlockBannerTimer > 0 {
		vector.DrawFilledRect(screen, 46, 38, 228, 26, color.RGBA{14, 22, 40, 235}, false)
		vector.StrokeRect(screen, 46, 38, 228, 26, 1, color.RGBA{80, 230, 255, 255}, false)
		ui.DrawText(screen, s.unlockBannerText, 60, 43, color.RGBA{255, 225, 60, 255})
		ui.DrawText(screen, "PRESS J TO FIRE!", 106, 52, color.RGBA{180, 240, 255, 255})
	}

	// 11. Paused Overlay
	if s.paused {
		vector.DrawFilledRect(screen, 0, 0, 320, 180, color.RGBA{10, 15, 25, 170}, false)
		ui.DrawText(screen, "PAUSED", 136, 80, color.RGBA{255, 240, 100, 255})
		ui.DrawText(screen, "PRESS P TO RESUME", 106, 96, color.RGBA{200, 210, 230, 255})
	}
}

func (s *PlayState) Exit() {}
