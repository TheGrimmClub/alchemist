package cmd

import (
	"fmt"
	"time"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/prompt"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Name the task you're about to work on",
	Long: `Records the task name and a short description before you write any code.

The task name becomes the title of your commit when you run 'alchemist brew',
which encourages you to decide what you're doing before you do it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}

		if existing, err := state.Load(); err == nil && existing.TaskName != "" {
			fmt.Printf("A task is already in progress: %q\n", existing.TaskName)
			if !prompt.Confirm("Replace it with a new task?", false) {
				return nil
			}
		}

		name := prompt.Line("Task name (becomes your commit title): ")
		for name == "" {
			fmt.Println("A task name is required.")
			name = prompt.Line("Task name: ")
		}
		desc := prompt.Line("Short description (optional): ")

		if err := state.Save(state.State{
			TaskName:    name,
			Description: desc,
			CreatedAt:   time.Now(),
		}); err != nil {
			return err
		}

		fmt.Printf("\n* Task started: %s\n", name)
		fmt.Println("Now write your code. When you're done, run 'alchemist brew' to commit.")
		return nil
	},
}
