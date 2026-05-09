package history

import (
	"errors"
	"time"
)

// TrendDirection indicates whether a metric is improving, degrading, or stable.
type TrendDirection string

const (
	TrendImproving TrendDirection = "improving"
	TrendDegrading TrendDirection = "degrading"
	TrendStable    TrendDirection = "stable"
)

// Trend holds trend analysis for a single job.
type Trend struct {
	JobName          string
	DurationTrend    TrendDirection
	SuccessRateTrend TrendDirection
	AvgDurationEarly time.Duration
	AvgDurationLate  time.Duration
	SuccessRateEarly float64
	SuccessRateLate  float64
}

// ComputeTrend compares the first half vs second half of the last n runs for a
// job and returns a Trend describing whether things are getting better or worse.
func ComputeTrend(path, jobName string, n int) (Trend, error) {
	if path == "" {
		return Trend{}, errors.New("history path must not be empty")
	}
	if jobName == "" {
		return Trend{}, errors.New("job name must not be empty")
	}
	if n < 2 {
		return Trend{}, errors.New("n must be at least 2")
	}

	store, err := New(path)
	if err != nil {
		return Trend{}, err
	}

	entries := store.entries[jobName]
	if len(entries) == 0 {
		return Trend{}, errors.New("no entries found for job: " + jobName)
	}

	// Take the last n entries (or all if fewer).
	if len(entries) > n {
		entries = entries[len(entries)-n:]
	}

	mid := len(entries) / 2
	early := entries[:mid]
	late := entries[mid:]

	earlyAvg, earlyRate := halfStats(early)
	lateAvg, lateRate := halfStats(late)

	return Trend{
		JobName:          jobName,
		DurationTrend:    direction(float64(earlyAvg), float64(lateAvg), true),
		SuccessRateTrend: direction(earlyRate, lateRate, false),
		AvgDurationEarly: earlyAvg,
		AvgDurationLate:  lateAvg,
		SuccessRateEarly: earlyRate,
		SuccessRateLate:  lateRate,
	}, nil
}

func halfStats(entries []Entry) (time.Duration, float64) {
	if len(entries) == 0 {
		return 0, 0
	}
	var total time.Duration
	var successes int
	for _, e := range entries {
		total += e.Duration
		if e.Error == "" {
			successes++
		}
	}
	return total / time.Duration(len(entries)), float64(successes) / float64(len(entries)) * 100
}

// direction returns a TrendDirection. lowerIsBetter applies to duration.
func direction(early, late float64, lowerIsBetter bool) TrendDirection {
	const threshold = 0.05 // 5% change required to be non-stable
	if early == 0 {
		return TrendStable
	}
	delta := (late - early) / early
	if delta > threshold {
		if lowerIsBetter {
			return TrendDegrading
		}
		return TrendImproving
	}
	if delta < -threshold {
		if lowerIsBetter {
			return TrendImproving
		}
		return TrendDegrading
	}
	return TrendStable
}
