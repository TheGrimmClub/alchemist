package cmd

import (
	"fmt"
	"os"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var brewCmd = &cobra.Command{
	Use:     "brew",
	Aliases: []string{"commit"},
	Short:   "Commit your work with a meaningful message",
	Long: `Stages your changes and commits them using the task name you set with
'alchemist start' as the commit title. The description (if any) becomes the
commit body.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		s, err := state.Load()
		if err != nil || s.TaskName == "" {
			return fmt.Errorf("no task in progress — run 'alchemist start' first")
		}
		m, err := tea.NewProgram(tui.NewBrewModel(s), tea.WithInput(os.Stdin)).Run()
		if err != nil {
			return err
		}
		return m.(tui.BrewModel).FinalErr
	},
}
