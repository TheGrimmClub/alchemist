// Package gitutil wraps go-git so the rest of Alchemist can stay readable.
// Stash operations fall back to the system git binary because go-git does not
// implement stash.
package gitutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func openRepo() (*git.Repository, error) {
	return git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
}

// Interactive runs a system git command wired to the terminal. Used for stash
// (not supported in go-git) and as a push fallback when SSH agent auth fails.
func Interactive(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Init creates a new git repository at path.
func Init(path string) error {
	_, err := git.PlainInit(path, false)
	return err
}

// IsRepo reports whether the current directory is inside a git repo.
func IsRepo() bool {
	_, err := openRepo()
	return err == nil
}

// EnsureRepo returns a friendly error if not inside a git repo.
func EnsureRepo() error {
	if !IsRepo() {
		return fmt.Errorf("not a git repository — run 'git init' or 'alchemist recipe' first")
	}
	return nil
}

// ChangedFiles lists tracked files with uncommitted changes (staged or unstaged).
func ChangedFiles() ([]string, error) {
	repo, err := openRepo()
	if err != nil {
		return nil, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := w.Status()
	if err != nil {
		return nil, err
	}
	var files []string
	for path, s := range status {
		if s.Worktree != git.Untracked && (s.Staging != git.Unmodified || s.Worktree != git.Unmodified) {
			files = append(files, path)
		}
	}
	sort.Strings(files)
	return files, nil
}

// UntrackedFiles lists files git is not yet tracking.
func UntrackedFiles() ([]string, error) {
	repo, err := openRepo()
	if err != nil {
		return nil, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := w.Status()
	if err != nil {
		return nil, err
	}
	var files []string
	for path, s := range status {
		if s.Worktree == git.Untracked {
			files = append(files, path)
		}
	}
	sort.Strings(files)
	return files, nil
}

// AddAll stages every change.
func AddAll() error {
	repo, err := openRepo()
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	return w.AddWithOptions(&git.AddOptions{All: true})
}

// authorFromConfig builds a git signature from local then global git config.
func authorFromConfig(repo *git.Repository) *object.Signature {
	name, email := "", ""
	if cfg, err := repo.Config(); err == nil {
		name = cfg.User.Name
		email = cfg.User.Email
	}
	if name == "" || email == "" {
		if cfg, err := config.LoadConfig(config.GlobalScope); err == nil {
			if name == "" {
				name = cfg.User.Name
			}
			if email == "" {
				email = cfg.User.Email
			}
		}
	}
	return &object.Signature{Name: name, Email: email, When: time.Now()}
}

// Commit creates a commit with the given message.
func Commit(message string) error {
	repo, err := openRepo()
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Commit(message, &git.CommitOptions{Author: authorFromConfig(repo)})
	return err
}

// Tag creates an annotated tag at HEAD.
func Tag(version, message string) error {
	repo, err := openRepo()
	if err != nil {
		return err
	}
	head, err := repo.Head()
	if err != nil {
		return err
	}
	_, err = repo.CreateTag(version, head.Hash(), &git.CreateTagOptions{
		Message: message,
		Tagger:  authorFromConfig(repo),
	})
	return err
}

// PushTags pushes commits and annotated tags. Tries SSH agent auth first;
// falls back to the system git binary so existing credential helpers work.
func PushTags() error {
	repo, err := openRepo()
	if err != nil {
		return Interactive("push", "--follow-tags")
	}
	if auth, err := gitssh.NewSSHAgentAuth("git"); err == nil {
		pushErr := repo.Push(&git.PushOptions{Auth: auth, FollowTags: true})
		if pushErr == nil || pushErr == git.NoErrAlreadyUpToDate {
			return nil
		}
	}
	return Interactive("push", "--follow-tags")
}

// Stash sets aside in-progress work (uses system git — no go-git stash API).
func Stash(message string) error {
	if message == "" {
		return Interactive("stash", "push")
	}
	return Interactive("stash", "push", "-m", message)
}

// StashPop brings the most recently stashed work back.
func StashPop() error { return Interactive("stash", "pop") }

// DiscardAll reverts all tracked files to their last committed state and
// unstages any staged changes.
func DiscardAll() error {
	repo, err := openRepo()
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	head, err := repo.Head()
	if err != nil {
		return err
	}
	return w.Reset(&git.ResetOptions{Commit: head.Hash(), Mode: git.HardReset})
}

// CleanDryRun returns the list of untracked files that Clean would remove.
func CleanDryRun() ([]string, error) {
	return UntrackedFiles()
}

// Clean removes all untracked files (and any directories that become empty).
func Clean() error {
	files, err := UntrackedFiles()
	if err != nil {
		return err
	}
	dirs := map[string]struct{}{}
	for _, f := range files {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
		for dir := filepath.Dir(f); dir != "."; dir = filepath.Dir(dir) {
			dirs[dir] = struct{}{}
		}
	}
	// Remove directories deepest-first so parents are cleaned up too.
	sorted := make([]string, 0, len(dirs))
	for d := range dirs {
		sorted = append(sorted, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(sorted)))
	for _, d := range sorted {
		_ = os.Remove(d) // no-op when directory is non-empty
	}
	return nil
}
