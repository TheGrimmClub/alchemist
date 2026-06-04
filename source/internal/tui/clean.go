package tui

import (
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cleanPhase int

const (
	cleanPhaseLoading cleanPhase = iota
	cleanPhaseNone
	cleanPhaseConfirm
	cleanPhaseCleaning
	cleanPhaseDone
	cleanPhaseAborted
	cleanPhaseError
)

type cleanLoadedMsg struct {
	files []string
	err   error
}

type cleanDoneMsg struct{ err error }

// CleanModel removes untracked files after confirmation.
type CleanModel struct {
	phase   cleanPhase
	files   []string
	spinner spinner.Model
	FinalErr error
}

func NewCleanModel() CleanModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	return CleanModel{phase: cleanPhaseLoading, spinner: s}
}

func (m CleanModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		files, err := gitutil.CleanDryRun()
		return cleanLoadedMsg{files: files, err: err}
	})
}

func (m CleanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case cleanLoadedMsg:
		if msg.err != nil {
			m.phase = cleanPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.files = msg.files
		if len(m.files) == 0 {
			m.phase = cleanPhaseNone
			return m, tea.Quit
		}
		m.phase = cleanPhaseConfirm
		return m, nil

	case cleanDoneMsg:
		if msg.err != nil {
			m.phase = cleanPhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.phase = cleanPhaseDone
		return m, tea.Quit

	case tea.KeyMsg:
		if m.phase != cleanPhaseConfirm {
			return m, nil
		}
		switch msg.String() {
		case "y", "Y":
			m.phase = cleanPhaseCleaning
			return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
				return cleanDoneMsg{err: gitutil.Clean()}
			})
		case "n", "N", "esc", "q":
			m.phase = cleanPhaseAborted
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m CleanModel) View() string {
	switch m.phase {
	case cleanPhaseLoading:
		return "\n" + m.spinner.View() + " Loading…\n"

	case cleanPhaseNone:
		return "\n" + hintStyle.Render("Nothing to clean — no untracked files.") + "\n"

	case cleanPhaseConfirm:
		s := "\n" + dangerStyle.Render("These untracked files will be deleted:") + "\n"
		s += renderFileList(m.files) + "\n"
		s += "Delete them permanently? " + hintStyle.Render("[y/N]") + " "
		return s

	case cleanPhaseCleaning:
		return "\n" + m.spinner.View() + " Cleaning…\n"

	case cleanPhaseDone:
		return "\n" + successStyle.Render("* Untracked files removed.") + "\n"

	case cleanPhaseAborted:
		return "\n" + hintStyle.Render("Cancelled — nothing was deleted.") + "\n"

	case cleanPhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}
