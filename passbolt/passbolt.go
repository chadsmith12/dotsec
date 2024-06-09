package passbolt

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/chadsmith12/dotsec/secrets"
	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
)

var (
	InvalidFolderErr = fmt.Errorf("failed to find folder")
)

type PassboltApi struct {
	server     string
	privateKey string
	password   string
	apiClient  *api.Client
	context    context.Context
}

type resourceResult struct {
	secretData secrets.SecretData
	err        error
}


// Initializes a new Passbolt Api with the context specified, with the credentails passed in.
// Returns an error if an error happens creating a client.
func NewClient(ctx context.Context, server, privateKey, password string) (*PassboltApi, error) {
	client, err := api.NewClient(nil, "", server, privateKey, password)
	if err != nil {
		return nil, fmt.Errorf("Creating Client: %w", err)
	}
	api := &PassboltApi{
		server:     server,
		privateKey: privateKey,
		password:   password,
		apiClient:  client,
		context:    ctx,
	}

	return api, nil
}

// Attempts to log the usar using the client.
func (client *PassboltApi) Login() error {
	return client.apiClient.Login(client.context)
}

// Checks to see if the user has a valid session
func (client *PassboltApi) ValidLogin() bool {
	return client.apiClient.CheckSession(client.context)
}

func (client *PassboltApi) GetSecretsByFolder(folderName string) ([]secrets.SecretData, error) {
	folder, err := client.GetFolderWithResources(folderName)
	secretData := make([]secrets.SecretData, 0)
	if err != nil {
		return secretData, err
	}

	client.populateSecrets(folder.ChildrenResources, &secretData)

	return secretData, nil
}

func (client *PassboltApi) GetFolderWithResources(folderName string) (api.Folder, error) {
	folders, err := client.apiClient.GetFolders(client.context, &api.GetFoldersOptions{
		FilterSearch:             folderName,
		ContainChildrenResources: true,
	})
	
	if err != nil {
		return api.Folder{}, err
	}

	for _, folder := range folders {
		if strings.EqualFold(folder.Name, folderName) {
			return folder, nil
		}
	}

	return api.Folder{}, InvalidFolderErr
}

func (client *PassboltApi) CreateSecretInFolder(folderId string, secret secrets.SecretData) error {
	_, err := helper.CreateResource(client.context, client.apiClient, folderId, secret.Key, "", "", secret.Value, "")

	return err
}

func (client *PassboltApi) UpdateSecret(resourceId string, secret secrets.SecretData) error {
	err := helper.UpdateResource(client.context, client.apiClient, resourceId, "", "", "", secret.Value, "")

	return err
}

func (client *PassboltApi) populateSecrets(resources []api.Resource, secrets *[]secrets.SecretData) {
	if len(resources) == 0 {
		return
	}
	ch := make(chan resourceResult)
	var wg sync.WaitGroup
	for _, resource := range resources {
		wg.Add(1)
		go client.downloadResource(resource, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		if result.err == nil {
			*secrets = append(*secrets, result.secretData)
		}
	}
}

func (client *PassboltApi) downloadResource(resource api.Resource, ch chan<- resourceResult, wg *sync.WaitGroup) {
	defer wg.Done()
	_, name, _, _, password, _, err := helper.GetResource(client.context, client.apiClient, resource.ID)
	if err != nil {
		secretData := secrets.SecretData{Key: "", Value: ""}
		ch <- resourceResult{secretData: secretData, err: err}
		return
	}

	ch <- resourceResult{secretData: secrets.SecretData{Key: name, Value: password}, err: nil}
}
