// Package state persists the current task between Alchemist commands.
// State lives in .alchemist/current.json inside the working repo.
package state

import (
	"encoding/json"
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

// Load reads the current task. It returns an error if none is saved.
func Load() (State, error) {
	var s State
	data, err := os.ReadFile(path())
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(data, &s)
	return s, err
}

// Clear removes the saved task (called once work is bottled).
func Clear() error {
	return os.Remove(path())
}
