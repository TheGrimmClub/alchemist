package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	dangerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
	hintStyle    = lipgloss.NewStyle().Faint(true)
	itemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	cursorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
)

func renderFileList(files []string) string {
	s := ""
	for _, f := range files {
		s += itemStyle.Render("  - "+f) + "\n"
	}
	return s
}
