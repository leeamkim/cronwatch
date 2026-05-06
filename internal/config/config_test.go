package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `{"jobs":[{"name":"backup","command":"/usr/bin/backup.sh","timeout":300000000000}]}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Name != "backup" {
		t.Errorf("expected name 'backup', got %q", cfg.Jobs[0].Name)
	}
	if cfg.Jobs[0].Timeout != 5*time.Minute {
		t.Errorf("expected 5m timeout, got %v", cfg.Jobs[0].Timeout)
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoJobs_ReturnsError(t *testing.T) {
	path := writeTemp(t, `{"jobs":[]}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty jobs")
	}
}

func TestLoad_MissingJobName_ReturnsError(t *testing.T) {
	path := writeTemp(t, `{"jobs":[{"command":"/usr/bin/backup.sh"}]}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing job name")
	}
}

func TestLoad_MissingCommand_ReturnsError(t *testing.T) {
	path := writeTemp(t, `{"jobs":[{"name":"backup"}]}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := writeTemp(t, `not json`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
