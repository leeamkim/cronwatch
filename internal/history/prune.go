package history

import (
	"fmt"
	"time"
)

// PruneOptions controls which entries are removed during a prune operation.
type PruneOptions struct {
	// OlderThan removes entries older than this duration.
	OlderThan time.Duration
	// JobName, if non-empty, restricts pruning to a single job.
	JobName string
	// DryRun reports what would be removed without modifying the store.
	DryRun bool
}

// PruneResult summarises the outcome of a prune operation.
type PruneResult struct {
	Removed int
	Retained int
}

// Prune removes history entries that match the given options.
// It returns a PruneResult describing how many entries were affected.
func (h *History) Prune(opts PruneOptions) (PruneResult, error) {
	if opts.OlderThan <= 0 {
		return PruneResult{}, fmt.Errorf("prune: OlderThan must be a positive duration")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	cutoff := time.Now().Add(-opts.OlderThan)
	var result PruneResult

	for job, entries := range h.store.Runs {
		if opts.JobName != "" && job != opts.JobName {
			continue
		}

		var kept []Entry
		for _, e := range entries {
			if e.StartedAt.Before(cutoff) {
				result.Removed++
			} else {
				kept = append(kept, e)
				result.Retained++
			}
		}

		if !opts.DryRun {
			if len(kept) == 0 {
				delete(h.store.Runs, job)
			} else {
				h.store.Runs[job] = kept
			}
		}
	}

	if opts.DryRun {
		return result, nil
	}

	return result, h.save()
}
