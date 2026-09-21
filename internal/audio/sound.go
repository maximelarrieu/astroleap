package audio

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"sync"

	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate = 44100
)

// SoundManager handles procedural audio synthesis, SFX, and BGM loops.
type SoundManager struct {
	ctx        *ebitenaudio.Context
	bgmPlayer  *ebitenaudio.Player
	bgmPlaying bool
	muted      bool
	mu         sync.Mutex

	// Pre-synthesized SFX PCM buffers
	jumpPCM     []byte
	thrustPCM   []byte
	stompPCM    []byte
	stompPCMs   [5][]byte
	crystalPCM  []byte
	hurtPCM     []byte
	winPCM      []byte
	gameoverPCM []byte
	laserPCM    []byte
	novaPCM       []byte
	powerupPCM    []byte
	typewriterPCM []byte
	bgmPCM        []byte
}

var (
	instance *SoundManager
	once     sync.Once
)

// Get returns the singleton SoundManager instance.
func Get() *SoundManager {
	once.Do(func() {
		ctx := ebitenaudio.NewContext(SampleRate)
		sm := &SoundManager{
			ctx: ctx,
		}
		sm.bakeSFX()
		instance = sm
	})
	return instance
}

func (s *SoundManager) SetMuted(m bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.muted = m
	if s.bgmPlayer != nil {
		if m {
			s.bgmPlayer.SetVolume(0)
		} else {
			s.bgmPlayer.SetVolume(0.55)
		}
	}
}

func (s *SoundManager) ToggleMute() bool {
	s.SetMuted(!s.muted)
	return s.muted
}

func (s *SoundManager) IsMuted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.muted
}

// PlaySFX plays a sound effect buffer.
func (s *SoundManager) PlaySFX(pcm []byte, vol float64) {
	if s.muted || len(pcm) == 0 {
		return
	}
	p := s.ctx.NewPlayerFromBytes(pcm)
	if p == nil {
		return
	}
	p.SetVolume(vol)
	p.Play()
}

func (s *SoundManager) PlayJump() {
	s.PlaySFX(s.jumpPCM, 0.45)
}

func (s *SoundManager) PlayThrust() {
	s.PlaySFX(s.thrustPCM, 0.25)
}

func (s *SoundManager) PlayStomp() {
	s.PlayStompCombo(1)
}

func (s *SoundManager) PlayStompCombo(combo int) {
	if s.muted {
		return
	}
	idx := combo - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(s.stompPCMs) {
		idx = len(s.stompPCMs) - 1
	}
	pcm := s.stompPCMs[idx]
	if len(pcm) == 0 {
		pcm = s.stompPCM
	}
	s.PlaySFX(pcm, 0.65)
}

func (s *SoundManager) PlayCrystal() {
	s.PlaySFX(s.crystalPCM, 0.5)
}

func (s *SoundManager) PlayHurt() {
	s.PlaySFX(s.hurtPCM, 0.6)
}

func (s *SoundManager) PlayWin() {
	s.PlaySFX(s.winPCM, 0.7)
}

func (s *SoundManager) PlayGameOver() {
	s.PlaySFX(s.gameoverPCM, 0.6)
}

func (s *SoundManager) PlayLaser() {
	s.PlaySFX(s.laserPCM, 0.5)
}

func (s *SoundManager) PlayNova() {
	s.PlaySFX(s.novaPCM, 0.65)
}

func (s *SoundManager) PlayPowerup() {
	s.PlaySFX(s.powerupPCM, 0.7)
}

func (s *SoundManager) PlayTypewriter() {
	s.PlaySFX(s.typewriterPCM, 0.25)
}

