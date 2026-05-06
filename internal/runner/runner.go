package runner

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/job"
	"github.com/user/cronwatch/internal/watcher"
)

// Runner executes configured jobs and reports results via a watcher.
type Runner struct {
	watcher *watcher.Watcher
}

// New creates a Runner backed by the given watcher.
func New(w *watcher.Watcher) (*Runner, error) {
	if w == nil {
		return nil, fmt.Errorf("runner: watcher must not be nil")
	}
	return &Runner{watcher: w}, nil
}

// Run executes the job described by cfg, monitors it, and returns any error.
func (r *Runner) Run(cfg config.JobConfig) error {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	run := job.NewRun(cfg.Name, timeout)

	cmd := exec.CommandContext(ctx, "sh", "-c", cfg.Command)
	err := cmd.Run()

	run.Finish(err)
	r.watcher.Watch(run)

	return err
}
