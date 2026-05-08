package history

import (
	"testing"
	"time"
)

func seedStreak(t *testing.T, entries []Entry) string {
	t.Helper()
	path := tempPath(t)
	store, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	store.Entries["myjob"] = entries
	if err := store.save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func makeEntry(status string, offset time.Duration) Entry {
	now := time.Now().Add(offset)
	return Entry{
		RunID:     now.Format(time.RFC3339Nano),
		StartedAt: now.Format(time.RFC3339),
		Status:    status,
	}
}

func TestComputeStreak_AllSuccess(t *testing.T) {
	path := seedStreak(t, []Entry{
		makeEntry("success", -3*time.Minute),
		makeEntry("success", -2*time.Minute),
		makeEntry("success", -1*time.Minute),
	})
	s, err := ComputeStreak(path, "myjob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count != 3 || s.Status != "success" {
		t.Errorf("expected 3 successes, got count=%d status=%s", s.Count, s.Status)
	}
}

func TestComputeStreak_BreaksOnFailure(t *testing.T) {
	path := seedStreak(t, []Entry{
		makeEntry("success", -4*time.Minute),
		makeEntry("failure", -3*time.Minute),
		makeEntry("failure", -2*time.Minute),
		makeEntry("failure", -1*time.Minute),
	})
	s, err := ComputeStreak(path, "myjob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count != 3 || s.Status != "failure" {
		t.Errorf("expected 3 failures, got count=%d status=%s", s.Count, s.Status)
	}
}

func TestComputeStreak_SingleEntry(t *testing.T) {
	path := seedStreak(t, []Entry{
		makeEntry("success", -1*time.Minute),
	})
	s, err := ComputeStreak(path, "myjob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count != 1 {
		t.Errorf("expected count 1, got %d", s.Count)
	}
}

func TestComputeStreak_NoEntries_ReturnsUnknown(t *testing.T) {
	path := seedStreak(t, []Entry{})
	s, err := ComputeStreak(path, "myjob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Status != "unknown" || s.Count != 0 {
		t.Errorf("expected unknown/0, got %s/%d", s.Status, s.Count)
	}
}

func TestComputeStreak_EmptyPath_ReturnsError(t *testing.T) {
	_, err := ComputeStreak("", "myjob")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestComputeStreak_EmptyJobName_ReturnsError(t *testing.T) {
	_, err := ComputeStreak("/tmp/x.json", "")
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}
