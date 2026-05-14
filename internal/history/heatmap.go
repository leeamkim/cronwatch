package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// HeatmapEntry holds the success/failure counts for a single day.
type HeatmapEntry struct {
	Date    string `json:"date"`
	Success int    `json:"success"`
	Failure int    `json:"failure"`
}

// Heatmap maps date strings (YYYY-MM-DD) to HeatmapEntry.
type Heatmap map[string]HeatmapEntry

// ComputeHeatmap builds a day-by-day activity heatmap for the given job
// over the past n days. Pass an empty jobName to aggregate all jobs.
func ComputeHeatmap(path, jobName string, days int) (Heatmap, error) {
	if path == "" {
		return nil, fmt.Errorf("history path must not be empty")
	}
	if days <= 0 {
		return nil, fmt.Errorf("days must be greater than zero")
	}

	store, err := load(path)
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	hm := make(Heatmap)

	for job, entries := range store {
		if jobName != "" && job != jobName {
			continue
		}
		for _, e := range entries {
			if e.StartedAt.Before(cutoff) {
				continue
			}
			day := e.StartedAt.UTC().Format("2006-01-02")
			ent := hm[day]
			ent.Date = day
			if e.Error == "" {
				ent.Success++
			} else {
				ent.Failure++
			}
			hm[day] = ent
		}
	}
	return hm, nil
}

// PrintHeatmap writes a simple ASCII heatmap to w (defaults to stdout).
func PrintHeatmap(hm Heatmap, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "%-12s %8s %8s %8s\n", "Date", "Success", "Failure", "Total")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 40))

	// Collect and sort dates.
	dates := make([]string, 0, len(hm))
	for d := range hm {
		dates = append(dates, d)
	}
	sortStrings(dates)

	for _, d := range dates {
		e := hm[d]
		total := e.Success + e.Failure
		fmt.Fprintf(w, "%-12s %8d %8d %8d\n", d, e.Success, e.Failure, total)
	}
}

// ExportHeatmapJSON serialises the heatmap as JSON to w (defaults to stdout).
func ExportHeatmapJSON(hm Heatmap, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(hm)
}

// sortStrings sorts a slice of strings in place (insertion sort — small slices).
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
