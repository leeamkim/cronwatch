package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/example/cronwatch/internal/history"
)

func pruneCmd(args []string) error {
	fs := flag.NewFlagSet("prune", flag.ContinueOnError)
	histPath := fs.String("history", defaultHistoryPath(), "path to history file")
	older := fs.Duration("older-than", 30*24*time.Hour, "remove entries older than this duration (e.g. 168h)")
	jobName := fs.String("job", "", "restrict pruning to a specific job name")
	dryRun := fs.Bool("dry-run", false, "report what would be removed without modifying the store")

	if err := fs.Parse(args); err != nil {
		return err
	}

	h, err := history.New(*histPath)
	if err != nil {
		return fmt.Errorf("prune: open history: %w", err)
	}

	opts := history.PruneOptions{
		OlderThan: *older,
		JobName:   *jobName,
		DryRun:    *dryRun,
	}

	result, err := h.Prune(opts)
	if err != nil {
		return fmt.Errorf("prune: %w", err)
	}

	w := os.Stdout
	if *dryRun {
		fmt.Fprintf(w, "[dry-run] would remove %d entr%s, retain %d\n",
			result.Removed, pluralSuffix(result.Removed), result.Retained)
	} else {
		fmt.Fprintf(w, "pruned %d entr%s, retained %d\n",
			result.Removed, pluralSuffix(result.Removed), result.Retained)
	}
	return nil
}

func pluralSuffix(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
