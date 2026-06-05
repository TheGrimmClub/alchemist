package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type startPhase int

const (
	startPhaseReplace startPhase = iota // confirm replacing existing task
	startPhaseInput                     // collect name + description
)

type startField int

const (
	startFieldName startField = iota
	startFieldDesc
)

// StartModel collects a task name (required) and optional description.
// If existingTask is non-empty, it first asks whether to replace it.
type StartModel struct {
	phase        startPhase
	existingTask string
	inputs       [2]textinput.Model
	focus        startField
	nameErr      string
	Aborted      bool
	FinalErr     error
}

func NewStartModel(existingTask string) StartModel {
	name := textinput.New()
	name.Placeholder = "e.g. Add user login"
	name.CharLimit = 72
	name.Width = 50

	desc := textinput.New()
	desc.Placeholder = "optional"
	desc.CharLimit = 120
	desc.Width = 50

	phase := startPhaseInput
	if existingTask != "" {
		phase = startPhaseReplace
	} else {
		name.Focus()
	}

	return StartModel{
		phase:        phase,
		existingTask: existingTask,
		inputs:       [2]textinput.Model{name, desc},
	}
}

func (m StartModel) Init() tea.Cmd {
	if m.phase == startPhaseInput {
		return textinput.Blink
	}
	return nil
}

func (m StartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.phase {

		case startPhaseReplace:
			switch msg.String() {
			case "y", "Y":
				m.phase = startPhaseInput
				m.inputs[startFieldName].Focus()
				return m, textinput.Blink
			case "n", "N", "esc", "q", "ctrl+c":
				m.Aborted = true
				return m, tea.Quit
			}

		case startPhaseInput:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.Aborted = true
				return m, tea.Quit

			case "tab", "shift+tab":
				if m.focus == startFieldName {
					if m.inputs[startFieldName].Value() == "" {
						m.nameErr = "A task name is required."
						return m, nil
					}
					m.nameErr = ""
					m.inputs[startFieldName].Blur()
					m.focus = startFieldDesc
					m.inputs[startFieldDesc].Focus()
				} else {
					m.inputs[startFieldDesc].Blur()
					m.focus = startFieldName
					m.inputs[startFieldName].Focus()
				}
				return m, textinput.Blink

			case "enter":
				if m.focus == startFieldName {
					if m.inputs[startFieldName].Value() == "" {
						m.nameErr = "A task name is required."
						return m, nil
					}
					m.nameErr = ""
					m.inputs[startFieldName].Blur()
					m.focus = startFieldDesc
					m.inputs[startFieldDesc].Focus()
					return m, textinput.Blink
				}
				if m.inputs[startFieldName].Value() == "" {
					m.nameErr = "A task name is required."
					m.inputs[startFieldDesc].Blur()
					m.focus = startFieldName
					m.inputs[startFieldName].Focus()
					return m, textinput.Blink
				}
				return m, tea.Quit
			}

			var cmd tea.Cmd
			m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
			return m, cmd
		}
	}

	if m.phase == startPhaseInput {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m StartModel) View() string {
	switch m.phase {
	case startPhaseReplace:
		s := "\n" + warnStyle.Render(fmt.Sprintf("A task is already in progress: %q", m.existingTask)) + "\n\n"
		s += "Replace it? " + hintStyle.Render("[y/N]") + " "
		return s

	case startPhaseInput:
		s := "\n" + titleStyle.Render("Start a new task") + "\n\n"
		s += labelStyle.Render("Task name") + " " + hintStyle.Render("(becomes your commit title)") + "\n"
		s += m.inputs[startFieldName].View() + "\n"
		if m.nameErr != "" {
			s += warnStyle.Render(m.nameErr) + "\n"
		}
		s += "\n"
		s += labelStyle.Render("Short description") + " " + hintStyle.Render("(optional)") + "\n"
		s += m.inputs[startFieldDesc].View() + "\n\n"
		if m.focus == startFieldName {
			s += hintStyle.Render("Enter/Tab to continue • Esc to cancel")
		} else {
			s += hintStyle.Render("Enter to save • Tab to go back • Esc to cancel")
		}
		s += "\n"
		return s
	}
	return ""
}

func (m StartModel) Name() string { return m.inputs[startFieldName].Value() }
func (m StartModel) Desc() string { return m.inputs[startFieldDesc].Value() }

func (m StartModel) String() string {
	return fmt.Sprintf("name=%q desc=%q", m.Name(), m.Desc())
}
