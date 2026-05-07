package history

import (
	"testing"
	"time"
)

func seedTagStore() map[string][]Entry {
	now := time.Now()
	return map[string][]Entry{
		"backup": {
			{Job: "backup", StartedAt: now.Add(-2 * time.Hour), Success: true, Tags: []string{"daily", "storage"}},
			{Job: "backup", StartedAt: now.Add(-1 * time.Hour), Success: false, Tags: []string{"daily"}},
		},
		"sync": {
			{Job: "sync", StartedAt: now.Add(-30 * time.Minute), Success: true, Tags: []string{"storage", "remote"}},
		},
		"report": {
			{Job: "report", StartedAt: now.Add(-10 * time.Minute), Success: true, Tags: []string{"weekly"}},
		},
	}
}

func TestTagFilter_MatchesSingleTag(t *testing.T) {
	store := seedTagStore()
	results, err := TagFilter(store, []string{"storage"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 entries, got %d", len(results))
	}
}

func TestTagFilter_MatchesMultipleTags(t *testing.T) {
	store := seedTagStore()
	results, err := TagFilter(store, []string{"daily", "weekly"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 entries, got %d", len(results))
	}
}

func TestTagFilter_NoMatchReturnsEmpty(t *testing.T) {
	store := seedTagStore()
	results, err := TagFilter(store, []string{"nonexistent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 entries, got %d", len(results))
	}
}

func TestTagFilter_EmptyTags_ReturnsError(t *testing.T) {
	store := seedTagStore()
	_, err := TagFilter(store, []string{})
	if err == nil {
		t.Error("expected error for empty tags, got nil")
	}
}

func TestTagSummary_CountsCorrectly(t *testing.T) {
	store := seedTagStore()
	summary := TagSummary(store)

	if summary["daily"] != 2 {
		t.Errorf("expected daily=2, got %d", summary["daily"])
	}
	if summary["storage"] != 2 {
		t.Errorf("expected storage=2, got %d", summary["storage"])
	}
	if summary["weekly"] != 1 {
		t.Errorf("expected weekly=1, got %d", summary["weekly"])
	}
	if summary["remote"] != 1 {
		t.Errorf("expected remote=1, got %d", summary["remote"])
	}
}

func TestTagSummary_EmptyStore(t *testing.T) {
	summary := TagSummary(map[string][]Entry{})
	if len(summary) != 0 {
		t.Errorf("expected empty summary, got %v", summary)
	}
}
