package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_EmptyPath_ReturnsError(t *testing.T) {
	_, err := history.New("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNew_NonExistentFile_CreatesEmptyStore(t *testing.T) {
	s, err := history.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.All(); len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestRecord_PersistsEntry(t *testing.T) {
	path := tempPath(t)
	s, err := history.New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := history.Entry{
		JobName:   "backup",
		StartedAt: time.Now(),
		Duration:  2 * time.Second,
		Success:   true,
	}
	if err := s.Record(entry); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	// Reload from disk
	s2, err := history.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	entries := s2.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(entries))
	}
	if entries[0].JobName != "backup" {
		t.Errorf("expected job name 'backup', got %q", entries[0].JobName)
	}
	if !entries[0].Success {
		t.Error("expected Success=true")
	}
}

func TestRecord_FailedEntry_StoresError(t *testing.T) {
	s, err := history.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := history.Entry{
		JobName: "cleanup",
		Success: false,
		TimedOut: true,
		Error:   "context deadline exceeded",
	}
	if err := s.Record(entry); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	entries := s.All()
	if entries[0].Error != "context deadline exceeded" {
		t.Errorf("unexpected error field: %q", entries[0].Error)
	}
	if !entries[0].TimedOut {
		t.Error("expected TimedOut=true")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatalf("setup error: %v", err)
	}
	_, err := history.New(path)
	if err == nil {
		t.Fatal("expected error for corrupt file, got nil")
	}
}
