package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/cmdcontext"
	"github.com/chadsmith12/dotsec/config"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/chadsmith12/dotsec/secrets"
	"github.com/passbolt/go-passbolt/api"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push foldername",
	Short: "Pushes alll the secrets into your file to the secret manager",
	Long: `Pushes the secrets from the folder specified and saves them to your secret manager folder. There are two types: dotnet or env.
		dotnet - Uses dotnet user-secrets to set the secrets in your dotnet projects secrets.json file.
		env - Saves the secrets to the .env file.

		If you do not specify the --project flag, then it will attempt to use your current working directory.
		You can specify the project directory for the secrets to try to be read `,
	Example: "dotsec push FolderName --project ./api",
	Run:     pushRun,
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("project", "p", "", "The path to the dotnet project to sync the secrets to. Default to the current directory. Only valid with --type dotnet.")
	pushCmd.Flags().StringP("file", "f", ".env", "The env file you want to save the secrets to. Default to .env in the current directory. Only valid with --type env.")
	pushCmd.Flags().String("type", "", "The type of secrets file you want to use. dotnet to use dotnet user-secrets or env to use a .env file.")
}

func pushRun(cmd *cobra.Command, args []string) {
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

	cmdCtx, err := cmdcontext.NewCommandContext(cmd, projectConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create command context: %v\n", err)
		os.Exit(1)
	}

	client, err := cmdCtx.UserClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get Passbolt client: %v\n", err)
		os.Exit(1)
	}
	folder, err := client.GetFolderWithResources(projectConfig.Folder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error - Using folder: %s - %v\n", folderName, err)
		os.Exit(1)
	}

	fetcher, err := cmdCtx.SecretsFetcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get secrets fetcher: %v\n", err)
		os.Exit(1)
	}

	secretsData, err := fetcher.FetchSecrets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error - Fetching Secrets: %v\n", err)
		os.Exit(1)
	}
	pushSecrets(secretsData, client, folder)
}

func pushSecrets(secretsData []secrets.SecretData, client *passbolt.PassboltApi, folder api.Folder) {
	for _, value := range secretsData {
		if id, ok := containsSecret(folder, value.Key); ok {
			err := client.UpdateSecret(id, value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error Updating Secret %s - %v", value.Key, err)
			}
		} else {
			err := client.CreateSecretInFolder(folder.ID, value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error Creating Secret %s -%v", value.Key, err)
			}
		}
	}
}

func containsSecret(folder api.Folder, key string) (string, bool) {
	for _, resource := range folder.ChildrenResources {
		if resource.Name == key {
			return resource.ID, true
		}
	}

	return "", false
}
