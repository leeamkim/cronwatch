package history

import (
	"testing"
	"time"
)

func TestCleanup_RemovesOldEntries(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	old := time.Now().Add(-48 * time.Hour)
	h.store.Jobs["backup"] = []Entry{
		{StartedAt: old, Success: true},
		{StartedAt: time.Now(), Success: true},
	}

	if err := h.Cleanup(CleanupOptions{MaxAge: 24 * time.Hour}); err != nil {
		t.Fatalf("Cleanup: %v", err)
	}

	if got := len(h.store.Jobs["backup"]); got != 1 {
		t.Errorf("expected 1 entry after cleanup, got %d", got)
	}
}

func TestCleanup_EnforcesMaxRuns(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now()
	entries := []Entry{
		{StartedAt: now.Add(-3 * time.Hour), Success: true},
		{StartedAt: now.Add(-2 * time.Hour), Success: false},
		{StartedAt: now.Add(-1 * time.Hour), Success: true},
		{StartedAt: now, Success: true},
	}
	h.store.Jobs["sync"] = entries

	if err := h.Cleanup(CleanupOptions{MaxRunsPerJob: 2}); err != nil {
		t.Fatalf("Cleanup: %v", err)
	}

	got := h.store.Jobs["sync"]
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	// Should retain the two most recent.
	if !got[0].StartedAt.Equal(now.Add(-1 * time.Hour)) {
		t.Errorf("unexpected oldest retained entry: %v", got[0].StartedAt)
	}
}

func TestCleanup_RemovesJobKeyWhenAllEntriesExpire(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	old := time.Now().Add(-72 * time.Hour)
	h.store.Jobs["prune-me"] = []Entry{
		{StartedAt: old, Success: true},
	}

	if err := h.Cleanup(CleanupOptions{MaxAge: 24 * time.Hour}); err != nil {
		t.Fatalf("Cleanup: %v", err)
	}

	if _, exists := h.store.Jobs["prune-me"]; exists {
		t.Error("expected job key to be removed when all entries expired")
	}
}
