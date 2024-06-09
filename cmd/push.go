package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chadsmith12/dotsec/dotnet"
	"github.com/chadsmith12/dotsec/env"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/chadsmith12/dotsec/secretsfetcher"
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
	Run: pushRun,
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("project", "p", "", "The path to the dotnet project to sync the secrets to. Default to the current directory. Only valid with --type dotnet.")
	pushCmd.Flags().StringP("file", "f", ".env", "The env file you want to save the secrets to. Default to .env in the current directory. Only valid with --type env.")
	pushCmd.Flags().String("type", "dotnet", "The type of secrets file you want to use. dotnet to use dotnet user-secrets or env to use a .env file.")
}

func pushRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Specify a folder to download secrets from")
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := getClient(ctx)
	folderName := args[0]
	folder, err := client.GetFolderWithResources(folderName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error - Using folder: %s - %v\n", folderName, err)
		os.Exit(1)
	}

	envType, _ := cmd.Flags().GetString("type")
	fetcher := getSecretsFetcher(cmd, envType)
	secretsData, err := fetcher.FetchSecrets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error - Fetching Secrets: %v\n", err)
	}
	pushSecrets(secretsData, client, folder)
}

func getSecretsFetcher(cmd *cobra.Command, envType string) secretsfetcher.SecretsFetcher {
	if envType == "dotnet" {
		project, _ := cmd.Flags().GetString("project")
		return dotnet.NewFetcher(project)
	}

	envFile, _ := cmd.Flags().GetString("file")
	return env.NewFetcher(envFile)
}

func pushSecrets(secretsData []passbolt.SecretData, client *passbolt.PassboltApi, folder api.Folder) {
	for _, value := range secretsData {
		if id, ok := containsSecret(folder, value.Key); ok {
			client.UpdateSecret(id, value)
		} else {
			client.CreateSecretInFolder(folder.ID, value)
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

func getClient(ctx context.Context) *passbolt.PassboltApi {
	server, keyFile, password := getConfiguration()
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

	return client
}
