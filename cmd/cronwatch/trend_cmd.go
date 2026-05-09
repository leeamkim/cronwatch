package main

import (
	"fmt"
	"os"
	"time"

	"github.com/user/cronwatch/internal/history"
	"github.com/spf13/cobra"
)

func trendCmd() *cobra.Command {
	var (
		hisPath string
		n       int
	)

	cmd := &cobra.Command{
		Use:   "trend <job>",
		Short: "Show performance trend for a job over recent runs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jobName := args[0]
			trend, err := history.ComputeTrend(hisPath, jobName, n)
			if err != nil {
				return fmt.Errorf("trend: %w", err)
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "Trend report for job: %s (last %d runs)\n", trend.JobName, n)
			fmt.Fprintf(w, "%-22s %-14s %-14s %s\n", "Metric", "Early half", "Late half", "Direction")
			fmt.Fprintf(w, "%-22s %-14s %-14s %s\n",
				"Avg duration",
				roundDur(trend.AvgDurationEarly),
				roundDur(trend.AvgDurationLate),
				formatDirection(trend.DurationTrend),
			)
			fmt.Fprintf(w, "%-22s %-14s %-14s %s\n",
				"Success rate",
				fmt.Sprintf("%.1f%%", trend.SuccessRateEarly),
				fmt.Sprintf("%.1f%%", trend.SuccessRateLate),
				formatDirection(trend.SuccessRateTrend),
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&hisPath, "history", defaultHistoryPath(), "path to history file")
	cmd.Flags().IntVar(&n, "runs", 20, "number of recent runs to analyse")
	return cmd
}

func roundDur(d time.Duration) string {
	return d.Round(time.Millisecond).String()
}

func formatDirection(dir history.TrendDirection) string {
	switch dir {
	case history.TrendImproving:
		return "✓ improving"
	case history.TrendDegrading:
		return "✗ degrading"
	default:
		return "~ stable"
	}
}

var _ = os.Stderr // keep import
