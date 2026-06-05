package tui

import (
	"fmt"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/state"
	tea "github.com/charmbracelet/bubbletea"
)

type lookLoadedMsg struct {
	task      state.State
	hasTask   bool
	changed   []string
	untracked []string
	err       error
}

// HmModel displays the current task and changed files.
type LookModel struct {
	task      state.State
	hasTask   bool
	changed   []string
	untracked []string
	loading   bool
	FinalErr  error
}

func NewLookModel() LookModel {
	return LookModel{loading: true}
}

func (m LookModel) Init() tea.Cmd {
	return func() tea.Msg {
		s, err := state.Load()
		if err != nil {
			return lookLoadedMsg{err: err}
		}
		changed, err := gitutil.ChangedFiles()
		if err != nil {
			return lookLoadedMsg{err: err}
		}
		untracked, err := gitutil.UntrackedFiles()
		if err != nil {
			return lookLoadedMsg{err: err}
		}
		return lookLoadedMsg{
			task:      s,
			hasTask:   s.TaskName != "",
			changed:   changed,
			untracked: untracked,
		}
	}
}

func (m LookModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case lookLoadedMsg:
		if msg.err != nil {
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.loading = false
		m.task = msg.task
		m.hasTask = msg.hasTask
		m.changed = msg.changed
		m.untracked = msg.untracked
		return m, nil
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m LookModel) View() string {
	if m.loading {
		return ""
	}

	s := "\n"
	if !m.hasTask {
		s += warnStyle.Render("Your cauldron is empty — no task in progress.") + "\n"
		s += hintStyle.Render("Run 'alchemist start' to begin.") + "\n"
	} else {
		s += titleStyle.Render("Current task") + "\n"
		s += fmt.Sprintf("  %s\n", m.task.TaskName)
		if m.task.Description != "" {
			s += fmt.Sprintf("  %s\n", hintStyle.Render(m.task.Description))
		}
	}
	s += "\n"

	all := append(m.changed, m.untracked...)
	if len(all) == 0 {
		s += hintStyle.Render("No ingredients in the cauldron — your project files are clean.") + "\n"
	} else {
		s += labelStyle.Render("Files with changes") + "\n"
		s += renderFileList(all)
		s += "\n"
		s += hintStyle.Render("alchemist brew   → commit  |  alchemist bottle → push") + "\n"
	}

	s += "\n" + hintStyle.Render("Press any key to exit.") + "\n"
	return s
}
