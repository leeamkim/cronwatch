package history

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// CompareResult holds the comparison between two jobs over a shared time window.
type CompareResult struct {
	JobA        string
	JobB        string
	WindowStart time.Time
	WindowEnd   time.Time
	StatsA      JobStats
	StatsB      JobStats
}

// CompareJobs computes stats for two jobs over the same time window and returns
// a CompareResult for side-by-side analysis.
func CompareJobs(path, jobA, jobB string, since time.Duration) (CompareResult, error) {
	if path == "" {
		return CompareResult{}, errors.New("history path must not be empty")
	}
	if jobA == "" || jobB == "" {
		return CompareResult{}, errors.New("both job names must be non-empty")
	}
	if since <= 0 {
		return CompareResult{}, errors.New("since duration must be positive")
	}

	statsA, err := ComputeStats(path, jobA)
	if err != nil {
		return CompareResult{}, fmt.Errorf("stats for %q: %w", jobA, err)
	}
	statsB, err := ComputeStats(path, jobB)
	if err != nil {
		return CompareResult{}, fmt.Errorf("stats for %q: %w", jobB, err)
	}

	now := time.Now()
	return CompareResult{
		JobA:        jobA,
		JobB:        jobB,
		WindowStart: now.Add(-since),
		WindowEnd:   now,
		StatsA:      statsA,
		StatsB:      statsB,
	}, nil
}

// PrintCompare writes a formatted side-by-side comparison to w.
// If w is nil, os.Stdout is used.
func PrintCompare(w io.Writer, r CompareResult) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Job Comparison  [%s  →  %s]\n",
		r.WindowStart.Format(time.RFC3339), r.WindowEnd.Format(time.RFC3339))
	fmt.Fprintf(w, "%-22s  %-20s  %-20s\n", "Metric", r.JobA, r.JobB)
	fmt.Fprintf(w, "%s\n", repeatChar('-', 66))
	fmt.Fprintf(w, "%-22s  %-20d  %-20d\n", "Total runs", r.StatsA.TotalRuns, r.StatsB.TotalRuns)
	fmt.Fprintf(w, "%-22s  %-20d  %-20d\n", "Successes", r.StatsA.Successes, r.StatsB.Successes)
	fmt.Fprintf(w, "%-22s  %-20d  %-20d\n", "Failures", r.StatsA.Failures, r.StatsB.Failures)
	fmt.Fprintf(w, "%-22s  %-20.1f  %-20.1f\n", "Success rate (%)", r.StatsA.SuccessRate, r.StatsB.SuccessRate)
	fmt.Fprintf(w, "%-22s  %-20s  %-20s\n", "Avg duration", r.StatsA.AvgDuration.Round(time.Millisecond), r.StatsB.AvgDuration.Round(time.Millisecond))
	fmt.Fprintf(w, "%-22s  %-20s  %-20s\n", "Max duration", r.StatsA.MaxDuration.Round(time.Millisecond), r.StatsB.MaxDuration.Round(time.Millisecond))
	fmt.Fprintf(w, "%-22s  %-20s  %-20s\n", "Last status", r.StatsA.LastStatus, r.StatsB.LastStatus)
}
