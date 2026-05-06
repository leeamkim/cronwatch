package watcher

import (
	"context"
	"time"

	"github.com/user/cronwatch/internal/job"
	"github.com/user/cronwatch/internal/notifier"
)

// Config holds configuration for a watched cron job.
type Config struct {
	JobName string
	Timeout time.Duration
	Notifier *notifier.Notifier
}

// Watcher monitors a cron job execution and sends alerts on failure or timeout.
type Watcher struct {
	cfg Config
}

// New creates a new Watcher with the given configuration.
func New(cfg Config) (*Watcher, error) {
	if cfg.JobName == "" {
		return nil, fmt.Errorf("watcher: job name must not be empty")
	}
	if cfg.Notifier == nil {
		return nil, fmt.Errorf("watcher: notifier must not be nil")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Minute
	}
	return &Watcher{cfg: cfg}, nil
}

// Watch runs fn under observation, alerting if it errors or exceeds the timeout.
func (w *Watcher) Watch(ctx context.Context, fn func() error) error {
	run := job.NewRun(w.cfg.JobName)

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- fn()
	}()

	select {
	case err := <-doneCh:
		run.Finish(err)
		if err != nil {
			_ = w.cfg.Notifier.Notify(run)
			return err
		}
		return nil
	case <-time.After(w.cfg.Timeout):
		run.Finish(nil)
		if run.IsTimedOut(w.cfg.Timeout) {
			_ = w.cfg.Notifier.Notify(run)
		}
		return fmt.Errorf("watcher: job %q exceeded timeout %s", w.cfg.JobName, w.cfg.Timeout)
	case <-ctx.Done():
		return ctx.Err()
	}
}
