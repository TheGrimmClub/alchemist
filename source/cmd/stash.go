package cmd

import (
	"os"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var stashResume bool

var stashCmd = &cobra.Command{
	Use:   "stash",
	Short: "Pause your work and set it aside to resume later",
	Long: `Saves your in-progress changes without committing, leaving a clean
working tree so you can switch gears. Run with --resume to bring the work back.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		m, err := tea.NewProgram(tui.NewStashModel(stashResume), tea.WithInput(os.Stdin)).Run()
		if err != nil {
			return err
		}
		return m.(tui.StashModel).FinalErr
	},
}

func init() {
	stashCmd.Flags().BoolVar(&stashResume, "resume", false, "bring back previously paused work")
}
