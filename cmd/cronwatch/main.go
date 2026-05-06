package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/notifier"
	"github.com/user/cronwatch/internal/runner"
	"github.com/user/cronwatch/internal/watcher"
)

func main() {
	configPath := flag.String("config", "cronwatch.json", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("cronwatch: %v", err)
	}

	var exitCode int
	for _, jobCfg := range cfg.Jobs {
		if err := runJob(jobCfg); err != nil {
			fmt.Fprintf(os.Stderr, "cronwatch: job %q failed: %v\n", jobCfg.Name, err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

func runJob(jobCfg config.JobConfig) error {
	n, err := notifier.New(jobCfg.Name, os.Stderr)
	if err != nil {
		return fmt.Errorf("create notifier: %w", err)
	}

	w, err := watcher.New(jobCfg.Name, n)
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}

	r, err := runner.New(w)
	if err != nil {
		return fmt.Errorf("create runner: %w", err)
	}

	return r.Run(jobCfg)
}
