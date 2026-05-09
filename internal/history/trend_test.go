package history

import (
	"testing"
	"time"
)

func seedTrendStore(t *testing.T) string {
	t.Helper()
	path := tempPath(t)
	store, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Early runs: slow and failing.
	for i := 0; i < 5; i++ {
		e := Entry{
			RunID:    fmt.Sprintf("early-%d", i),
			JobName:  "backup",
			Duration: 10 * time.Second,
			Error:    "timeout",
			At:       time.Now().Add(-time.Duration(10-i) * time.Hour),
		}
		store.entries["backup"] = append(store.entries["backup"], e)
	}
	// Late runs: fast and succeeding.
	for i := 0; i < 5; i++ {
		e := Entry{
			RunID:    fmt.Sprintf("late-%d", i),
			JobName:  "backup",
			Duration: 2 * time.Second,
			Error:    "",
			At:       time.Now().Add(-time.Duration(5-i) * time.Hour),
		}
		store.entries["backup"] = append(store.entries["backup"], e)
	}
	if err := store.save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestComputeTrend_ImprovingDuration(t *testing.T) {
	path := seedTrendStore(t)
	trend, err := ComputeTrend(path, "backup", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trend.DurationTrend != TrendImproving {
		t.Errorf("expected DurationTrend=improving, got %s", trend.DurationTrend)
	}
}

func TestComputeTrend_ImprovingSuccessRate(t *testing.T) {
	path := seedTrendStore(t)
	trend, err := ComputeTrend(path, "backup", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trend.SuccessRateTrend != TrendImproving {
		t.Errorf("expected SuccessRateTrend=improving, got %s", trend.SuccessRateTrend)
	}
}

func TestComputeTrend_EmptyPath_ReturnsError(t *testing.T) {
	_, err := ComputeTrend("", "backup", 10)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestComputeTrend_EmptyJobName_ReturnsError(t *testing.T) {
	path := seedTrendStore(t)
	_, err := ComputeTrend(path, "", 10)
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestComputeTrend_NLessThanTwo_ReturnsError(t *testing.T) {
	path := seedTrendStore(t)
	_, err := ComputeTrend(path, "backup", 1)
	if err == nil {
		t.Fatal("expected error for n < 2")
	}
}

func TestComputeTrend_StableWhenSimilar(t *testing.T) {
	path := tempPath(t)
	store, _ := New(path)
	for i := 0; i < 6; i++ {
		store.entries["noop"] = append(store.entries["noop"], Entry{
			RunID:    fmt.Sprintf("r%d", i),
			JobName:  "noop",
			Duration: 5 * time.Second,
			Error:    "",
			At:       time.Now(),
		})
	}
	_ = store.save()
	trend, err := ComputeTrend(path, "noop", 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trend.DurationTrend != TrendStable {
		t.Errorf("expected stable, got %s", trend.DurationTrend)
	}
}
