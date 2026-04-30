package job

import "time"

// Status represents the current state of a cron job execution.
type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusTimeout Status = "timeout"
	StatusRunning Status = "running"
)

// Job represents a monitored cron job.
type Job struct {
	ID          string
	Name        string
	Schedule    string
	TimeoutSecs int
	LastRunAt   *time.Time
	LastStatus  Status
	LastError   string
	DurationMs  int64
}

// Run represents a single execution record of a job.
type Run struct {
	JobID      string
	StartedAt  time.Time
	FinishedAt *time.Time
	Status     Status
	Error      string
	DurationMs int64
}

// NewRun creates a new Run for the given job ID.
func NewRun(jobID string) *Run {
	return &Run{
		JobID:     jobID,
		StartedAt: time.Now().UTC(),
		Status:    StatusRunning,
	}
}

// Finish marks the run as completed with the given status and optional error message.
func (r *Run) Finish(status Status, errMsg string) {
	now := time.Now().UTC()
	r.FinishedAt = &now
	r.Status = status
	r.Error = errMsg
	r.DurationMs = now.Sub(r.StartedAt).Milliseconds()
}

// IsTimedOut returns true if the run has exceeded the job's timeout threshold.
func (r *Run) IsTimedOut(timeoutSecs int) bool {
	if r.FinishedAt != nil {
		return false
	}
	return time.Since(r.StartedAt).Seconds() > float64(timeoutSecs)
}
