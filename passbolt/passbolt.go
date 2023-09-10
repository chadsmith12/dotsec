package passbolt

import (
	"context"
	"fmt"

	"github.com/passbolt/go-passbolt/api"
)

type PassboltApi struct {
	server string
	privateKey string
	password string
	apiClient *api.Client
	context context.Context
}

// Initializes a new Passbolt Api with the context specified, with the credentails passed in.
// Returns an error if an error happens creating a client.
func NewClient(ctx context.Context, server, privateKey, password string) (*PassboltApi, error) {
	client, err := api.NewClient(nil, "", server, privateKey, password)
	if err != nil {
		return nil, fmt.Errorf("Creating Client: %w", err)
	}
	api := &PassboltApi {
		server: server,
		privateKey: privateKey,
		password: password,
		apiClient: client,
		context: ctx,
	}

	return api, nil
}

// Attempts to log the user using the client.
func (client *PassboltApi) Login() error {
	return client.apiClient.Login(client.context)
}
