// Package gitutil wraps the git command line so the rest of Alchemist
// can stay readable.
package gitutil

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// run executes a git command and returns trimmed combined output.
func run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// Interactive runs git wired to the terminal (stdin/stdout/stderr) so the
// student sees git's own output and prompts.
func Interactive(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsRepo reports whether the current directory is inside a git repo.
func IsRepo() bool {
	_, err := run("rev-parse", "--is-inside-work-tree")
	return err == nil
}

// EnsureRepo returns a friendly error if not inside a git repo.
func EnsureRepo() error {
	if !IsRepo() {
		return fmt.Errorf("not a git repository — run 'git init' or 'alchemist recipe' first")
	}
	return nil
}

// ChangedFiles lists tracked files with uncommitted changes.
func ChangedFiles() ([]string, error) {
	out, err := run("status", "--porcelain")
	if err != nil {
		return nil, err
	}
	return parseStatus(out, false), nil
}

// UntrackedFiles lists files git is not yet tracking.
func UntrackedFiles() ([]string, error) {
	out, err := run("status", "--porcelain")
	if err != nil {
		return nil, err
	}
	return parseStatus(out, true), nil
}

func parseStatus(out string, wantUntracked bool) []string {
	var files []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		code := line[:2]
		name := strings.TrimSpace(line[3:])
		isUntracked := code == "??"
		if wantUntracked == isUntracked {
			files = append(files, name)
		}
	}
	return files
}

// AddAll stages every change.
func AddAll() error { return Interactive("add", "-A") }

// Commit creates a commit with the given message.
func Commit(message string) error { return Interactive("commit", "-m", message) }

// Tag creates an annotated tag.
func Tag(version, message string) error {
	return Interactive("tag", "-a", version, "-m", message)
}

// PushTags pushes commits and any annotated tags they point to.
func PushTags() error { return Interactive("push", "--follow-tags") }

// Stash sets aside in-progress work.
func Stash(message string) error {
	if message == "" {
		return Interactive("stash", "push")
	}
	return Interactive("stash", "push", "-m", message)
}

// StashPop brings the most recently stashed work back.
func StashPop() error { return Interactive("stash", "pop") }

// DiscardAll reverts tracked files to their last committed state.
func DiscardAll() error {
	_ = Interactive("restore", "--staged", ".")
	return Interactive("restore", ".")
}

// CleanDryRun previews which untracked files would be removed.
func CleanDryRun() (string, error) { return run("clean", "-nd") }

// Clean removes untracked files and directories.
func Clean() error { return Interactive("clean", "-fd") }
