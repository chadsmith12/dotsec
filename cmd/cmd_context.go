package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chadsmith12/dotsec/dotnet"
	"github.com/chadsmith12/dotsec/env"
	"github.com/chadsmith12/dotsec/input"
	"github.com/chadsmith12/dotsec/passbolt"
	"github.com/chadsmith12/dotsec/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configuration struct {
	server string
	privateKey string
	password string
}

type CommandContext struct {
	secretsType string
	configuration *Configuration
	client *passbolt.PassboltApi
	cmd *cobra.Command
}

// Initializes a new CommandContext which gives you context into the how the current command is being ran.
// This is meant to be used accross commands that need to access Passbolt and the secrets.
// This assumes that tool has already been configured correctly. If not then it will crash with an exit code of 1.
func NewCommandContext(cmd *cobra.Command) *CommandContext {
	envType, _ := cmd.Flags().GetString("type")
	return &CommandContext{
		secretsType: envType,
		configuration: getConfiguration(),
		cmd: cmd,
	}	
}

// Gets the secret fetcher we are going to use get the secrets from this environment type
func (cmdContext *CommandContext) SecretsFetcher() secrets.SecretsFetcher {
	if cmdContext.secretsType == "dotnet" {
		project, _ := cmdContext.cmd.Flags().GetString("project") 
		return dotnet.NewFetcher(project)
	}
	
	envFile, _ := cmdContext.cmd.Flags().GetString("file")
	return env.NewFetcher(envFile)
}

// Gets the secrets setter we are going to use to set the secrets from this environment type
func (cmdContext *CommandContext) SecretsSetter() secrets.SecretsSetter {
	if cmdContext.secretsType == "dotnet" {
		project, _ := cmdContext.cmd.Flags().GetString("project")
		return dotnet.NewSetter(project)
	}

	envFile, _ := cmdContext.cmd.Flags().GetString("file")
	return env.NewSetter(envFile)
}

// Attempts to get and initialize the passbolt api client with the user logged in
func (cmdContext *CommandContext) UserClient(ctx context.Context) *passbolt.PassboltApi {
	if cmdContext.client != nil && cmdContext.client.ValidLogin() {
		return cmdContext.client
	}
	if cmdContext.client != nil {
		err := cmdContext.client.Login()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to login the Password with error: %v\n", err)
			os.Exit(1)
		}
	}

	password := cmdContext.Password()
	keyData, err := os.ReadFile(cmdContext.configuration.privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read key file with error: %v\n", err)
	}

	client, err := passbolt.NewClient(ctx, cmdContext.configuration.server, string(keyData), password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Passbolt client with error: %v\n", err)
		os.Exit(1)
	}

	cmdContext.client = client
	err = cmdContext.client.Login()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to login the Password with error: %v\n", err)
		os.Exit(1)
	}

	return cmdContext.client
}

// Attempts to get the password to unlock the users private key.
// First checks to see if we have it stored from viper in the configuration.
// If not then we will immediately prompt the user for their password.
// If the prompt fails for some reason, we immediately quit the application as we can't do anything else.
func (cmdContext *CommandContext) Password() string {
	password := cmdContext.configuration.password
	var err error
	if password == "" {
		password, err = input.PromptUser("Master Password: ", true)
		fmt.Println()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed getting the password to unlock private key with error - %v\n", err)
		os.Exit(1)
	}

	return password
}

func getConfiguration() (*Configuration) {
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

	return &Configuration{
		server: server,
		privateKey: privateKey,
		password: password,
	}
}
