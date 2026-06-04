package cmd

import (
	"fmt"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/state"

	"github.com/spf13/cobra"
)

var hmCmd = &cobra.Command{
	Use:   "hm",
	Short: "Check your work!",
	Long: `Shows your tasks changes:

  - the title comes from the task you named with 'alchemist start'
  - the files are shown by git

The info is shown on the command line.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}

		s, err := state.Load()
		if err != nil || s.TaskName == "" {
			fmt.Printf("your caldron is empty - no task in progress\n — run 'alchemist start'")
		}

		fmt.Printf("\nTask: %q\n\n", s.TaskName)
		changed, err := gitutil.ChangedFiles()
		if err != nil {
			return err
		}
		untracked, err := gitutil.UntrackedFiles()
		if err != nil {
			return err
		}
		files := append(changed, untracked...)
		if len(files) == 0 {
			fmt.Println("No ingredients in the cauldron\n— your project files are clean.")
			return nil
		}

		fmt.Println("Files with changes:")
		for _, f := range files {
			fmt.Printf("  - %q\n", f)
		}
		fmt.Println()

		fmt.Println("Run 'alchemist brew' when you're tasks is implemented to add and commit.")
		fmt.Println("Run 'alchemist bottle' when you're ready to tag and push.")
		return nil
	},
}
