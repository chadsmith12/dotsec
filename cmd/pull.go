package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/dotnet"
	"github.com/chadsmith12/dotsec/env"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validEnvironments = []string {"dotnet", "env"}

// pullCmd represents the sync command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls down the secrets for a folder from passbolt",
	Long: `Pulls down the secrets from the folder specified and saves them to your projects secrets file. There are two types: dotnet or env.
		dotnet - Uses dotnet user-secrets to set the secrets in your dotnet projects secrets.json file.
		env - Saves the secrets to the .env file. 

		If you do not specify the --project flag, then it will attempt to use your current working directory.
		You can specify the project directory for the secrets to try to be set.
	
		When using dotnet user-secrets your project will first be initialized to work with user-secrets.
		When using env a file will be created and/or replaced with the secrets downloaded.

		Example: dotsec pull "SecretsFolder" --project ./projects/testProject/
				 dotnet pull "SecretsFolder" --type env --file ".env" --project ./projects/testProject`,
	Run: pullRun,
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringP("project", "p", "", "The path to the dotnet project to sync the secrets to. Default to the current directory. Only valid with --type dotnet.")
	pullCmd.Flags().StringP("file", "f", ".env", "The env file you want to save the secrets to. Default to .env in the current directory. Only valid with --type env.")
	pullCmd.Flags().String("type", "dotnet", "The type of secrets file you want to use. dotnet to use dotnet user-secrets or env to use a .env file.")

}

func pullRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Specify a folder to download secrets from")
		os.Exit(1)
	}
	folderName := args[0]
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

	secrets, err := client.GetSecretsByFolder(folderName);
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve folder: ", err)
		os.Exit(1)
	}
	
	project, _ := cmd.Flags().GetString("project")
	envType, _ := cmd.Flags().GetString("type")
	if envType == "dotnet" {
		fmt.Fprintln(os.Stderr, "Using type of dotnet")
		dotnet.InitSecrets(project)
		for _, secret := range secrets {
			dotnet.SetSecret(project, secret.Key, secret.Value)
		}
	} else {
		envFile, _ := cmd.Flags().GetString("file")
		err := env.SetSecrets(project, envFile, secrets)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to set secrets: ", err)
			os.Exit(1)
		}
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
