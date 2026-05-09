package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func seedTrendPrintStore(t *testing.T) string {
	t.Helper()
	path := tempPath(t)
	store, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 4; i++ {
		store.entries["sync"] = append(store.entries["sync"], Entry{
			RunID:    fmt.Sprintf("r%d", i),
			JobName:  "sync",
			Duration: time.Duration(i+1) * time.Second,
			Error:    "",
			At:       time.Now(),
		})
	}
	_ = store.save()
	return path
}

func TestPrintTrend_ContainsHeader(t *testing.T) {
	path := seedTrendPrintStore(t)
	var buf bytes.Buffer
	if err := PrintTrend(path, 4, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Job") {
		t.Error("expected header row containing 'Job'")
	}
}

func TestPrintTrend_ContainsJobName(t *testing.T) {
	path := seedTrendPrintStore(t)
	var buf bytes.Buffer
	_ = PrintTrend(path, 4, &buf)
	if !strings.Contains(buf.String(), "sync") {
		t.Error("expected job name 'sync' in output")
	}
}

func TestPrintTrend_EmptyPath_ReturnsError(t *testing.T) {
	if err := PrintTrend("", 10, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestPrintTrend_NilWriter_UsesStdout(t *testing.T) {
	path := seedTrendPrintStore(t)
	// Should not panic.
	if err := PrintTrend(path, 4, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintTrend_EmptyStore_PrintsMessage(t *testing.T) {
	path := tempPath(t)
	_, _ = New(path) // creates empty store
	var buf bytes.Buffer
	_ = PrintTrend(path, 10, &buf)
	if !strings.Contains(buf.String(), "No history") {
		t.Error("expected 'No history' message for empty store")
	}
}
