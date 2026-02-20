package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chadsmith12/dotsec/cmdcontext"
	"github.com/chadsmith12/dotsec/config"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/chadsmith12/dotsec/secrets"
	"github.com/passbolt/go-passbolt/api"
	"github.com/spf13/cobra"
)

var prune bool
var force bool

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
	pushCmd.Flags().BoolVar(&prune, "prune", false, "Delete secrets in Passbolt that don't exist locally")
	pushCmd.Flags().BoolVar(&force, "force", false, "Skip dirty check warning")
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

	if projectConfig.Team == "" {
		fmt.Fprintln(os.Stderr, "Error: Team is required. Set it in .dotsecrc or use --team flag.")
		os.Exit(1)
	}

	group, err := client.GetGroup(projectConfig.Team)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get team '%s': %v\n", projectConfig.Team, err)
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

	if !force {
		extraSecrets := findExtraSecrets(folder, secretsData)
		if len(extraSecrets) > 0 {
			fmt.Fprintf(os.Stderr, "Warning: Passbolt has %d secrets not in your local file:\n  - %s\nRun 'dotsec pull' first, or use --force to override.\n", len(extraSecrets), strings.Join(extraSecrets, "\n  - "))
			os.Exit(1)
		}
	}

	if prune {
		secretsToDelete := findSecretsToDelete(folder, secretsData)
		if len(secretsToDelete) > 0 {
			fmt.Printf("Pruning %d secrets from Passbolt:\n  - %s\n", len(secretsToDelete), strings.Join(secretsToDelete, "\n  - "))
			for _, resource := range folder.ChildrenResources {
				if containsSecretName(secretsData, resource.Name) {
					continue
				}
				if err := client.DeleteSecret(resource.ID); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to delete secret %s: %v\n", resource.Name, err)
				}
			}
		}
	}

	pushSecrets(secretsData, client, folder, group.ID)
}

func findExtraSecrets(folder api.Folder, localSecrets []secrets.SecretData) []string {
	localKeys := make(map[string]bool)
	for _, s := range localSecrets {
		localKeys[s.Key] = true
	}

	var extra []string
	for _, resource := range folder.ChildrenResources {
		if !localKeys[resource.Name] {
			extra = append(extra, resource.Name)
		}
	}
	return extra
}

func findSecretsToDelete(folder api.Folder, localSecrets []secrets.SecretData) []string {
	localKeys := make(map[string]bool)
	for _, s := range localSecrets {
		localKeys[s.Key] = true
	}

	var toDelete []string
	for _, resource := range folder.ChildrenResources {
		if !localKeys[resource.Name] {
			toDelete = append(toDelete, resource.Name)
		}
	}
	return toDelete
}

func containsSecretName(secretsData []secrets.SecretData, name string) bool {
	for _, s := range secretsData {
		if s.Key == name {
			return true
		}
	}
	return false
}

func pushSecrets(secretsData []secrets.SecretData, client *passbolt.PassboltApi, folder api.Folder, groupId string) {
	for _, value := range secretsData {
		if id, ok := containsSecret(folder, value.Key); ok {
			err := client.UpdateSecret(id, value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error Updating Secret %s - %v\n", value.Key, err)
			}
		} else {
			resourceId, err := client.CreateSecretInFolder(folder.ID, value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error Creating Secret %s - %v\n", value.Key, err)
				continue
			}

			shareErrors := client.ShareResourcesWithGroup([]string{resourceId}, groupId)
			if len(shareErrors) > 0 {
				fmt.Fprintf(os.Stderr, "Warning: Secret %s created but failed to share with group: %v\n", value.Key, shareErrors[0])
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
