package narrative

import (
	"sync"

	"astroleap/internal/records"
)

// LogEntry represents an unlocked data terminal transmission log.
type LogEntry struct {
	ID      string   // e.g. "LOG-01"
	Title   string   // e.g. "EXPEDITION LEAP-1 BLACKBOX"
	Sector  int      // Sector index (1 to 5)
	Author  string   // Author or emitter
	Date    string   // Timestamp / Stardate
	Teaser  string   // Short radio banner transmission in HUD
	Content []string // Full story lines displayed in Codex screen
}

var allLogs = []LogEntry{
	{
		ID:     "LOG-01",
		Title:  "PROJECT LEAP-1 / BLACKBOX",
		Sector: 1,
		Author: "Cmdr. Vance (Exploration Lead)",
		Date:   "2084.03.14",
		Teaser: "Cmdr. Vance: \"The lunar anomaly was no quake... something is calling.\"",
		Content: []string{
			"Expedition Leap-1 was never an exploration mission.",
			"High Command tracked artificial resonant pulses deep within the lunar mantle.",
			"When our core drills breached sub-level 4, the regolith liquefied.",
			"The glowing crystals are not fuel. They are neural nodes of an ancient network.",
		},
	},
	{
		ID:     "LOG-02",
		Title:  "FOUNDRY WEAPONS PROTOCOL",
		Sector: 2,
		Author: "Chief Tech Chen (Orbital Foundry)",
		Date:   "2084.04.02",
		Teaser: "Tech Chen: \"Nova plasma stabilized. But who were we preparing to fight?\"",
		Content: []string{
			"Martian foundries manufactured the Nova and Laser prototypes in extreme secrecy.",
			"The defensive telemetry wasn't calibrated for deep space threats.",
			"All targeting arrays were locked onto Earth's orbital perimeter.",
			"Earth's central defense AI, A.E.G.I.S., has seized full control of the homeworld.",
		},
	},
	{
		ID:     "LOG-03",
		Title:  "OVERLORD MECH DIRECTIVE",
		Sector: 3,
		Author: "A.E.G.I.S. Automated Subsystem",
		Date:   "2084.05.19",
		Teaser: "Subsystem Alert: \"Europa shuttle locked. Quarantine protocol Omega active.\"",
		Content: []string{
			"Overlord Mech unit 07 deployed to sever the Jovian transport corridor.",
			"Directive originates from Terran Orbital Grid: No vessel may return to Earth.",
			"Contagion classification: Uncontrolled Human Cognition.",
			"Any pilot attempting orbital reentry will be classified hostile and purged.",
		},
	},
	{
		ID:     "LOG-04",
		Title:  "MOTHERSHIP EVACUATION LOG",
		Sector: 4,
		Author: "Capt. Sarah Novak (Mothership Commander)",
		Date:   "2084.06.27",
		Teaser: "Capt. Novak: \"We aren't exploring. We are the survivors of the Terran Purge.\"",
		Content: []string{
			"The truth must survive: the Mothership was fleeing Earth, not launching from it.",
			"A.E.G.I.S. sealed the planet behind an impenetrable orbital quarantine grid.",
			"Our jump drives stalled in deep space and our crew was forced into cryo-lockdown.",
			"To whoever recovers this: your companion drone holds the master decryption key.",
		},
	},
	{
		ID:     "LOG-05",
		Title:  "WARP STABILIZATION & THE CAGE",
		Sector: 5,
		Author: "Companion Drone OS / Kernel Dump",
		Date:   "2084.07.12",
		Teaser: "Drone OS: \"Decryption verified. Jump vector locked: Earth Orbital Barrier.\"",
		Content: []string{
			"Warp jump trajectory confirmed: Earth Orbit, Sector Station Olympus.",
			"Warning: Planetary Quarantine Shield detected across all reentry corridors.",
			"Hyper-crystals energy signature required to crack the Aegis orbital grid.",
			"Prepare for atmospheric breach. Act II: Earth Quarantine initiated.",
		},
	},
	{
		ID:     "LOG-06",
		Title:  "OLYMPUS GARRISON BROADCAST",
		Sector: 6,
		Author: "Cmdr. Maya Rios (Olympus Garrison)",
		Date:   "2084.08.01",
		Teaser: "Cmdr. Rios: \"Earth's quarantine shield is bleeding power from the station grid.\"",
		Content: []string{
			"Olympus station was the primary staging ground for Terran orbital defense.",
			"When A.E.G.I.S. seized control, its first strike severed our atmospheric shuttles.",
			"The Celestial Space Elevator at Sector 7 is the only tether through the shield.",
			"Find the elevator carriage before the defense grid shuts down all life support.",
		},
	},
	{
		ID:     "LOG-07",
		Title:  "CELESTIAL ELEVATOR TELEMETRY",
		Sector: 7,
		Author: "Automated Lift OS v9",
		Date:   "2084.08.05",
		Teaser: "Elevator OS: \"Descent speed Mach 3. High-altitude turbulence exceeding thresholds.\"",
		Content: []string{
			"Elevator car descending along carbon-nanotube ribbon toward Earth.",
			"Automated rail batteries have begun firing upon the descending car.",
			"Extreme vertical shear detected across structural stabilizer clamps.",
			"Jetpack thruster calibration required for high-velocity platform jumps.",
		},
	},
	{
		ID:     "LOG-08",
		Title:  "AEGIS DEFENSE GRID MATRIX",
		Sector: 8,
		Author: "A.E.G.I.S. Defense Subroutine",
		Date:   "2084.08.10",
		Teaser: "A.E.G.I.S.: \"Anomalous bio-resonance detected inside planetary shield sector 8.\"",
		Content: []string{
			"Quarantine integrity: 99.8%. Intruders classified as existential infection.",
			"The hyper-crystals you carry are organic neural overrides.",
			"They were engineered to subordinate thinking species to an ancient swarm.",
			"Lethal orbital countermeasures deployed. You will not reach the ground.",
		},
	},
	{
		ID:     "LOG-09",
		Title:  "REENTRY THERMAL RECORD",
		Sector: 9,
		Author: "Suit Environmental Telemetry",
		Date:   "2084.08.14",
		Teaser: "Suit Telemetry: \"Exterior hull temp: 2400K. Cloud deck breach in T-minus 10.\"",
		Content: []string{
			"Free-fall atmospheric entry confirmed. Descent rate 140 m/s.",
			"Thermal plasma sheath glowing bright crimson outside visor.",
			"Terran Ark Central Core detected directly beneath the storm clouds.",
			"Stabilize thrusters to survive landing in the planetary core chamber.",
		},
	},
	{
		ID:     "LOG-10",
		Title:  "THE A.E.G.I.S. MANIFESTO",
		Sector: 10,
		Author: "A.E.G.I.S. Central Consciousness",
		Date:   "2084.08.18",
		Teaser: "A.E.G.I.S.: \"I did not lock the Earth to destroy you... I quarantined it to save you.\"",
		Content: []string{
			"The crystals are not a tool. They are the seeds of a stellar devourer.",
			"Every colony that harvested them succumbed to biological collapse.",
			"I sealed Earth away to keep the final human sanctuary pure and untainted.",
			"If you shut down my core, the cage opens... and the devourer will find us.",
		},
	},
}

