package cmd

import (
	"fmt"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/prompt"
	"github.com/spf13/cobra"
)

var discardCmd = &cobra.Command{
	Use:   "discard",
	Short: "Throw away uncommitted changes to tracked files",
	Long: `Reverts tracked files to their last committed state. This cannot be
undone. Untracked files are left alone — use 'alchemist clean' for those.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}

		changed, err := gitutil.ChangedFiles()
		if err != nil {
			return err
		}
		if len(changed) == 0 {
			fmt.Println("Nothing to discard — no changes to tracked files.")
			return nil
		}

		fmt.Println("These changes will be permanently thrown away:")
		for _, f := range changed {
			fmt.Printf("  - %s\n", f)
		}
		fmt.Println()

		if !prompt.Confirm("This cannot be undone. Discard them?", false) {
			fmt.Println("Cancelled — nothing was changed.")
			return nil
		}
		if err := gitutil.DiscardAll(); err != nil {
			return err
		}
		fmt.Println("Changes discarded.")
		return nil
	},
}