// StartBGM begins playing the ambient cosmic soundtrack in a seamless loop.
func (s *SoundManager) StartBGM() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bgmPlaying || len(s.bgmPCM) == 0 {
		return
	}

	reader := bytes.NewReader(s.bgmPCM)
	loop := ebitenaudio.NewInfiniteLoop(reader, int64(len(s.bgmPCM)))
	player, err := s.ctx.NewPlayer(loop)
	if err != nil {
		return
	}
	s.bgmPlayer = player
	if s.muted {
		s.bgmPlayer.SetVolume(0)
	} else {
		s.bgmPlayer.SetVolume(0.5)
	}
	s.bgmPlayer.Play()
	s.bgmPlaying = true
}

func (s *SoundManager) StopBGM() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bgmPlayer != nil {
		s.bgmPlayer.Close()
		s.bgmPlayer = nil
		s.bgmPlaying = false
	}
}

// ============================================================================
// DSP SYNTHESIS HELPERS
// ============================================================================

func encodeStereo16(samples [][2]float64) []byte {
	buf := make([]byte, len(samples)*4)
	for i, s := range samples {
		l := math.Max(-1.0, math.Min(1.0, s[0]))
		r := math.Max(-1.0, math.Min(1.0, s[1]))
		iL := int16(l * 32767.0)
		iR := int16(r * 32767.0)
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(iL))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(iR))
	}
	return buf
}

