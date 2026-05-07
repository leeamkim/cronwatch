package history

import (
	"errors"
	"testing"
)

func TestRecordRetry_PersistsEntry(t *testing.T) {
	p := tempPath(t)

	err := RecordRetry(p, "backup", "run-001", 1, errors.New("timeout"), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := GetRetries(p, "backup", "run-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 retry entry, got %d", len(entries))
	}
	if entries[0].Error != "timeout" {
		t.Errorf("expected error 'timeout', got %q", entries[0].Error)
	}
	if entries[0].Success {
		t.Error("expected success=false")
	}
}

func TestRecordRetry_MultipleAttempts(t *testing.T) {
	p := tempPath(t)

	_ = RecordRetry(p, "deploy", "run-42", 1, errors.New("fail"), false)
	_ = RecordRetry(p, "deploy", "run-42", 2, errors.New("fail again"), false)
	_ = RecordRetry(p, "deploy", "run-42", 3, nil, true)

	entries, err := GetRetries(p, "deploy", "run-42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 retry entries, got %d", len(entries))
	}
	if !entries[2].Success {
		t.Error("expected final attempt to succeed")
	}
}

func TestRecordRetry_EmptyPath_ReturnsError(t *testing.T) {
	err := RecordRetry("", "myjob", "run-1", 1, nil, true)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRecordRetry_EmptyJobName_ReturnsError(t *testing.T) {
	p := tempPath(t)
	err := RecordRetry(p, "", "run-1", 1, nil, true)
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestRecordRetry_InvalidAttempt_ReturnsError(t *testing.T) {
	p := tempPath(t)
	err := RecordRetry(p, "myjob", "run-1", 0, nil, true)
	if err == nil {
		t.Fatal("expected error for attempt < 1")
	}
}

func TestGetRetries_NoEntries_ReturnsEmpty(t *testing.T) {
	p := tempPath(t)

	entries, err := GetRetries(p, "nonexistent", "run-99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetRetries_IsolatesRunID(t *testing.T) {
	p := tempPath(t)

	_ = RecordRetry(p, "sync", "run-A", 1, errors.New("err"), false)
	_ = RecordRetry(p, "sync", "run-B", 1, errors.New("err"), false)

	entries, err := GetRetries(p, "sync", "run-A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry for run-A, got %d", len(entries))
	}
}
