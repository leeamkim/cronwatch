package history

import (
	"testing"
	"time"
)

func seedBaseline(t *testing.T, entries []Entry) string {
	t.Helper()
	path := tempPath(t)
	store, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	store.data["job"] = entries
	if err := store.save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestComputeBaseline_AvgDuration(t *testing.T) {
	path := seedBaseline(t, []Entry{
		{Status: "success", Duration: 4 * time.Second},
		{Status: "success", Duration: 6 * time.Second},
	})

	b, err := ComputeBaseline(path, "job", 0)
	if err != nil {
		t.Fatalf("ComputeBaseline: %v", err)
	}
	if b.AvgDuration != 5*time.Second {
		t.Errorf("expected avg 5s, got %s", b.AvgDuration)
	}
}

func TestComputeBaseline_SuccessRate(t *testing.T) {
	path := seedBaseline(t, []Entry{
		{Status: "success", Duration: time.Second},
		{Status: "failure", Duration: time.Second},
		{Status: "success", Duration: time.Second},
		{Status: "success", Duration: time.Second},
	})

	b, err := ComputeBaseline(path, "job", 0)
	if err != nil {
		t.Fatalf("ComputeBaseline: %v", err)
	}
	if b.SuccessRate != 0.75 {
		t.Errorf("expected 0.75 success rate, got %f", b.SuccessRate)
	}
}

func TestComputeBaseline_LimitsToN(t *testing.T) {
	path := seedBaseline(t, []Entry{
		{Status: "failure", Duration: 10 * time.Second},
		{Status: "success", Duration: 2 * time.Second},
		{Status: "success", Duration: 4 * time.Second},
	})

	b, err := ComputeBaseline(path, "job", 2)
	if err != nil {
		t.Fatalf("ComputeBaseline: %v", err)
	}
	if b.SampleSize != 2 {
		t.Errorf("expected sample size 2, got %d", b.SampleSize)
	}
	if b.SuccessRate != 1.0 {
		t.Errorf("expected 100%% success rate from last 2, got %f", b.SuccessRate)
	}
}

func TestComputeBaseline_EmptyPath_ReturnsError(t *testing.T) {
	_, err := ComputeBaseline("", "job", 0)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestComputeBaseline_UnknownJob_ReturnsError(t *testing.T) {
	path := seedBaseline(t, []Entry{})
	_, err := ComputeBaseline(path, "ghost", 0)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestBaseline_IsAnomalous(t *testing.T) {
	b := &Baseline{AvgDuration: 10 * time.Second}

	if !b.IsAnomalous(25*time.Second, 2.0) {
		t.Error("expected anomalous for 25s vs avg 10s at 2x threshold")
	}
	if b.IsAnomalous(15*time.Second, 2.0) {
		t.Error("expected not anomalous for 15s vs avg 10s at 2x threshold")
	}
}

func TestBaseline_IsAnomalous_ZeroAvg(t *testing.T) {
	// A zero average duration should never be considered anomalous regardless
	// of the observed duration, since there is no meaningful baseline to compare
	// against.
	b := &Baseline{AvgDuration: 0}

	if b.IsAnomalous(5*time.Second, 2.0) {
		t.Error("expected not anomalous when avg duration is zero")
	}
}
