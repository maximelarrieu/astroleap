# Game Design Document: AstroLeap (Lunar Odyssey)

> **Game Title**: AstroLeap: Lunar Odyssey  
> **Target Genre**: 2D Low-Gravity Platformer  
> **Target Platform**: Desktop (Linux / macOS / Windows) & WebAssembly (WASM)  
> **Target Aspect Ratio & Resolution**: 16:9 Widescreen (`320x180` Virtual Pixel Canvas)  
> **Author / Lead Designer**: User & Antigravity  

---

## 1. Executive Summary & Elevator Pitch

- **Elevator Pitch**: *Super Mario Bros* on the Moon—guide an astronaut across hazardous lunar terrain with floaty low-gravity physics and a limited jetpack thruster, stomping alien critters and collecting glowing energy crystals on a quest to reach the rescue lunar lander.
- **Core Inspiration**: *Super Mario Bros* meets *Moon Patrol* & *Jetpack Joyride* low-gravity jump mechanics.
- **Target Audience / Mood**: Nostalgic retro arcade platforming with floaty, satisfying aerial control and vibrant sci-fi cosmic aesthetics.

---

## 2. Core Gameplay Loop & Mechanics

- **Primary Gameplay Loop**:
  1. **Traverse & Jump**: Run across lunar craters, platforms, and floating space rocks using floaty moon gravity.
  2. **Thrust & Hover**: Ignite the jetpack thruster meter in mid-air to extend jumps, glide across chasms, or position above enemies.
  3. **Stomp & Avoid**: Stomp on lunar aliens (Moon Slimes, Eyeballs) to bounce high and eliminate them, or avoid hazards (cosmic spikes, acid pools, abysses).
  4. **Collect & Recharge**: Collect glowing Energy Crystals to earn score, refill thrusters, and earn extra lives.
  5. **Stage Exit**: Reach the Lunar Lander module at the end of the stage to board and blast off to the next zone.

- **Sector & Mission Progression (5 Stages)**:
  1. **Sector 1: Lunar Outpost** (Moon regolith, floaty jumping tutorial, Laser Blaster unlock, escape lander).
  2. **Sector 2: Phobos Ridge** (Martian crimson canyon, Nova Cannon secret unlock, moving platforms across spike chasm).
  3. **Sector 3: Europa Ice Core** (Europa glacial caverns, Boss Overlord Mech encounter guarding transport shuttle).
  4. **Sector 4: Mothership Corridor** (Derelict human mothership interior, security lockdown; recover the holographic **Security Keycard** to unlock the bulk-head airlock door leading to the engine core).
  5. **Sector 5: Reactor Bay** (Spaceship reactor engine room in critical failure; recover 3 **Warp Repair Cores** across moving pistons and coolant vents to repair the ship, engage warp drive, and return to Earth!).

- **Core Mechanics**:
  - **Moon Gravity**: Lower gravity constant ($g \approx 0.18\times$ standard) allowing majestic, high-arc jumps.
  - **Jetpack Thruster Gauge**: Holding jump mid-air activates thruster boost with particle sparks. Thruster fuel depletes smoothly and auto-recharges when standing on solid ground or grabbing blue fuel canisters.
  - **Moving Platform Carriage & Physics**: Mechanical platforms moving between waypoints with Euclidean distance rate timing. Riders are locked smoothly to platform movement without sinking or jitter, with full jump clearance and wall-collision prevention.
  - **Keycard & Airlock Security**: In Sector 4, the exit airlock is locked behind an active security barrier. Collecting the golden digital keycard deactivates the barrier.
  - **Reactor Repair System**: In Sector 5, the warp engine console remains locked until all 3 critical repair cores (Coolant Regulator, Plasma Stabilizer, Hyper-Flux Rod) are retrieved and installed.
  - **Enemy Stomping**: Landing on top of aliens knocks them flat with a satisfying bounce impulse and pop sound effect.
  - **Unlockable Weapons**:
    - **Laser Blaster**: Unlocked by collecting the holographic Laser Pod at Tile X=20. Fires high-speed cyan photon energy beams that blast through aliens horizontally.
    - **Nova Cannon**: Unlocked by reaching the secret high platform at Tile X=55 with the jetpack. Fires bouncing golden plasma stars that ricochet across lunar floors and walls before exploding.
    - **Active Weapon HUD & Switching**: Players can cycle between unlocked weapons seamlessly using `Q` or gamepad shoulder buttons.
  - **Win Condition**: Board the Lunar Lander, clear the Mothership Corridor with the Security Key, and fully repair the Reactor Core in Sector 5 to engage Warp Drive and return home!
  - **Loss Condition**: Running out of Astronaut Oxygen/Hearts (3 hits) or falling into the lunar/ship chasms.

---

## 3. Controls & Input Mapping Scheme

