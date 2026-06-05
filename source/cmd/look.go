package cmd

import (
	"os"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var lookCmd = &cobra.Command{
	Use:     "look",
	Aliases: []string{"hm"},
	Short:   "Check your current task and changed files",
	Long: `Shows the task you named with 'alchemist start' and lists all files
that have changed since your last commit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		m, err := tea.NewProgram(tui.NewLookModel(), tea.WithInput(os.Stdin)).Run()
		if err != nil {
			return err
		}
		return m.(tui.LookModel).FinalErr
	},
}
