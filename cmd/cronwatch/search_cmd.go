package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/cronwatch/internal/history"
)

func searchCmd(args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	jobName := fs.String("job", "", "filter by job name")
	status := fs.String("status", "", "filter by status: success or failure")
	sinceDur := fs.Duration("since", 0, "only show entries newer than this duration (e.g. 24h)")
	maxResults := fs.Int("n", 20, "maximum number of results to display")
	histPath := fs.String("history", defaultHistoryPath(), "path to history file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *status != "" && *status != "success" && *status != "failure" {
		return fmt.Errorf("--status must be \"success\" or \"failure\"")
	}

	h, err := history.New(*histPath)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}

	store := h.All()

	opts := history.SearchOptions{
		JobName:    *jobName,
		Status:     *status,
		MaxResults: *maxResults,
	}
	if *sinceDur > 0 {
		opts.Since = time.Now().Add(-*sinceDur)
	}

	results := history.Search(store, opts)

	if len(results) == 0 {
		fmt.Fprintln(os.Stdout, "no matching entries found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-30s %10s %s\n", "JOB", "STARTED", "DURATION", "STATUS")
	for _, e := range results {
		status := "ok"
		if e.Error != "" {
			status = "FAIL: " + e.Error
		}
		fmt.Fprintf(os.Stdout, "%-20s %-30s %9dms %s\n",
			e.JobName,
			e.StartedAt.Format(time.RFC3339),
			e.DurationMs,
			status,
		)
	}
	return nil
}
