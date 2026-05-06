package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// JobConfig holds configuration for a single monitored cron job.
type JobConfig struct {
	Name    string        `json:"name"`
	Command string        `json:"command"`
	Timeout time.Duration `json:"timeout"`
	Email   string        `json:"email,omitempty"`
}

// Config holds the full cronwatch configuration.
type Config struct {
	Jobs []JobConfig `json:"jobs"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("config: no jobs defined")
	}
	for i, j := range c.Jobs {
		if j.Name == "" {
			return fmt.Errorf("config: job[%d]: name is required", i)
		}
		if j.Command == "" {
			return fmt.Errorf("config: job[%d] %q: command is required", i, j.Name)
		}
	}
	return nil
}
