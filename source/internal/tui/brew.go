package tui

import (
	"strings"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type brewPhase int

const (
	brewPhaseLoading brewPhase = iota
	brewPhaseNoChanges
	brewPhaseConfirm
	brewPhaseStaging
	brewPhaseCommitting
	brewPhaseDone
	brewPhaseAborted
	brewPhaseError
)

type brewLoadedMsg struct {
	files []string
	err   error
}

type brewStagedMsg struct{ err error }
type brewCommittedMsg struct{ err error }

// BrewModel stages changes and commits them with the active task's name.
type BrewModel struct {
	phase    brewPhase
	task     state.State
	files    []string
	spinner  spinner.Model
	FinalErr error
}

func NewBrewModel(task state.State) BrewModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	return BrewModel{phase: brewPhaseLoading, task: task, spinner: s}
}

func (m BrewModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		changed, err := gitutil.ChangedFiles()
		if err != nil {
			return brewLoadedMsg{err: err}
		}
		untracked, err := gitutil.UntrackedFiles()
		if err != nil {
			return brewLoadedMsg{err: err}
		}
		return brewLoadedMsg{files: append(changed, untracked...)}
	})
}

func (m BrewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case brewLoadedMsg:
		if msg.err != nil {
			m.phase = brewPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.files = msg.files
		if len(m.files) == 0 {
			m.phase = brewPhaseNoChanges
			return m, tea.Quit
		}
		m.phase = brewPhaseConfirm
		return m, nil

	case brewStagedMsg:
		if msg.err != nil {
			m.phase = brewPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		commitMsg := buildCommitMessage(m.task)
		m.phase = brewPhaseCommitting
		return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
			return brewCommittedMsg{err: gitutil.Commit(commitMsg)}
		})

	case brewCommittedMsg:
		if msg.err != nil {
			m.phase = brewPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.phase = brewPhaseDone
		return m, tea.Quit

	case tea.KeyMsg:
		if m.phase != brewPhaseConfirm {
			return m, nil
		}
		switch msg.String() {
		case "y", "Y":
			m.phase = brewPhaseStaging
			return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
				return brewStagedMsg{err: gitutil.AddAll()}
			})
		case "n", "N", "esc", "q":
			m.phase = brewPhaseAborted
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m BrewModel) View() string {
	switch m.phase {
	case brewPhaseLoading:
		return "\n" + m.spinner.View() + " Loading…\n"

	case brewPhaseNoChanges:
		return "\n" + hintStyle.Render("Nothing to commit — your working tree is clean.") + "\n"

	case brewPhaseConfirm:
		s := "\n" + titleStyle.Render("Brew: commit your work") + "\n"
		s += hintStyle.Render("Task: "+m.task.TaskName) + "\n\n"
		s += labelStyle.Render("Files with changes") + "\n"
		s += renderFileList(m.files)
		s += "\nStage and commit all of these? " + hintStyle.Render("[y/N]") + " "
		return s

	case brewPhaseStaging:
		return "\n" + m.spinner.View() + " Staging changes…\n"

	case brewPhaseCommitting:
		return "\n" + m.spinner.View() + " Committing…\n"

	case brewPhaseDone:
		return "\n" + successStyle.Render("* Brewed! Committed as: "+m.task.TaskName) + "\n" +
			hintStyle.Render("Run 'alchemist bottle' when you're ready to tag and push.") + "\n"

	case brewPhaseAborted:
		return "\n" + hintStyle.Render("Cancelled — nothing was committed.") + "\n"

	case brewPhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}

func buildCommitMessage(s state.State) string {
	msg := s.TaskName
	if s.Description != "" {
		msg += "\n\n" + strings.TrimSpace(s.Description)
	}
	return msg
}
