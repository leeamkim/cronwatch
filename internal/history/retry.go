package history

import (
	"fmt"
	"time"
)

// RetryEntry records a retry attempt for a job run.
type RetryEntry struct {
	JobName   string    `json:"job_name"`
	RunID     string    `json:"run_id"`
	Attempt   int       `json:"attempt"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
	Succeeded bool      `json:"succeeded"`
}

// RecordRetry appends a retry entry to the history store for the given job.
func RecordRetry(path, jobName, runID string, attempt int, err error, succeeded bool) error {
	if path == "" {
		return fmt.Errorf("history path must not be empty")
	}
	if jobName == "" {
		return fmt.Errorf("job name must not be empty")
	}
	if runID == "" {
		return fmt.Errorf("run ID must not be empty")
	}
	if attempt < 1 {
		return fmt.Errorf("attempt must be >= 1, got %d", attempt)
	}

	s, loadErr := New(path)
	if loadErr != nil {
		return fmt.Errorf("load store: %w", loadErr)
	}

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	retryKey := jobName + ":retries"
	entry := Entry{
		RunID:     runID + fmt.Sprintf("#retry%d", attempt),
		JobName:   jobName,
		StartedAt: time.Now(),
		FinishedAt: func() *time.Time { t := time.Now(); return &t }(),
		Success:   succeeded,
		Error:     errMsg,
		Tags:      []string{"retry", fmt.Sprintf("attempt:%d", attempt)},
	}

	s.mu.Lock()
	s.data[retryKey] = append(s.data[retryKey], entry)
	s.mu.Unlock()

	return s.save()
}

// GetRetries returns all retry entries for the given job and original run ID.
func GetRetries(path, jobName, runID string) ([]Entry, error) {
	if path == "" {
		return nil, fmt.Errorf("history path must not be empty")
	}

	s, err := New(path)
	if err != nil {
		return nil, fmt.Errorf("load store: %w", err)
	}

	retryKey := jobName + ":retries"
	s.mu.RLock()
	all := s.data[retryKey]
	s.mu.RUnlock()

	var results []Entry
	prefix := runID + "#retry"
	for _, e := range all {
		if len(e.RunID) >= len(prefix) && e.RunID[:len(prefix)] == prefix {
			results = append(results, e)
		}
	}
	return results, nil
}
