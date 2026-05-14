package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/cronwatch/internal/history"
)

func heatmapCmd(args []string) error {
	fs := flag.NewFlagSet("heatmap", flag.ContinueOnError)
	histPath := fs.String("history", defaultHistoryPath(), "path to history file")
	job := fs.String("job", "", "job name to filter (empty = all jobs)")
	days := fs.Int("days", 30, "number of past days to include")
	jsonOut := fs.Bool("json", false, "output as JSON instead of table")

	if err := fs.Parse(args); err != nil {
		return err
	}

	hm, err := history.ComputeHeatmap(*histPath, *job, *days)
	if err != nil {
		return fmt.Errorf("heatmap: %w", err)
	}

	if len(hm) == 0 {
		fmt.Fprintln(os.Stdout, "no data found for the given filters")
		return nil
	}

	if *jsonOut {
		return history.ExportHeatmapJSON(hm, os.Stdout)
	}

	history.PrintHeatmap(hm, os.Stdout)
	return nil
}
