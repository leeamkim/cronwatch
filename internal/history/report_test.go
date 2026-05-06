package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func seedStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now()
	_ = s.Record("backup", Entry{StartedAt: now.Add(-2 * time.Minute), FinishedAt: now.Add(-time.Minute), Error: ""})
	_ = s.Record("backup", Entry{StartedAt: now.Add(-time.Minute), FinishedAt: now, Error: "exit status 1"})
	_ = s.Record("sync", Entry{StartedAt: now.Add(-30 * time.Second), FinishedAt: now, TimedOut: true, Error: "timed out"})
	return s
}

func TestSummarise_CountsCorrectly(t *testing.T) {
	s := seedStore(t)
	sums := Summarise(s)

	byName := make(map[string]Summary, len(sums))
	for _, sum := range sums {
		byName[sum.JobName] = sum
	}

	backup, ok := byName["backup"]
	if !ok {
		t.Fatal("expected summary for 'backup'")
	}
	if backup.TotalRuns != 2 {
		t.Errorf("TotalRuns: got %d, want 2", backup.TotalRuns)
	}
	if backup.Failures != 1 {
		t.Errorf("Failures: got %d, want 1", backup.Failures)
	}
	if backup.LastStatus != "FAILED" {
		t.Errorf("LastStatus: got %s, want FAILED", backup.LastStatus)
	}

	sync, ok := byName["sync"]
	if !ok {
		t.Fatal("expected summary for 'sync'")
	}
	if sync.Timeouts != 1 {
		t.Errorf("Timeouts: got %d, want 1", sync.Timeouts)
	}
	if sync.LastStatus != "TIMEOUT" {
		t.Errorf("LastStatus: got %s, want TIMEOUT", sync.LastStatus)
	}
}

func TestPrintReport_ContainsHeaders(t *testing.T) {
	s := seedStore(t)
	var buf bytes.Buffer
	PrintReport(s, &buf)
	out := buf.String()
	for _, header := range []string{"JOB", "TOTAL", "FAILURES", "TIMEOUTS", "LAST RUN", "LAST STATUS"} {
		if !strings.Contains(out, header) {
			t.Errorf("output missing header %q", header)
		}
	}
}

func TestPrintReport_ContainsJobNames(t *testing.T) {
	s := seedStore(t)
	var buf bytes.Buffer
	PrintReport(s, &buf)
	out := buf.String()
	for _, name := range []string{"backup", "sync"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing job name %q", name)
		}
	}
}

func TestPrintReport_NilWriter_UsesStdout(t *testing.T) {
	s := seedStore(t)
	// Should not panic when w is nil.
	PrintReport(s, nil)
}
