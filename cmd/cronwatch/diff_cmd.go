package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/cronwatch/internal/history"
)

func diffCmd(args []string) error {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	histPath := fs.String("history", defaultHistoryPath(), "path to history file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: cronwatch diff [flags] <job-name>")
	}

	jobName := fs.Arg(0)

	d, err := history.Diff(*histPath, jobName)
	if err != nil {
		return fmt.Errorf("diff: %w", err)
	}

	history.PrintDiff(os.Stdout, d)
	return nil
}
