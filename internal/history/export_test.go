package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestExportCSV_ContainsHeader(t *testing.T) {
	path := tempPath(t)
	h, _ := New(path)

	_ = h.Record("backup", time.Now(), 1200, true, "")

	var buf bytes.Buffer
	if err := ExportCSV(path, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "job,started_at,duration_ms,success,error") {
		t.Errorf("expected CSV header in output, got:\n%s", out)
	}
}

func TestExportCSV_ContainsEntries(t *testing.T) {
	path := tempPath(t)
	h, _ := New(path)

	start := time.Now().UTC().Truncate(time.Second)
	_ = h.Record("sync", start, 500, true, "")
	_ = h.Record("sync", start, 800, false, "exit status 1")

	var buf bytes.Buffer
	if err := ExportCSV(path, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "sync") {
		t.Errorf("expected job name 'sync' in output")
	}
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output")
	}
	if !strings.Contains(out, "false") {
		t.Errorf("expected failed entry in output")
	}
}

func TestExportCSV_EmptyPath_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV("", &buf); err == nil {
		t.Error("expected error for empty path, got nil")
	}
}

func TestExportCSV_NilWriter_UsesStdout(t *testing.T) {
	path := tempPath(t)
	_, _ = New(path)

	// Should not panic or error when writer is nil (falls back to stdout)
	if err := ExportCSV(path, nil); err != nil {
		t.Fatalf("unexpected error with nil writer: %v", err)
	}
}
