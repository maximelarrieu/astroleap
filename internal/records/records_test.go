package records

import "testing"

func TestFormatTime(t *testing.T) {
	if got := FormatTime(0); got != "--:--.-" {
		t.Errorf("expected '--:--.-', got %s", got)
	}

	if got := FormatTime(65.4); got != "01:05.4" {
		t.Errorf("expected '01:05.4', got %s", got)
	}

	if got := FormatTime(185.0); got != "03:05.0" {
		t.Errorf("expected '03:05.0', got %s", got)
	}
}

func TestEvaluateRank(t *testing.T) {
	// S Rank
	if rank := EvaluateRank(14000, 180.0, 9); rank != "S" {
		t.Errorf("expected S rank, got %s", rank)
	}

	// A Rank
	if rank := EvaluateRank(9000, 240.0, 5); rank != "A" {
		t.Errorf("expected A rank, got %s", rank)
	}

	// B Rank
	if rank := EvaluateRank(6000, 360.0, 3); rank != "B" {
		t.Errorf("expected B rank, got %s", rank)
	}

	// C Rank
	if rank := EvaluateRank(3000, 500.0, 1); rank != "C" {
		t.Errorf("expected C rank, got %s", rank)
	}
}

func TestSubmitRunRecords(t *testing.T) {
	res := SubmitRun(15000, 150.0, 7)
	rec := Get()

	if rec.BestScore < 15000 {
		t.Errorf("expected best score >= 15000, got %d", rec.BestScore)
	}

	if rec.BestTime <= 0 || rec.BestTime > 150.0 {
		t.Errorf("expected best time <= 150.0, got %f", rec.BestTime)
	}

	if rec.TotalMedals < 7 {
		t.Errorf("expected total medals >= 7, got %d", rec.TotalMedals)
	}

	if !res.NewHighScore && rec.BestScore != 15000 {
		t.Errorf("expected HighScore notification")
	}
}
