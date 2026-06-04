package tui

import (
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type discardPhase int

const (
	discardPhaseLoading discardPhase = iota
	discardPhaseNoChanges
	discardPhaseConfirm
	discardPhaseDiscarding
	discardPhaseDone
	discardPhaseAborted
	discardPhaseError
)

type discardLoadedMsg struct {
	files []string
	err   error
}

type discardDoneMsg struct{ err error }

// DiscardModel reverts all tracked file changes after confirmation.
type DiscardModel struct {
	phase   discardPhase
	files   []string
	spinner spinner.Model
	FinalErr error
}

func NewDiscardModel() DiscardModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	return DiscardModel{phase: discardPhaseLoading, spinner: s}
}

func (m DiscardModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		files, err := gitutil.ChangedFiles()
		return discardLoadedMsg{files: files, err: err}
	})
}

func (m DiscardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case discardLoadedMsg:
		if msg.err != nil {
			m.phase = discardPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.files = msg.files
		if len(m.files) == 0 {
			m.phase = discardPhaseNoChanges
			return m, tea.Quit
		}
		m.phase = discardPhaseConfirm
		return m, nil

	case discardDoneMsg:
		if msg.err != nil {
			m.phase = discardPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.phase = discardPhaseDone
		return m, tea.Quit

	case tea.KeyMsg:
		if m.phase != discardPhaseConfirm {
			return m, nil
		}
		switch msg.String() {
		case "y", "Y":
			m.phase = discardPhaseDiscarding
			return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
				return discardDoneMsg{err: gitutil.DiscardAll()}
			})
		case "n", "N", "esc", "q":
			m.phase = discardPhaseAborted
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m DiscardModel) View() string {
	switch m.phase {
	case discardPhaseLoading:
		return "\n" + m.spinner.View() + " Loading…\n"

	case discardPhaseNoChanges:
		return "\n" + hintStyle.Render("Nothing to discard — no changes to tracked files.") + "\n"

	case discardPhaseConfirm:
		s := "\n" + dangerStyle.Render("These changes will be permanently thrown away:") + "\n"
		s += renderFileList(m.files) + "\n"
		s += "This cannot be undone. Discard them? " + hintStyle.Render("[y/N]") + " "
		return s

	case discardPhaseDiscarding:
		return "\n" + m.spinner.View() + " Discarding…\n"

	case discardPhaseDone:
		return "\n" + successStyle.Render("* Changes discarded.") + "\n"

	case discardPhaseAborted:
		return "\n" + hintStyle.Render("Cancelled — nothing was changed.") + "\n"

	case discardPhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}
