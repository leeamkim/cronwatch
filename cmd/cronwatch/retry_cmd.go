package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/example/cronwatch/internal/history"
)

func retryCmd(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("retry", flag.ContinueOnError)
	histPath := fs.String("history", defaultHistoryPath(), "path to history file")
	jobName := fs.String("job", "", "job name to query retries for (required)")
	runID := fs.String("run", "", "run ID to query retries for (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *jobName == "" {
		return fmt.Errorf("--job is required")
	}
	if *runID == "" {
		return fmt.Errorf("--run is required")
	}

	if out == nil {
		out = os.Stdout
	}

	entries, err := history.GetRetries(*histPath, *jobName, *runID)
	if err != nil {
		return fmt.Errorf("get retries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintf(out, "No retry entries found for job %q run %q\n", *jobName, *runID)
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ATTEMPT\tTIMESTAMP\tSUCCEEDED\tERROR")
	for i, e := range entries {
		succeeded := "no"
		if e.Success {
			succeeded = "yes"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			i+1,
			e.StartedAt.Format(time.RFC3339),
			succeeded,
			e.Error,
		)
	}
	return w.Flush()
}
