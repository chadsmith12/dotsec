/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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

	password, err := input.PromptUser("Master Password: ", true)
	if err != nil {
		log.Fatalf("Error getting users master password: %v", err)
	}

	fmt.Println("")
	viper.Set("server", server)
	viper.Set("privateKey", privateKey)
	viper.Set("password", password)
	
	configErr := viper.SafeWriteConfig()
	if configErr != nil {
		fmt.Fprintf(os.Stderr, "Error Writing config: %v", configErr)
		os.Exit(1)
	}
}
