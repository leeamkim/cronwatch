package history

import (
	"testing"
	"time"
)

func seedAnnotateStore(t *testing.T) (string, *History) {
	t.Helper()
	p := tempPath(t)
	h, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	h.store["backup"] = []Entry{
		{RunID: "run-1", StartedAt: time.Now(), Success: true},
		{RunID: "run-2", StartedAt: time.Now(), Success: false, Error: "exit 1"},
	}
	if err := h.save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p, h
}

func TestAnnotate_AddsNoteToEntry(t *testing.T) {
	p, _ := seedAnnotateStore(t)

	if err := Annotate(p, "backup", "run-2", "investigated: disk full"); err != nil {
		t.Fatalf("Annotate: %v", err)
	}

	annotated, err := GetAnnotations(p, "backup")
	if err != nil {
		t.Fatalf("GetAnnotations: %v", err)
	}
	if len(annotated) != 1 {
		t.Fatalf("expected 1 annotated entry, got %d", len(annotated))
	}
	if annotated[0].RunID != "run-2" {
		t.Errorf("expected run-2, got %s", annotated[0].RunID)
	}
	if annotated[0].Note != "investigated: disk full" {
		t.Errorf("unexpected note: %s", annotated[0].Note)
	}
}

func TestAnnotate_UnknownJob_ReturnsError(t *testing.T) {
	p, _ := seedAnnotateStore(t)
	err := Annotate(p, "nonexistent", "run-1", "note")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestAnnotate_UnknownRunID_ReturnsError(t *testing.T) {
	p, _ := seedAnnotateStore(t)
	err := Annotate(p, "backup", "run-999", "note")
	if err == nil {
		t.Fatal("expected error for unknown run ID")
	}
}

func TestAnnotate_EmptyNote_ReturnsError(t *testing.T) {
	p, _ := seedAnnotateStore(t)
	err := Annotate(p, "backup", "run-1", "")
	if err == nil {
		t.Fatal("expected error for empty note")
	}
}

func TestAnnotate_EmptyPath_ReturnsError(t *testing.T) {
	err := Annotate("", "backup", "run-1", "note")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetAnnotations_NoAnnotations_ReturnsEmpty(t *testing.T) {
	p, _ := seedAnnotateStore(t)
	annotated, err := GetAnnotations(p, "backup")
	if err != nil {
		t.Fatalf("GetAnnotations: %v", err)
	}
	if len(annotated) != 0 {
		t.Errorf("expected 0 annotations before any are added, got %d", len(annotated))
	}
}

func TestGetAnnotations_EmptyJobName_ReturnsError(t *testing.T) {
	p, _ := seedAnnotateStore(t)
	_, err := GetAnnotations(p, "")
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}
