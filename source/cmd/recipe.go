package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/TheGrimmClub/alchemist/internal/prompt"
	"github.com/spf13/cobra"
)

type templateFile struct {
	path    string
	content string
}

type recipeTemplate struct {
	name        string
	description string
	files       []templateFile
}

// Shared ignore block so every scaffolded project hides Alchemist's task state.
const ignoreAlchemist = ".alchemist/\n"

var recipes = map[string]recipeTemplate{
	"python": {
		name:        "python",
		description: "A minimal Python project",
		files: []templateFile{
			{"main.py", "def main():\n    print(\"Hello from Alchemist!\")\n\n\nif __name__ == \"__main__\":\n    main()\n"},
			{"requirements.txt", ""},
			{"README.md", "# {{project}}\n\nA Python project scaffolded with Alchemist.\n\n## Run\n\n    python main.py\n"},
			{".gitignore", "__pycache__/\n*.pyc\n.venv/\nvenv/\n" + ignoreAlchemist},
		},
	},
	"go": {
		name:        "go",
		description: "A minimal Go project",
		files: []templateFile{
			{"main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello from Alchemist!\")\n}\n"},
			{"README.md", "# {{project}}\n\nA Go project scaffolded with Alchemist.\n\n## Run\n\n    go run .\n"},
			{".gitignore", "/bin/\n*.exe\n" + ignoreAlchemist},
		},
	},
	"blank": {
		name:        "blank",
		description: "Just a README and .gitignore",
		files: []templateFile{
			{"README.md", "# {{project}}\n\nScaffolded with Alchemist.\n"},
			{".gitignore", ignoreAlchemist},
		},
	},
}

var recipeCmd = &cobra.Command{
	Use:   "recipe [template]",
	Short: "Scaffold a new project from a template",
	Long:  "Creates a starter project and initializes git. Templates: python, go, blank.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) == 1 {
			name = args[0]
		} else {
			fmt.Println("Available templates:")
			for _, key := range sortedRecipeKeys() {
				fmt.Printf("  %-8s %s\n", key, recipes[key].description)
			}
			name = prompt.LineDefault("Choose a template", "python")
		}

		tmpl, ok := recipes[name]
		if !ok {
			return fmt.Errorf("unknown template %q (try: python, go, blank)", name)
		}

		project := prompt.LineDefault("Project name", "my-project")
		if err := os.MkdirAll(project, 0o755); err != nil {
			return err
		}

		for _, f := range tmpl.files {
			content := strings.ReplaceAll(f.content, "{{project}}", project)
			dest := filepath.Join(project, f.path)
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
				return err
			}
			fmt.Printf("  created %s\n", dest)
		}

		if err := os.Chdir(project); err == nil {
			_ = gitutil.Interactive("init", "-q")
		}

		fmt.Printf("\n* Project %q ready (template: %s)\n", project, tmpl.name)
		fmt.Printf("Next:\n  cd %s\n  alchemist start\n", project)
		return nil
	},
}

func sortedRecipeKeys() []string {
	keys := make([]string, 0, len(recipes))
	for k := range recipes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
