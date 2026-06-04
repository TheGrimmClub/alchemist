package cmd

import (
	"fmt"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/prompt"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/spf13/cobra"
)

var bottleCmd = &cobra.Command{
	Use:   "bottle",
	Short: "Finalize your work: optional version tag and push",
	Long: `Pushes your committed work. Optionally creates an annotated version
tag first. Once everything is pushed, the current task is cleared so you can
start fresh.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}

		if changed, _ := gitutil.ChangedFiles(); len(changed) > 0 {
			fmt.Println("You still have uncommitted changes.")
			if prompt.Confirm("Commit them with 'brew' first?", true) {
				fmt.Println("Run 'alchemist brew', then 'alchemist bottle' again.")
				return nil
			}
		}

		if prompt.Confirm("Tag this as a release?", false) {
			version := prompt.Line("Version tag (e.g. v1.0.0): ")
			if version != "" {
				msg := prompt.LineDefault("Tag message", "Release "+version)
				if err := gitutil.Tag(version, msg); err != nil {
					return err
				}
				fmt.Printf("Tagged %s\n", version)
			}
		}

		if err := gitutil.PushTags(); err != nil {
			return err
		}

		_ = state.Clear()
		fmt.Println("\n* Bottled and pushed! Nice work.")
		return nil
	},
}