- **Primary Input Devices**: Keyboard, Gamepad, Touch (Virtual Controls for WASM mobile).
- **Default Control Scheme**:

| Logical Action | Keyboard | Gamepad Button | Touch / Mouse |
| :--- | :--- | :--- | :--- |
| `ActionMoveLeft` | `A` or `Left Arrow` | D-Pad Left / Left Stick | Virtual Left Zone |
| `ActionMoveRight` | `D` or `Right Arrow` | D-Pad Right / Left Stick | Virtual Center Zone |
| `ActionJump` | `Space` / `W` / `Up Arrow` | Button South (`A` / `X`) | Virtual Right Zone |
| `ActionJetpack` | Hold `Space` / `W` mid-air | Hold Button South (`A` / `X`) | Hold Right Zone |
| `ActionShoot` | `J` / `Z` / `F` / Left Click | Button West (`X` / `Square`) | Screen Click / Tap |
| `ActionSwitch` | `Q` / `Tab` / `E` | Button North (`Y`) / `R1` | Top HUD Weapon Tap |
| `ActionPause` | `Escape` / `P` | Start Button | Pause Icon Button |
| `ActionRestart` | `R` (on Game Over) | Select / Back | Tap Screen |

---

## 4. Visual Style & Asset Strategy

- **Aesthetic Direction**: 16-bit "Retro-HD" Pixel Art (Crisp silhouettes, cosmic color ramps, deep space nebulae, starfield parallax).
- **Asset Pipeline**:
  - **Generative AI Sprites & Art (`nano-banana` / Gemini Image)**:
    - Astronaut Player Character: White space suit, gold visor, jetpack thruster.
    - Alien Critters: Bouncing lunar slime/blob, floating cyclops alien.
    - Environment: Moon dust terrain, crater blocks, glowing energy crystals, rescue lander.
    - Deep Space Parallax: Distant Earth, starry cosmos, purple nebula dust.
  - **Procedural Pure-Code FX (`procedural-art`)**:
    - Zero-dependency engine rendering for rocket thrust flame particles, landing dust clouds, starfield sparkle, and HUD status bars.

---

## 5. Audio & Soundscape Strategy

- **Background Music (BGM)**:
  - **Style**: Atmospheric space-synth / cosmic chiptune (44.1 kHz stereo) with floating arpeggios, driving synthwave bass, and ethereal pads.
  - **Generation Engine**: High-fidelity AI music (`lyria`) paired with built-in zero-dependency polyphonic FM synthesis (`procedural-composer`) so the game has great audio out of the box.
- **Sound Effects (SFX)**:
  - Float Jump (low-pass frequency sweep)
  - Thruster Hiss (band-pass filtered white noise burst)
  - Alien Stomp (punchy square wave pop + pitch bend)
  - Crystal Pickup (sparkling dual-sine chime)
  - Rocket Blastoff / Victory Fanfare (ascending chord stinger)
  - Player Hurt (crunchy noise/saw shock)

---

## 6. Game State Sequence & HUD Layout

- **Scene Progression**:
  `Boot (Company Logo)` $\rightarrow$ `Title Screen` $\rightarrow$ `Gameplay` $\rightarrow$ `Level Victory / Game Over` $\rightarrow$ `Title Screen`
- **HUD & UI Overlay (16:9 Virtual 320x180 Canvas)**:
  - **Top-Left**: Astronaut Hearts (3 units) & Score
  - **Top-Center**: Energy Crystals collected (`💎 x 00`)
  - **Top-Right**: Jetpack Thruster Fuel Gauge (Cyan/Orange segmented bar)
  - **Game Over / Victory**: Centered modal with final score and restart prompt.

---

## 7. Technical Scope & Architecture Notes (Ebitengine v2)

- **Engine**: Ebitengine v2 (`github.com/hajimehoshi/ebiten/v2`)
- **Resolution**: Fixed internal virtual canvas `320x180` scaled to window with crisp nearest-neighbor sampling.
- **Physics**: Sub-pixel delta time ($dt$) fixed 60 FPS loop with AABB sweep collision detection for tilemaps and dynamic entities.
- **Directory Structure (`internal/`)**:
  - `cmd/game/main.go`: Desktop window entrypoint.
  - `internal/game/`: Top-level `ebiten.Game` implementation.
  - `internal/state/`: Scene state machine (`TitleState`, `PlayState`, `GameOverState`, `WinState`).
  - `internal/entity/`: Player, Enemies, Crystals, Lander.
  - `internal/level/`: Tilemap parser, tile collision grid, parallax background.
  - `internal/input/`: Cross-platform input mapping (Keyboard, Gamepad, Touch).
  - `internal/audio/`: Sound manager (SFX synthesis + BGM player).
  - `internal/ui/`: Retro HUD rendering, thruster meter, font rasterizer.
  - `internal/assets/`: Embedded sprite sheets and audio assets (`//go:embed`).
