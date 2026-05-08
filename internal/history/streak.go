package history

import (
	"fmt"
	"sort"
)

// Streak describes a consecutive run of successes or failures for a job.
type Streak struct {
	JobName string
	Status  string // "success" or "failure"
	Count   int
	Latest  string // RFC3339 timestamp of the most recent run in the streak
}

// ComputeStreak returns the current consecutive streak (success or failure)
// for the given job by inspecting its recorded history entries.
func ComputeStreak(path, jobName string) (Streak, error) {
	if path == "" {
		return Streak{}, fmt.Errorf("history path must not be empty")
	}
	if jobName == "" {
		return Streak{}, fmt.Errorf("job name must not be empty")
	}

	store, err := New(path)
	if err != nil {
		return Streak{}, fmt.Errorf("open store: %w", err)
	}

	entries := store.Entries[jobName]
	if len(entries) == 0 {
		return Streak{JobName: jobName, Status: "unknown", Count: 0}, nil
	}

	// Sort entries ascending by timestamp so the most recent is last.
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].StartedAt < sorted[j].StartedAt
	})

	latest := sorted[len(sorted)-1]
	currentStatus := latest.Status
	count := 0

	for i := len(sorted) - 1; i >= 0; i-- {
		if sorted[i].Status == currentStatus {
			count++
		} else {
			break
		}
	}

	return Streak{
		JobName: jobName,
		Status:  currentStatus,
		Count:   count,
		Latest:  latest.StartedAt,
	}, nil
}
