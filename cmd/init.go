package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chadsmith12/dotsec/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type step int

const (
	stepFolder step = iota
	stepType
	stepPath
	stepConfirm
	stepDone
)

type model struct {
	step      step
	input     string
	cursor    int
	done      bool
	cancelled bool
	err       error

	folder       string
	secretType   string
	path         string
	typeOptions  []string
	selectedType int
}

func initialModel() model {
	return model{
		step:         stepFolder,
		input:        "",
		cursor:       0,
		done:         false,
		cancelled:    false,
		err:          nil,
		folder:       "",
		secretType:   "",
		path:         "",
		typeOptions:  []string{"dotnet", "env"},
		selectedType: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case stepFolder:
			return m.updateFolder(msg)
		case stepType:
			return m.updateType(msg)
		case stepPath:
			return m.updatePath(msg)
		case stepConfirm:
			return m.updateConfirm(msg)
		case stepDone:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) updateFolder(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.cancelled = true
		return m, tea.Quit
	case "enter":
		if strings.TrimSpace(m.input) != "" {
			m.folder = strings.TrimSpace(m.input)
			m.step = stepType
			m.input = ""
		}
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		m.input += msg.String()
	}
	return m, nil
}

func (m model) updateType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.cancelled = true
		return m, tea.Quit
	case "up", "k":
		if m.selectedType > 0 {
			m.selectedType--
		}
	case "down", "j":
		if m.selectedType < len(m.typeOptions)-1 {
			m.selectedType++
		}
	case "enter":
		m.secretType = m.typeOptions[m.selectedType]
		m.step = stepPath
		m.input = ""
		if m.secretType == "env" {
			m.input = ".env"
		}
	}
	return m, nil
}

func (m model) updatePath(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.cancelled = true
		return m, tea.Quit
	case "enter":
		m.path = strings.TrimSpace(m.input)
		if m.path == "" {
			if m.secretType == "env" {
				m.path = ".env"
			} else {
				m.path = "."
			}
		}
		m.step = stepConfirm
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		m.input += msg.String()
	}
	return m, nil
}

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.cancelled = true
		return m, tea.Quit
	case "y", "Y", "enter":
		m.step = stepDone
		m.done = true
		return m, tea.Quit
	case "n", "N":
		m.cancelled = true
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	switch m.step {
	case stepFolder:
		return m.viewFolder()
	case stepType:
		return m.viewType()
	case stepPath:
		return m.viewPath()
	case stepConfirm:
		return m.viewConfirm()
	case stepDone:
		return "✅ Configuration saved to .dotsecrc\n"
	}
	return ""
}

func (m model) viewFolder() string {
	return fmt.Sprintf("🔐 dotsec init\n\n📁 Passbolt folder name: %s\n\n(Press Enter to continue, Ctrl+C to cancel)", m.input+"█")
}

func (m model) viewType() string {
	s := "🔐 dotsec init\n\n🎯 Select secret type:\n\n"

	for i, option := range m.typeOptions {
		cursor := " "
		if i == m.selectedType {
			cursor = "▶"
		}

		description := ""
		if option == "dotnet" {
			description = " - .NET user-secrets"
		} else {
			description = " - Environment files"
		}

		s += fmt.Sprintf("%s %s%s\n", cursor, option, description)
	}

	s += "\n(Use arrow keys or j/k to select, Enter to continue, Ctrl+C to cancel)"
	return s
}

func (m model) viewPath() string {
	var prompt string
	var placeholder string

	if m.secretType == "dotnet" {
		prompt = "📁 Project path (. for current directory):"
		placeholder = "."
	} else {
		prompt = "📄 Environment file path:"
		placeholder = ".env"
	}

	input := m.input
	if input == "" {
		input = placeholder
	}

	return fmt.Sprintf("🔐 dotsec init\n\n%s %s\n\n(Press Enter to continue, Ctrl+C to cancel)", prompt, input+"█")
}

func (m model) viewConfirm() string {
	return fmt.Sprintf(`🔐 dotsec init

📋 Configuration Summary:
   📁 Folder: %s
   🎯 Type: %s
   📄 Path: %s

💾 Save this configuration? (y/N)`, m.folder, m.secretType, m.path)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .dotsecrc file",
	Long: `Initializes a new .dotsecrc file that can be used for project configuration to allow using
	dotsec without additional arguments or flags.`,

	Run: initRun,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRun(cmd *cobra.Command, args []string) {
	m := initialModel()
	program := tea.NewProgram(m)
	finalModel, err := program.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	result := finalModel.(model)
	if result.cancelled {
		fmt.Println("Cancelled.")
		return
	}

	if result.done {
		err := config.WriteProjectConfigWithData(result.folder, result.secretType, result.path)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
