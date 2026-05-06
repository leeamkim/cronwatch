package history

import (
	"fmt"
	"time"
)

// Retention defines how long to keep history entries.
type Retention struct {
	MaxAge  time.Duration
	MaxRuns int
}

// DefaultRetention keeps 30 days of history and at most 500 runs per job.
var DefaultRetention = Retention{
	MaxAge:  30 * 24 * time.Hour,
	MaxRuns: 500,
}

// Cleanup removes entries from the store that exceed the given retention
// policy. It returns the number of entries removed.
func (s *Store) Cleanup(r Retention) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-r.MaxAge)
	removed := 0

	for job, entries := range s.data {
		var kept []Entry
		for _, e := range entries {
			if e.StartedAt.After(cutoff) {
				kept = append(kept, e)
			} else {
				removed++
			}
		}

		// Enforce max runs per job (keep the most recent ones).
		if r.MaxRuns > 0 && len(kept) > r.MaxRuns {
			excess := len(kept) - r.MaxRuns
			removed += excess
			kept = kept[excess:]
		}

		if len(kept) == 0 {
			delete(s.data, job)
		} else {
			s.data[job] = kept
		}
	}

	if err := s.persist(); err != nil {
		return removed, fmt.Errorf("cleanup: persist: %w", err)
	}

	return removed, nil
}
