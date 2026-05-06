package history

import "time"

// CleanupOptions controls how old entries are pruned from the store.
type CleanupOptions struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// MaxRunsPerJob keeps at most this many recent runs per job. Zero means unlimited.
	MaxRunsPerJob int
}

// Cleanup prunes the in-memory store according to opts and persists the result.
func (h *History) Cleanup(opts CleanupOptions) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()

	for job, entries := range h.store.Jobs {
		filtered := entries[:0]

		for _, e := range entries {
			if opts.MaxAge > 0 && now.Sub(e.StartedAt) > opts.MaxAge {
				continue
			}
			filtered = append(filtered, e)
		}

		if opts.MaxRunsPerJob > 0 && len(filtered) > opts.MaxRunsPerJob {
			// Keep the most recent N entries (entries are stored oldest-first).
			filtered = filtered[len(filtered)-opts.MaxRunsPerJob:]
		}

		if len(filtered) == 0 {
			delete(h.store.Jobs, job)
		} else {
			h.store.Jobs[job] = filtered
		}
	}

	return h.persist()
}
