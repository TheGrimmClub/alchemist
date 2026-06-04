package cmd

import (
	"fmt"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
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

		if stashResume {
			fmt.Println("Bringing your paused work back...")
			return gitutil.StashPop()
		}

		fmt.Println("Pausing your work and setting it aside...")
		if err := gitutil.Stash(""); err != nil {
			return err
		}
		fmt.Println("Done — working tree is clean. Run 'alchemist stash --resume' to bring it back.")
		return nil
	},
}

func init() {
	stashCmd.Flags().BoolVar(&stashResume, "resume", false, "bring back previously paused work")
}
