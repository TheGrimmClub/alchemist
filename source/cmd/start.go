package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"next"},
	Short:   "Name the task you're about to work on",
	Long: `Records the task name and a short description before you write any code.

The task name becomes the title of your commit when you run 'alchemist brew',
which encourages you to decide what you're doing before you do it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}
		existing, _ := state.Load()
		result, err := tea.NewProgram(
			tui.NewStartModel(existing.TaskName),
			tea.WithInput(os.Stdin),
		).Run()
		if err != nil {
			return err
		}
		m := result.(tui.StartModel)
		if m.FinalErr != nil {
			return m.FinalErr
		}
		if m.Aborted || m.Name() == "" {
			return nil
		}
		if err := state.Save(state.State{
			TaskName:    m.Name(),
			Description: m.Desc(),
			CreatedAt:   time.Now(),
		}); err != nil {
			return err
		}
		fmt.Printf("Task started: %s\n", m.Name())
		return nil
	},
}
