package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"cronwatch/internal/history"
)

func compareCmd() *cobra.Command {
	var (
		historyPath string
		sinceStr    string
	)

	cmd := &cobra.Command{
		Use:   "compare <jobA> <jobB>",
		Short: "Compare statistics for two jobs side by side",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if historyPath == "" {
				historyPath = defaultHistoryPath()
			}

			since, err := time.ParseDuration(sinceStr)
			if err != nil {
				return fmt.Errorf("invalid --since value %q: %w", sinceStr, err)
			}

			result, err := history.CompareJobs(historyPath, args[0], args[1], since)
			if err != nil {
				return fmt.Errorf("compare: %w", err)
			}

			history.PrintCompare(os.Stdout, result)
			return nil
		},
	}

	cmd.Flags().StringVar(&historyPath, "history", "", "path to history file (default: ~/.cronwatch/history.json)")
	cmd.Flags().StringVar(&sinceStr, "since", "168h", "time window to compare over (e.g. 24h, 7d expressed as hours)")

	return cmd
}
