// Package editor opens text in the user's preferred editor.
package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// WriteTempFile writes content to a new temp file and returns its path.
// The caller is responsible for removing the file.
func WriteTempFile(content string) (string, error) {
	tmp, err := os.CreateTemp("", "alchemist-*.txt")
	if err != nil {
		return "", err
	}
	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}
	tmp.Close()
	return tmp.Name(), nil
}

// Command returns an exec.Cmd that opens path in the preferred editor.
// Use with tea.ExecProcess to suspend the TUI while editing.
func Command(path string) *exec.Cmd {
	return exec.Command(pick(), path)
}

// Edit writes content to a temp file, opens it, and returns the edited text.
func Edit(content string) (string, error) {
	path, err := WriteTempFile(content)
	if err != nil {
		return "", err
	}
	defer os.Remove(path)
	cmd := Command(path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not open editor %q: %w", pick(), err)
	}
	edited, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(edited), nil
}

func pick() string {
	if _, err := exec.LookPath("micro"); err == nil {
		return "micro"
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "nano"
}
