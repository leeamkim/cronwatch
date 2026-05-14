package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func seedHeatmap(t *testing.T, entries map[string][]Entry) string {
	t.Helper()
	p := tempPath(t)
	store := make(store)
	for job, ee := range entries {
		store[job] = ee
	}
	if err := save(p, store); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return p
}

func TestComputeHeatmap_CountsDays(t *testing.T) {
	now := time.Now().UTC()
	p := seedHeatmap(t, map[string][]Entry{
		"backup": {
			{RunID: "1", StartedAt: now, Error: ""},
			{RunID: "2", StartedAt: now, Error: "fail"},
			{RunID: "3", StartedAt: now.AddDate(0, 0, -1), Error: ""},
		},
	})

	hm, err := ComputeHeatmap(p, "backup", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	today := now.Format("2006-01-02")
	if hm[today].Success != 1 {
		t.Errorf("expected 1 success today, got %d", hm[today].Success)
	}
	if hm[today].Failure != 1 {
		t.Errorf("expected 1 failure today, got %d", hm[today].Failure)
	}
	if len(hm) != 2 {
		t.Errorf("expected 2 days, got %d", len(hm))
	}
}

func TestComputeHeatmap_AllJobs(t *testing.T) {
	now := time.Now().UTC()
	p := seedHeatmap(t, map[string][]Entry{
		"job-a": {{RunID: "1", StartedAt: now, Error: ""}},
		"job-b": {{RunID: "2", StartedAt: now, Error: "err"}},
	})

	hm, err := ComputeHeatmap(p, "", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	today := now.Format("2006-01-02")
	if hm[today].Success+hm[today].Failure != 2 {
		t.Errorf("expected 2 total entries, got %d", hm[today].Success+hm[today].Failure)
	}
}

func TestComputeHeatmap_ExcludesOldEntries(t *testing.T) {
	now := time.Now().UTC()
	p := seedHeatmap(t, map[string][]Entry{
		"job": {
			{RunID: "old", StartedAt: now.AddDate(0, 0, -30), Error: ""},
			{RunID: "new", StartedAt: now, Error: ""},
		},
	})

	hm, err := ComputeHeatmap(p, "job", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hm) != 1 {
		t.Errorf("expected 1 day, got %d", len(hm))
	}
}

func TestComputeHeatmap_EmptyPath_ReturnsError(t *testing.T) {
	_, err := ComputeHeatmap("", "job", 7)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestComputeHeatmap_ZeroDays_ReturnsError(t *testing.T) {
	_, err := ComputeHeatmap("/tmp/x.json", "job", 0)
	if err == nil {
		t.Fatal("expected error for zero days")
	}
}

func TestPrintHeatmap_ContainsHeaders(t *testing.T) {
	hm := Heatmap{
		"2024-01-01": {Date: "2024-01-01", Success: 3, Failure: 1},
	}
	var buf bytes.Buffer
	PrintHeatmap(hm, &buf)
	out := buf.String()
	for _, h := range []string{"Date", "Success", "Failure", "Total"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestPrintHeatmap_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic.
	PrintHeatmap(Heatmap{}, nil)
}
