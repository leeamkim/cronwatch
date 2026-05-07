package history

import (
	"bytes"
	"testing"
	"time"
)

func seedDiffStore(t *testing.T) (string, *Store) {
	t.Helper()
	path := tempPath(t)
	store, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return path, store
}

func TestDiff_StatusChanged(t *testing.T) {
	path, store := seedDiffStore(t)
	now := time.Now()

	store.data["backup"] = []Entry{
		{RunID: "r1", Status: "success", Duration: 2 * time.Second, StartedAt: now.Add(-2 * time.Hour)},
		{RunID: "r2", Status: "failure", Duration: 3 * time.Second, StartedAt: now.Add(-1 * time.Hour)},
	}
	if err := store.save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	d, err := Diff(path, "backup")
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if !d.StatusChanged {
		t.Error("expected StatusChanged=true")
	}
	if d.DurationDelta != time.Second {
		t.Errorf("expected DurationDelta=1s, got %s", d.DurationDelta)
	}
}

func TestDiff_NoChange(t *testing.T) {
	path, store := seedDiffStore(t)
	now := time.Now()

	store.data["sync"] = []Entry{
		{RunID: "r1", Status: "success", Duration: 5 * time.Second, StartedAt: now.Add(-2 * time.Hour), Output: "ok"},
		{RunID: "r2", Status: "success", Duration: 5 * time.Second, StartedAt: now.Add(-1 * time.Hour), Output: "ok"},
	}
	if err := store.save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	d, err := Diff(path, "sync")
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if d.StatusChanged {
		t.Error("expected StatusChanged=false")
	}
	if d.OutputChanged {
		t.Error("expected OutputChanged=false")
	}
}

func TestDiff_InsufficientRuns_ReturnsError(t *testing.T) {
	path, store := seedDiffStore(t)
	store.data["lonely"] = []Entry{
		{RunID: "r1", Status: "success"},
	}
	if err := store.save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	_, err := Diff(path, "lonely")
	if err == nil {
		t.Fatal("expected error for fewer than 2 runs")
	}
}

func TestDiff_EmptyPath_ReturnsError(t *testing.T) {
	_, err := Diff("", "job")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestDiff_EmptyJobName_ReturnsError(t *testing.T) {
	path, _ := seedDiffStore(t)
	_, err := Diff(path, "")
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestPrintDiff_ContainsJobName(t *testing.T) {
	d := &DiffResult{
		JobName:       "nightly",
		RunA:          Entry{Status: "success", Duration: time.Second, StartedAt: time.Now().Add(-2 * time.Hour)},
		RunB:          Entry{Status: "failure", Duration: 2 * time.Second, StartedAt: time.Now().Add(-time.Hour)},
		StatusChanged: true,
		DurationDelta: time.Second,
	}
	var buf bytes.Buffer
	PrintDiff(&buf, d)
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("nightly")) {
		t.Errorf("expected job name in output, got: %s", out)
	}
	if !bytes.Contains([]byte(out), []byte("Status changed")) {
		t.Errorf("expected status change notice in output, got: %s", out)
	}
}
