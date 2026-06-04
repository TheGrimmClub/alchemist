// Package editor opens text in the user's editor and returns the result.
package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Edit writes content to a temp file, opens it in an editor, and returns
// the edited text. It prefers micro, then $EDITOR, then an OS default.
func Edit(content string) (string, error) {
	tmp, err := os.CreateTemp("", "alchemist-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()

	editor := pick()
	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not open editor %q: %w", editor, err)
	}

	edited, err := os.ReadFile(tmp.Name())
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
