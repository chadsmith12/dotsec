/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/dotnet"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: syncRun,
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringP("project", "p", "", "The path to the dotnet project to sync the secrets to. Default to the current directory")
	syncCmd.Flags().String("folder", "", "The folder we want to pull the secrets from")
	syncCmd.MarkFlagRequired("folder")
}

func syncRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	server, keyFile, password := checkConfiguration()
	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read key file: ", err)
		os.Exit(1)
	}

	client, err := passbolt.NewClient(ctx, server, string(keyData), password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create api client: ", err)
		os.Exit(1)
	}

	err = client.Login()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to login to Passbolt: ", err)
		os.Exit(1)
	}

	folderName, err := cmd.Flags().GetString("folder")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find the folder flag. Did you use --folder? %s\n", err)
		os.Exit(1)
	}
	secrets, err := client.GetSecretsByFolder(folderName);
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve folder: ", err)
		os.Exit(1)
	}
	
	project, _ := cmd.Flags().GetString("project")
	for _, secret := range secrets {
		dotnet.SetSecret(project, secret.Key, secret.Value)
	}
}

func checkConfiguration() (string, string, string) {
	server:= viper.GetViper().GetString("server")
	if server == "" {
		fmt.Fprint(os.Stderr, "Server is not configured. Run the configure command, use the --server flag, or environment variable to set the server\n")
		os.Exit(1)
	}

	privateKey := viper.GetViper().GetString("privateKey")
	if privateKey == "" {
		fmt.Fprint(os.Stderr, "privateKey is not configured. Run the configure command, use the --privateKey flag, or environment variable to set it to a valid private key file to load.\n")
		os.Exit(1)
	}

	password := viper.GetViper().GetString("password")
	if password == "" {
		fmt.Fprint(os.Stderr, "password is not configured. Run the configure command, use the --password flag, or the environment variable to set the master password to use.\n")
		os.Exit(1)
	}

	return server, privateKey, password
}
