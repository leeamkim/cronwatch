package runner_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/notifier"
	"github.com/user/cronwatch/internal/runner"
	"github.com/user/cronwatch/internal/watcher"
)

func newTestRunner(t *testing.T, buf *bytes.Buffer) *runner.Runner {
	t.Helper()
	n, err := notifier.New("test-job", buf)
	if err != nil {
		t.Fatal(err)
	}
	w, err := watcher.New("test-job", n)
	if err != nil {
		t.Fatal(err)
	}
	r, err := runner.New(w)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestRun_SuccessfulCommand_NoAlert(t *testing.T) {
	var buf bytes.Buffer
	r := newTestRunner(t, &buf)
	cfg := config.JobConfig{Name: "test-job", Command: "true", Timeout: 5 * time.Second}
	if err := r.Run(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no alert output, got: %s", buf.String())
	}
}

func TestRun_FailingCommand_SendsAlert(t *testing.T) {
	var buf bytes.Buffer
	r := newTestRunner(t, &buf)
	cfg := config.JobConfig{Name: "test-job", Command: "false", Timeout: 5 * time.Second}
	_ = r.Run(cfg)
	if !strings.Contains(buf.String(), "test-job") {
		t.Errorf("expected alert containing job name, got: %s", buf.String())
	}
}

func TestNew_NilWatcher_ReturnsError(t *testing.T) {
	_, err := runner.New(nil)
	if err == nil {
		t.Fatal("expected error for nil watcher")
	}
}

func TestRun_DefaultTimeout_Applied(t *testing.T) {
	var buf bytes.Buffer
	r := newTestRunner(t, &buf)
	cfg := config.JobConfig{Name: "test-job", Command: "true"} // zero timeout
	if err := r.Run(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
