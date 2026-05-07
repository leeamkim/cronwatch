package history

import (
	"testing"
	"time"
)

func seedSearch() map[string][]Entry {
	now := time.Now()
	return map[string][]Entry{
		"backup": {
			{JobName: "backup", StartedAt: now.Add(-2 * time.Hour), DurationMs: 100, Error: ""},
			{JobName: "backup", StartedAt: now.Add(-1 * time.Hour), DurationMs: 120, Error: "timeout"},
		},
		"sync": {
			{JobName: "sync", StartedAt: now.Add(-30 * time.Minute), DurationMs: 50, Error: ""},
		},
	}
}

func TestSearch_NoFilters_ReturnsAll(t *testing.T) {
	store := seedSearch()
	results := Search(store, SearchOptions{})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestSearch_FilterByJobName(t *testing.T) {
	store := seedSearch()
	results := Search(store, SearchOptions{JobName: "backup"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.JobName != "backup" {
			t.Errorf("unexpected job name: %s", r.JobName)
		}
	}
}

func TestSearch_FilterByStatus_Failure(t *testing.T) {
	store := seedSearch()
	results := Search(store, SearchOptions{Status: "failure"})
	if len(results) != 1 {
		t.Fatalf("expected 1 failure, got %d", len(results))
	}
	if results[0].Error == "" {
		t.Error("expected entry with error")
	}
}

func TestSearch_FilterBySince(t *testing.T) {
	store := seedSearch()
	cutoff := time.Now().Add(-90 * time.Minute)
	results := Search(store, SearchOptions{Since: cutoff})
	if len(results) != 2 {
		t.Fatalf("expected 2 results after cutoff, got %d", len(results))
	}
}

func TestSearch_MaxResults(t *testing.T) {
	store := seedSearch()
	results := Search(store, SearchOptions{MaxResults: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ResultsOrderedByStartedAtDesc(t *testing.T) {
	store := seedSearch()
	results := Search(store, SearchOptions{})
	for i := 1; i < len(results); i++ {
		if results[i].StartedAt.After(results[i-1].StartedAt) {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}