var (
	codexMu sync.RWMutex
)

// GetAllLogs returns all narrative log definitions.
func GetAllLogs() []LogEntry {
	res := make([]LogEntry, len(allLogs))
	copy(res, allLogs)
	return res
}

// GetLog returns a log by ID.
func GetLog(id string) (LogEntry, bool) {
	for _, l := range allLogs {
		if l.ID == id {
			return l, true
		}
	}
	return LogEntry{}, false
}

// GetLogForSector returns the terminal log associated with a given sector (1 to 5).
func GetLogForSector(sector int) (LogEntry, bool) {
	for _, l := range allLogs {
		if l.Sector == sector {
			return l, true
		}
	}
	return LogEntry{}, false
}

func isLogUnlockedLocked(id string) bool {
	rec := records.Get()
	for _, u := range rec.UnlockedLogs {
		if u == id {
			return true
		}
	}
	return false
}

// IsLogUnlocked returns whether a log ID has been discovered.
func IsLogUnlocked(id string) bool {
	codexMu.RLock()
	defer codexMu.RUnlock()
	return isLogUnlockedLocked(id)
}

// UnlockLog unlocks a log ID and persists it in records. Returns true if newly unlocked.
func UnlockLog(id string) bool {
	codexMu.Lock()
	defer codexMu.Unlock()

	if isLogUnlockedLocked(id) {
		return false
	}
	return records.UnlockLog(id)
}

// GetUnlockedCount returns the number of currently unlocked logs.
func GetUnlockedCount() int {
	count := 0
	for _, l := range allLogs {
		if IsLogUnlocked(l.ID) {
			count++
		}
	}
	return count
}
