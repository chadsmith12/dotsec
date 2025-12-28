package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/chadsmith12/dotsec/colors"
	"github.com/chadsmith12/dotsec/config"
	"github.com/chadsmith12/dotsec/input"
	"github.com/spf13/cobra"
)

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
	var folder string
	for folder == "" {
		val, err := input.PromptUser(colors.Yellow("Passbolt Folder Name: "), false)
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		folder = strings.TrimSpace(val)
		if folder == "" {
			fmt.Println(colors.Red("Folder name is required"))
		}
	}

	var secretType string
	for secretType == "" {
		val, err := input.PromptUser(colors.Yellow("Secret type (dotnet/env) [env]: "), false)
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		val = strings.TrimSpace(val)
		if val == "" {
			secretType = "env"
		} else if val == "dotnet" || val == "env" {
			secretType = val
		} else {
			fmt.Println(colors.Red("Invalid type. Please enter 'dotnet' or 'env'"))
		}
	}

	defaultPath := ".env"
	if secretType == "dotnet" {
		defaultPath = "."
	}
	val, err := input.PromptUser(colors.Yellow(fmt.Sprintf("Path [%s]: ", defaultPath)), false)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	path := strings.TrimSpace(val)
	if path == "" {
		path = defaultPath
	}

	fmt.Printf(colors.Cyan("\nConfiguration:\n Folder: %s\n Type: %s\n Path: %s\n\n"), folder, secretType, path)
	confirm, _ := input.PromptUser(colors.Yellow("Save to .dotsecrc? [Y/n]: "), false)
	if strings.ToLower(strings.TrimSpace(confirm)) == "n" {
		fmt.Println("Cancelled")
		return
	}

	if err := config.WriteProjectConfigWithData(folder, secretType, path); err != nil {
		log.Fatalf("Error saving config: %v", err)
	}

	fmt.Println(colors.Green("Configuration saved to .dotsecrc"))
}
