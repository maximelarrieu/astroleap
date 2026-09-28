package narrative

import (
	"testing"

	"astroleap/internal/records"
)

func TestCodexLogs(t *testing.T) {
	logs := GetAllLogs()
	if len(logs) != 5 {
		t.Fatalf("expected 5 logs, got %d", len(logs))
	}

	for i := 1; i <= 5; i++ {
		l, ok := GetLogForSector(i)
		if !ok {
			t.Errorf("expected log for sector %d", i)
		}
		if l.Sector != i {
			t.Errorf("expected sector %d, got %d", i, l.Sector)
		}
		if l.Title == "" || l.Teaser == "" || len(l.Content) == 0 {
			t.Errorf("log for sector %d has empty fields", i)
		}
	}
}

func TestUnlockLog(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	// Reset records singleton for testing
	records.UnlockLog("TEST-INIT")

	if !UnlockLog("LOG-01") {
		t.Errorf("expected LOG-01 to be unlocked first time")
	}

	if UnlockLog("LOG-01") {
		t.Errorf("expected LOG-01 to NOT unlock a second time")
	}

	if !IsLogUnlocked("LOG-01") {
		t.Errorf("expected LOG-01 to be recognized as unlocked")
	}

	if IsLogUnlocked("LOG-99") {
		t.Errorf("LOG-99 should not be unlocked")
	}
}
