package state

import (
	"os"
	"path/filepath"
	"testing"
)

// Each test runs in its own temp dir so .alchemist/ is isolated.
func TestLoad_MissingFileReturnsEmptyState(t *testing.T) {
	t.Chdir(t.TempDir())

	s, err := Load()
	if err != nil {
		t.Fatalf("expected nil error for a missing state file, got %v", err)
	}
	if s.TaskName != "" {
		t.Fatalf("expected empty TaskName, got %q", s.TaskName)
	}
}

func TestLoad_ValidFileReturnsTask(t *testing.T) {
	t.Chdir(t.TempDir())

	want := State{TaskName: "My Task", Description: "desc"}
	if err := Save(want); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got.TaskName != want.TaskName || got.Description != want.Description {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", got, want)
	}
}

func TestLoad_CorruptJSONReturnsError(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("setup mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, fileName), []byte("{not valid json"), 0o644); err != nil {
		t.Fatalf("setup write failed: %v", err)
	}

	if _, err := Load(); err == nil {
		t.Fatal("expected an error for corrupt JSON, got nil")
	}
}

func TestClear_MissingFileIsNoOp(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := Clear(); err != nil {
		t.Fatalf("expected nil error clearing a missing state file, got %v", err)
	}
}

func TestClear_RemovesExistingFile(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := Save(State{TaskName: "x"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if err := Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	if _, err := os.Stat(path()); !os.IsNotExist(err) {
		t.Fatalf("expected state file to be gone, stat err = %v", err)
	}
}
