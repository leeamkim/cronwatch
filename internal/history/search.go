package history

import (
	"strings"
	"time"
)

// SearchOptions defines filters for querying history entries.
type SearchOptions struct {
	JobName  string
	Status   string // "success", "failure", or "" for all
	Since    time.Time
	MaxResults int
}

// Search returns history entries matching the given options.
func Search(store map[string][]Entry, opts SearchOptions) []Entry {
	var results []Entry

	for jobName, entries := range store {
		if opts.JobName != "" && !strings.EqualFold(jobName, opts.JobName) {
			continue
		}
		for _, e := range entries {
			if !opts.Since.IsZero() && e.StartedAt.Before(opts.Since) {
				continue
			}
			if opts.Status != "" {
				switch opts.Status {
				case "success":
					if e.Error != "" {
						continue
					}
				case "failure":
					if e.Error == "" {
						continue
					}
				}
			}
			results = append(results, e)
		}
	}

	// Sort by StartedAt descending
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].StartedAt.After(results[j-1].StartedAt); j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	if opts.MaxResults > 0 && len(results) > opts.MaxResults {
		results = results[:opts.MaxResults]
	}

	return results
}
