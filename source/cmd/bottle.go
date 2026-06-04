package cmd

import (
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var bottleCmd = &cobra.Command{
	Use:   "bottle",
	Short: "Finalize your work: optional version tag and push",
	Long: `Pushes your committed work. Optionally creates an annotated version
tag first. Once everything is pushed, the current task is cleared.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		m, err := tea.NewProgram(tui.NewBottleModel()).Run()
		if err != nil {
			return err
		}
		return m.(tui.BottleModel).FinalErr
	},
}
