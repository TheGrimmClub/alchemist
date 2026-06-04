package tui

import (
	"os/exec"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type bottlePhase int

const (
	bottlePhaseConfirmPush bottlePhase = iota
	bottlePhaseConfirmTag
	bottlePhaseTagVersion
	bottlePhaseTagMessage
	bottlePhaseTagging
	bottlePhasePushing
	bottlePhaseDone
	bottlePhaseAborted
	bottlePhaseError
)

type bottleTaggedMsg struct{ err error }
type bottlePushedMsg struct{ err error }

// BottleModel optionally tags then pushes, clearing the task state on success.
type BottleModel struct {
	phase    bottlePhase
	spinner  spinner.Model
	inputs   [2]textinput.Model // [0]=version, [1]=tag message
	focus    int
	FinalErr error
}

func NewBottleModel() BottleModel {
	ver := textinput.New()
	ver.Placeholder = "e.g. v1.0.0"
	ver.CharLimit = 30
	ver.Width = 30
	ver.Focus()

	msg := textinput.New()
	msg.Placeholder = "Release …"
	msg.CharLimit = 72
	msg.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	return BottleModel{
		phase:   bottlePhaseConfirmPush,
		spinner: s,
		inputs:  [2]textinput.Model{ver, msg},
	}
}

func (m BottleModel) Init() tea.Cmd { return nil }

func (m BottleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case bottleTaggedMsg:
		if msg.err != nil {
			m.phase = bottlePhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.phase = bottlePhasePushing
		pushCmd := exec.Command("git", "push", "--follow-tags")
		return m, tea.ExecProcess(pushCmd, func(err error) tea.Msg {
			return bottlePushedMsg{err: err}
		})

	case bottlePushedMsg:
		if msg.err != nil {
			m.phase = bottlePhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		_ = state.Clear()
		m.phase = bottlePhaseDone
		return m, tea.Quit

	case tea.KeyMsg:
		switch m.phase {

		case bottlePhaseConfirmPush:
			switch msg.String() {
			case "y", "Y":
				m.phase = bottlePhaseConfirmTag
			case "n", "N", "esc", "q":
				m.phase = bottlePhaseAborted
				return m, tea.Quit
			}

		case bottlePhaseConfirmTag:
			switch msg.String() {
			case "y", "Y":
				m.phase = bottlePhaseTagVersion
				return m, textinput.Blink
			case "n", "N":
				m.phase = bottlePhasePushing
				pushCmd := exec.Command("git", "push", "--follow-tags")
				return m, tea.ExecProcess(pushCmd, func(err error) tea.Msg {
					return bottlePushedMsg{err: err}
				})
			case "esc", "q":
				m.phase = bottlePhaseAborted
				return m, tea.Quit
			}

		case bottlePhaseTagVersion:
			switch msg.String() {
			case "enter":
				if m.inputs[0].Value() == "" {
					return m, nil
				}
				if m.inputs[1].Placeholder == "Release …" {
					m.inputs[1].Placeholder = "Release " + m.inputs[0].Value()
				}
				m.inputs[0].Blur()
				m.focus = 1
				m.inputs[1].Focus()
				m.phase = bottlePhaseTagMessage
				return m, textinput.Blink
			case "esc":
				m.phase = bottlePhaseAborted
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.inputs[0], cmd = m.inputs[0].Update(msg)
				return m, cmd
			}

		case bottlePhaseTagMessage:
			switch msg.String() {
			case "enter":
				tagMsg := m.inputs[1].Value()
				if tagMsg == "" {
					tagMsg = m.inputs[1].Placeholder
				}
				version := m.inputs[0].Value()
				m.inputs[1].Blur()
				m.phase = bottlePhaseTagging
				return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
					return bottleTaggedMsg{err: gitutil.Tag(version, tagMsg)}
				})
			case "esc":
				m.phase = bottlePhaseAborted
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.inputs[1], cmd = m.inputs[1].Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m BottleModel) View() string {
	switch m.phase {
	case bottlePhaseConfirmPush:
		return "\n" + titleStyle.Render("Bottle: push your work") + "\n\n" +
			"Push your commits? " + hintStyle.Render("[y/N]") + " "

	case bottlePhaseConfirmTag:
		return "\n" + titleStyle.Render("Bottle: push your work") + "\n\n" +
			"Tag this as a release? " + hintStyle.Render("[y/N]") + " "

	case bottlePhaseTagVersion:
		return "\n" + labelStyle.Render("Version tag") + "\n" +
			m.inputs[0].View() + "\n\n" +
			hintStyle.Render("Enter to continue • Esc to cancel") + "\n"

	case bottlePhaseTagMessage:
		return "\n" + labelStyle.Render("Tag message") + "\n" +
			m.inputs[1].View() + "\n\n" +
			hintStyle.Render("Enter to create tag • Esc to cancel") + "\n"

	case bottlePhaseTagging:
		return "\n" + m.spinner.View() + " Creating tag…\n"

	case bottlePhasePushing:
		return "\n" + m.spinner.View() + " Pushing…\n"

	case bottlePhaseDone:
		return "\n" + successStyle.Render("* Bottled and pushed! Nice work.") + "\n"

	case bottlePhaseAborted:
		return "\n" + hintStyle.Render("Cancelled — nothing was pushed.") + "\n"

	case bottlePhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}
