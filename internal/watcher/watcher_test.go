package watcher_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/notifier"
	"github.com/user/cronwatch/internal/watcher"
)

func newTestNotifier(t *testing.T, buf *bytes.Buffer) *notifier.Notifier {
	t.Helper()
	n, err := notifier.New(buf)
	if err != nil {
		t.Fatalf("notifier.New: %v", err)
	}
	return n
}

func TestWatch_SuccessfulJob_NoAlert(t *testing.T) {
	var buf bytes.Buffer
	w, _ := watcher.New(watcher.Config{
		JobName:  "backup",
		Timeout:  5 * time.Second,
		Notifier: newTestNotifier(t, &buf),
	})

	err := w.Watch(context.Background(), func() error { return nil })
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no alert output, got: %s", buf.String())
	}
}

func TestWatch_FailingJob_SendsAlert(t *testing.T) {
	var buf bytes.Buffer
	w, _ := watcher.New(watcher.Config{
		JobName:  "backup",
		Timeout:  5 * time.Second,
		Notifier: newTestNotifier(t, &buf),
	})

	jobErr := errors.New("disk full")
	err := w.Watch(context.Background(), func() error { return jobErr })
	if !errors.Is(err, jobErr) {
		t.Fatalf("expected job error, got %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected alert to be written, got empty output")
	}
}

func TestWatch_TimedOutJob_SendsAlert(t *testing.T) {
	var buf bytes.Buffer
	w, _ := watcher.New(watcher.Config{
		JobName:  "slow-job",
		Timeout:  50 * time.Millisecond,
		Notifier: newTestNotifier(t, &buf),
	})

	err := w.Watch(context.Background(), func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if buf.Len() == 0 {
		t.Error("expected timeout alert to be written")
	}
}

func TestNew_EmptyJobName_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	_, err := watcher.New(watcher.Config{
		JobName:  "",
		Notifier: newTestNotifier(t, &buf),
	})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestNew_NilNotifier_ReturnsError(t *testing.T) {
	_, err := watcher.New(watcher.Config{JobName: "backup", Notifier: nil})
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}
