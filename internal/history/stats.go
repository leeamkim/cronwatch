package history

import (
	"fmt"
	"time"
)

// JobStats holds aggregated statistics for a single job.
type JobStats struct {
	JobName      string
	TotalRuns    int
	SuccessCount int
	FailureCount int
	AvgDuration  time.Duration
	MaxDuration  time.Duration
	LastRun      time.Time
	LastStatus   string
}

// SuccessRate returns the percentage of successful runs.
func (s JobStats) SuccessRate() float64 {
	if s.TotalRuns == 0 {
		return 0
	}
	return float64(s.SuccessCount) / float64(s.TotalRuns) * 100
}

// String returns a human-readable summary of the job stats.
func (s JobStats) String() string {
	return fmt.Sprintf(
		"%s: runs=%d success=%.1f%% avg=%s max=%s last=%s (%s)",
		s.JobName,
		s.TotalRuns,
		s.SuccessRate(),
		s.AvgDuration.Round(time.Millisecond),
		s.MaxDuration.Round(time.Millisecond),
		s.LastRun.Format(time.RFC3339),
		s.LastStatus,
	)
}

// ComputeStats returns aggregated statistics for each job in the store.
func ComputeStats(s *Store) map[string]JobStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]JobStats, len(s.data))

	for job, entries := range s.data {
		if len(entries) == 0 {
			continue
		}

		stats := JobStats{JobName: job}
		var totalDuration time.Duration

		for _, e := range entries {
			stats.TotalRuns++
			if e.Error == "" {
				stats.SuccessCount++
			} else {
				stats.FailureCount++
			}
			totalDuration += e.Duration
			if e.Duration > stats.MaxDuration {
				stats.MaxDuration = e.Duration
			}
			if e.StartedAt.After(stats.LastRun) {
				stats.LastRun = e.StartedAt
				if e.Error == "" {
					stats.LastStatus = "ok"
				} else {
					stats.LastStatus = "failed"
				}
			}
		}

		if stats.TotalRuns > 0 {
			stats.AvgDuration = totalDuration / time.Duration(stats.TotalRuns)
		}

		result[job] = stats
	}

	return result
}
