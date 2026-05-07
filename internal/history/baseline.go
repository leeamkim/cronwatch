package history

import (
	"fmt"
	"time"
)

// Baseline holds the reference metrics for a job derived from its history.
type Baseline struct {
	JobName        string
	AvgDuration    time.Duration
	SuccessRate    float64
	SampleSize     int
}

// ComputeBaseline calculates a baseline from the last n successful runs of a job.
// If n <= 0, all available runs are used.
func ComputeBaseline(path, jobName string, n int) (*Baseline, error) {
	if path == "" {
		return nil, fmt.Errorf("history path must not be empty")
	}
	if jobName == "" {
		return nil, fmt.Errorf("job name must not be empty")
	}

	store, err := New(path)
	if err != nil {
		return nil, fmt.Errorf("open history: %w", err)
	}

	entries := store.data[jobName]
	if len(entries) == 0 {
		return nil, fmt.Errorf("no history found for job %q", jobName)
	}

	if n > 0 && len(entries) > n {
		entries = entries[len(entries)-n:]
	}

	var totalDur time.Duration
	var successes int
	for _, e := range entries {
		totalDur += e.Duration
		if e.Status == "success" {
			successes++
		}
	}

	count := len(entries)
	return &Baseline{
		JobName:     jobName,
		AvgDuration: time.Duration(int64(totalDur) / int64(count)),
		SuccessRate: float64(successes) / float64(count),
		SampleSize:  count,
	}, nil
}

// IsAnomalous returns true when a run's duration exceeds the baseline average
// by more than the given threshold multiplier (e.g. 2.0 = twice as long).
func (b *Baseline) IsAnomalous(d time.Duration, threshold float64) bool {
	if b.AvgDuration == 0 || threshold <= 0 {
		return false
	}
	return float64(d) > float64(b.AvgDuration)*threshold
}
