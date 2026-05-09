package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/cronwatch/internal/history"
)

func alertCmd(args []string) error {
	fs := flag.NewFlagSet("alert", flag.ContinueOnError)
	historyPath := fs.String("history", defaultHistoryPath(), "path to history file")
	configPath := fs.String("config", "", "path to alert rules YAML config (required)")
	quiet := fs.Bool("quiet", false, "suppress output, use exit code only")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *configPath == "" {
		return fmt.Errorf("--config is required")
	}

	rules, err := history.LoadAlertRules(*configPath)
	if err != nil {
		return fmt.Errorf("load alert rules: %w", err)
	}

	breaches, err := history.CheckAlerts(*historyPath, rules)
	if err != nil {
		return fmt.Errorf("check alerts: %w", err)
	}

	if len(breaches) == 0 {
		if !*quiet {
			fmt.Println("No alert thresholds breached.")
		}
		return nil
	}

	if !*quiet {
		fmt.Fprintf(os.Stderr, "ALERT: %d threshold%s breached:\n", len(breaches), pluralSuffix(len(breaches)))
		for _, b := range breaches {
			fmt.Fprintf(os.Stderr, "  [%s] %s\n", b.JobName, b.Reason)
		}
	}

	os.Exit(1)
	return nil
}
