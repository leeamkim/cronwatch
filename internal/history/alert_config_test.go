package history

import (
	"os"
	"testing"
	"time"
)

func writeAlertConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "alert-cfg-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadAlertRules_ValidConfig(t *testing.T) {
	path := writeAlertConfig(t, `
alert_rules:
  - job: backup
    max_fail_streak: 3
    min_success_rate: 0.9
    max_avg_duration: 30s
  - job: sync
    max_fail_streak: 2
`)
	rules, err := LoadAlertRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].JobName != "backup" {
		t.Errorf("expected job 'backup', got %q", rules[0].JobName)
	}
	if rules[0].MaxAvgDuration != 30*time.Second {
		t.Errorf("expected 30s, got %s", rules[0].MaxAvgDuration)
	}
	if rules[0].MinSuccessRate != 0.9 {
		t.Errorf("expected 0.9, got %f", rules[0].MinSuccessRate)
	}
}

func TestLoadAlertRules_EmptyPath_ReturnsError(t *testing.T) {
	_, err := LoadAlertRules("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoadAlertRules_MissingJobName_ReturnsError(t *testing.T) {
	path := writeAlertConfig(t, `
alert_rules:
  - max_fail_streak: 3
`)
	_, err := LoadAlertRules(path)
	if err == nil {
		t.Fatal("expected error for missing job name")
	}
}

func TestLoadAlertRules_InvalidDuration_ReturnsError(t *testing.T) {
	path := writeAlertConfig(t, `
alert_rules:
  - job: nightly
    max_avg_duration: notaduration
`)
	_, err := LoadAlertRules(path)
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestLoadAlertRules_MissingFile_ReturnsError(t *testing.T) {
	_, err := LoadAlertRules("/nonexistent/path/alert.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
