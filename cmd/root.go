/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotsec",
	Short: "Quickly download secrets for your projects locally from Passbolt",
	Long: `dotsec allows developers to quickly download their secrets and place them in their dotnet projects secrets.json file or their .env file for local development.

	The tool is designed to be used by teams where it can be difficult to keep and share development secrets across the team. Store them in the Passbolt password manager and quickly download and place them in your secrets.json file or a .env file. When using dotnet dotsec will use run the dotnet user-secrets command, and using a .env will parse and save the secrets to a .env file.`,
}

// This is called by main.main().
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file for dotsec to read information from.")
	rootCmd.PersistentFlags().String("server", "", "Passbolt Server to use (https://passbolt.example.com)")
	rootCmd.PersistentFlags().String("privateKey", "", "Passbolt User Private Key")
	rootCmd.PersistentFlags().String("password", "", "Passbolt User Password")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("privateKey", rootCmd.PersistentFlags().Lookup("privateKey"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// find the XDG config directory
		configDir := xdg.ConfigHome
		configDir = filepath.Join(configDir, "dotsec")
		os.MkdirAll(configDir, 0700)
		viper.SetConfigPermissions(os.FileMode(0600))
		viper.AddConfigPath(configDir)
		viper.SetConfigType("json")
		viper.SetConfigName(".config")
	}
	
	// read in the environment variables that match and use those
	viper.SetEnvPrefix("dotsec")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Failed to read config with error: %v", err)
			os.Exit(1)
		}
	}
}

