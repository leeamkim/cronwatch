package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func seedCompareStore(t *testing.T) string {
	t.Helper()
	p := tempPath(t)
	h, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now()
	for i := 0; i < 4; i++ {
		_ = h.Record(Entry{
			JobName:   "alpha",
			RunID:     fmt.Sprintf("a-%d", i),
			StartedAt: now.Add(-time.Duration(i) * time.Minute),
			Duration:  2 * time.Second,
			Success:   true,
		})
	}
	for i := 0; i < 3; i++ {
		success := i%2 == 0
		_ = h.Record(Entry{
			JobName:   "beta",
			RunID:     fmt.Sprintf("b-%d", i),
			StartedAt: now.Add(-time.Duration(i) * time.Minute),
			Duration:  5 * time.Second,
			Success:   success,
		})
	}
	return p
}

func TestCompareJobs_ReturnsStats(t *testing.T) {
	p := seedCompareStore(t)
	r, err := CompareJobs(p, "alpha", "beta", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.StatsA.TotalRuns != 4 {
		t.Errorf("expected 4 runs for alpha, got %d", r.StatsA.TotalRuns)
	}
	if r.StatsB.TotalRuns != 3 {
		t.Errorf("expected 3 runs for beta, got %d", r.StatsB.TotalRuns)
	}
}

func TestCompareJobs_EmptyPath_ReturnsError(t *testing.T) {
	_, err := CompareJobs("", "alpha", "beta", time.Hour)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCompareJobs_EmptyJobName_ReturnsError(t *testing.T) {
	p := seedCompareStore(t)
	_, err := CompareJobs(p, "", "beta", time.Hour)
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestCompareJobs_ZeroDuration_ReturnsError(t *testing.T) {
	p := seedCompareStore(t)
	_, err := CompareJobs(p, "alpha", "beta", 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestPrintCompare_ContainsHeaders(t *testing.T) {
	p := seedCompareStore(t)
	r, err := CompareJobs(p, "alpha", "beta", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	PrintCompare(&buf, r)
	out := buf.String()
	for _, want := range []string{"alpha", "beta", "Total runs", "Success rate", "Avg duration"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestPrintCompare_NilWriter_UsesStdout(t *testing.T) {
	p := seedCompareStore(t)
	r, err := CompareJobs(p, "alpha", "beta", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// should not panic
	PrintCompare(nil, r)
}
