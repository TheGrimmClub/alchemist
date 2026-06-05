package cmd

import (
	"fmt"
	"os"

	"github.com/TheGrimmClub/alchemist/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var validTemplates = []string{"python", "go", "blank"}

var recipeCmd = &cobra.Command{
	Use:   "recipe [template]",
	Short: "Scaffold a new project from a template",
	Long:  "Creates a starter project and initializes git. Templates: python, go, blank.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		template := ""
		if len(args) == 1 {
			template = args[0]
			if !isValidTemplate(template) {
				return fmt.Errorf("unknown template %q — try: python, go, blank", template)
			}
		}
		m, err := tea.NewProgram(tui.NewRecipeModel(template), tea.WithInput(os.Stdin)).Run()
		if err != nil {
			return err
		}
		return m.(tui.RecipeModel).FinalErr
	},
}

func isValidTemplate(name string) bool {
	for _, t := range validTemplates {
		if t == name {
			return true
		}
	}
	return false
}
