package history

import (
	"errors"
	"time"
)

// Annotation holds a user-supplied note attached to a specific history entry.
type Annotation struct {
	JobName   string    `json:"job_name"`
	RunID     string    `json:"run_id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// Annotate attaches a note to the history entry identified by jobName and runID.
// The entry must already exist in the store. The annotation is persisted by
// re-saving the store.
func Annotate(path, jobName, runID, note string) error {
	if path == "" {
		return errors.New("annotate: path must not be empty")
	}
	if jobName == "" {
		return errors.New("annotate: job name must not be empty")
	}
	if runID == "" {
		return errors.New("annotate: run ID must not be empty")
	}
	if note == "" {
		return errors.New("annotate: note must not be empty")
	}

	h, err := New(path)
	if err != nil {
		return err
	}

	entries, ok := h.store[jobName]
	if !ok {
		return errors.New("annotate: job not found: " + jobName)
	}

	found := false
	for i := range entries {
		if entries[i].RunID == runID {
			entries[i].Note = note
			found = true
			break
		}
	}
	if !found {
		return errors.New("annotate: run ID not found: " + runID)
	}

	h.store[jobName] = entries
	return h.save()
}

// GetAnnotations returns all entries for jobName that have a non-empty note.
func GetAnnotations(path, jobName string) ([]Entry, error) {
	if path == "" {
		return nil, errors.New("get annotations: path must not be empty")
	}
	if jobName == "" {
		return nil, errors.New("get annotations: job name must not be empty")
	}

	h, err := New(path)
	if err != nil {
		return nil, err
	}

	var annotated []Entry
	for _, e := range h.store[jobName] {
		if e.Note != "" {
			annotated = append(annotated, e)
		}
	}
	return annotated, nil
}
