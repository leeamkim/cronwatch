package history

import (
	"fmt"
	"io"
	"os"
)

// PrintTrend writes a human-readable trend table for all jobs to w.
// If w is nil, os.Stdout is used.
func PrintTrend(path string, n int, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	if path == "" {
		return fmt.Errorf("history path must not be empty")
	}

	store, err := New(path)
	if err != nil {
		return err
	}

	if len(store.entries) == 0 {
		fmt.Fprintln(w, "No history entries found.")
		return nil
	}

	fmt.Fprintf(w, "%-20s %-12s %-12s %-12s %s\n",
		"Job", "Dur (early)", "Dur (late)", "SR (early)", "SR (late)")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 72))

	for jobName := range store.entries {
		t, err := ComputeTrend(path, jobName, n)
		if err != nil {
			// Skip jobs with insufficient data.
			continue
		}
		fmt.Fprintf(w, "%-20s %-12s %-12s %-12s %s\n",
			jobName,
			t.AvgDurationEarly.Round(1e6).String(),
			t.AvgDurationLate.Round(1e6).String(),
			fmt.Sprintf("%.1f%%", t.SuccessRateEarly),
			fmt.Sprintf("%.1f%%", t.SuccessRateLate),
		)
	}
	return nil
}

func repeatChar(c rune, n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = c
	}
	return string(out)
}
