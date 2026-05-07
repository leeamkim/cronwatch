package history

import (
	"testing"
	"time"
)

func TestPrune_RemovesOldEntries(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now()
	old := now.Add(-48 * time.Hour)

	h.store.Runs["backup"] = []Entry{
		{StartedAt: old, Success: true},
		{StartedAt: now, Success: true},
	}

	res, err := h.Prune(PruneOptions{OlderThan: 24 * time.Hour})
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if res.Removed != 1 {
		t.Errorf("Removed = %d, want 1", res.Removed)
	}
	if res.Retained != 1 {
		t.Errorf("Retained = %d, want 1", res.Retained)
	}
	if len(h.store.Runs["backup"]) != 1 {
		t.Errorf("entries after prune = %d, want 1", len(h.store.Runs["backup"]))
	}
}

func TestPrune_DryRun_DoesNotModify(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	old := time.Now().Add(-72 * time.Hour)
	h.store.Runs["sync"] = []Entry{
		{StartedAt: old, Success: false},
	}

	res, err := h.Prune(PruneOptions{OlderThan: 24 * time.Hour, DryRun: true})
	if err != nil {
		t.Fatalf("Prune dry run: %v", err)
	}
	if res.Removed != 1 {
		t.Errorf("Removed = %d, want 1", res.Removed)
	}
	if len(h.store.Runs["sync"]) != 1 {
		t.Error("dry run should not remove entries")
	}
}

func TestPrune_SpecificJob_OnlyAffectsThatJob(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	old := time.Now().Add(-48 * time.Hour)
	h.store.Runs["jobA"] = []Entry{{StartedAt: old, Success: true}}
	h.store.Runs["jobB"] = []Entry{{StartedAt: old, Success: true}}

	_, err = h.Prune(PruneOptions{OlderThan: 24 * time.Hour, JobName: "jobA"})
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if _, ok := h.store.Runs["jobA"]; ok {
		t.Error("jobA should have been pruned")
	}
	if len(h.store.Runs["jobB"]) != 1 {
		t.Error("jobB should be untouched")
	}
}

func TestPrune_ZeroDuration_ReturnsError(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = h.Prune(PruneOptions{OlderThan: 0})
	if err == nil {
		t.Error("expected error for zero OlderThan")
	}
}
