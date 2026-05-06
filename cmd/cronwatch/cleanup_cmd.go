package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/cronwatch/internal/history"
)

// cleanupCmd parses flags and prunes the history store.
func cleanupCmd(args []string) error {
	fs := flag.NewFlagSet("cleanup", flag.ContinueOnError)
	maxAgeDays := fs.Int("max-age-days", 30, "remove entries older than this many days (0 = disabled)")
	maxRuns := fs.Int("max-runs", 100, "keep at most this many recent runs per job (0 = unlimited)")
	historyPath := fs.String("history", defaultHistoryPath(), "path to history file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	h, err := history.New(*historyPath)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}

	opts := history.CleanupOptions{
		MaxRunsPerJob: *maxRuns,
	}
	if *maxAgeDays > 0 {
		opts.MaxAge = time.Duration(*maxAgeDays) * 24 * time.Hour
	}

	if err := h.Cleanup(opts); err != nil {
		return fmt.Errorf("cleanup: %w", err)
	}

	fmt.Fprintf(os.Stdout, "history cleaned (max-age-days=%d, max-runs=%d)\n", *maxAgeDays, *maxRuns)
	return nil
}

func defaultHistoryPath() string {
	if p := os.Getenv("CRONWATCH_HISTORY"); p != "" {
		return p
	}
	return "/var/lib/cronwatch/history.json"
}