func (s *SoundManager) bakeSFX() {
	// 1. Jump SFX: Floaty rising frequency sweep
	{
		dur := 0.22
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			freq := 240.0 + (t/dur)*360.0 // 240Hz to 600Hz
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			env := math.Sin((t / dur) * math.Pi) // Smooth hump envelope
			v := (math.Sin(phase) + 0.3*math.Sin(phase*2.0)) * env * 0.7
			samples[i] = [2]float64{v, v}
		}
		s.jumpPCM = encodeStereo16(samples)
	}

	// 2. Thruster SFX: Warm noise hiss burst
	{
		dur := 0.12
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		rng := rand.New(rand.NewSource(42))
		last := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			noise := (rng.Float64()*2.0 - 1.0)
			// Simple low-pass filter
			last = last*0.82 + noise*0.18
			env := 1.0 - (t / dur)
			v := last * env * 0.55
			samples[i] = [2]float64{v, v}
		}
		s.thrustPCM = encodeStereo16(samples)
	}

	// 3. Stomp SFX: Escalating pitch punchy low-end thud (combo ladder)
	{
		dur := 0.18
		numSamples := int(float64(SampleRate) * dur)
		basePitches := []float64{320.0, 392.0, 480.0, 600.0, 750.0}
		for step, baseFreq := range basePitches {
			samples := make([][2]float64, numSamples)
			phase := 0.0
			for i := 0; i < numSamples; i++ {
				t := float64(i) / float64(SampleRate)
				freq := baseFreq * math.Exp(-t*18.0) // Rapid pitch decay
				phase += 2.0 * math.Pi * freq / float64(SampleRate)
				env := math.Exp(-t * 14.0)
				// Square-ish wave for punch
				sq := math.Sin(phase)
				if sq > 0 {
					sq = 1.0
				} else {
					sq = -1.0
				}
				v := sq * env * 0.65
				samples[i] = [2]float64{v, v}
			}
			s.stompPCMs[step] = encodeStereo16(samples)
		}
		s.stompPCM = s.stompPCMs[0]
	}

	// 4. Crystal SFX: Sparkly dual-tone bell chime
	{
		dur := 0.35
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		p1, p2 := 0.0, 0.0
		f1, f2 := 987.77, 1318.51 // B5 and E6 major interval
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			p1 += 2.0 * math.Pi * f1 / float64(SampleRate)
			p2 += 2.0 * math.Pi * f2 / float64(SampleRate)
			env := math.Exp(-t * 9.0)
			v := (math.Sin(p1)*0.6 + math.Sin(p2)*0.4) * env * 0.7
			samples[i] = [2]float64{v, v}
		}
		s.crystalPCM = encodeStereo16(samples)
	}

	// 5. Hurt SFX: Harsh saw tone + noise
	{
		dur := 0.28
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		rng := rand.New(rand.NewSource(99))
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			freq := 160.0 - (t/dur)*80.0
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			saw := (2.0 * (phase/(2.0*math.Pi) - math.Floor(phase/(2.0*math.Pi)+0.5)))
			noise := (rng.Float64()*2.0 - 1.0) * 0.4
			env := 1.0 - (t / dur)
			v := (saw*0.6 + noise*0.4) * env * 0.8
			samples[i] = [2]float64{v, v}
		}
		s.hurtPCM = encodeStereo16(samples)
	}

	// 6. Win SFX: Ascending fanfare chord (C5 - E5 - G5 - C6)
	{
		notes := []float64{523.25, 659.25, 783.99, 1046.50}
		durPerNote := 0.12
		totDur := durPerNote * 4.0
		numSamples := int(float64(SampleRate) * totDur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			noteIdx := int(t / durPerNote)
			if noteIdx > 3 {
				noteIdx = 3
			}
			freq := notes[noteIdx]
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			noteT := math.Mod(t, durPerNote)
			env := math.Sin((noteT / durPerNote) * math.Pi)
			if noteIdx == 3 {
				env = math.Exp(-noteT * 3.0)
			}
			v := math.Sin(phase) * env * 0.7
			samples[i] = [2]float64{v, v}
		}
		s.winPCM = encodeStereo16(samples)
	}

	// 7. GameOver SFX: Sad descending tones
	{
		notes := []float64{392.00, 369.99, 349.23, 293.66}
		durPerNote := 0.2
		totDur := durPerNote * 4.0
		numSamples := int(float64(SampleRate) * totDur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			noteIdx := int(t / durPerNote)
			if noteIdx > 3 {
				noteIdx = 3
			}
			freq := notes[noteIdx]
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			noteT := math.Mod(t, durPerNote)
			env := math.Exp(-noteT * 4.0)
			v := (math.Sin(phase) + 0.2*math.Sin(phase*2)) * env * 0.65
			samples[i] = [2]float64{v, v}
		}
		s.gameoverPCM = encodeStereo16(samples)
	}

	// 8. Laser Blaster SFX: Fast downward futuristic frequency sweep
	{
		dur := 0.14
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			// Frequency sweeps down from 1200Hz to 280Hz
			freq := 1200.0 * math.Exp(-t*16.0)
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			env := 1.0 - (t / dur)
			// Square wave with high harmonic bite
			sq := math.Sin(phase)
			if sq > 0 {
				sq = 0.8
			} else {
				sq = -0.8
			}
			v := (sq*0.65 + math.Sin(phase*2.0)*0.35) * env * 0.6
			samples[i] = [2]float64{v, v}
		}
		s.laserPCM = encodeStereo16(samples)
	}

	// 9. Nova Cannon SFX: Resonant plasma discharge + low boom
	{
		dur := 0.26
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			freq := 380.0 * math.Exp(-t*10.0) + 60.0
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			env := math.Exp(-t * 8.0)
			// Triangle / Sine mix with sub-bass resonance
			tri := (2.0/math.Pi)*math.Asin(math.Sin(phase))
			v := (tri*0.5 + math.Sin(phase*0.5)*0.5) * env * 0.75
			samples[i] = [2]float64{v, v}
		}
		s.novaPCM = encodeStereo16(samples)
	}

	// 10. Weapon Unlock / Power-Up Fanfare: Bright ascending 4-tone chime
	{
		notes := []float64{659.25, 830.61, 987.77, 1318.51} // E5, G#5, B5, E6
		durPerNote := 0.09
		totDur := durPerNote * 4.0
		numSamples := int(float64(SampleRate) * totDur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			noteIdx := int(t / durPerNote)
			if noteIdx > 3 {
				noteIdx = 3
			}
			freq := notes[noteIdx]
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			noteT := math.Mod(t, durPerNote)
			env := math.Sin((noteT / durPerNote) * math.Pi)
			if noteIdx == 3 {
				env = math.Exp(-noteT * 2.5)
			}
			v := (math.Sin(phase)*0.7 + math.Sin(phase*2.0)*0.3) * env * 0.75
			samples[i] = [2]float64{v, v}
		}
		s.powerupPCM = encodeStereo16(samples)
	}

	// 11. Typewriter SFX: Gentle retro computer terminal blip
	{
		dur := 0.025
		numSamples := int(float64(SampleRate) * dur)
		samples := make([][2]float64, numSamples)
		phase := 0.0
		freq := 1760.0 // Soft high blip
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			phase += 2.0 * math.Pi * freq / float64(SampleRate)
			env := math.Exp(-t * 90.0)
			v := math.Sin(phase) * env * 0.18
			samples[i] = [2]float64{v, v}
		}
		s.typewriterPCM = encodeStereo16(samples)
	}

	// 12. Background Music (BGM): 8-second seamless looping space-synth track
	{
		bpm := 120.0
		secPerBeat := 60.0 / bpm
		totalBeats := 16.0 // 8 seconds total loop
		totDur := totalBeats * secPerBeat
		numSamples := int(float64(SampleRate) * totDur)
		samples := make([][2]float64, numSamples)

		// Chord progression: Am -> F -> C -> G (Root freq in Hz)
		bassNotes := []float64{110.0, 87.31, 130.81, 98.0} // A2, F2, C3, G2
		arpeggios := [][]float64{
			{220.0, 261.63, 329.63, 440.0}, // Am (A3, C4, E4, A4)
			{174.61, 220.0, 261.63, 349.23}, // F  (F3, A3, C4, F4)
			{261.63, 329.63, 392.0, 523.25}, // C  (C4, E4, G4, C5)
			{196.0, 246.94, 293.66, 392.0},  // G  (G3, B3, D4, G4)
		}

		bassPhase := 0.0
		arpPhase := 0.0

		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(SampleRate)
			barIdx := int(t / (4.0 * secPerBeat)) % 4
			barT := math.Mod(t, 4.0*secPerBeat)

			// Bassline: Pulse wave on beats 1 and 3
			curBassFreq := bassNotes[barIdx]
			bassPhase += 2.0 * math.Pi * curBassFreq / float64(SampleRate)
			beatInBar := math.Mod(barT, secPerBeat)
			bassEnv := math.Exp(-beatInBar * 3.5)
			// Triangle/Soft pulse
			bassVal := math.Sin(bassPhase) * 0.45 * bassEnv

			// Arpeggiator: 16th notes (4 notes per beat)
			stepInBeat := int(barT / (secPerBeat / 4.0))
			arpNoteIdx := stepInBeat % 4
			curArpFreq := arpeggios[barIdx][arpNoteIdx]
			arpPhase += 2.0 * math.Pi * curArpFreq / float64(SampleRate)
			stepT := math.Mod(barT, secPerBeat/4.0)
			arpEnv := math.Exp(-stepT * 12.0)
			arpVal := math.Sin(arpPhase) * 0.28 * arpEnv

			// Ambient stereo pad
			padFreq := arpeggios[barIdx][0] * 1.5
			padVal := math.Sin(2.0*math.Pi*padFreq*t) * 0.12

			// Hi-hat tick every 8th note
			tickEnv := math.Exp(-math.Mod(t, secPerBeat/2.0) * 45.0)
			rngVal := (rand.Float64()*2.0 - 1.0) * 0.06 * tickEnv

			left := (bassVal*0.7 + arpVal*0.9 + padVal*0.6 + rngVal) * 0.5
			right := (bassVal*0.7 + arpVal*0.6 + padVal*0.9 + rngVal) * 0.5

			samples[i] = [2]float64{left, right}
		}
		s.bgmPCM = encodeStereo16(samples)
	}
}
