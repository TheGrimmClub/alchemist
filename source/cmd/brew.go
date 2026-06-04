package cmd

import (
	"fmt"
	"strings"

	"github.com/TheGrimmClub/alchemist/internal/editor"
	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/prompt"
	"github.com/TheGrimmClub/alchemist/internal/state"
	"github.com/spf13/cobra"
)

var brewCmd = &cobra.Command{
	Use:   "brew",
	Short: "Commit your work with a detailed, editable message",
	Long: `Stages your changes and drafts a commit message:

  - the title comes from the task you named with 'alchemist start'
  - the body is built from your one-line summary and the changed files

The draft opens in your editor (micro by default) so you can refine it
before the commit is made.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitutil.EnsureRepo(); err != nil {
			return err
		}

		s, err := state.Load()
		if err != nil || s.TaskName == "" {
			return fmt.Errorf("no task in progress — run 'alchemist start' first")
		}

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
			fmt.Println("Nothing to commit — your working tree is clean.")
			return nil
		}

		fmt.Println("Files with changes:")
		for _, f := range files {
			fmt.Printf("  - %s\n", f)
		}
		fmt.Println()

		if !prompt.Confirm("Stage all of these?", true) {
			fmt.Println("Stage what you want with 'git add', then run 'alchemist brew' again.")
			return nil
		}
		if err := gitutil.AddAll(); err != nil {
			return err
		}

		summary := prompt.Line("In one line, what did you do? ")

		edited, err := editor.Edit(buildMessage(s, summary, files))
		if err != nil {
			return err
		}
		message := stripComments(edited)
		if strings.TrimSpace(message) == "" {
			return fmt.Errorf("empty commit message — aborting")
		}

		if err := gitutil.Commit(message); err != nil {
			return err
		}
		fmt.Println("\n* Brewed! Your work is committed.")
		fmt.Println("Run 'alchemist bottle' when you're ready to tag and push.")
		return nil
	},
}

func buildMessage(s state.State, summary string, files []string) string {
	var b strings.Builder
	b.WriteString(s.TaskName)
	b.WriteString("\n\n")
	if summary != "" {
		b.WriteString(summary)
		b.WriteString("\n\n")
	}
	if s.Description != "" {
		b.WriteString(s.Description)
		b.WriteString("\n\n")
	}
	b.WriteString("Files changed:\n")
	for _, f := range files {
		b.WriteString("  - ")
		b.WriteString(f)
		b.WriteString("\n")
	}
	b.WriteString("\n# The first line is your commit title. Lines starting with '#' are ignored.\n")
	return b.String()
}

func stripComments(s string) string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
