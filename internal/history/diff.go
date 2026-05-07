package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// DiffResult holds the comparison between two runs of the same job.
type DiffResult struct {
	JobName     string
	RunA        Entry
	RunB        Entry
	StatusChanged bool
	DurationDelta time.Duration
	OutputChanged bool
}

// Diff compares the two most recent runs of a job and returns a DiffResult.
func Diff(path, jobName string) (*DiffResult, error) {
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
	if len(entries) < 2 {
		return nil, fmt.Errorf("job %q has fewer than 2 recorded runs", jobName)
	}

	a := entries[len(entries)-2]
	b := entries[len(entries)-1]

	return &DiffResult{
		JobName:       jobName,
		RunA:          a,
		RunB:          b,
		StatusChanged: a.Status != b.Status,
		DurationDelta: b.Duration - a.Duration,
		OutputChanged: a.Output != b.Output,
	}, nil
}

// PrintDiff writes a human-readable diff summary to w.
// If w is nil, os.Stdout is used.
func PrintDiff(w io.Writer, d *DiffResult) {
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, "Diff for job: %s\n", d.JobName)
	fmt.Fprintf(w, "  Run A (%s): status=%s duration=%s\n",
		d.RunA.StartedAt.Format(time.RFC3339), d.RunA.Status, d.RunA.Duration)
	fmt.Fprintf(w, "  Run B (%s): status=%s duration=%s\n",
		d.RunB.StartedAt.Format(time.RFC3339), d.RunB.Status, d.RunB.Duration)

	if d.StatusChanged {
		fmt.Fprintf(w, "  [!] Status changed: %s -> %s\n", d.RunA.Status, d.RunB.Status)
	} else {
		fmt.Fprintf(w, "  Status unchanged: %s\n", d.RunA.Status)
	}

	if d.DurationDelta != 0 {
		fmt.Fprintf(w, "  Duration delta: %+.3fs\n", d.DurationDelta.Seconds())
	}

	if d.OutputChanged {
		fmt.Fprintln(w, "  [!] Output changed between runs")
	}
}
