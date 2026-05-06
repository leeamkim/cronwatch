package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Summary holds aggregated statistics for a single job.
type Summary struct {
	JobName     string
	Total       int
	Failures    int
	Timeouts    int
	AvgDuration time.Duration
}

// Summarise computes per-job statistics from the provided entries.
func Summarise(entries []Entry) []Summary {
	type agg struct {
		total, failures, timeouts int
		totalDur                  time.Duration
	}
	m := make(map[string]*agg)
	for _, e := range entries {
		a, ok := m[e.JobName]
		if !ok {
			a = &agg{}
			m[e.JobName] = a
		}
		a.total++
		a.totalDur += e.Duration
		if !e.Success {
			a.failures++
		}
		if e.TimedOut {
			a.timeouts++
		}
	}

	summaries := make([]Summary, 0, len(m))
	for name, a := range m {
		avg := time.Duration(0)
		if a.total > 0 {
			avg = a.totalDur / time.Duration(a.total)
		}
		summaries = append(summaries, Summary{
			JobName:     name,
			Total:       a.total,
			Failures:    a.failures,
			Timeouts:    a.timeouts,
			AvgDuration: avg,
		})
	}
	return summaries
}

// PrintReport writes a formatted summary table to w.
func PrintReport(w io.Writer, entries []Entry) error {
	summaries := Summarise(entries)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "JOB\tRUNS\tFAILURES\tTIMEOUTS\tAVG DURATION")
	for _, s := range summaries {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%s\n",
			s.JobName, s.Total, s.Failures, s.Timeouts, s.AvgDuration.Round(time.Millisecond))
	}
	return tw.Flush()
}
