package history

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func seedAlertStore(t *testing.T, entries map[string][]Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "alert-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestCheckAlerts_EmptyPath_ReturnsError(t *testing.T) {
	_, err := CheckAlerts("", []AlertRule{{JobName: "job"}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCheckAlerts_NoRules_ReturnsNil(t *testing.T) {
	path := seedAlertStore(t, map[string][]Entry{})
	breaches, err := CheckAlerts(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(breaches) != 0 {
		t.Fatalf("expected 0 breaches, got %d", len(breaches))
	}
}

func TestCheckAlerts_FailStreak_Breach(t *testing.T) {
	entries := map[string][]Entry{
		"backup": {
			{Status: "success", Duration: time.Second},
			{Status: "failure", Duration: time.Second},
			{Status: "failure", Duration: time.Second},
			{Status: "failure", Duration: time.Second},
		},
	}
	path := seedAlertStore(t, entries)
	breaches, err := CheckAlerts(path, []AlertRule{{JobName: "backup", MaxFailStreak: 3}})
	if err != nil {
		t.Fatal(err)
	}
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].JobName != "backup" {
		t.Errorf("unexpected job name: %s", breaches[0].JobName)
	}
}

func TestCheckAlerts_SuccessRate_Breach(t *testing.T) {
	entries := map[string][]Entry{
		"sync": {
			{Status: "success"}, {Status: "failure"}, {Status: "failure"}, {Status: "failure"},
		},
	}
	path := seedAlertStore(t, entries)
	breaches, err := CheckAlerts(path, []AlertRule{{JobName: "sync", MinSuccessRate: 0.8}})
	if err != nil {
		t.Fatal(err)
	}
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
}

func TestCheckAlerts_AvgDuration_Breach(t *testing.T) {
	entries := map[string][]Entry{
		"report": {
			{Status: "success", Duration: 10 * time.Second},
			{Status: "success", Duration: 20 * time.Second},
		},
	}
	path := seedAlertStore(t, entries)
	breaches, err := CheckAlerts(path, []AlertRule{{JobName: "report", MaxAvgDuration: 5 * time.Second}})
	if err != nil {
		t.Fatal(err)
	}
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
}

func TestCheckAlerts_NoBreach_ReturnsEmpty(t *testing.T) {
	entries := map[string][]Entry{
		"nightly": {
			{Status: "success", Duration: time.Second},
			{Status: "success", Duration: time.Second},
		},
	}
	path := seedAlertStore(t, entries)
	rules := []AlertRule{{JobName: "nightly", MaxFailStreak: 3, MinSuccessRate: 0.5, MaxAvgDuration: time.Minute}}
	breaches, err := CheckAlerts(path, rules)
	if err != nil {
		t.Fatal(err)
	}
	if len(breaches) != 0 {
		t.Fatalf("expected no breaches, got %d", len(breaches))
	}
}
