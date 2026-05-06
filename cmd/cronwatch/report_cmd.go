package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/cronwatch/internal/history"
)

// reportCmd handles the "report" sub-command, printing a summary of all
// recorded job runs from the history store.
func reportCmd(args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	historyPath := fs.String("history", "cronwatch-history.json", "path to history file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("report: %w", err)
	}

	s, err := history.New(*historyPath)
	if err != nil {
		return fmt.Errorf("report: open history: %w", err)
	}

	history.PrintReport(s, os.Stdout)
	return nil
}
