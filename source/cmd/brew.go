package cmd

import (
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var brewCmd = &cobra.Command{
	Use:   "brew",
	Short: "Commit your work with a detailed, editable message",
	Long: `Stages your changes and drafts a commit message:

  - the title comes from the task you named with 'alchemist start'
  - the body lists the changed files

The draft opens in your editor so you can refine it before committing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		m, err := tea.NewProgram(tui.NewBrewModel()).Run()
		if err != nil {
			return err
		}
		return m.(tui.BrewModel).FinalErr
	},
}
