package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alchemist",
	Short: "Alchemist — a task-first git workflow for students",
	Long: `Alchemist guides you through a task-first git workflow.

  recipe   Scaffold a new project from a template
  start    Name the task you're about to work on
  brew     Commit your work with a detailed, editable message
  bottle   Finalize: optional version tag and push
  discard  Throw away uncommitted changes to tracked files
  stash    Pause work and set it aside to resume later
  clean    Remove untracked files (build artifacts, stray files)
  hm       Find what you where actively working on.

Typical flow:  start  ->  (write code)  ->  brew  ->  bottle`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		recipeCmd,
		startCmd,
		brewCmd,
		bottleCmd,
		discardCmd,
		stashCmd,
		cleanCmd,
		hmCmd,
	)
}
