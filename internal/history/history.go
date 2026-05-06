package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single job run record.
type Entry struct {
	JobName   string        `json:"job_name"`
	StartedAt time.Time     `json:"started_at"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	TimedOut  bool          `json:"timed_out"`
	Error     string        `json:"error,omitempty"`
}

// Store persists job run history to a JSON file.
type Store struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New creates a new Store backed by the given file path.
// Existing entries are loaded from disk if the file exists.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("history: path must not be empty")
	}
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("history: loading existing entries: %w", err)
	}
	return s, nil
}

// Record appends an entry and flushes to disk.
func (s *Store) Record(e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	return s.flush()
}

// All returns a copy of all recorded entries.
func (s *Store) All() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshalling entries: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}
