// Package state persists the current task between Alchemist commands.
// State lives in .alchemist/current.json inside the working repo.
package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	dir      = ".alchemist"
	fileName = "current.json"
)

// State is the task a student is currently working on.
type State struct {
	TaskName    string    `json:"task_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func path() string {
	return filepath.Join(dir, fileName)
}

// Save writes the current task to disk.
func Save(s State) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path(), data, 0o644)
}

// Load reads the current task. A missing state file is not an error — it
// simply means no task has been started yet, so an empty State is returned.
// Real errors (corrupt JSON, permission problems) are still surfaced.
func Load() (State, error) {
	var s State
	data, err := os.ReadFile(path())
	if errors.Is(err, os.ErrNotExist) {
		return s, nil // no task started yet — empty state, not an error
	}
	if err != nil {
		return s, err // real error (permissions, …)
	}
	return s, json.Unmarshal(data, &s)
}

// Clear removes the saved task (called once work is bottled). It is
// idempotent: clearing when no task exists is a no-op, not an error.
func Clear() error {
	err := os.Remove(path())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
