package cmd

import (
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove untracked files (build artifacts, stray files)",
	Long: `Deletes files git is not tracking, such as build output and stray
files. Tracked files are left alone — use 'alchemist discard' for those.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		m, err := tea.NewProgram(tui.NewCleanModel()).Run()
		if err != nil {
			return err
		}
		return m.(tui.CleanModel).FinalErr
	},
}
