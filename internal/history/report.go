package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Summary holds aggregated stats for a single job.
type Summary struct {
	JobName    string
	TotalRuns  int
	Failures   int
	Timeouts   int
	LastRun    time.Time
	LastStatus string
}

// Summarise returns a Summary for each job found in the store.
func Summarise(s *Store) []Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var summaries []Summary
	for job, entries := range s.data {
		sum := Summary{JobName: job, TotalRuns: len(entries)}
		for _, e := range entries {
			if e.Error != "" {
				sum.Failures++
			}
			if e.TimedOut {
				sum.Timeouts++
			}
			if e.FinishedAt.After(sum.LastRun) {
				sum.LastRun = e.FinishedAt
				if e.Error != "" {
					sum.LastStatus = "FAILED"
				} else if e.TimedOut {
					sum.LastStatus = "TIMEOUT"
				} else {
					sum.LastStatus = "OK"
				}
			}
		}
		summaries = append(summaries, sum)
	}
	return summaries
}

// PrintReport writes a human-readable report table to w.
// If w is nil, os.Stdout is used.
func PrintReport(s *Store, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "JOB\tTOTAL\tFAILURES\tTIMEOUTS\tLAST RUN\tLAST STATUS")
	fmt.Fprintln(tw, "---\t-----\t--------\t--------\t--------\t-----------")
	for _, sum := range Summarise(s) {
		lastRun := "-"
		if !sum.LastRun.IsZero() {
			lastRun = sum.LastRun.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%s\t%s\n",
			sum.JobName, sum.TotalRuns, sum.Failures, sum.Timeouts, lastRun, sum.LastStatus)
	}
	tw.Flush()
}
