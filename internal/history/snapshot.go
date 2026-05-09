package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the current state of all job history at a point in time.
type Snapshot struct {
	CapturedAt time.Time            `json:"captured_at"`
	Entries    map[string][]Entry   `json:"entries"`
}

// TakeSnapshot reads the store at path and writes a JSON snapshot to dest.
func TakeSnapshot(path, dest string) error {
	if path == "" {
		return fmt.Errorf("history path must not be empty")
	}
	if dest == "" {
		return fmt.Errorf("snapshot destination path must not be empty")
	}

	store, err := load(path)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Entries:    store,
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create snapshot file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("encode snapshot: %w", err)
	}

	return nil
}

// LoadSnapshot reads a previously written snapshot from path.
func LoadSnapshot(path string) (*Snapshot, error) {
	if path == "" {
		return nil, fmt.Errorf("snapshot path must not be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}

	return &snap, nil
}
