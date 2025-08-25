package cmd

import (
	"fmt"
	"log"
	"os"

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
	step step
	input string
	cursor int
	done bool
	cancelled bool
	err error
}

func initialModel() model {
	return model{
		step: stepFolder,
		input: "",
		cursor: 0,
		done: false,
		cancelled: false,
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Folder to pull from?"

	return s
}

// configureCmd represents the configure command
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
	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	err := config.WriteProjectConfig()
	if err != nil {
		log.Fatalf("%v", err)
	}
	// server, err := input.PromptUser("Server (https://passbolt.example.com): ", false)
	// if err != nil {
	// 	log.Fatalf("Error getting server: %v", err)
	// }
	//
	// privateKey, err := input.PromptUser("Path to Private Key: ", false)
	// if err != nil {
	// 	log.Fatalf("Error getting path to the private key: %v", err)
	// }
	//
	// password, err := input.PromptUser("Master Password (leave blank to ask each time): ", true)
	// if err != nil {
	// 	log.Fatalf("Error getting users master password: %v", err)
	// }
	//
	// fmt.Println("")
	// viper.Set("server", server)
	// viper.Set("privateKey", privateKey)
	// viper.Set("password", password)
	//
	// saveConfigFile()
}

