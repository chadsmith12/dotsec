package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/chadsmith12/dotsec/input"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure saves the server auth details to the Config file.",
	Long: `Configure saves the server auth details to the config file to make working with the tool easier and quicker.
	If no flags or environment variables are used, then dotsec will use the config file created.`,

	Run: configureRun,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func configureRun(cmd *cobra.Command, args []string) {
	server, err := input.PromptUser("Server (https://passbolt.example.com): ", false)
	if err != nil {
		log.Fatalf("Error getting server: %v", err)
	}

	privateKey, err := input.PromptUser("Path to Private Key: ", false)
	if err != nil {
		log.Fatalf("Error getting path to the private key: %v", err)
	}

	password, err := input.PromptUser("Master Password (leave blank to ask each time): ", true)
	if err != nil {
		log.Fatalf("Error getting users master password: %v", err)
	}

	fmt.Println("")
	viper.Set("server", server)
	viper.Set("privateKey", privateKey)
	viper.Set("password", password)

	saveConfigFile()
}

func saveConfigFile() {
	// first try to just save the config in general.
	// if there is an error we will then try to see if it's possible to save to a differnet path
	err := viper.SafeWriteConfig()
	if err != nil {
		trySaveConfigAs(err)
	}
}

func trySaveConfigAs(configError error) {
	if _, ok := configError.(viper.ConfigFileAlreadyExistsError); !ok {
		fmt.Fprintf(os.Stderr, "error writing config: %v\n", configError)
		os.Exit(1)
	}

	currentConfigFile := viper.ConfigFileUsed()
	filePath, promptError := input.PromptUser(fmt.Sprintf("Save config file as (%v): ", currentConfigFile), false)
	if promptError != nil {
		fmt.Fprintf(os.Stderr, "error getting new path for config file: %v\n", promptError)
		os.Exit(1)
	}

	if filePath == "" {
		filePath = currentConfigFile
	}
	err := viper.WriteConfigAs(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error saving config to new path: %v\n", err)
		os.Exit(1)
	}
}
