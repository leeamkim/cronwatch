package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func seedSnapshot(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "history.json")
	h, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = h.Record("backup", Entry{RunID: "r1", StartedAt: time.Now(), Duration: 2 * time.Second, Success: true})
	_ = h.Record("backup", Entry{RunID: "r2", StartedAt: time.Now(), Duration: 3 * time.Second, Success: false, Error: "exit 1"})
	_ = h.Record("deploy", Entry{RunID: "r3", StartedAt: time.Now(), Duration: 1 * time.Second, Success: true})
	return path
}

func TestTakeSnapshot_CreatesFile(t *testing.T) {
	src := seedSnapshot(t)
	dest := filepath.Join(t.TempDir(), "snap.json")

	if err := TakeSnapshot(src, dest); err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("snapshot file not created: %v", err)
	}
}

func TestTakeSnapshot_EmptyPath_ReturnsError(t *testing.T) {
	if err := TakeSnapshot("", "/tmp/snap.json"); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestTakeSnapshot_EmptyDest_ReturnsError(t *testing.T) {
	if err := TakeSnapshot("/some/path", ""); err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestLoadSnapshot_ReturnsEntries(t *testing.T) {
	src := seedSnapshot(t)
	dest := filepath.Join(t.TempDir(), "snap.json")

	if err := TakeSnapshot(src, dest); err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	snap, err := LoadSnapshot(dest)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
	if len(snap.Entries["backup"]) != 2 {
		t.Errorf("expected 2 backup entries, got %d", len(snap.Entries["backup"]))
	}
	if len(snap.Entries["deploy"]) != 1 {
		t.Errorf("expected 1 deploy entry, got %d", len(snap.Entries["deploy"]))
	}
}

func TestLoadSnapshot_EmptyPath_ReturnsError(t *testing.T) {
	if _, err := LoadSnapshot(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoadSnapshot_MissingFile_ReturnsError(t *testing.T) {
	if _, err := LoadSnapshot("/nonexistent/snap.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}
