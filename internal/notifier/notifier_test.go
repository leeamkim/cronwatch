package notifier_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/notifier"
)

func TestNotify_WritesFormattedAlert(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)

	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := n.Notify(notifier.Alert{
		JobName:   "backup-db",
		AlertType: notifier.AlertFailure,
		Message:   "exit status 1",
		OccuredAt: at,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "backup-db") {
		t.Errorf("expected job name in output, got: %s", output)
	}
	if !strings.Contains(output, string(notifier.AlertFailure)) {
		t.Errorf("expected alert type in output, got: %s", output)
	}
	if !strings.Contains(output, "exit status 1") {
		t.Errorf("expected message in output, got: %s", output)
	}
	if !strings.Contains(output, "2024-06-01T12:00:00Z") {
		t.Errorf("expected timestamp in output, got: %s", output)
	}
}

func TestNotify_TimeoutAlert(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)

	err := n.Notify(notifier.Alert{
		JobName:   "send-report",
		AlertType: notifier.AlertTimeout,
		Message:   "exceeded 30s limit",
		OccuredAt: time.Now(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), string(notifier.AlertTimeout)) {
		t.Errorf("expected TIMEOUT in output, got: %s", buf.String())
	}
}

func TestNotify_EmptyJobName_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)

	err := n.Notify(notifier.Alert{
		AlertType: notifier.AlertFailure,
		Message:   "something failed",
	})

	if err == nil {
		t.Fatal("expected error for empty job name, got nil")
	}
}

func TestNew_NilWriter_UsesStderr(t *testing.T) {
	// Should not panic
	n := notifier.New(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
