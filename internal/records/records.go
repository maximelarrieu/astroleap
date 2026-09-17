package records

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type GameRecords struct {
	BestScore   int     `json:"best_score"`
	BestTime    float64 `json:"best_time"` // in seconds
	TotalMedals int     `json:"total_medals"`
	Completions int     `json:"completions"`
}

type RecordResult struct {
	NewHighScore bool
	NewBestTime  bool
}

var (
	currentRecords GameRecords
	recordsOnce    sync.Once
	recordsMu      sync.Mutex
)

func getRecordFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".astroleap_records.json"
	}
	return filepath.Join(home, ".astroleap_records.json")
}

func load() {
	currentRecords = GameRecords{
		BestScore:   0,
		BestTime:    0,
		TotalMedals: 0,
		Completions: 0,
	}

	data, err := os.ReadFile(getRecordFilePath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &currentRecords)
}

// Get returns the current all-time best records.
func Get() GameRecords {
	recordsOnce.Do(load)
	recordsMu.Lock()
	defer recordsMu.Unlock()
	return currentRecords
}

// SubmitRun updates records if new high score or best speedrun time was achieved.
func SubmitRun(score int, timeSeconds float64, medals int) RecordResult {
	recordsOnce.Do(load)
	recordsMu.Lock()
	defer recordsMu.Unlock()

	res := RecordResult{}

	if score > currentRecords.BestScore {
		currentRecords.BestScore = score
		res.NewHighScore = true
	}

	if currentRecords.BestTime <= 0 || (timeSeconds > 0 && timeSeconds < currentRecords.BestTime) {
		currentRecords.BestTime = timeSeconds
		res.NewBestTime = true
	}

	if medals > currentRecords.TotalMedals {
		currentRecords.TotalMedals = medals
	}
	currentRecords.Completions++

	// Save to disk asynchronously/safely
	data, err := json.Marshal(currentRecords)
	if err == nil {
		_ = os.WriteFile(getRecordFilePath(), data, 0644)
	}

	return res
}

// FormatTime converts seconds into MM:SS.S string format.
func FormatTime(seconds float64) string {
	if seconds <= 0 {
		return "--:--.-"
	}
	mins := int(seconds) / 60
	secs := int(seconds) % 60
	tenths := int(seconds*10) % 10

	minStr := string(rune('0'+mins/10)) + string(rune('0'+mins%10))
	secStr := string(rune('0'+secs/10)) + string(rune('0'+secs%10))
	tenthStr := string(rune('0' + tenths))

	return minStr + ":" + secStr + "." + tenthStr
}

// EvaluateRank determines mission completion rank (S, A, B, C).
func EvaluateRank(score int, timeSeconds float64, medals int) string {
	// S-Rank: under 3m30s, score >= 12000, or all 9 medals found
	if (timeSeconds < 210.0 && score >= 12000) || medals >= 9 {
		return "S"
	}
	if timeSeconds < 300.0 && score >= 8000 {
		return "A"
	}
	if timeSeconds < 420.0 && score >= 5000 {
		return "B"
	}
	return "C"
}
