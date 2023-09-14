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

// pullCmd represents the sync command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls down the secrets for a folder from passbolt",
	Long: `Pulls down the secrets from the folder specified and saves them to your projects secrets file. 

		If you do not specify the --project flag, then it will attempt to run dotnet user-secrets in your current working directory.
		You can specify the project directory to run the dotnet user-secrets in.

		Example: dotsec pull "SecretsFolder" --project ./projects/testProject/`,
	Run: syncRun,
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringP("project", "p", "", "The path to the dotnet project to sync the secrets to. Default to the current directory")
}

func syncRun(cmd *cobra.Command, args []string) {
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
