package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/cmdcontext"
	"github.com/chadsmith12/dotsec/config"
	"github.com/spf13/cobra"
)

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
}

func pullRun(cmd *cobra.Command, args []string) {
	folderName := ""
	if len(args) > 0 {
		folderName = args[0]
	}
	projectConfig, err := config.LoadProjectConfig(cmd, folderName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmdContext, err := cmdcontext.NewCommandContext(cmd, projectConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create command context: %v\n", err)
		os.Exit(1)
	}

	client, err := cmdContext.UserClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get Passbolt client: %v\n", err)
		os.Exit(1)
	}

	secrets, err := client.GetSecretsByFolder(projectConfig.Folder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve folder: ", err)
		os.Exit(1)
	}

	setter, err := cmdContext.SecretsSetter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get secrets setter: %v\n", err)
		os.Exit(1)
	}

	err = setter.SetSecrets(secrets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set secrets: %v\n", err)
		os.Exit(1)
	}
}
