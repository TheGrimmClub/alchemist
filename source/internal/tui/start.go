package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type startField int

const (
	startFieldName startField = iota
	startFieldDesc
)

// StartModel collects a task name (required) and optional description.
type StartModel struct {
	inputs   [2]textinput.Model
	focus    startField
	nameErr  string
	Aborted  bool
	FinalErr error
}

func NewStartModel() StartModel {
	name := textinput.New()
	name.Placeholder = "e.g. Add user login"
	name.Focus()
	name.CharLimit = 72
	name.Width = 50

	desc := textinput.New()
	desc.Placeholder = "optional"
	desc.CharLimit = 120
	desc.Width = 50

	return StartModel{inputs: [2]textinput.Model{name, desc}}
}

func (m StartModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			// Submit from desc field
			if m.inputs[startFieldName].Value() == "" {
				m.nameErr = "A task name is required."
				m.inputs[startFieldDesc].Blur()
				m.focus = startFieldName
				m.inputs[startFieldName].Focus()
				return m, textinput.Blink
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m StartModel) View() string {
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

func (m StartModel) Name() string { return m.inputs[startFieldName].Value() }
func (m StartModel) Desc() string { return m.inputs[startFieldDesc].Value() }

func (m StartModel) String() string {
	return fmt.Sprintf("name=%q desc=%q", m.Name(), m.Desc())
}
