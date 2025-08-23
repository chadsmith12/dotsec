package dotnet

import (
	"github.com/chadsmith12/dotsec/secrets"
)

type DotNetFetcher struct {
	project string
}

func NewFetcher(project string) DotNetFetcher {
	return DotNetFetcher{project: project}
}

func (fetcher DotNetFetcher) FetchSecrets() ([]secrets.SecretData, error) {
	stdOut, err := ListSecrets(fetcher.project)
	if err != nil {
		return []secrets.SecretData{}, err
	}
	values, err := ParseSecrets(stdOut)
	if err != nil {
		return []secrets.SecretData{}, err
	}

	secretsData := secrets.SecretDataFromSlice(values)
	return secretsData, nil
}
