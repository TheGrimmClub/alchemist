package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TheGrimmClub/alchemist/internal/gitutil"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type recipePhase int

const (
	recipePhaseSelect  recipePhase = iota // choose template
	recipePhaseNaming                     // enter project name
	recipePhaseCreating                   // creating files
	recipePhaseDone
	recipePhaseError
)

type recipeCreatedMsg struct {
	project string
	tmpl    string
	created []string
	err     error
}

type recipeTemplate struct {
	name        string
	description string
	files       []recipeFile
}

type recipeFile struct {
	path    string
	content string
}

const ignoreAlchemist = ".alchemist/\n"

var templates = map[string]recipeTemplate{
	"python": {
		name:        "python",
		description: "A minimal Python project",
		files: []recipeFile{
			{"main.py", "def main():\n    print(\"Hello from Alchemist!\")\n\n\nif __name__ == \"__main__\":\n    main()\n"},
			{"requirements.txt", ""},
			{"README.md", "# {{project}}\n\nA Python project scaffolded with Alchemist.\n\n## Run\n\n    python main.py\n"},
			{".gitignore", "__pycache__/\n*.pyc\n.venv/\nvenv/\n" + ignoreAlchemist},
		},
	},
	"go": {
		name:        "go",
		description: "A minimal Go project",
		files: []recipeFile{
			{"main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello from Alchemist!\")\n}\n"},
			{"README.md", "# {{project}}\n\nA Go project scaffolded with Alchemist.\n\n## Run\n\n    go run .\n"},
			{".gitignore", "/bin/\n*.exe\n" + ignoreAlchemist},
		},
	},
	"blank": {
		name:        "blank",
		description: "Just a README and .gitignore",
		files: []recipeFile{
			{"README.md", "# {{project}}\n\nScaffolded with Alchemist.\n"},
			{".gitignore", ignoreAlchemist},
		},
	},
}

func sortedTemplateKeys() []string {
	keys := make([]string, 0, len(templates))
	for k := range templates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// RecipeModel lets the user pick a template and name a project.
type RecipeModel struct {
	phase    recipePhase
	keys     []string
	cursor   int
	projInput textinput.Model
	result   recipeCreatedMsg
	FinalErr error
}

func NewRecipeModel(template string) RecipeModel {
	input := textinput.New()
	input.Placeholder = "my-project"
	input.CharLimit = 64
	input.Width = 40

	m := RecipeModel{
		keys:      sortedTemplateKeys(),
		projInput: input,
	}

	if template != "" {
		for i, k := range m.keys {
			if k == template {
				m.cursor = i
				break
			}
		}
		m.phase = recipePhaseNaming
		m.projInput.Focus()
	}
	return m
}

func (m RecipeModel) Init() tea.Cmd {
	if m.phase == recipePhaseNaming {
		return textinput.Blink
	}
	return nil
}

func (m RecipeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case recipeCreatedMsg:
		if msg.err != nil {
			m.phase = recipePhaseError
			m.FinalErr = msg.err
			return m, tea.Quit
		}
		m.result = msg
		m.phase = recipePhaseDone
		return m, tea.Quit

	case tea.KeyMsg:
		switch m.phase {

		case recipePhaseSelect:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.keys)-1 {
					m.cursor++
				}
			case "enter", " ":
				m.phase = recipePhaseNaming
				m.projInput.Focus()
				return m, textinput.Blink
			case "ctrl+c", "esc", "q":
				return m, tea.Quit
			}

		case recipePhaseNaming:
			switch msg.String() {
			case "enter":
				project := m.projInput.Value()
				if project == "" {
					project = "my-project"
				}
				tmplKey := m.keys[m.cursor]
				tmpl := templates[tmplKey]
				m.phase = recipePhaseCreating
				return m, func() tea.Msg {
					return createProject(project, tmpl)
				}
			case "esc":
				m.phase = recipePhaseSelect
				m.projInput.Blur()
				return m, nil
			default:
				var cmd tea.Cmd
				m.projInput, cmd = m.projInput.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m RecipeModel) View() string {
	switch m.phase {
	case recipePhaseSelect:
		s := "\n" + titleStyle.Render("Recipe: scaffold a new project") + "\n\n"
		s += labelStyle.Render("Choose a template") + "\n\n"
		for i, k := range m.keys {
			tmpl := templates[k]
			cursor := "  "
			if i == m.cursor {
				cursor = cursorStyle.Render("> ")
			}
			line := fmt.Sprintf("%-8s %s", tmpl.name, hintStyle.Render(tmpl.description))
			if i == m.cursor {
				line = itemStyle.Render(fmt.Sprintf("%-8s", tmpl.name)) + " " + hintStyle.Render(tmpl.description)
			}
			s += cursor + line + "\n"
		}
		s += "\n" + hintStyle.Render("↑/↓ to choose • Enter to select • Esc to quit") + "\n"
		return s

	case recipePhaseNaming:
		tmplKey := m.keys[m.cursor]
		s := "\n" + titleStyle.Render("Recipe: scaffold a new project") + "\n"
		s += hintStyle.Render("Template: "+tmplKey) + "\n\n"
		s += labelStyle.Render("Project name") + "\n"
		s += m.projInput.View() + "\n\n"
		s += hintStyle.Render("Enter to create • Esc to go back") + "\n"
		return s

	case recipePhaseCreating:
		return "\n" + hintStyle.Render("Creating project…") + "\n"

	case recipePhaseDone:
		s := "\n" + successStyle.Render(fmt.Sprintf("* Project %q ready (template: %s)", m.result.project, m.result.tmpl)) + "\n"
		for _, f := range m.result.created {
			s += itemStyle.Render("  - "+f) + "\n"
		}
		s += "\n" + hintStyle.Render(fmt.Sprintf("Next:\n  cd %s\n  alchemist start", m.result.project)) + "\n"
		return s

	case recipePhaseError:
		return "\n" + warnStyle.Render("Error: "+m.FinalErr.Error()) + "\n"

	default:
		return ""
	}
}

func createProject(project string, tmpl recipeTemplate) recipeCreatedMsg {
	if err := os.MkdirAll(project, 0o755); err != nil {
		return recipeCreatedMsg{err: err}
	}
	var created []string
	for _, f := range tmpl.files {
		content := strings.ReplaceAll(f.content, "{{project}}", project)
		dest := filepath.Join(project, f.path)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return recipeCreatedMsg{err: err}
		}
		if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
			return recipeCreatedMsg{err: err}
		}
		created = append(created, dest)
	}
	if err := gitutil.Init(project); err != nil {
		return recipeCreatedMsg{err: err}
	}
	return recipeCreatedMsg{project: project, tmpl: tmpl.name, created: created}
}
