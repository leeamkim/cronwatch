package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

// ExportCSV writes all history entries to w in CSV format.
// If w is nil, os.Stdout is used.
func ExportCSV(path string, w io.Writer) error {
	if path == "" {
		return fmt.Errorf("history path must not be empty")
	}

	if w == nil {
		w = os.Stdout
	}

	h, err := New(path)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}

	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"job", "started_at", "duration_ms", "success", "error"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for job, entries := range h.store.Runs {
		for _, e := range entries {
			errMsg := ""
			if e.Error != "" {
				errMsg = e.Error
			}
			row := []string{
				job,
				e.StartedAt.Format(time.RFC3339),
				fmt.Sprintf("%d", e.DurationMs),
				fmt.Sprintf("%t", e.Success),
				errMsg,
			}
			if err := cw.Write(row); err != nil {
				return fmt.Errorf("write csv row: %w", err)
			}
		}
	}

	return nil
}
