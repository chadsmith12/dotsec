package cmdcontext

import (
	"context"
	"errors"
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
	server     string
	privateKey string
	password   string
}

type CommandContext struct {
	secretsType   string
	configuration *Configuration
	client        *passbolt.PassboltApi
	cmd           *cobra.Command
}

// Initializes a new CommandContext which gives you context into the how the current command is being ran.
// This is meant to be used accross commands that need to access Passbolt and the secrets.
// Returns an error if the configuration is invalid or required flags are missing.
func NewCommandContext(cmd *cobra.Command) (*CommandContext, error) {
	envType, err := cmd.Flags().GetString("type")
	if err != nil {
		return nil, fmt.Errorf("failed to get type flag: %w", err)
	}
	if envType == "" {
		return nil, errors.New("type flag is required")
	}

	configuration, err := getConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	return &CommandContext{
		secretsType:   envType,
		configuration: configuration,
		cmd:           cmd,
	}, nil
}

// Gets the secret fetcher we are going to use get the secrets from this environment type
func (cmdContext *CommandContext) SecretsFetcher() (secrets.SecretsFetcher, error) {
	switch cmdContext.secretsType {
	case "dotnet":
		project, err := cmdContext.cmd.Flags().GetString("project")
		if err != nil {
			project = ""
		}
		return dotnet.NewFetcher(project), nil
	case "env":
		envFile, err := cmdContext.cmd.Flags().GetString("file")
		if err != nil {
			envFile = ".env"
		}
		return env.NewFetcher(envFile), nil
	default:
		return nil, fmt.Errorf("unsupported secrets type: %s", cmdContext.secretsType)
	}
}

// Gets the secrets setter we are going to use to set the secrets from this environment type
func (cmdContext *CommandContext) SecretsSetter() (secrets.SecretsSetter, error) {
	switch cmdContext.secretsType {
	case "dotnet":
		project, err := cmdContext.cmd.Flags().GetString("project")
		if err != nil {
			project = ""
		}
		return dotnet.NewSetter(project), nil
	case "env":
		envFile, err := cmdContext.cmd.Flags().GetString("file")
		if err != nil {
			envFile = ".env"
		}
		return env.NewSetter(envFile), nil
	default:
		return nil, fmt.Errorf("unsupported secrets type: %s", cmdContext.secretsType)
	}
}

// Attempts to get and initialize the passbolt api client with the user logged in
func (cmdContext *CommandContext) UserClient(ctx context.Context) (*passbolt.PassboltApi, error) {
	if cmdContext.client != nil && cmdContext.client.ValidLogin() {
		return cmdContext.client, nil
	}
	if cmdContext.client != nil {
		err := cmdContext.client.Login()
		if err != nil {
			return nil, fmt.Errorf("failed to re-login to Passbolt: %w", err)
		}
		return cmdContext.client, nil
	}

	password, err := cmdContext.Password()
	if err != nil {
		return nil, fmt.Errorf("failed to get password: %w", err)
	}

	keyData, err := os.ReadFile(cmdContext.configuration.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	client, err := passbolt.NewClient(ctx, cmdContext.configuration.server, string(keyData), password)
	if err != nil {
		return nil, fmt.Errorf("failed to create Passbolt client: %w", err)
	}

	err = client.Login()
	if err != nil {
		return nil, fmt.Errorf("failed to login to Passbolt: %w", err)
	}

	cmdContext.client = client
	return cmdContext.client, nil
}

// Attempts to get the password to unlock the users private key.
// First checks to see if we have it stored from viper in the configuration.
// If not then we will immediately prompt the user for their password.
func (cmdContext *CommandContext) Password() (string, error) {
	password := cmdContext.configuration.password
	if password != "" {
		return password, nil
	}

	password, err := input.PromptUser("Master Password: ", true)
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}
	fmt.Println()
	return password, nil
}

func getConfiguration() (*Configuration, error) {
	server := viper.GetViper().GetString("server")
	if server == "" {
		return nil, errors.New("server not configured - use configure command, --server flag, or environment variable")
	}

	privateKey := viper.GetViper().GetString("privateKey")
	if privateKey == "" {
		return nil, errors.New("privateKey not configured - use configure command, --privateKey flag, or environment variable")
	}

	password := viper.GetViper().GetString("password")

	return &Configuration{
		server:     server,
		privateKey: privateKey,
		password:   password,
	}, nil
}
