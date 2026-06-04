package tui

import (
	"os"
	"strings"

	"github.com/TheGrimmClub/alchemist/internal/editor"
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
	task      state.State
	hasTask   bool
	files     []string
	err       error
}

type brewStagedMsg struct{ err error }

type brewEditorResultMsg struct {
	content string
	err     error
}

type brewCommittedMsg struct{ err error }

// BrewModel stages changes, opens the editor, and commits.
type BrewModel struct {
	phase    brewPhase
	task     state.State
	hasTask  bool
	files    []string
	message  string
	tmpFile  string
	spinner  spinner.Model
	FinalErr error
}

func NewBrewModel() BrewModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	return BrewModel{phase: brewPhaseLoading, spinner: s}
}

func (m BrewModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		s, _ := state.Load()
		changed, err := gitutil.ChangedFiles()
		if err != nil {
			return brewLoadedMsg{err: err}
		}
		untracked, err := gitutil.UntrackedFiles()
		if err != nil {
			return brewLoadedMsg{err: err}
		}
		return brewLoadedMsg{
			task:    s,
			hasTask: s.TaskName != "",
			files:   append(changed, untracked...),
		}
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
		m.task = msg.task
		m.hasTask = msg.hasTask
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
		var writeErr error
		m.tmpFile, writeErr = editor.WriteTempFile(buildCommitMessage(m.task, m.files))
		if writeErr != nil {
			m.phase = brewPhaseError
			m.FinalErr = writeErr
			return m, tea.Quit
		}
		editorCmd := editor.Command(m.tmpFile)
		return m, tea.ExecProcess(editorCmd, func(err error) tea.Msg {
			content, readErr := os.ReadFile(m.tmpFile)
			os.Remove(m.tmpFile)
			if err != nil {
				return brewEditorResultMsg{err: err}
			}
			return brewEditorResultMsg{content: string(content), err: readErr}
		})

	case brewEditorResultMsg:
		if msg.err != nil {
			m.phase = brewPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.message = stripCommitComments(msg.content)
		if strings.TrimSpace(m.message) == "" {
			m.phase = brewPhaseAborted
			return m, tea.Quit
		}
		m.phase = brewPhaseCommitting
		return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
			return brewCommittedMsg{err: gitutil.Commit(m.message)}
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
		return "\n" + hintStyle.Render("Nothing in the cauldron — your working tree is clean.") + "\n"

	case brewPhaseConfirm:
		s := "\n"
		if !m.hasTask {
			s += warnStyle.Render("No task in progress — run 'alchemist start' first.") + "\n\n"
		} else {
			s += titleStyle.Render("Brew: commit your work") + "\n"
			s += hintStyle.Render("Task: "+m.task.TaskName) + "\n\n"
		}
		s += labelStyle.Render("Files with changes") + "\n"
		s += renderFileList(m.files)
		s += "\n"
		s += "Bottle all of these files? " + hintStyle.Render("[y/N]") + " "
		return s

	case brewPhaseStaging:
		return "\n" + m.spinner.View() + " Staging changes…\n"

	case brewPhaseCommitting:
		return "\n" + m.spinner.View() + " Committing…\n"

	case brewPhaseDone:
		return "\n" + successStyle.Render("* Brewed! Your work is committed.") + "\n" +
			hintStyle.Render("Run 'alchemist bottle' when you're ready to tag and push.") + "\n"

	case brewPhaseAborted:
		return "\n" + hintStyle.Render("Cancelled — nothing was committed.") + "\n"

	case brewPhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}

func buildCommitMessage(s state.State, files []string) string {
	var b strings.Builder
	b.WriteString(s.TaskName)
	b.WriteString("\n\n")
	if s.Description != "" {
		b.WriteString(s.Description)
		b.WriteString("\n\n")
	}
	b.WriteString("Files changed:\n")
	for _, f := range files {
		b.WriteString("  - ")
		b.WriteString(f)
		b.WriteString("\n")
	}
	b.WriteString("\n# The first line is your commit title. Lines starting with '#' are ignored.\n")
	return b.String()
}

func stripCommitComments(s string) string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
