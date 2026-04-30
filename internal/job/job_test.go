package job

import (
	"testing"
	"time"
)

func TestNewRun(t *testing.T) {
	run := NewRun("job-123")

	if run.JobID != "job-123" {
		t.Errorf("expected JobID 'job-123', got '%s'", run.JobID)
	}
	if run.Status != StatusRunning {
		t.Errorf("expected status Running, got '%s'", run.Status)
	}
	if run.FinishedAt != nil {
		t.Error("expected FinishedAt to be nil for a new run")
	}
}

func TestRun_Finish(t *testing.T) {
	run := NewRun("job-456")
	time.Sleep(5 * time.Millisecond)
	run.Finish(StatusSuccess, "")

	if run.Status != StatusSuccess {
		t.Errorf("expected status Success, got '%s'", run.Status)
	}
	if run.FinishedAt == nil {
		t.Error("expected FinishedAt to be set after Finish")
	}
	if run.DurationMs <= 0 {
		t.Errorf("expected positive DurationMs, got %d", run.DurationMs)
	}
}

func TestRun_Finish_WithError(t *testing.T) {
	run := NewRun("job-789")
	run.Finish(StatusFailed, "exit code 1")

	if run.Status != StatusFailed {
		t.Errorf("expected status Failed, got '%s'", run.Status)
	}
	if run.Error != "exit code 1" {
		t.Errorf("expected error 'exit code 1', got '%s'", run.Error)
	}
}

func TestRun_IsTimedOut(t *testing.T) {
	run := NewRun("job-timeout")

	if run.IsTimedOut(60) {
		t.Error("fresh run should not be timed out with 60s timeout")
	}

	// Backdate start time to simulate a long-running job
	run.StartedAt = time.Now().UTC().Add(-10 * time.Second)
	if !run.IsTimedOut(5) {
		t.Error("run started 10s ago should be timed out with 5s timeout")
	}
}

func TestRun_IsTimedOut_WhenFinished(t *testing.T) {
	run := NewRun("job-done")
	run.StartedAt = time.Now().UTC().Add(-30 * time.Second)
	run.Finish(StatusSuccess, "")

	if run.IsTimedOut(5) {
		t.Error("finished run should never be considered timed out")
	}
}
