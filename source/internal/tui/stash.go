package tui

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type stashPhase int

const (
	stashPhaseConfirm stashPhase = iota
	stashPhaseRunning
	stashPhaseDone
	stashPhaseAborted
	stashPhaseError
)

type stashDoneMsg struct{ err error }

// StashModel pauses or resumes work via git stash.
// git stash needs the real terminal (credential/passphrase prompts), so it
// runs via tea.ExecProcess which suspends the TUI.
type StashModel struct {
	resume  bool
	phase   stashPhase
	FinalErr error
}

func NewStashModel(resume bool) StashModel {
	return StashModel{resume: resume}
}

func (m StashModel) Init() tea.Cmd { return nil }

func (m StashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stashDoneMsg:
		if msg.err != nil {
			m.phase = stashPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.phase = stashPhaseDone
		return m, tea.Quit

	case tea.KeyMsg:
		if m.phase != stashPhaseConfirm {
			return m, nil
		}
		switch msg.String() {
		case "y", "Y":
			m.phase = stashPhaseRunning
			var args []string
			if m.resume {
				args = []string{"stash", "pop"}
			} else {
				args = []string{"stash", "push"}
			}
			stashCmd := exec.Command("git", args...)
			return m, tea.ExecProcess(stashCmd, func(err error) tea.Msg {
				return stashDoneMsg{err: err}
			})
		case "n", "N", "esc", "q":
			m.phase = stashPhaseAborted
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m StashModel) View() string {
	switch m.phase {
	case stashPhaseConfirm:
		if m.resume {
			return "\n" + titleStyle.Render("Resume stashed work") + "\n\n" +
				"Bring your paused work back? " + hintStyle.Render("[y/N]") + " "
		}
		return "\n" + titleStyle.Render("Stash: pause your work") + "\n\n" +
			"Pause and stash your current changes? " + hintStyle.Render("[y/N]") + " "

	case stashPhaseRunning:
		return "" // ExecProcess owns the terminal

	case stashPhaseDone:
		if m.resume {
			return "\n" + successStyle.Render("* Work resumed.") + "\n"
		}
		return "\n" + successStyle.Render("* Work stashed. Run 'alchemist stash --resume' to bring it back.") + "\n"

	case stashPhaseAborted:
		return "\n" + hintStyle.Render("Cancelled.") + "\n"

	case stashPhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}
