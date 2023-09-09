/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotsec",
	Short: "Quickly download secrets for your dotnet applications",
	Long: `dotsec allows developers to quickly download their secrets and place them in their dotnet projects secrets.json file for local development.

	The tool is designed to be used by teams where it can be difficult to keep and share development secrets across the team. Store them in the Passbolt password manager and quickly download and place them in your secrets.json file`,
}

// This is called by main.main().
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dotsec.yaml)")
}


