package cmd

import (
	"fmt"
	"strings"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/prompt"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove untracked files (build artifacts, stray files)",
	Long: `Deletes files git is not tracking, such as build output and stray
files. Tracked files and their changes are left alone — use 'alchemist discard'
for those.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}

		preview, err := gitutil.CleanDryRun()
		if err != nil {
			return err
		}
		if strings.TrimSpace(preview) == "" {
			fmt.Println("Nothing to clean — no untracked files.")
			return nil
		}

		fmt.Println("These untracked files will be deleted:")
		for _, line := range strings.Split(preview, "\n") {
			name := strings.TrimPrefix(strings.TrimSpace(line), "Would remove ")
			fmt.Printf("  - %s\n", name)
		}
		fmt.Println()

		if !prompt.Confirm("Delete them permanently?", false) {
			fmt.Println("Cancelled — nothing was deleted.")
			return nil
		}
		if err := gitutil.Clean(); err != nil {
			return err
		}
		fmt.Println("Untracked files removed.")
		return nil
	},
}
