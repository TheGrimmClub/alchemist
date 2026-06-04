package cmd

import (
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var discardCmd = &cobra.Command{
	Use:   "discard",
	Short: "Throw away uncommitted changes to tracked files",
	Long: `Reverts tracked files to their last committed state. This cannot be
undone. Untracked files are left alone — use 'alchemist clean' for those.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		m, err := tea.NewProgram(tui.NewDiscardModel()).Run()
		if err != nil {
			return err
		}
		return m.(tui.DiscardModel).FinalErr
	},
}
