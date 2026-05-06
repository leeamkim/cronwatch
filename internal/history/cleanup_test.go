package history

import (
	"testing"
	"time"
)

func TestCleanup_RemovesOldEntries(t *testing.T) {
	path := tempPath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now()
	old := now.Add(-60 * 24 * time.Hour) // 60 days ago

	s.mu.Lock()
	s.data["backup"] = []Entry{
		{StartedAt: old, Duration: time.Second, ExitCode: 0},
		{StartedAt: now, Duration: time.Second, ExitCode: 0},
	}
	s.mu.Unlock()

	removed, err := s.Cleanup(DefaultRetention)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}

	entries := s.Get("backup")
	if len(entries) != 1 {
		t.Errorf("expected 1 remaining entry, got %d", len(entries))
	}
}

func TestCleanup_EnforcesMaxRuns(t *testing.T) {
	path := tempPath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now()
	s.mu.Lock()
	for i := 0; i < 10; i++ {
		s.data["sync"] = append(s.data["sync"], Entry{
			StartedAt: now.Add(-time.Duration(i) * time.Minute),
			Duration:  time.Second,
			ExitCode:  0,
		})
	}
	s.mu.Unlock()

	retention := Retention{MaxAge: 30 * 24 * time.Hour, MaxRuns: 5}
	removed, err := s.Cleanup(retention)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 5 {
		t.Errorf("expected 5 removed, got %d", removed)
	}
	if got := len(s.Get("sync")); got != 5 {
		t.Errorf("expected 5 remaining entries, got %d", got)
	}
}

func TestCleanup_RemovesJobKeyWhenAllEntriesExpire(t *testing.T) {
	path := tempPath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	old := time.Now().Add(-60 * 24 * time.Hour)
	s.mu.Lock()
	s.data["stale"] = []Entry{
		{StartedAt: old, Duration: time.Second, ExitCode: 0},
	}
	s.mu.Unlock()

	_, err = s.Cleanup(DefaultRetention)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if entries := s.Get("stale"); len(entries) != 0 {
		t.Errorf("expected job key to be removed, got %d entries", len(entries))
	}
}
